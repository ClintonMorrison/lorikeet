package server

import "encoding/json"

type TypedError string

const (
	ERROR_BAD_REQUEST         TypedError = "BAD_REQUEST"
	ERROR_INVALID_USER_NAME   TypedError = "INVALID_USER_NAME"
	ERROR_SERVER_ERROR        TypedError = "SERVER_ERROR"
	ERROR_INVALID_CREDENTIALS TypedError = "INVALID_CREDENTIALS"
	ERROR_TOO_MANY_REQUESTS   TypedError = "TOO_MANY_REQUESTS"
)

func (t TypedError) Error() string {
	return string(t)
}

func NewErrorResponse(code int, msg string) ApiResponse {
	var body, _ = json.Marshal(ErrorBody{msg})

	return ApiResponse{code, emptyHeaders, body, msg}
}

var badRequestResponse = NewErrorResponse(400, "Invalid request.")
var usernameTakenResponse = NewErrorResponse(400, "Username already taken.")
var unauthorizedResponse = NewErrorResponse(401, "Invalid username or password.")
var tooManyRequestsResponse = NewErrorResponse(429, "Too many failed attempts. Try again in a few hours.")
var serverErrorResponse = NewErrorResponse(500, "Server error. Please try again later.")

func responseForError(err error) ApiResponse {
	switch err {
	case ERROR_BAD_REQUEST:
		return badRequestResponse
	case ERROR_INVALID_USER_NAME:
		return usernameTakenResponse
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
