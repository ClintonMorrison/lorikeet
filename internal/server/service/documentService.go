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
	available, err := s.repo.IsUsernameAvailable(auth)
	s.logError(err)

	if !available {
		return errors.ALREADY_EXISTS
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

	_, err = s.repo.CreateUser(auth, []byte(document))
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

	user, err := s.repo.GetUser(auth)
	if err != nil {
		return errors.INVALID_CREDENTIALS
	}

	_, err = s.repo.UpdateUser(user, model.UserUpdate{Document: []byte(document)})
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

	user, err := s.repo.GetUser(auth)
	if err != nil {
		return nil, errors.INVALID_CREDENTIALS
	}

	return user.Document, nil
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

	user, err := s.repo.GetUser(auth)
	if err != nil {
		return errors.INVALID_CREDENTIALS
	}

	err = s.repo.DeleteUser(user)
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

	user, err := s.repo.GetUser(auth)
	if err != nil {
		return "", errors.INVALID_CREDENTIALS
	}

	user, err = s.repo.UpdateUser(user, model.UserUpdate{Password: newPassword, Document: []byte(document)})
	if err != nil {
		s.logError(err)
		return "", errors.SERVER_ERROR
	}

	// Grant session for new user
	session, err := s.sessionTable.Grant(user.Username, user.Auth.Ip, user.Auth.Password)
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
