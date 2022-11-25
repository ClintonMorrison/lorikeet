package main

import (
	"fmt"
	"os"

	"github.com/ClintonMorrison/lorikeet/internal/config"
	"github.com/ClintonMorrison/lorikeet/internal/server"
)

func main() {
	recaptchaSecret := os.Getenv("LORIKEET_RECAPTCHA_SECRET")

	if recaptchaSecret == "" {
		fmt.Println("Environment variable 'LORIKEET_RECAPTCHA_SECRET' is not set")
		os.Exit(1)
	}

	server.Run(
		config.DATA_PATH,
		config.SERVER_ADDRESS,
		os.Getenv("LORIKEET_RECAPTCHA_SECRET"),
		config.LOG_PATH,
		config.REQUEST_LOG_PATH,
		config.ERROR_LOG_PATH)
}
