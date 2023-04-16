package utils

import "net/http"

func GetIpFromRequest(r *http.Request) string {
	return r.Header.Get("X-Forwarded-For")
}
