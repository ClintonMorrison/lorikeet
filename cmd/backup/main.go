package main

import (
	"github.com/ClintonMorrison/lorikeet/internal/storage"
	"time"
	"fmt"
	"os"
	"github.com/ClintonMorrison/lorikeet/internal/config"
)

func main() {
	// Make directory for backups
	err := os.Mkdir(config.BACKUP_PATH, 0700)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	timeStamp := time.Now().Format(time.RFC3339)
	outputPath := fmt.Sprintf("%s/backup-%s.tar.gz", config.BACKUP_PATH, timeStamp)

	err = storage.ToTarball(config.DATA_PATH, outputPath)
	if err != nil {
		panic(err)
	}
}
