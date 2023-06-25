package repository

import (
	"encoding/json"
	"fmt"

	"github.com/ClintonMorrison/lorikeet/internal/model"
	"github.com/ClintonMorrison/lorikeet/internal/storage"
)

type UserMetadataRepository struct {
	dataPath string
}

func (r *UserMetadataRepository) resourceForMetadata(auth model.Auth) storage.FileResource {
	fileName := fmt.Sprintf("%s/metadata.json", pathForUserFolder(r.dataPath, auth))
	return storage.NewFileResource(fileName)
}

// Get loads metadata for the given user
func (r *UserMetadataRepository) Get(auth model.Auth) (model.UserMetadata, error) {
	metadata := &model.UserMetadata{}

	resource := r.resourceForMetadata(auth)
	data, err := resource.Read()
	if err != nil {
		return *metadata, err
	}

	err = json.Unmarshal(data, metadata)
	if err != nil {
		return *metadata, err
	}

	return *metadata, nil
}

// CreateOrUpdate updates metadata for the given user
func (r *UserMetadataRepository) CreateOrUpdate(auth model.Auth, metadata model.UserMetadata) error {
	resource := r.resourceForMetadata(auth)

	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	err = resource.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// Remove deletes saved metadata for the given user
func (r *UserMetadataRepository) Remove(auth model.Auth) error {
	resource := r.resourceForMetadata(auth)
	return resource.Remove()
}
