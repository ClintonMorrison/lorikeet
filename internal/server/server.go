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
		lockoutTable:    lockoutTable,
		sessionTable:    sessionTable,
		errorLogger:     errorLogger,
		lockByUser:      make(map[string]*sync.RWMutex),
	}
	sessionService := &SessionService{
		recaptchaClient: recaptchaClient,
		documentService: documentService,
		sessionTable:    sessionTable,
		lockoutTable:    lockoutTable,
		errorLogger:     errorLogger,
	}

	repository.createDataDirectory()

	documentController := NewDocumentController(documentService, requestLogger)
	http.HandleFunc(documentApiPath, documentController.handle)

	sessionController := NewSessionController(sessionService, requestLogger)
	http.HandleFunc(sessionApiPath, sessionController.handle)

	fmt.Printf("Listening on http://localhost%s\n", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
