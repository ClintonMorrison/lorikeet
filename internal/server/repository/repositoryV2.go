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

type V2 struct {
	dataPath string
}

func NewRepositoryV2(dataPath string) *V1 {
	return &V1{dataPath}
}

func (r *V2) CreateDataDirectory() {
	err := storage.CreateDirectory(r.dataPath)

	if err != nil {
		panic(err)
	}
}

func (r *V2) pathForUserFolder(auth model.Auth) string {
	return fmt.Sprintf("%s/%s", r.dataPath, auth.Username)
}

func (r *V2) pathForDocument(auth model.Auth, salt []byte) (string, error) {
	signature, err := auth.Signature(salt)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s.data.txt", r.pathForUserFolder(auth), signature), nil
}

func (r *V2) pathForSalt(auth model.Auth) string {
	return fmt.Sprintf("%s/salt.txt", r.pathForUserFolder(auth))
}

//
// Salt Files
//
func (r *V2) saltFileExists(auth model.Auth) (bool, error) {
	filename := r.pathForSalt(auth)
	return storage.FileExists(filename)
}

func (r *V2) writeSaltFile(auth model.Auth) ([]byte, error) {
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

func (r *V2) readSaltFile(auth model.Auth) ([]byte, error) {
	fileName := r.pathForSalt(auth)
	salt, err := ioutil.ReadFile(fileName)
	if err != nil {
		return salt, err
	}

	return salt, nil
}

func (r *V2) deleteSaltFile(auth model.Auth) error {
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
func (r *V2) documentExists(auth model.Auth, salt []byte) (bool, error) {
	filename, err := r.pathForDocument(auth, salt)
	if err != nil {
		return false, err
	}

	return storage.FileExists(filename)
}

func (r *V2) writeDocument(data []byte, auth model.Auth, salt []byte) error {
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

func (r *V2) readDocument(auth model.Auth, salt []byte) ([]byte, error) {
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

func (r *V2) deleteDocument(auth model.Auth, salt []byte) error {
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

func (r *V2) moveDocument(currentAuth model.Auth, salt []byte, newAuth model.Auth) error {
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
func (r *V2) IsUsernameAvailable(auth model.Auth) (bool, error) {
	exists, err := r.saltFileExists(auth)
	return !exists, err
}

func (r *V2) CreateUser(auth model.Auth, document []byte) (*model.User, error) {
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

func (r *V2) GetUser(auth model.Auth) (*model.User, error) {
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
		Username: auth.Username,
		Auth:     auth,
		Metadata: model.UserMetadata{StorageVersion: 2},
		Salt:     salt,
		Document: document,
	}, nil

}

func (r *V2) DeleteUser(user *model.User) error {
	err := r.deleteDocument(user.Auth, user.Salt)
	if err != nil {
		return err
	}

	err = r.deleteSaltFile(user.Auth)
	if err != nil {
		return err
	}

	return nil
}

func (r *V2) UpdateUser(user *model.User, update model.UserUpdate) (*model.User, error) {
	updatedUser := user

	if len(update.Password) > 0 {
		newAuth := model.Auth{
			Ip:       user.Auth.Ip,
			Username: user.Auth.Username,
			Password: update.Password,
		}
		err := r.moveDocument(user.Auth, user.Salt, newAuth)
		if err != nil {
			return nil, err
		}

		updatedUser.Auth = newAuth
	}

	if len(update.Document) > 0 {
		err := r.writeDocument(update.Document, updatedUser.Auth, updatedUser.Salt)
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
