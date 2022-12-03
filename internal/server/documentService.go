package server

import (
	"log"
	"regexp"
	"sync"
)

type DocumentService struct {
	repo            *Repository
	recaptchaClient *RecaptchaClient
	sessionTable    *SessionTable
	errorLogger     *log.Logger

	lockByUser map[string]*sync.RWMutex
	lockMux    sync.RWMutex
}

var usernameMatchesRegex = regexp.MustCompile(`^[a-zA-Z0-9 \@\.\!\-\_\$\+]+$`).MatchString

func isUsernameValid(username string) bool {
	return usernameMatchesRegex(username)
}

func (s *DocumentService) checkUserNameFree(auth Auth) error {
	exists, err := s.repo.saltFileExists(auth)
	s.logError(err)

	if exists {
		return ERROR_ALREADY_EXISTS
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
	session, err := s.sessionTable.GetSession(context.sessionToken, context.username, context.ip)
	if err != nil {
		return Auth{}, ERROR_INVALID_CREDENTIALS
	}

	return context.ToAuth(session.DecryptToken), nil
}

func (s *DocumentService) checkAuth(auth Auth) ([]byte, error) {
	salt, err := s.saltForUser(auth)
	if err != nil {
		return nil, err
	}

	err = s.checkDocumentExists(auth, salt)
	if err != nil {
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

func (s *DocumentService) CreateDocument(context RequestContext, document string, recaptchaResponse string) (string, error) {
	auth := context.ToAuth(context.password)

	// Validate username
	if !isUsernameValid(auth.username) {
		return "", ERROR_USERNAME_INVALID
	}

	// Validate recaptcha
	recaptchaValid := s.recaptchaClient.Verify(recaptchaResponse, auth.ip)
	if !recaptchaValid {
		s.errorLogger.Println("Recaptcha was not valid")
		return "", ERROR_INVALID_CREDENTIALS
	}

	userMux := s.getLockForUser(auth.username)
	userMux.Lock()
	defer userMux.Unlock()

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
		return "", err
	}
	if err != nil {
		s.logError(err)
		return "", ERROR_SERVER_ERROR
	}

	return session.SessionToken, nil
}

func (s *DocumentService) UpdateDocument(context RequestContext, document string) error {
	// Validate username
	if !isUsernameValid(context.username) {
		return ERROR_INVALID_CREDENTIALS
	}

	// Validate session
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
	// Validate username
	if !isUsernameValid(context.username) {
		return nil, ERROR_INVALID_CREDENTIALS
	}

	auth, err := s.authFromSession(context)
	if err != nil {
		return nil, err
	}

	userMux := s.getLockForUser(auth.username)
	userMux.RLock()
	defer userMux.RUnlock()

	salt, err := s.checkAuth(auth)
	if err != nil {
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
	// Validate username
	if !isUsernameValid(context.username) {
		return ERROR_INVALID_CREDENTIALS
	}

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
	// Validate username
	if !isUsernameValid(context.username) {
		return "", ERROR_INVALID_CREDENTIALS
	}

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
