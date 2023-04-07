package controller

import (
	"encoding/json"

	"github.com/ClintonMorrison/lorikeet/internal/errors"
)

func NewErrorResponse(code int, msg string) ApiResponse {
	var body, _ = json.Marshal(ErrorBody{msg})

	return ApiResponse{code, emptyHeaders, body, msg}
}

var badRequestResponse = NewErrorResponse(400, "Invalid request.")
var usernameTakenResponse = NewErrorResponse(400, "Username already taken.")
var unauthorizedResponse = NewErrorResponse(401, "Invalid username or password.")
var usernameInvalidResponse = NewErrorResponse(422, "Username can only contain letters, numbers, or certain symbols (. @ ! $ + - _)")
var tooManyRequestsResponse = NewErrorResponse(429, "Too many failed attempts. Try again in a few hours.")
var serverErrorResponse = NewErrorResponse(500, "Server errors. Please try again later.")

func responseForError(err error) ApiResponse {
	switch err {
	case errors.BAD_REQUEST:
		return badRequestResponse
	case errors.ALREADY_EXISTS:
		return usernameTakenResponse
	case errors.INVALID_CREDENTIALS:
		return unauthorizedResponse
	case errors.USERNAME_INVALID:
		return usernameInvalidResponse
	case errors.TOO_MANY_REQUESTS:
		return tooManyRequestsResponse
	case errors.SERVER_ERROR:
	default:
		return serverErrorResponse
	}

	return serverErrorResponse
}
