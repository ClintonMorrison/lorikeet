package server

import (
	"log"

	"github.com/ClintonMorrison/lorikeet/internal/model"
)

type SessionService struct {
	recaptchaClient *RecaptchaClient
	documentService *DocumentService
	sessionTable    *SessionTable
	errorLogger     *log.Logger
}

// GrantSession returns a new session token, or an error
func (s *SessionService) GrantSession(auth model.Auth, recaptchaResponse string) (string, error) {
	// Validate recaptcha
	recaptchaValid := s.recaptchaClient.Verify(recaptchaResponse, auth.Ip)
	if !recaptchaValid {
		s.errorLogger.Println("Recaptcha in grant session request was not valid")
		return "", ERROR_INVALID_CREDENTIALS
	}

	// Validate auth
	_, err := s.documentService.checkAuth(auth)
	if err != nil {
		s.errorLogger.Println("Auth in grant session request was not valid")
		return "", ERROR_INVALID_CREDENTIALS
	}

	// Grant session (proves user passed recaptcha with valid auth)
	session, err := s.sessionTable.Grant(auth.Username, auth.Ip, auth.Password)
	if err != nil {
		s.errorLogger.Println("Error granting user session")
		return "", err
	}

	return session.SessionToken, nil
}

// RevokeSession deletes an existing session
func (s *SessionService) RevokeSession(token string, username string, ip string) error {
	err := s.sessionTable.RevokeSession(token, username)
	if err != nil {
		s.errorLogger.Println("Unable to revoke session")
		return err
	}

	return nil
}
