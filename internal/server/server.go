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

// corsMiddleware wraps an http.HandlerFunc and adds CORS headers for local development
func corsMiddleware(handler http.HandlerFunc, localDev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if localDev {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Handle preflight OPTIONS requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		handler(w, r)
	}
}

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
		localDev,
	)
	repository := repository.NewUserRepository(dataPath)
	lockoutTable := lockout.NewTable()
	sessionTable := session.NewTable()
	userLockTable := service.NewUserLockTable()
	documentService := service.NewDocumentService(
		repository,
		recaptchaClient,
		sessionTable,
		userLockTable,
		errorLogger,
	)
	sessionService := service.NewSessionService(
		recaptchaClient,
		repository,
		sessionTable,
		userLockTable,
		errorLogger,
	)

	repository.InitialSetup()

	documentController := controller.NewDocumentController(cookieHelper, documentService, lockoutTable, requestLogger)
	http.HandleFunc(documentApiPath, corsMiddleware(documentController.Handle, localDev))

	sessionController := controller.NewSessionController(cookieHelper, sessionService, lockoutTable, requestLogger)
	http.HandleFunc(sessionApiPath, corsMiddleware(sessionController.Handle, localDev))

	if localDev {
		debugLogger.Println("CORS enabled for http://localhost:3000")
	}

	fmt.Printf("Listening on http://localhost%s\n", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
