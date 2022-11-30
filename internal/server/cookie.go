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

const sessionCookieName = "session"
const sessionCookieLifespan = time.Second * 60 // TODO

func SetCookieHeader(name string, value string, lifespan time.Duration) ResponseHeader {
	// TODO: add "Secure" in prod environment?
	maxAge := int(lifespan.Seconds())
	return ResponseHeader{
		"Set-Cookie",
		fmt.Sprintf("%s=%s; SameSite=Strict; Max-Age=%d; HttpOnly; Path=/;",
			name, url.QueryEscape(value), maxAge),
	}
}

func GetSessionToken(request *http.Request) string {
	cookies := ParseCookies(request)
	return cookies[sessionCookieName]
}

func SetSessionCookieHeader(sessionToken string) ResponseHeader {
	return SetCookieHeader(sessionCookieName, sessionToken, sessionCookieLifespan)
}

func ClearSessionCookieHeader() ResponseHeader {
	return SetCookieHeader(sessionCookieName, "", 0)
}
