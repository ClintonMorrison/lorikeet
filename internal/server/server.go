package server

import (
	"fmt"
	"net/http"

	"github.com/ClintonMorrison/lorikeet/internal/server/controller"
	"github.com/ClintonMorrison/lorikeet/internal/server/lockout"
	"github.com/ClintonMorrison/lorikeet/internal/server/recaptcha"
	"github.com/ClintonMorrison/lorikeet/internal/server/repository"
	"github.com/ClintonMorrison/lorikeet/internal/server/service"
	"github.com/ClintonMorrison/lorikeet/internal/server/session"
	"github.com/ClintonMorrison/lorikeet/internal/storage"
	"github.com/ClintonMorrison/lorikeet/internal/utils"
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
	requestLogger, err := utils.CreateLogger(requestLogFilename, "[REQUEST] ")
	if err != nil {
		panic(err)
	}

	// Error logger
	errorLogger, err := utils.CreateLogger(errorLogFilename, "[ERROR] ")
	if err != nil {
		panic(err)
	}

	// Debug logger
	debugLogger, err := utils.CreateLogger(debugLogFilename, "[DEBUG] ")
	if err != nil {
		panic(err)
	}

	if localDev {
		debugLogger.Println("Server running in local dev mode")
	}

	cookieHelper := controller.NewCookieHelper(localDev)
	recaptchaClient := recaptcha.NewClient(
		debugLogger,
		recaptchaSecret,
	)
	repository := repository.NewRepositoryV1(dataPath)
	lockoutTable := lockout.NewTable()
	sessionTable := session.NewTable()
	documentService := service.NewDocumentService(
		repository,
		recaptchaClient,
		sessionTable,
		errorLogger,
	)
	sessionService := service.NewSessionService(
		recaptchaClient,
		repository,
		sessionTable,
		errorLogger,
	)

	repository.CreateDataDirectory()

	documentController := controller.NewDocumentController(cookieHelper, documentService, lockoutTable, requestLogger)
	http.HandleFunc(documentApiPath, documentController.Handle)

	sessionController := controller.NewSessionController(cookieHelper, sessionService, lockoutTable, requestLogger)
	http.HandleFunc(sessionApiPath, sessionController.Handle)

	fmt.Printf("Listening on http://localhost%s\n", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
