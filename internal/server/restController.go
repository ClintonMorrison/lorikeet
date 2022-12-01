package server

import (
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
	Code     int
	Headers  []ResponseHeader
	Body     []byte
	ErrorMsg string
}

type MethodHandler func(ApiRequest) ApiResponse

type RestController struct {
	lockoutTable  *LockoutTable
	requestLogger *log.Logger
	Get           MethodHandler
	Post          MethodHandler
	Put           MethodHandler
	Delete        MethodHandler
}

type ErrorBody struct {
	Error string `json:"error"`
}

var emptyBody = make([]byte, 0)
var emptyHeaders = make([]ResponseHeader, 0)

func (c *RestController) runMethodHandler(w http.ResponseWriter, r *http.Request, context RequestContext, handler MethodHandler) ApiResponse {
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
		Context: context,
		Body:    body,
	}

	response := handler(request)

	return response
}

func (c *RestController) checkLockoutAndHandle(w http.ResponseWriter, r *http.Request, context RequestContext) ApiResponse {
	if !c.lockoutTable.ShouldAllow(context.ip, context.username) {
		return tooManyRequestsResponse
	}

	response := ApiResponse{}

	// Call handler based on method
	switch r.Method {
	case "GET":
		response = c.runMethodHandler(w, r, context, c.Get)
	case "PUT":
		response = c.runMethodHandler(w, r, context, c.Put)
	case "POST":
		response = c.runMethodHandler(w, r, context, c.Post)
	case "DELETE":
		response = c.runMethodHandler(w, r, context, c.Delete)
	default:
		response = badRequestResponse
	}

	if response.Code >= 400 {
		c.lockoutTable.LogFailure(context.ip, context.username)
	}

	return response
}

func (c *RestController) Handle(w http.ResponseWriter, r *http.Request) {
	context := ParseBasicContext(r)

	response := c.checkLockoutAndHandle(w, r, context)

	c.logRequest(r, response, context.username)

	w.Header().Add("Content-Type", "application/json")
	for _, header := range response.Headers {
		w.Header().Add(header.Name, header.Value)
	}

	w.WriteHeader(response.Code)
	w.Write(response.Body)
}

func (c *RestController) logRequest(r *http.Request, response ApiResponse, username string) {
	ip := r.Header.Get("X-Forwarded-For")
	result := "OK"
	if response.ErrorMsg != "" {
		result = response.ErrorMsg
	}

	name := username

	c.requestLogger.Printf(
		"%s %s | %d [%s] | %s | %s\n",
		r.Method, r.RequestURI,
		response.Code, result, name,
		ip)
}
