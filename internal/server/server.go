package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ClintonMorrison/lorikeet/internal/server/repository"
	"github.com/ClintonMorrison/lorikeet/internal/storage"
)

const documentApiPath = "/api/document"
const sessionApiPath = "/api/session"

func Run(
	dataPath string,
	address string,
	recaptchaSecret string,
	localDev bool,
	logPath string,
	requestLogFilename string,
	errorLogFilename string,
	debugLogFilename string) {

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

	// Debug logger
	debugLogger, err := createLogger(debugLogFilename, "[DEBUG] ")
	if err != nil {
		panic(err)
	}

	if localDev {
		debugLogger.Println("Server running in local dev mode")
	}

	cookieHelper := &CookieHelper{localDev}
	recaptchaClient := &RecaptchaClient{
		debugLogger: debugLogger,
		secret:      recaptchaSecret,
	}
	repository := repository.NewRepositoryV1(dataPath)
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

	repository.CreateDataDirectory()

	documentController := NewDocumentController(cookieHelper, documentService, lockoutTable, requestLogger)
	http.HandleFunc(documentApiPath, documentController.Handle)

	sessionController := NewSessionController(cookieHelper, sessionService, lockoutTable, requestLogger)
	http.HandleFunc(sessionApiPath, sessionController.Handle)

	fmt.Printf("Listening on http://localhost%s\n", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
