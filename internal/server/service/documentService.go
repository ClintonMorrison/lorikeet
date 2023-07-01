package service

import (
	"log"
	"regexp"

	"github.com/ClintonMorrison/lorikeet/internal/errors"
	"github.com/ClintonMorrison/lorikeet/internal/model"
	"github.com/ClintonMorrison/lorikeet/internal/server/recaptcha"
	"github.com/ClintonMorrison/lorikeet/internal/server/repository"
	"github.com/ClintonMorrison/lorikeet/internal/server/session"
)

type Document struct {
	Data                 []byte
	Salt                 []byte
	StorageVersion       int
	ClientEncryptVersion int
}

func adaptUserToDocument(user *model.User) Document {
	return Document{
		Data:                 user.Document,
		Salt:                 user.ClientSalt,
		StorageVersion:       user.Metadata.StorageVersion,
		ClientEncryptVersion: user.Metadata.ClientEncryptVersion,
	}
}

type DocumentService struct {
	repo            repository.UserRepository
	recaptchaClient *recaptcha.Client
	sessionTable    *session.Table
	userLockTable   *UserLockTable
	errorLogger     *log.Logger
}

func NewDocumentService(
	repository repository.UserRepository,
	recaptchaClient *recaptcha.Client,
	sessionTable *session.Table,
	userLockTable *UserLockTable,
	errorLogger *log.Logger,
) *DocumentService {
	return &DocumentService{
		repo:            repository,
		recaptchaClient: recaptchaClient,
		sessionTable:    sessionTable,
		userLockTable:   userLockTable,
		errorLogger:     errorLogger,
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

func (s *DocumentService) CreateDocument(context model.RequestContext, document string, recaptchaResponse string) (Document, string, error) {
	auth := context.ToAuth(context.Password)

	// Validate username
	if !isUsernameValid(auth.Username) {
		return Document{}, "", errors.USERNAME_INVALID
	}

	// Validate recaptcha
	recaptchaValid := s.recaptchaClient.Verify(recaptchaResponse, auth.Ip)
	if !recaptchaValid {
		s.errorLogger.Println("Recaptcha was not valid")
		return Document{}, "", errors.INVALID_CREDENTIALS
	}

	// Validate user name is free
	s.userLockTable.Lock(auth.Username)
	defer s.userLockTable.Unlock(auth.Username)

	err := s.checkUserNameFree(auth)
	if err != nil {
		return Document{}, "", err
	}

	// Create the user
	user, err := s.repo.CreateUser(auth, []byte(document))
	if err != nil {
		s.logError(err)
		return Document{}, "", errors.SERVER_ERROR
	}

	// Grant session for new user
	session, err := s.sessionTable.Grant(auth.Username, auth.Ip, auth.Password)
	if err != nil {
		s.errorLogger.Println("Error granting user session")
		return Document{}, "", err
	}
	if err != nil {
		s.logError(err)
		return Document{}, "", errors.SERVER_ERROR
	}

	return adaptUserToDocument(user), session.SessionToken, nil
}

func (s *DocumentService) UpdateDocument(context model.RequestContext, document string) (Document, error) {
	// Validate username
	if !isUsernameValid(context.Username) {
		return Document{}, errors.INVALID_CREDENTIALS
	}

	// Validate session
	auth, err := s.authFromSession(context)
	if err != nil {
		return Document{}, err
	}

	s.userLockTable.Lock(auth.Username)
	defer s.userLockTable.Unlock(auth.Username)

	user, err := s.repo.GetUser(auth)
	if err != nil {
		return Document{}, errors.INVALID_CREDENTIALS
	}

	user, err = s.repo.UpdateUser(user, model.UserUpdate{Document: []byte(document)})
	if err != nil {
		s.logError(err)
		return Document{}, errors.SERVER_ERROR
	}

	return adaptUserToDocument(user), nil
}

func (s *DocumentService) MigrateDocument(context model.RequestContext) (Document, error) {
	// Validate session
	auth, err := s.authFromSession(context)
	if err != nil {
		return Document{}, err
	}

	s.userLockTable.Lock(auth.Username)
	defer s.userLockTable.Unlock(auth.Username)

	user, err := s.repo.GetUser(auth)
	if err != nil {
		return Document{}, errors.INVALID_CREDENTIALS
	}

	user, err = s.repo.MigrateUser(user)
	if err != nil {
		return Document{}, errors.SERVER_ERROR
	}

	return adaptUserToDocument(user), nil
}

func (s *DocumentService) GetDocument(context model.RequestContext) (Document, error) {
	// Validate username
	if !isUsernameValid(context.Username) {
		return Document{}, errors.INVALID_CREDENTIALS
	}

	auth, err := s.authFromSession(context)
	if err != nil {
		return Document{}, err
	}

	s.userLockTable.Lock(auth.Username)
	defer s.userLockTable.Unlock(auth.Username)

	user, err := s.repo.GetUser(auth)
	if err != nil {
		return Document{}, errors.INVALID_CREDENTIALS
	}

	document := adaptUserToDocument(user)

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

	s.userLockTable.Lock(auth.Username)
	defer s.userLockTable.Unlock(auth.Username)

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

func (s *DocumentService) UpdateDocumentAndPassword(context model.RequestContext, newPassword string, document string) (Document, string, error) {
	// Validate username
	if !isUsernameValid(context.Username) {
		return Document{}, "", errors.INVALID_CREDENTIALS
	}

	auth, err := s.authFromSession(context)
	if err != nil {
		return Document{}, "", err
	}

	s.userLockTable.Lock(auth.Username)
	defer s.userLockTable.Unlock(auth.Username)

	user, err := s.repo.GetUser(auth)
	if err != nil {
		return Document{}, "", errors.INVALID_CREDENTIALS
	}

	user, err = s.repo.UpdateUser(user, model.UserUpdate{Password: newPassword, Document: []byte(document)})
	if err != nil {
		s.logError(err)
		return Document{}, "", errors.SERVER_ERROR
	}

	// Grant session for new user
	session, err := s.sessionTable.Grant(user.Username, user.Auth.Ip, user.Auth.Password)
	if err != nil {
		s.errorLogger.Println("Error granting user session")
		return Document{}, "", err
	}

	if err != nil {
		s.logError(err)
		return Document{}, "", errors.SERVER_ERROR
	}

	return adaptUserToDocument(user), session.SessionToken, nil
}

func (s *DocumentService) logError(err error) {
	if err != nil {
		s.errorLogger.Printf("%s\n", err.Error())
	}
}
