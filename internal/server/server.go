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

	repository := &Repository{dataPath}
	lockoutTable := NewLockoutTable()
	service := &DocumentService{
		repo:         repository,
		lockoutTable: lockoutTable,
		errorLogger:  errorLogger,
		lockByUser:   make(map[string]*sync.RWMutex),
	}

	repository.createDataDirectory()

	documentController := DocumentController{service, requestLogger}
	http.HandleFunc(documentApiPath, documentController.handle)

	sessionController := SessionController{service, requestLogger}
	http.HandleFunc(sessionApiPath, sessionController.handle)

	fmt.Printf("Listening on http://localhost%s\n", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
