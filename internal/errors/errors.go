package errors

type TypedError string

const (
	BAD_REQUEST         TypedError = "BAD_REQUEST"
	ALREADY_EXISTS      TypedError = "ERROR_ALREADY_EXISTS"
	SERVER_ERROR        TypedError = "SERVER_ERROR"
	INVALID_CREDENTIALS TypedError = "INVALID_CREDENTIALS"
	USERNAME_INVALID    TypedError = "INVALID_USERNAME"
	TOO_MANY_REQUESTS   TypedError = "TOO_MANY_REQUESTS"
)

func (t TypedError) Error() string {
	return string(t)
}
