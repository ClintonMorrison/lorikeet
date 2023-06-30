package controller

import (
	"encoding/json"
	"log"

	"github.com/ClintonMorrison/lorikeet/internal/errors"
	"github.com/ClintonMorrison/lorikeet/internal/server/lockout"
	"github.com/ClintonMorrison/lorikeet/internal/server/service"
)

type DocumentResponse struct {
	Document             string `json:"document"`             // encrypted password data
	Salt                 string `json:"salt"`                 // salt for client to use
	ClientStorageVersion int    `json:"clientStorageVersion"` // specifies which client encrypt algorithm was used
	ServerStorageVersion int    `json:"serverStorageVersion"` // specifies how data is stored on sever (1 = legacy, 2 = new)
}

type DocumentRequest struct {
	Password        string `json:"password,omitempty"`
	Document        string `json:"document"`
	RecaptchaResult string `json:"recaptchaResult"`
}

func NewDocumentController(
	cookieHelper *CookieHelper,
	service *service.DocumentService,
	lockoutTable *lockout.Table,
	requestLogger *log.Logger) RestController {
	// GET /document
	var get MethodHandler = func(request ApiRequest) ApiResponse {
		document, err := service.GetDocument(request.Context)
		if err != nil {
			return responseForError(err)
		}

		headers := make([]ResponseHeader, 0)

		responseBody, err := marshalResponse(DocumentResponse{
			Document:             string(document.Data),
			Salt:                 string(document.Salt),
			ClientStorageVersion: document.ClientVersion,
			ServerStorageVersion: document.ServerVersion,
		})

		if err != nil {
			return responseForError(err)
		}

		return ApiResponse{200, headers, responseBody, ""}
	}

	// POST /document
	var post MethodHandler = func(request ApiRequest) ApiResponse {
		parsedBody, err := parseDocumentRequestBody(request.Body)
		if err != nil {
			return responseForError(err)
		}

		sessionToken, err := service.CreateDocument(request.Context, parsedBody.Document, parsedBody.RecaptchaResult)
		if err != nil {
			return responseForError(err)
		}

		headers := make([]ResponseHeader, 0)
		headers = append(headers, cookieHelper.SetSessionCookieHeader(sessionToken))

		return ApiResponse{201, headers, emptyBody, ""}
	}

	// PUT /document
	var put MethodHandler = func(request ApiRequest) ApiResponse {
		parsedBody, err := parseDocumentRequestBody(request.Body)
		if err != nil {
			return responseForError(err)
		}

		// Update password if password was given
		if len(parsedBody.Password) > 0 {
			sessionToken, err := service.UpdateDocumentAndPassword(request.Context, parsedBody.Password, parsedBody.Document)
			if err != nil {
				return responseForError(err)
			}

			headers := make([]ResponseHeader, 0)
			headers = append(headers, cookieHelper.SetSessionCookieHeader(sessionToken))

			return ApiResponse{202, headers, emptyBody, ""}
		}

		err = service.UpdateDocument(request.Context, parsedBody.Document)
		if err != nil {
			return responseForError(err)
		}

		return ApiResponse{202, emptyHeaders, emptyBody, ""}
	}

	// DELETE /document
	var delete MethodHandler = func(request ApiRequest) ApiResponse {
		err := service.DeleteDocument(request.Context)
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
		Get:           get,
		Post:          post,
		Put:           put,
		Delete:        delete,
	}
}

func parseDocumentRequestBody(body []byte) (*DocumentRequest, error) {
	documentRequest := &DocumentRequest{}
	err := json.Unmarshal(body, documentRequest)
	if err != nil {
		return nil, errors.BAD_REQUEST
	}

	return documentRequest, nil
}

func marshalResponse(response DocumentResponse) ([]byte, error) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return emptyBody, errors.SERVER_ERROR
	}
	return jsonResponse, nil
}
