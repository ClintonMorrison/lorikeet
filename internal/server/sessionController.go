package server

import (
	"encoding/json"
	"log"
)

type SessionRequest struct {
	DecryptToken    string `json:"decryptToken"`
	RecaptchaResult string `json:"recaptchaResult"`
}

func NewSessionController(service *SessionService, lockoutTable *LockoutTable, requestLogger *log.Logger) RestController {
	// POST /session
	var post MethodHandler = func(request ApiRequest) ApiResponse {
		sessionRequest, err := parseSessionRequestBody(request.Body)
		if err != nil {
			return badRequestResponse
		}

		auth := Auth{
			username: request.Context.username,
			password: sessionRequest.DecryptToken,
			ip:       request.Context.ip,
		}

		token, err := service.GrantSession(auth, sessionRequest.RecaptchaResult)
		if err != nil {
			return responseForError(err)
		}

		headers := make([]ResponseHeader, 0)
		headers = append(headers, SetSessionCookieHeader(token))

		return ApiResponse{201, headers, emptyBody}
	}

	// DELETE /session
	var delete MethodHandler = func(request ApiRequest) ApiResponse {
		err := service.RevokeSession(request.Context.sessionToken, request.Context.username, request.Context.ip)
		if err != nil {
			return responseForError(err)
		}

		headers := make([]ResponseHeader, 0)
		headers = append(headers, ClearSessionCookieHeader())

		return ApiResponse{204, headers, emptyBody}
	}

	return RestController{
		lockoutTable:  lockoutTable,
		requestLogger: requestLogger,
		Post:          post,
		Delete:        delete,
	}
}

func parseSessionRequestBody(body []byte) (*SessionRequest, error) {
	sessionRequest := &SessionRequest{}
	err := json.Unmarshal(body, sessionRequest)
	if err != nil {
		return nil, ERROR_BAD_REQUEST
	}

	return sessionRequest, nil
}
