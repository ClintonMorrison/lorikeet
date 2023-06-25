package repository

import (
	"fmt"

	"github.com/ClintonMorrison/lorikeet/internal/model"
	"github.com/ClintonMorrison/lorikeet/internal/storage"
	"github.com/ClintonMorrison/lorikeet/internal/utils"
)

type SaltRepository struct {
	dataPath string
}

func (r *SaltRepository) pathForSalt(auth model.Auth) string {
	return fmt.Sprintf("%s/salt.txt", pathForUserFolder(r.dataPath, auth))
}

func (r *SaltRepository) resourceForSalt(auth model.Auth) storage.FileResource {
	return storage.NewFileResource(r.pathForSalt(auth))
}

// Create generates a new salt file
func (r *SaltRepository) Create(auth model.Auth) ([]byte, error) {
	resource := r.resourceForSalt(auth)

	// Generate salt
	salt, err := utils.MakeSalt()
	if err != nil {
		return salt, err
	}

	// Save salt in file
	err = resource.Write(salt)
	if err != nil {
		return salt, err
	}

	return salt, nil
}

// Create generates a new salt file
func (r *SaltRepository) Exists(auth model.Auth) (bool, error) {
	resource := r.resourceForSalt(auth)
	return resource.Exists()
}

// Get loads an existing salt file
func (r *SaltRepository) Get(auth model.Auth) ([]byte, error) {
	resource := r.resourceForSalt(auth)
	return resource.Read()
}

// Delete removes an existing salt file
func (r *SaltRepository) Remove(auth model.Auth) error {
	resource := r.resourceForSalt(auth)
	return resource.Remove()
}
