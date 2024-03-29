package service

import (
	"log"
	"time"

	"github.com/ClintonMorrison/lorikeet/internal/errors"
	"github.com/ClintonMorrison/lorikeet/internal/model"
	"github.com/ClintonMorrison/lorikeet/internal/server/recaptcha"
	"github.com/ClintonMorrison/lorikeet/internal/server/repository"
	"github.com/ClintonMorrison/lorikeet/internal/server/session"
)

type SessionService struct {
	recaptchaClient *recaptcha.Client
	repository      *repository.UserRepository
	sessionTable    *session.Table
	userLockTable   *UserLockTable
	errorLogger     *log.Logger
}

func NewSessionService(
	recaptchaClient *recaptcha.Client,
	repository *repository.UserRepository,
	sessionTable *session.Table,
	userLockTable *UserLockTable,
	errorLogger *log.Logger,

) *SessionService {
	return &SessionService{
		recaptchaClient,
		repository,
		sessionTable,
		userLockTable,
		errorLogger,
	}
}

// GrantSession returns a new session token, or an error
func (s *SessionService) GrantSession(auth model.Auth, recaptchaResponse string) (string, *model.User, error) {
	// Validate recaptcha
	recaptchaValid := s.recaptchaClient.Verify(recaptchaResponse, auth.Ip)
	if !recaptchaValid {
		s.errorLogger.Println("Recaptcha in grant session request was not valid")
		return "", nil, errors.INVALID_CREDENTIALS
	}

	// Lock on username
	s.userLockTable.Lock(auth.Username)
	defer s.userLockTable.Unlock(auth.Username)

	// Validate auth
	user, err := s.repository.GetUser(auth)
	if err != nil {
		s.errorLogger.Println("Auth in grant session request was not valid")
		return "", nil, errors.INVALID_CREDENTIALS
	}

	// Update last accessed time
	user, err = s.repository.UpdateUser(user, model.UserUpdate{LastAccessTime: time.Now()})
	if err != nil {
		s.errorLogger.Println("Failed to update lastAccessTime on user")
		return "", nil, errors.INVALID_CREDENTIALS
	}

	// Grant session (proves user passed recaptcha with valid auth)
	session, err := s.sessionTable.Grant(auth.Username, auth.Ip, auth.Password)
	if err != nil {
		s.errorLogger.Println("Error granting user session")
		return "", nil, err
	}

	return session.SessionToken, user, nil
}

// RevokeSession deletes an existing session
func (s *SessionService) RevokeSession(token string, username string, ip string) error {
	s.userLockTable.Lock(username)
	defer s.userLockTable.Unlock(username)

	err := s.sessionTable.RevokeSession(token, username)
	if err != nil {
		s.errorLogger.Println("Unable to revoke session")
		return err
	}

	return nil
}
