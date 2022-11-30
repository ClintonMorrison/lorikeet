package server

import (
	"encoding/json"
	"log"
)

type SessionRequest struct {
	Username        string `json:"username"`
	DecryptToken    string `json:"decryptToken"`
	RecaptchaResult string `json:"recaptchaResult"`
}

func NewSessionController(service *SessionService, requestLogger *log.Logger) RestController {
	// `POST /session`
	var Post MethodHandler = func(request ApiRequest) ApiResponse {
		sessionRequest, err := parseSessionRequestBody(request.Body)
		if err != nil {
			return badRequestResponse
		}

		auth := Auth{
			username: sessionRequest.Username,
			password: sessionRequest.DecryptToken,
			ip:       request.Context.ip,
		}

		token, err := service.GrantSession(auth, sessionRequest.RecaptchaResult)
		if err != nil {
			return responseForSessionError(err)
		}

		headers := make([]ResponseHeader, 0)
		headers = append(headers, SetSessionCookieHeader(token))

		return ApiResponse{201, headers, emptyBody}
	}

	return RestController{
		requestLogger: requestLogger,
		Post:          Post,
	}
}

var fallbackErrorJSON, _ = json.Marshal(internalServerError)

func responseForSessionError(err error) ApiResponse {
	switch err {
	case ERROR_BAD_REQUEST:
	case ERROR_INVALID_USER_NAME:
		return badRequestResponse
	case ERROR_INVALID_CREDENTIALS:
		return unauthorizedResponse
	case ERROR_TOO_MANY_REQUESTS:
		return tooManyRequestsResponse
	case ERROR_SERVER_ERROR:
	default:
		return serverErrorResponse
	}

	return serverErrorResponse
}

func parseSessionRequestBody(body []byte) (*SessionRequest, error) {
	sessionRequest := &SessionRequest{}
	err := json.Unmarshal(body, sessionRequest)
	if err != nil {
		return nil, ERROR_BAD_REQUEST
	}

	return sessionRequest, nil
}
