package service

import (
	"log"
	"regexp"
	"sync"

	"github.com/ClintonMorrison/lorikeet/internal/errors"
	"github.com/ClintonMorrison/lorikeet/internal/model"
	"github.com/ClintonMorrison/lorikeet/internal/server/recaptcha"
	"github.com/ClintonMorrison/lorikeet/internal/server/repository"
	"github.com/ClintonMorrison/lorikeet/internal/server/session"
)

type DocumentService struct {
	repo            *repository.V1
	recaptchaClient *recaptcha.Client
	sessionTable    *session.Table
	errorLogger     *log.Logger

	lockByUser map[string]*sync.RWMutex
	lockMux    sync.RWMutex
}

func NewDocumentService(
	repository *repository.V1,
	recaptchaClient *recaptcha.Client,
	sessionTable *session.Table,
	errorLogger *log.Logger,
) *DocumentService {
	return &DocumentService{
		repo:            repository,
		recaptchaClient: recaptchaClient,
		sessionTable:    sessionTable,
		errorLogger:     errorLogger,
		lockByUser:      make(map[string]*sync.RWMutex),
	}
}

var usernameMatchesRegex = regexp.MustCompile(`^[a-zA-Z0-9 \@\.\!\-\_\$\+]+$`).MatchString

func isUsernameValid(username string) bool {
	return usernameMatchesRegex(username)
}

func (s *DocumentService) checkUserNameFree(auth model.Auth) error {
	exists, err := s.repo.SaltFileExists(auth)
	s.logError(err)

	if exists {
		return errors.ALREADY_EXISTS
	}

	return nil
}

func (s *DocumentService) saltForUser(auth model.Auth) ([]byte, error) {
	salt, err := s.repo.ReadSaltFile(auth)

	if err != nil {
		s.logError(err)
		return nil, errors.INVALID_CREDENTIALS
	}

	return salt, nil
}

func (s *DocumentService) checkDocumentExists(auth model.Auth, salt []byte) error {
	exists, err := s.repo.DocumentExists(auth, salt)
	s.logError(err)

	if !exists {
		return errors.INVALID_CREDENTIALS
	}

	return nil
}

func (s *DocumentService) authFromSession(context model.RequestContext) (model.Auth, error) {
	session, err := s.sessionTable.GetSession(context.SessionToken, context.Username, context.Ip)
	if err != nil {
		return model.Auth{}, errors.INVALID_CREDENTIALS
	}

	return context.ToAuth(session.DecryptToken), nil
}

func (s *DocumentService) checkAuth(auth model.Auth) ([]byte, error) {
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

func (s *DocumentService) createSalt(auth model.Auth) ([]byte, error) {
	salt, err := s.repo.WriteSaltFile(auth)
	s.logError(err)

	if err != nil {
		s.logError(err)
		return nil, errors.SERVER_ERROR
	}

	return salt, nil
}

func (s *DocumentService) CreateDocument(context model.RequestContext, document string, recaptchaResponse string) (string, error) {
	auth := context.ToAuth(context.Password)

	// Validate username
	if !isUsernameValid(auth.Username) {
		return "", errors.USERNAME_INVALID
	}

	// Validate recaptcha
	recaptchaValid := s.recaptchaClient.Verify(recaptchaResponse, auth.Ip)
	if !recaptchaValid {
		s.errorLogger.Println("Recaptcha was not valid")
		return "", errors.INVALID_CREDENTIALS
	}

	userMux := s.getLockForUser(auth.Username)
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

	err = s.repo.WriteDocument([]byte(document), auth, salt)
	if err != nil {
		s.logError(err)
		return "", errors.SERVER_ERROR
	}

	// Grant session for new user
	session, err := s.sessionTable.Grant(auth.Username, auth.Ip, auth.Password)
	if err != nil {
		s.errorLogger.Println("Error granting user session")
		return "", err
	}
	if err != nil {
		s.logError(err)
		return "", errors.SERVER_ERROR
	}

	return session.SessionToken, nil
}

func (s *DocumentService) UpdateDocument(context model.RequestContext, document string) error {
	// Validate username
	if !isUsernameValid(context.Username) {
		return errors.INVALID_CREDENTIALS
	}

	// Validate session
	auth, err := s.authFromSession(context)
	if err != nil {
		return err
	}

	userMux := s.getLockForUser(auth.Username)
	userMux.Lock()
	defer userMux.Unlock()

	salt, err := s.checkAuth(auth)
	if err != nil {
		return err
	}

	err = s.repo.WriteDocument([]byte(document), auth, salt)
	if err != nil {
		s.logError(err)
		return errors.SERVER_ERROR
	}

	return nil
}

func (s *DocumentService) GetDocument(context model.RequestContext) ([]byte, error) {
	// Validate username
	if !isUsernameValid(context.Username) {
		return nil, errors.INVALID_CREDENTIALS
	}

	auth, err := s.authFromSession(context)
	if err != nil {
		return nil, err
	}

	userMux := s.getLockForUser(auth.Username)
	userMux.RLock()
	defer userMux.RUnlock()

	salt, err := s.checkAuth(auth)
	if err != nil {
		return nil, err
	}

	document, err := s.repo.ReadDocument(auth, salt)
	if err != nil {
		s.logError(err)
		return nil, errors.SERVER_ERROR
	}

	return document, nil
}

func (s *DocumentService) DeleteDocument(context model.RequestContext) error {
	// Validate username
	if !isUsernameValid(context.Username) {
		return errors.INVALID_CREDENTIALS
	}

	auth, err := s.authFromSession(context)
	if err != nil {
		return err
	}

	userMux := s.getLockForUser(auth.Username)
	userMux.Lock()
	defer userMux.Unlock()

	salt, err := s.checkAuth(auth)
	if err != nil {
		return err
	}

	err = s.repo.DeleteDocument(auth, salt)
	if err != nil {
		s.logError(err)
		return errors.SERVER_ERROR
	}

	err = s.repo.DeleteSaltFile(auth)
	if err != nil {
		s.logError(err)
		return errors.SERVER_ERROR
	}

	return nil
}

func (s *DocumentService) UpdateDocumentAndPassword(context model.RequestContext, newPassword string, document string) (string, error) {
	// Validate username
	if !isUsernameValid(context.Username) {
		return "", errors.INVALID_CREDENTIALS
	}

	auth, err := s.authFromSession(context)
	if err != nil {
		return "", err
	}

	userMux := s.getLockForUser(auth.Username)
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

	newAuth := model.Auth{Username: auth.Username, Password: string(newPassword), Ip: auth.Ip}
	err = s.repo.MoveDocument(auth, salt, newAuth)
	if err != nil {
		s.logError(err)
		return "", errors.SERVER_ERROR
	}

	err = s.repo.WriteDocument([]byte(document), newAuth, salt)
	if err != nil {
		s.logError(err)
		return "", errors.SERVER_ERROR
	}

	// Grant session for new user
	session, err := s.sessionTable.Grant(newAuth.Username, newAuth.Ip, newAuth.Password)
	if err != nil {
		s.errorLogger.Println("Error granting user session")
		return "", err
	}
	if err != nil {
		s.logError(err)
		return "", errors.SERVER_ERROR
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
