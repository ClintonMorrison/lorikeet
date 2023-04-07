package repository

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ClintonMorrison/lorikeet/internal/model"
	"github.com/ClintonMorrison/lorikeet/internal/storage"
	"github.com/ClintonMorrison/lorikeet/internal/utils"
)

type V1 struct {
	dataPath string
}

func NewRepositoryV1(dataPath string) *V1 {
	return &V1{dataPath}
}

func (r *V1) CreateDataDirectory() {
	err := storage.CreateDirectory(r.dataPath)

	if err != nil {
		panic(err)
	}
}

func (r *V1) pathForDocument(auth model.Auth, salt []byte) (string, error) {
	signature, err := auth.Signature(salt)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s.txt", r.dataPath, signature), nil
}

func (r *V1) pathForSalt(auth model.Auth) string {
	return fmt.Sprintf("%s/%s.salt.txt", r.dataPath, auth.Username)
}

//
// Salt Files
//
func (r *V1) SaltFileExists(auth model.Auth) (bool, error) {
	filename := r.pathForSalt(auth)
	return storage.FileExists(filename)
}

func (r *V1) WriteSaltFile(auth model.Auth) ([]byte, error) {
	fileName := r.pathForSalt(auth)

	salt, err := utils.MakeSalt()
	if err != nil {
		return salt, err
	}

	err = ioutil.WriteFile(fileName, salt, 0644)
	if err != nil {
		return salt, err
	}

	return salt, nil
}

func (r *V1) ReadSaltFile(auth model.Auth) ([]byte, error) {
	fileName := r.pathForSalt(auth)
	salt, err := ioutil.ReadFile(fileName)
	if err != nil {
		return salt, err
	}

	return salt, nil
}

func (r *V1) DeleteSaltFile(auth model.Auth) error {
	fileName := r.pathForSalt(auth)

	err := os.Remove(fileName)
	if err != nil {
		return err
	}

	return nil
}

//
// Document Files
//
func (r *V1) DocumentExists(auth model.Auth, salt []byte) (bool, error) {
	filename, err := r.pathForDocument(auth, salt)
	if err != nil {
		return false, err
	}

	return storage.FileExists(filename)
}

func (r *V1) WriteDocument(data []byte, auth model.Auth, salt []byte) error {
	filename, err := r.pathForDocument(auth, salt)
	if err != nil {
		return err
	}

	saltedPassword, err := auth.SaltedPassword(salt)
	if err != nil {
		return err
	}

	encrypted, err := utils.Encrypt(data, []byte(saltedPassword))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, encrypted, 0644)
	if err != nil {
		return nil
	}

	return nil
}

func (r *V1) ReadDocument(auth model.Auth, salt []byte) ([]byte, error) {
	data := make([]byte, 0)
	filename, err := r.pathForDocument(auth, salt)
	if err != nil {
		return data, err
	}

	data, err = ioutil.ReadFile(filename)
	if err != nil {
		return data, err
	}

	saltedPassword, err := auth.SaltedPassword(salt)

	data, err = ioutil.ReadFile(filename)
	if err != nil {
		return data, err
	}

	decrypted := utils.Decrypt(data, []byte(saltedPassword))

	return decrypted, nil
}

func (r *V1) DeleteDocument(auth model.Auth, salt []byte) error {
	fileName, err := r.pathForDocument(auth, salt)
	if err != nil {
		return err
	}

	err = os.Remove(fileName)
	if err != nil {
		return err
	}

	return nil
}

func (r *V1) MoveDocument(currentAuth model.Auth, salt []byte, newAuth model.Auth) error {
	currentFilename, err := r.pathForDocument(currentAuth, salt)
	if err != nil {
		return err
	}

	newFilename, err := r.pathForDocument(newAuth, salt)
	if err != nil {
		return err
	}

	err = os.Rename(currentFilename, newFilename)
	if err != nil {
		return err
	}

	return nil
}
