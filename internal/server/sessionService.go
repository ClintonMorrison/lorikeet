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
		s.errorLogger.Println("Recaptcha in grant session request was not valid")
		s.lockoutTable.logFailure(auth.ip, auth.username)
		return "", ERROR_INVALID_CREDENTIALS
	}

	// Validate auth
	_, err := s.documentService.checkAuth(auth)
	if err != nil {
		s.errorLogger.Println("Auth in grant session request was not valid")
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

// RevokeSession deletes an existing session
func (s *SessionService) RevokeSession(token string, username string, ip string) error {
	// Check request isn't blocked
	if !s.lockoutTable.shouldAllow(ip, username) {
		s.lockoutTable.logFailure(ip, username)
		return ERROR_TOO_MANY_REQUESTS
	}

	// Revoke session
	err := s.sessionTable.RevokeSession(token, username)
	if err != nil {
		s.errorLogger.Println("Unable to revoke session")
		s.lockoutTable.logFailure(ip, username)
		return err
	}

	return nil
}