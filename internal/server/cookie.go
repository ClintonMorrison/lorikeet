package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func ParseCookies(request *http.Request) map[string]string {
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

func SetCookie(w http.ResponseWriter, name string, value string, lifespan time.Duration) {
	// TODO: add "Secure" in prod environment?
	maxAge := int(lifespan.Seconds())
	w.Header().Add("Set-Cookie", fmt.Sprintf("%s=%s; SameSite=Strict; Max-Age=%d; HttpOnly; Path=/;",
		name, url.QueryEscape(value), maxAge))
}

const sessionCookieName = "session"
const sessionCookieLifespan = time.Second * 60 // TODO

func GetSessionToken(request *http.Request) string {
	cookies := ParseCookies(request)
	return cookies[sessionCookieName]
}

func SetSessionCookie(w http.ResponseWriter, sessionToken string) {
	SetCookie(w, sessionCookieName, sessionToken, sessionCookieLifespan)
}

func ClearSessionCookie(w http.ResponseWriter) {
	SetCookie(w, sessionCookieName, "", 0)
}
