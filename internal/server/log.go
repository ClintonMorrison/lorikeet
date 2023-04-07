package server

import (
	"io"
	"log"
	"os"
)

func createLogger(filename string, prefix string) (*log.Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0700)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	writer := io.MultiWriter(os.Stdout, file)
	return log.New(writer, prefix, log.LstdFlags), nil
}
