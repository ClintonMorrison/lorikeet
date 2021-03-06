package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ClintonMorrison/lorikeet/internal/storage"
)

const documentApiPath = "/api/document"

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
	service := &Service{
		repo:         repository,
		lockoutTable: lockoutTable,
		errorLogger:  errorLogger,
		lockByUser:   make(map[string]*sync.RWMutex),
	}
	controller := Controller{service, requestLogger}

	repository.createDataDirectory()

	http.HandleFunc(documentApiPath, controller.handleDocument)
	fmt.Printf("Listening on http://localhost%s\n", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
