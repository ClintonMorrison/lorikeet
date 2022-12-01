package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ClintonMorrison/lorikeet/internal/storage"
)

const documentApiPath = "/api/document"
const sessionApiPath = "/api/session"

func Run(
	dataPath string,
	address string,
	recaptchaSecret string,
	logPath string,
	requestLogFilename string,
	errorLogFilename string) {

	storage.CreateDirectory(logPath)

	// Request logger
	requestLogger, err := createLogger(requestLogFilename, "[REQUEST] ")
	if err != nil {
		panic(err)
	}

	// Error logger
	errorLogger, err := createLogger(errorLogFilename, "[ERROR] ")
	if err != nil {
		panic(err)
	}

	recaptchaClient := &RecaptchaClient{
		secret: recaptchaSecret,
	}
	repository := &Repository{dataPath}
	lockoutTable := NewLockoutTable()
	sessionTable := NewSessionTable()
	documentService := &DocumentService{
		repo:            repository,
		recaptchaClient: recaptchaClient,
		sessionTable:    sessionTable,
		errorLogger:     errorLogger,
		lockByUser:      make(map[string]*sync.RWMutex),
	}
	sessionService := &SessionService{
		recaptchaClient: recaptchaClient,
		documentService: documentService,
		sessionTable:    sessionTable,
		errorLogger:     errorLogger,
	}

	repository.createDataDirectory()

	documentController := NewDocumentController(documentService, lockoutTable, requestLogger)
	http.HandleFunc(documentApiPath, documentController.Handle)

	sessionController := NewSessionController(sessionService, lockoutTable, requestLogger)
	http.HandleFunc(sessionApiPath, sessionController.Handle)

	fmt.Printf("Listening on http://localhost%s\n", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
