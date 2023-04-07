package controller

import (
	"encoding/json"
	"log"

	"github.com/ClintonMorrison/lorikeet/internal/errors"
	"github.com/ClintonMorrison/lorikeet/internal/model"
	"github.com/ClintonMorrison/lorikeet/internal/server/lockout"
	"github.com/ClintonMorrison/lorikeet/internal/server/service"
)

type SessionRequest struct {
	DecryptToken    string `json:"decryptToken"`
	RecaptchaResult string `json:"recaptchaResult"`
}

func NewSessionController(
	cookieHelper *CookieHelper,
	service *service.SessionService,
	lockoutTable *lockout.Table,
	requestLogger *log.Logger) RestController {
	// POST /session
	var post MethodHandler = func(request ApiRequest) ApiResponse {
		sessionRequest, err := parseSessionRequestBody(request.Body)
		if err != nil {
			return badRequestResponse
		}

		auth := model.Auth{
			Username: request.Context.Username,
			Password: sessionRequest.DecryptToken,
			Ip:       request.Context.Ip,
		}

		token, err := service.GrantSession(auth, sessionRequest.RecaptchaResult)
		if err != nil {
			return responseForError(err)
		}

		headers := make([]ResponseHeader, 0)
		headers = append(headers, cookieHelper.SetSessionCookieHeader(token))

		return ApiResponse{201, headers, emptyBody, ""}
	}

	// DELETE /session
	var delete MethodHandler = func(request ApiRequest) ApiResponse {
		err := service.RevokeSession(request.Context.SessionToken, request.Context.Username, request.Context.Ip)
		if err != nil {
			return responseForError(err)
		}

		headers := make([]ResponseHeader, 0)
		headers = append(headers, cookieHelper.ClearSessionCookieHeader())

		return ApiResponse{204, headers, emptyBody, ""}
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
		return nil, errors.BAD_REQUEST
	}

	return sessionRequest, nil
}
