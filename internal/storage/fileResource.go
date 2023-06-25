package storage

import (
	"fmt"
	"io/ioutil"
	"os"
)

// A helper for handling files
type FileResource struct {
	fileName string
}

func NewFileResource(fileName string) FileResource {
	return FileResource{fileName: fileName}
}

func (fr *FileResource) Exists() (bool, error) {
	return FileExists(fr.fileName)
}

func (fr *FileResource) Write(data []byte) error {
	return ioutil.WriteFile(fr.fileName, data, 0644)
}

func (fr *FileResource) Read() ([]byte, error) {
	return ioutil.ReadFile(fr.fileName)
}

func (fr *FileResource) Remove() error {
	return os.Remove(fr.fileName)
}

func (fr *FileResource) MoveTo(newFileName string) error {
	exists, err := fr.Exists()
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("cannot move resource because is does not exist")
	}

	err = os.Rename(fr.fileName, newFileName)
	if err == nil {
		fr.fileName = newFileName
	}

	return err
}
