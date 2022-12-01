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

func GetSessionToken(request *http.Request) string {
	cookies := ParseCookies(request)
	return cookies[sessionCookieName]
}

type CookieHelper struct {
	localDev bool
}

const sessionCookieName = "session"

func (ch *CookieHelper) SetCookieHeader(name string, value string, lifespan time.Duration) ResponseHeader {
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
	return ch.SetCookieHeader(sessionCookieName, sessionToken, sessionLifespan)
}

func (ch *CookieHelper) ClearSessionCookieHeader() ResponseHeader {
	return ch.SetCookieHeader(sessionCookieName, "", 0)
}
