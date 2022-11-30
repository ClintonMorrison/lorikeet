package server

import (
	"log"
)

type SessionService struct {
	recaptchaClient *RecaptchaClient
	documentService *DocumentService
	sessionTable    *SessionTable
	lockoutTable    *LockoutTable
	errorLogger     *log.Logger
}

// GrantSession returns a new session token, or an error
func (s *SessionService) GrantSession(auth Auth, recaptchaResponse string) (string, error) {
	// Check request isn't blocked
	if !s.lockoutTable.shouldAllow(auth.ip, auth.username) {
		s.lockoutTable.logFailure(auth.ip, auth.username)
		return "", ERROR_TOO_MANY_REQUESTS
	}

	// Validate recaptcha
	recaptchaValid := s.recaptchaClient.Verify(recaptchaResponse, auth.ip)
	if !recaptchaValid {
		s.errorLogger.Println("Recaptcha was not valid")
		s.lockoutTable.logFailure(auth.ip, auth.username)
		return "", ERROR_INVALID_CREDENTIALS
	}

	// Validate auth
	_, err := s.documentService.checkAuth(auth)
	if err != nil {
		s.errorLogger.Println("Auth in session request was not valid")
		s.lockoutTable.logFailure(auth.ip, auth.username)
		return "", ERROR_INVALID_CREDENTIALS
	}

	// Grant session (proves user passed recaptcha with valid auth)
	session, err := s.sessionTable.Grant(auth.username, auth.ip, auth.password)
	if err != nil {
		s.errorLogger.Println("Error granting user session")
		s.lockoutTable.logFailure(auth.ip, auth.username)
		return "", err
	}

	return session.SessionToken, nil
}
