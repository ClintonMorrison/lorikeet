package repository

import (
	"fmt"

	"github.com/ClintonMorrison/lorikeet/internal/model"
	"github.com/ClintonMorrison/lorikeet/internal/storage"
	"github.com/ClintonMorrison/lorikeet/internal/utils"
)

type DocumentRepository struct {
	dataPath string
}

func (r *DocumentRepository) pathForDocument(auth model.Auth, salt []byte) (string, error) {
	signature, err := auth.Signature(salt)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s.data.txt", pathForUserFolder(r.dataPath, auth), signature), nil
}

func (r *DocumentRepository) resourceForDocument(auth model.Auth, salt []byte) (storage.FileResource, error) {
	fileName, err := r.pathForDocument(auth, salt)
	if err != nil {
		return storage.FileResource{}, err
	}

	return storage.NewFileResource(fileName), nil
}

// CreateOrUpdate writes the given document
func (r *DocumentRepository) CreateOrUpdate(data []byte, auth model.Auth, salt []byte) error {
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
		return err
	}

	return nil
}

func (r *DocumentRepository) Get(auth model.Auth, salt []byte) ([]byte, error) {
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

	decrypted, err := utils.Decrypt(data, []byte(saltedPassword))
	if err != nil {
		return data, err
	}

	return decrypted, nil
}

func (r *DocumentRepository) Exists(auth model.Auth, salt []byte) (bool, error) {
	resource, err := r.resourceForDocument(auth, salt)
	if err != nil {
		return false, err
	}

	return resource.Exists()
}

func (r *DocumentRepository) Move(currentAuth model.Auth, salt []byte, newAuth model.Auth) error {
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

func (r *DocumentRepository) Remove(auth model.Auth, salt []byte) error {
	resource, err := r.resourceForDocument(auth, salt)
	if err != nil {
		return err
	}

	return resource.Remove()
}
