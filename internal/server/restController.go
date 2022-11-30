package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type ResponseHeader struct {
	Name  string
	Value string
}

type ApiRequest struct {
	Headers http.Header
	Context RequestContext
	Body    []byte
}

type ApiResponse struct {
	Code    int
	Headers []ResponseHeader
	Body    []byte
}

type MethodHandler func(ApiRequest) ApiResponse // TODO: error?

type RestController struct {
	requestLogger *log.Logger
	Get           MethodHandler
	Post          MethodHandler
	Put           MethodHandler
	Delete        MethodHandler
	Options       MethodHandler
}

type ErrorBody struct {
	Error string `json:"error"`
}

func NewErrorResponse(code int, msg string) ApiResponse {
	var body, _ = json.Marshal(ErrorBody{msg})

	return ApiResponse{code, make([]ResponseHeader, 0), body}
}

var emptyBody = make([]byte, 0)
var badRequestResponse = NewErrorResponse(400, "Invalid request.")
var unauthorizedResponse = NewErrorResponse(403, "Invalid username or password.")
var tooManyRequestsResponse = NewErrorResponse(429, "Too many failed attempts. Try again in a few hours.")
var serverErrorResponse = NewErrorResponse(500, "Server error. Please try again later.")

func (c *RestController) runMethodHandler(r *http.Request, w http.ResponseWriter, handler MethodHandler) ApiResponse {
	// This resource does not support the request method
	if handler == nil {
		return badRequestResponse
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return badRequestResponse
	}

	request := ApiRequest{
		Headers: r.Header,
		Context: ParseBasicContext(r),
		Body:    body,
	}

	response := handler(request)

	c.logRequest(r, response.Code, request.Context.username)
	return response
}

func (c *RestController) handle(w http.ResponseWriter, r *http.Request) {

	response := ApiResponse{}

	// Call handler based on method
	switch r.Method {
	case "GET":
		response = c.runMethodHandler(r, w, c.Get)
	case "PUT":
		response = c.runMethodHandler(r, w, c.Put)
	case "POST":
		response = c.runMethodHandler(r, w, c.Post)
	case "DELETE":
		response = c.runMethodHandler(r, w, c.Delete)
	default:
		response = badRequestResponse
	}

	w.Header().Add("Content-Type", "application/json")
	for _, header := range response.Headers {
		w.Header().Add(header.Name, header.Value)
	}

	w.WriteHeader(response.Code)
	w.Write(response.Body)
}

func (c *RestController) logRequest(r *http.Request, responseCode int, username string) {
	ip := r.Header.Get("X-Forwarded-For")
	result := "OK"
	if responseCode >= 400 {
		result = "ERROR"
	}

	name := username

	c.requestLogger.Printf(
		"%s %s | %d [%s] | %s | %s\n",
		r.Method, r.RequestURI,
		responseCode, result, name,
		ip)
}
