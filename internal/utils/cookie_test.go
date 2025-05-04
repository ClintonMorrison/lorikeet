package utils

import (
	"net/http"
	"testing"
)

func TestParseCookies(t *testing.T) {
	request := &http.Request{
		Header: http.Header{
			"Cookie": []string{"session=abc123;\n other=xyz"},
		},
	}

	result := ParseCookies(request)

	if len(result) != 2 {
		t.Errorf("Expected 2 cookies, got %d", len(result))
	}

	if result["session"] != "abc123" {
		t.Errorf("Expected session cookie to be 'abc123', got '%s'", result["session"])
	}

	if result["other"] != "xyz" {
		t.Errorf("Expected other cookie to be 'xyz', got '%s'", result["other"])
	}
}
