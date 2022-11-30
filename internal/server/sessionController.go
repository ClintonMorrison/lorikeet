package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type SessionResponse struct {
	Code  int    `json:"-"`
	Error string `json:"error,omitempty"`
}

type SessionRequest struct {
	Username        string `json:"username"`
	DecryptToken    string `json:"decryptToken"`
	RecaptchaResult string `json:"recaptchaResult"`
}

type SessionController struct {
	service       *SessionService
	requestLogger *log.Logger
}

var invalidSessionRequestResponse = SessionResponse{400, "Invalid request."}
var invalidSessionCredentialsResponse = SessionResponse{403, "Invalid username or password."}
var tooManySessionRequestsResponse = SessionResponse{429, "Too many failed attempts. Try again in a few hours."}
var internalSessionServerError = SessionResponse{500, "Server error. Please try again later."}

var fallbackErrorJSON, _ = json.Marshal(internalServerError)

func responseForSessionError(err error) SessionResponse {
	switch err {
	case ERROR_BAD_REQUEST:
	case ERROR_INVALID_USER_NAME:
		return invalidSessionRequestResponse
	case ERROR_INVALID_CREDENTIALS:
		return invalidSessionCredentialsResponse
	case ERROR_TOO_MANY_REQUESTS:
		return tooManySessionRequestsResponse
	case ERROR_SERVER_ERROR:
	default:
		return internalSessionServerError
	}

	return internalSessionServerError
}

func parseSessionRequestBody(body []byte) (*SessionRequest, error) {
	sessionRequest := &SessionRequest{}
	err := json.Unmarshal(body, sessionRequest)
	if err != nil {
		return nil, ERROR_BAD_REQUEST
	}

	return sessionRequest, nil
}

func (c *SessionController) postSession(w http.ResponseWriter, context RequestContext, sessionRequest SessionRequest) SessionResponse {
	token, err := c.service.GrantSession(context.ToAuth(sessionRequest.DecryptToken), sessionRequest.RecaptchaResult)
	if err != nil {
		return responseForSessionError(err)
	}

	SetSessionCookie(w, token)

	return SessionResponse{201, ""}
}

func (c *SessionController) parseRequestOrWriteError(w http.ResponseWriter, r *http.Request) (RequestContext, *SessionRequest) {
	// Read auth headers
	context := ParseBasicContext(r)

	// Read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.writeResponse(r, w, invalidSessionRequestResponse, context)
		return context, nil
	}

	// Parse body
	sessionRequest, err := parseSessionRequestBody(body)
	if err != nil {
		c.writeResponse(r, w, invalidSessionRequestResponse, context)
		return context, nil
	}

	// Add user data to context
	context.username = sessionRequest.Username

	return context, sessionRequest
}

func (c *SessionController) handle(w http.ResponseWriter, r *http.Request) {
	context, sessionRequest := c.parseRequestOrWriteError(w, r)
	if sessionRequest == nil {
		return
	}

	var response SessionResponse

	// Call handler based on method
	switch r.Method {
	case "POST":
		response = c.postSession(w, context, *sessionRequest)
	default:
		response = invalidSessionRequestResponse
	}

	c.writeResponse(r, w, response, context)
}

func (c *SessionController) logRequest(r *http.Request, response SessionResponse, context RequestContext) {
	ip := r.Header.Get("X-Forwarded-For")
	result := "OK"
	if response.Error != "" {
		result = response.Error
	}

	name := context.username

	c.requestLogger.Printf(
		"%s %s | %d [%s] | %s | %s\n",
		r.Method, r.RequestURI,
		response.Code, result, name,
		ip)
}

func (c *SessionController) writeResponse(r *http.Request, w http.ResponseWriter, response SessionResponse, context RequestContext) {
	c.logRequest(r, response, context)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(response.Code)

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(500)
		w.Write(fallbackErrorJSON)
		return
	}

	w.Write(jsonResponse)
}
