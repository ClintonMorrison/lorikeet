package utils

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

func FormatCookie(name string, value string, lifespan time.Duration, secure bool) string {
	maxAge := int(lifespan.Seconds())

	parts := make([]string, 0)
	parts = append(parts, fmt.Sprintf("%s=%s", name, url.QueryEscape(value)))
	parts = append(parts, "SameSite=Strict")
	parts = append(parts, fmt.Sprintf("Max-Age=%d", maxAge))
	parts = append(parts, "HttpOnly")
	parts = append(parts, "Path=/")

	if secure {
		parts = append(parts, "Secure")
	}

	return strings.Join(parts, "; ")
}
