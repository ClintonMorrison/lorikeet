package server

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

// how long session is valid
const sessionLifespan = time.Hour * 24

type Session struct {
	SessionToken string
	DecryptToken string
	Username     string
	Ip           string
	IssuedAt     time.Time
	ExpiresAt    time.Time
}

type SessionTable struct {
	sessionByToken map[string]Session
	mux            sync.RWMutex
}

func generateToken(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	// Note that err == nil only if we read len(bytes) bytes.
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), err
}

func NewSessionTable() *SessionTable {
	return &SessionTable{
		sessionByToken: make(map[string]Session, 0)}
}

func (st *SessionTable) Grant(username string, ip string, decryptToken string) (*Session, error) {
	st.mux.Lock()
	defer st.mux.Unlock()

	st.purgeExpiredSessions()

	// Generate new token
	token, err := generateToken(256)
	if err != nil {
		return nil, ERROR_SERVER_ERROR
	}

	// Make sure it doesn't match an existing session
	_, alreadyExists := st.sessionByToken[token]
	if alreadyExists {
		return nil, ERROR_SERVER_ERROR
	}

	session := Session{
		SessionToken: token,
		DecryptToken: decryptToken,
		Username:     username,
		Ip:           ip,
		IssuedAt:     time.Now(),
		ExpiresAt:    time.Now().Add(sessionLifespan),
	}

	// Add session to table
	st.sessionByToken[token] = session

	return &session, nil
}

func (st *SessionTable) IsValid(token string, username string, ip string) bool {
	session, err := st.GetSession(token, username, ip)
	if err != nil {
		return false
	}

	return session.SessionToken == token
}

func (st *SessionTable) GetSession(token string, username string, ip string) (*Session, error) {
	st.mux.Lock()
	defer st.mux.Unlock()

	st.purgeExpiredSessions()

	// Get session from map
	session, exists := st.sessionByToken[token]
	if !exists {
		return nil, ERROR_INVALID_CREDENTIALS
	}

	// Make sure username matches
	if session.Username != username {
		return nil, ERROR_INVALID_CREDENTIALS
	}

	// Make sure token matches
	if session.SessionToken != token {
		return nil, ERROR_INVALID_CREDENTIALS
	}

	return &session, nil
}

func (st *SessionTable) RevokeSession(token string, username string) error {
	st.mux.Lock()
	defer st.mux.Unlock()

	st.purgeExpiredSessions()

	// Get session from map
	session, exists := st.sessionByToken[token]
	if !exists {
		return ERROR_INVALID_CREDENTIALS
	}

	// Make sure username matches
	if session.Username != username {
		return ERROR_INVALID_CREDENTIALS
	}

	// Make sure token matches
	if session.SessionToken != token {
		return ERROR_INVALID_CREDENTIALS
	}

	// Remove the session
	delete(st.sessionByToken, token)

	return nil
}

func (st *SessionTable) purgeExpiredSessions() {
	now := time.Now()

	for token, session := range st.sessionByToken {
		if now.After(session.ExpiresAt) {
			delete(st.sessionByToken, token)
		}
	}
}
