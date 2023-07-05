package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ClintonMorrison/lorikeet/internal/config"
	"github.com/ClintonMorrison/lorikeet/internal/storage"
)

func main() {

	tarballPath := flag.String("file", "", "tarball to extract backup from")
	flag.Parse()

	if *tarballPath == "" {
		flag.Usage()
		return
	}

	// Ensure backup exists
	exists, err := storage.FileExists(*tarballPath)
	if !exists {
		panic("Tarball file does not exist")
	}

	// Move old data directory
	timeStamp := time.Now().Format(time.RFC3339)
	tmpDataPath := fmt.Sprintf("%s/data-old-%s", config.BASE_PATH, timeStamp)
	os.Rename(config.DATA_PATH, tmpDataPath)

	// Make new directory for data
	err = os.Mkdir(config.DATA_PATH, 0700)
	if err != nil {
		panic(err)
	}

	// Extract backup to data directory
	err = storage.FromTarball(config.DATA_PATH, *tarballPath)
	if err != nil {
		panic(err)
	}
}
