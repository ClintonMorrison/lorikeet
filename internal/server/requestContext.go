package server

import (
	"net/http"
	"strings"

	"github.com/ClintonMorrison/lorikeet/internal/model"
)

type RequestContext struct {
	username     string
	ip           string
	password     string // password from basic auth field
	sessionToken string
}

func (rc RequestContext) ToAuth(decryptToken string) model.Auth {
	return model.Auth{
		Username: rc.username,
		Password: decryptToken,
		Ip:       rc.ip,
	}
}

func ParseBasicContext(r *http.Request) RequestContext {
	username, password, _ := r.BasicAuth()

	username = strings.TrimSpace(strings.ToLower(username))
	ip := r.Header.Get("X-Forwarded-For")

	cookies := ParseCookies(r)
	sesionToken := cookies["session"]

	return RequestContext{username, ip, password, sesionToken}
}
