package server

import (
	"fmt"
	"log"
	"sync"
)

type DocumentService struct {
	repo         *Repository
	lockoutTable *LockoutTable
	sessionTable *SessionTable
	errorLogger  *log.Logger

	lockByUser map[string]*sync.RWMutex
	lockMux    sync.RWMutex
}

func (s *DocumentService) checkUserNameFree(auth Auth) error {
	allowed := s.lockoutTable.shouldAllow(auth.ip, "")
	if !allowed {
		return ERROR_TOO_MANY_REQUESTS
	}

	exists, err := s.repo.saltFileExists(auth)
	s.logError(err)

	if exists {
		s.lockoutTable.logFailure(auth.ip, "")
		return ERROR_INVALID_USER_NAME
	}

	return nil
}

func (s *DocumentService) saltForUser(auth Auth) ([]byte, error) {
	salt, err := s.repo.readSaltFile(auth)

	if err != nil {
		s.logError(err)
		return nil, ERROR_INVALID_CREDENTIALS
	}

	return salt, nil
}

func (s *DocumentService) checkDocumentExists(auth Auth, salt []byte) error {
	exists, err := s.repo.documentExists(auth, salt)
	s.logError(err)

	if !exists {
		return ERROR_INVALID_CREDENTIALS
	}

	return nil
}

func (s *DocumentService) authFromSession(context RequestContext) (Auth, error) {
	// TODO: check if lockout table should allow?
	fmt.Println("Looking up session " + context.sessionToken)
	session, err := s.sessionTable.GetSession(context.sessionToken, context.username, context.ip)
	if err != nil {
		s.lockoutTable.logFailure(context.ip, context.username)
		return Auth{}, ERROR_INVALID_CREDENTIALS
	}

	return context.ToAuth(session.DecryptToken), nil
}

func (s *DocumentService) checkAuth(auth Auth) ([]byte, error) {
	allowed := s.lockoutTable.shouldAllow(auth.ip, auth.username)
	if !allowed {
		return nil, ERROR_TOO_MANY_REQUESTS
	}

	// TODO: validate session here

	salt, err := s.saltForUser(auth)
	if err != nil {
		s.lockoutTable.logFailure(auth.ip, auth.username)
		return nil, err
	}

	err = s.checkDocumentExists(auth, salt)
	if err != nil {
		s.lockoutTable.logFailure(auth.ip, auth.username)
		return nil, err
	}

	return salt, nil
}

func (s *DocumentService) createSalt(auth Auth) ([]byte, error) {
	salt, err := s.repo.writeSaltFile(auth)
	s.logError(err)

	if err != nil {
		s.logError(err)
		return nil, ERROR_SERVER_ERROR
	}

	return salt, nil
}

func (s *DocumentService) CreateDocument(context RequestContext, document string) (string, error) {
	auth := context.ToAuth(context.password)

	// TODO: check recaptcha??

	userMux := s.getLockForUser(auth.username)
	userMux.Lock()
	defer userMux.Unlock()

	// TODO: more validation on username?
	err := s.checkUserNameFree(auth)
	if err != nil {
		return "", err
	}

	salt, err := s.createSalt(auth)
	if err != nil {
		return "", err
	}

	err = s.repo.writeDocument([]byte(document), auth, salt)
	if err != nil {
		s.logError(err)
		return "", ERROR_SERVER_ERROR
	}

	// Grant session for new user
	session, err := s.sessionTable.Grant(auth.username, auth.ip, auth.password)
	if err != nil {
		s.errorLogger.Println("Error granting user session")
		s.lockoutTable.logFailure(auth.ip, auth.username)
		return "", err
	}
	if err != nil {
		s.logError(err)
		return "", ERROR_SERVER_ERROR
	}

	return session.SessionToken, nil
}

func (s *DocumentService) UpdateDocument(context RequestContext, document string) error {
	auth, err := s.authFromSession(context)
	if err != nil {
		return err
	}

	userMux := s.getLockForUser(auth.username)
	userMux.Lock()
	defer userMux.Unlock()

	salt, err := s.checkAuth(auth)
	if err != nil {
		return err
	}

	err = s.repo.writeDocument([]byte(document), auth, salt)
	if err != nil {
		s.logError(err)
		return ERROR_SERVER_ERROR
	}

	return nil
}

func (s *DocumentService) GetDocument(context RequestContext) ([]byte, error) {
	auth, err := s.authFromSession(context)
	if err != nil {
		fmt.Println("Was not able to get auth from session: " + err.Error())
		return nil, err
	}

	userMux := s.getLockForUser(auth.username)
	userMux.RLock()
	defer userMux.RUnlock()

	salt, err := s.checkAuth(auth)
	if err != nil {
		fmt.Println("Auth does not appear to be valid: " + err.Error())
		return nil, err
	}

	document, err := s.repo.readDocument(auth, salt)
	if err != nil {
		s.logError(err)
		return nil, ERROR_SERVER_ERROR
	}

	return document, nil
}

func (s *DocumentService) DeleteDocument(context RequestContext) error {
	auth, err := s.authFromSession(context)
	if err != nil {
		return err
	}

	userMux := s.getLockForUser(auth.username)
	userMux.Lock()
	defer userMux.Unlock()

	salt, err := s.checkAuth(auth)
	if err != nil {
		return err
	}

	err = s.repo.deleteDocument(auth, salt)
	if err != nil {
		s.logError(err)
		return ERROR_SERVER_ERROR
	}

	err = s.repo.deleteSaltFile(auth)
	if err != nil {
		s.logError(err)
		return ERROR_SERVER_ERROR
	}

	return nil
}

func (s *DocumentService) UpdateDocumentAndPassword(context RequestContext, newPassword string, document string) (string, error) {
	auth, err := s.authFromSession(context)
	if err != nil {
		return "", err
	}

	userMux := s.getLockForUser(auth.username)
	userMux.Lock()
	defer userMux.Unlock()

	salt, err := s.saltForUser(auth)
	if err != nil {
		return "", err
	}

	err = s.checkDocumentExists(auth, salt)
	if err != nil {
		return "", err
	}

	newAuth := Auth{auth.username, string(newPassword), auth.ip}
	err = s.repo.moveDocument(auth, salt, newAuth)
	if err != nil {
		s.logError(err)
		return "", ERROR_SERVER_ERROR
	}

	err = s.repo.writeDocument([]byte(document), newAuth, salt)
	if err != nil {
		s.logError(err)
		return "", ERROR_SERVER_ERROR
	}

	// Grant session for new user
	session, err := s.sessionTable.Grant(newAuth.username, newAuth.ip, newAuth.password)
	if err != nil {
		s.errorLogger.Println("Error granting user session")
		s.lockoutTable.logFailure(auth.ip, auth.username)
		return "", err
	}
	if err != nil {
		s.logError(err)
		return "", ERROR_SERVER_ERROR
	}

	return session.SessionToken, nil
}

func (s *DocumentService) logError(err error) {
	if err != nil {
		s.errorLogger.Printf("%s\n", err.Error())
	}
}

func (s *DocumentService) getLockForUser(username string) *sync.RWMutex {
	s.lockMux.Lock()
	defer s.lockMux.Unlock()

	if s.lockByUser[username] == nil {
		s.lockByUser[username] = &sync.RWMutex{}
	}

	return s.lockByUser[username]
}
