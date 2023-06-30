package repository

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ClintonMorrison/lorikeet/internal/errors"
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

func (r *V1) InitialSetup() {
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
func (r *V1) saltFileExists(auth model.Auth) (bool, error) {
	filename := r.pathForSalt(auth)
	return storage.FileExists(filename)
}

func (r *V1) writeSaltFile(auth model.Auth) ([]byte, error) {
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

func (r *V1) readSaltFile(auth model.Auth) ([]byte, error) {
	fileName := r.pathForSalt(auth)
	salt, err := ioutil.ReadFile(fileName)
	if err != nil {
		return salt, err
	}

	return salt, nil
}

func (r *V1) deleteSaltFile(auth model.Auth) error {
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
func (r *V1) documentExists(auth model.Auth, salt []byte) (bool, error) {
	filename, err := r.pathForDocument(auth, salt)
	if err != nil {
		return false, err
	}

	return storage.FileExists(filename)
}

func (r *V1) writeDocument(data []byte, auth model.Auth, salt []byte) error {
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

func (r *V1) readDocument(auth model.Auth, salt []byte) ([]byte, error) {
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

func (r *V1) deleteDocument(auth model.Auth, salt []byte) error {
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

func (r *V1) moveDocument(currentAuth model.Auth, salt []byte, newAuth model.Auth) error {
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

// NEW INTERFACE
func (r *V1) IsUsernameAvailable(auth model.Auth) (bool, error) {
	exists, err := r.saltFileExists(auth)
	return !exists, err
}

func (r *V1) CreateUser(auth model.Auth, document []byte) (*model.User, error) {
	salt, err := r.writeSaltFile(auth)
	if err != nil {
		return nil, err
	}

	err = r.writeDocument(document, auth, salt)
	if err != nil {
		return nil, err
	}

	user, err := r.GetUser(auth)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *V1) GetUser(auth model.Auth) (*model.User, error) {
	salt, err := r.readSaltFile(auth)
	if err != nil {
		return nil, err
	}

	exists, err := r.documentExists(auth, salt)
	if !exists {
		return nil, errors.INVALID_CREDENTIALS
	}

	document, err := r.readDocument(auth, salt)
	if err != nil {
		return nil, err
	}

	return &model.User{
		Username:   auth.Username,
		Auth:       auth,
		Metadata:   model.UserMetadata{StorageVersion: 1},
		ServerSalt: salt,
		Document:   document,
	}, nil

}

func (r *V1) DeleteUser(user *model.User) error {
	err := r.deleteDocument(user.Auth, user.ServerSalt)
	if err != nil {
		return err
	}

	err = r.deleteSaltFile(user.Auth)
	if err != nil {
		return err
	}

	return nil
}

func (r *V1) UpdateUser(user *model.User, update model.UserUpdate) (*model.User, error) {
	updatedUser := user

	if len(update.Password) > 0 {
		newAuth := model.Auth{
			Ip:       user.Auth.Ip,
			Username: user.Auth.Username,
			Password: update.Password,
		}
		err := r.moveDocument(user.Auth, user.ServerSalt, newAuth)
		if err != nil {
			return nil, err
		}

		updatedUser.Auth = newAuth
	}

	if len(update.Document) > 0 {
		err := r.writeDocument(update.Document, updatedUser.Auth, updatedUser.ServerSalt)
		if err != nil {
			return nil, err
		}
	}

	updatedUser, err := r.GetUser(updatedUser.Auth)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// TODO: add way to get sign up and last accessed dates from legacy user files?
// Could be based on last modified dates
// https://www.includehelp.com/golang/print-the-last-modified-time-of-an-existing-file.aspx#:~:text=In%20the%20main()%20function,using%20the%20ModTime()%20function.
