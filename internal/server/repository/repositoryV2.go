package repository

import (
	"fmt"

	"github.com/ClintonMorrison/lorikeet/internal/errors"
	"github.com/ClintonMorrison/lorikeet/internal/model"
	"github.com/ClintonMorrison/lorikeet/internal/storage"
	"github.com/ClintonMorrison/lorikeet/internal/utils"
)

type V2 struct {
	dataPath string
}

func NewRepositoryV2(baseDataPath string) *V1 {
	dataPath := fmt.Sprintf("%s/v2", baseDataPath)
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

func (r *V2) resourceForDocument(auth model.Auth, salt []byte) (storage.FileResource, error) {
	fileName, err := r.pathForDocument(auth, salt)
	if err != nil {
		return storage.FileResource{}, err
	}

	return storage.NewFileResource(fileName), nil
}

func (r *V2) pathForSalt(auth model.Auth) string {
	return fmt.Sprintf("%s/salt.txt", r.pathForUserFolder(auth))
}

func (r *V2) resourceForSalt(auth model.Auth) storage.FileResource {
	return storage.NewFileResource(r.pathForSalt(auth))
}

//
// Salt Files
//
func (r *V2) writeSaltFile(auth model.Auth) ([]byte, error) {
	resource := r.resourceForSalt(auth)

	salt, err := utils.MakeSalt()
	if err != nil {
		return salt, err
	}

	err = resource.Write(salt)
	if err != nil {
		return salt, err
	}

	return salt, nil
}

//
// Document Files
//
func (r *V2) writeDocument(data []byte, auth model.Auth, salt []byte) error {
	resource, err := r.resourceForDocument(auth, salt)
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

	err = resource.Write(encrypted)
	if err != nil {
		return nil
	}

	return nil
}

func (r *V2) readDocument(auth model.Auth, salt []byte) ([]byte, error) {
	data := make([]byte, 0)
	resource, err := r.resourceForDocument(auth, salt)
	if err != nil {
		return data, err
	}

	data, err = resource.Read()
	if err != nil {
		return data, err
	}

	saltedPassword, err := auth.SaltedPassword(salt)
	if err != nil {
		return data, err
	}

	decrypted := utils.Decrypt(data, []byte(saltedPassword))

	return decrypted, nil
}

func (r *V2) moveDocument(currentAuth model.Auth, salt []byte, newAuth model.Auth) error {
	resource, err := r.resourceForDocument(currentAuth, salt)
	if err != nil {
		return err
	}

	newFilename, err := r.pathForDocument(newAuth, salt)
	if err != nil {
		return err
	}

	err = resource.MoveTo(newFilename)
	if err != nil {
		return err
	}

	return nil
}

// NEW INTERFACE
func (r *V2) IsUsernameAvailable(auth model.Auth) (bool, error) {
	resource := r.resourceForSalt(auth)
	exists, err := resource.Exists()
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
	saltResource := r.resourceForSalt(auth)
	salt, err := saltResource.Read()
	if err != nil {
		return nil, err
	}

	documentResource, err := r.resourceForDocument(auth, salt)
	if err != nil {
		return nil, err
	}

	exists, err := documentResource.Exists()
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
	documentResource, err := r.resourceForDocument(user.Auth, user.Salt)
	if err != nil {
		return err
	}

	saltResource := r.resourceForSalt(user.Auth)

	err = documentResource.Remove()
	if err != nil {
		return err
	}

	err = saltResource.Remove()
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
