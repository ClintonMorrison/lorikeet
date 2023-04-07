package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ClintonMorrison/lorikeet/internal/model"
	"github.com/ClintonMorrison/lorikeet/internal/server/session"
)

func parseCookies(request *http.Request) map[string]string {
	result := make(map[string]string, 0)
	if request == nil {
		return result
	}

	items := strings.Split(request.Header.Get("Cookie"), ";")

	for _, item := range items {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.ToLower(parts[0])
		val, err := url.QueryUnescape(parts[1])
		if err != nil {
			continue
		}

		result[key] = val
	}

	return result
}

type CookieHelper struct {
	localDev bool
}

const sessionCookieName = "session"

func (ch *CookieHelper) setCookieHeader(name string, value string, lifespan time.Duration) ResponseHeader {
	maxAge := int(lifespan.Seconds())

	parts := make([]string, 0)
	parts = append(parts, fmt.Sprintf("%s=%s", name, url.QueryEscape(value)))
	parts = append(parts, "SameSite=Strict")
	parts = append(parts, fmt.Sprintf("Max-Age=%d", maxAge))
	parts = append(parts, "HttpOnly")
	parts = append(parts, "Path=/")

	if !ch.localDev {
		parts = append(parts, "Secure")
	}

	return ResponseHeader{
		"Set-Cookie",
		strings.Join(parts, "; "),
	}
}

func (ch *CookieHelper) SetSessionCookieHeader(sessionToken string) ResponseHeader {
	return ch.setCookieHeader(sessionCookieName, sessionToken, session.Lifespan)
}

func (ch *CookieHelper) ClearSessionCookieHeader() ResponseHeader {
	return ch.setCookieHeader(sessionCookieName, "", 0)
}

func ParseBasicContext(r *http.Request) model.RequestContext {
	username, password, _ := r.BasicAuth()

	username = strings.TrimSpace(strings.ToLower(username))
	ip := r.Header.Get("X-Forwarded-For")

	cookies := parseCookies(r)
	sesionToken := cookies["session"]

	return model.RequestContext{Username: username, Ip: ip, Password: password, SessionToken: sesionToken}
}
