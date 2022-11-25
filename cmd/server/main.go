package main

import (
	"github.com/ClintonMorrison/lorikeet/internal/config"
	"github.com/ClintonMorrison/lorikeet/internal/server"
)

func main() {
	server.Run(
		config.DATA_PATH,
		config.SERVER_ADDRESS,
		"", // TODO: use environment variable DO NOT COMMIT
		config.LOG_PATH,
		config.REQUEST_LOG_PATH,
		config.ERROR_LOG_PATH)
}
