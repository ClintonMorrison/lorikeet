package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type DocumentResponse struct {
	Code     int    `json:"-"`
	Error    string `json:"error,omitempty"`
	Document string `json:"document,omitempty"`
}

type PasswordRequest struct {
	Password string `json:"password"`
	Document string `json:"document"`
}

type DocumentRequest struct {
	Password string `json:"password,omitempty"`
	Document string `json:"document"`
}

type DocumentController struct {
	service       *DocumentService
	requestLogger *log.Logger
}

var invalidRequestResponse = DocumentResponse{400, "Invalid request.", ""}
var usernameTakenResponse = DocumentResponse{400, "Username already taken.", ""}
var invalidCredentialsResponse = DocumentResponse{401, "Invalid user or credentials.", ""}
var tooManyRequestsResponse = DocumentResponse{429, "Too many failed attempts. Try again in a few hours.", ""}
var internalServerError = DocumentResponse{500, "Server error. Please try again later.", ""}
var fallbackLoginErrorJSON, _ = json.Marshal(internalServerError)

func responseForError(err error) DocumentResponse {
	switch err {
	case ERROR_BAD_REQUEST:
		return invalidRequestResponse
	case ERROR_INVALID_USER_NAME:
		return usernameTakenResponse
	case ERROR_INVALID_CREDENTIALS:
		return invalidCredentialsResponse
	case ERROR_TOO_MANY_REQUESTS:
		return tooManyRequestsResponse
	case ERROR_SERVER_ERROR:
	default:
		return internalServerError
	}

	return internalServerError
}

//
// Document API
//
func parseDocumentRequestBody(body []byte) (*DocumentRequest, error) {
	documentRequest := &DocumentRequest{}
	err := json.Unmarshal(body, documentRequest)
	if err != nil {
		return nil, ERROR_BAD_REQUEST
	}

	return documentRequest, nil
}

func (c *DocumentController) postDocument(w http.ResponseWriter, context RequestContext, body []byte) DocumentResponse {
	request, err := parseDocumentRequestBody(body)
	if err != nil {
		return responseForError(err)
	}

	sessionToken, err := c.service.CreateDocument(context, request.Document)
	if err != nil {
		return responseForError(err)
	}

	SetSessionCookie(w, sessionToken)

	return DocumentResponse{201, "", ""}
}

func (c *DocumentController) putDocument(w http.ResponseWriter, context RequestContext, body []byte) DocumentResponse {
	request, err := parseDocumentRequestBody(body)
	if err != nil {
		return responseForError(err)
	}

	// Update password if password was given
	if len(request.Password) > 0 {
		sessionToken, err := c.service.UpdateDocumentAndPassword(context, request.Password, request.Document)
		if err != nil {
			return responseForError(err)
		}
		fmt.Println("Returning session token! " + sessionToken + "/")
		SetSessionCookie(w, sessionToken)

		return DocumentResponse{202, "", ""}
	}

	err = c.service.UpdateDocument(context, request.Document)
	if err != nil {
		return responseForError(err)
	}

	return DocumentResponse{202, "", ""}
}

func (c *DocumentController) getDocument(w http.ResponseWriter, context RequestContext) DocumentResponse {
	document, err := c.service.GetDocument(context)
	if err != nil {
		return responseForError(err)
	}

	return DocumentResponse{200, "", string(document)}
}

func (c *DocumentController) deleteDocument(w http.ResponseWriter, context RequestContext) DocumentResponse {
	err := c.service.DeleteDocument(context)
	if err != nil {
		return responseForError(err)
	}

	ClearSessionCookie(w)

	return DocumentResponse{204, "", ""}
}

func (c *DocumentController) parseRequestOrWriteError(w http.ResponseWriter, r *http.Request) (RequestContext, []byte) {
	// Read auth headers
	context := ParseBasicContext(r)
	if context.username == "" || context.ip == "" {
		c.writeResponse(r, w, invalidRequestResponse, context)
		return RequestContext{}, nil
	}

	// Read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.writeResponse(r, w, invalidRequestResponse, context)
		return RequestContext{}, nil
	}

	return context, body
}

func (c *DocumentController) handle(w http.ResponseWriter, r *http.Request) {
	context, body := c.parseRequestOrWriteError(w, r)

	var response DocumentResponse

	// Call handler based on method
	switch r.Method {
	case "GET":
		response = c.getDocument(w, context)
	case "PUT":
		response = c.putDocument(w, context, body)
	case "POST":
		response = c.postDocument(w, context, body)
	case "DELETE":
		response = c.deleteDocument(w, context)
	default:
		response = invalidRequestResponse
	}

	c.writeResponse(r, w, response, context)
}

func (c *DocumentController) logRequest(r *http.Request, response DocumentResponse, context RequestContext) {
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

func (c *DocumentController) writeResponse(r *http.Request, w http.ResponseWriter, response DocumentResponse, context RequestContext) {
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
