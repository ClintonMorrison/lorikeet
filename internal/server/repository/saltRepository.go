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

type UserSalts struct {
	ClientSalt []byte
	ServerSalt []byte
}

func (r *SaltRepository) pathForServerSalt(auth model.Auth) string {
	return fmt.Sprintf("%s/server.salt.txt", pathForUserFolder(r.dataPath, auth))
}

func (r *SaltRepository) resourceForServerSalt(auth model.Auth) storage.FileResource {
	return storage.NewFileResource(r.pathForServerSalt(auth))
}

func (r *SaltRepository) pathForClientSalt(auth model.Auth) string {
	return fmt.Sprintf("%s/client.salt.txt", pathForUserFolder(r.dataPath, auth))
}

func (r *SaltRepository) resourceForClientSalt(auth model.Auth) storage.FileResource {
	return storage.NewFileResource(r.pathForServerSalt(auth))
}

// Create generates a new salt file
func (r *SaltRepository) Create(auth model.Auth) (UserSalts, error) {
	salts := UserSalts{}

	// Generate server salt
	serverResource := r.resourceForServerSalt(auth)
	serverSalt, err := utils.MakeSalt()
	if err != nil {
		return UserSalts{}, err
	}

	err = serverResource.Write(serverSalt)
	if err != nil {
		return salts, err
	}

	// Generate client salt
	clientResource := r.resourceForClientSalt(auth)
	clientSalt, err := utils.MakeSalt()
	if err != nil {
		return salts, err
	}

	err = clientResource.Write(clientSalt)
	if err != nil {
		return salts, err
	}

	return UserSalts{ClientSalt: clientSalt, ServerSalt: serverSalt}, nil
}

// Create generates new salt files
func (r *SaltRepository) Exists(auth model.Auth) (bool, error) {
	serverResource := r.resourceForServerSalt(auth)
	serverExists, err := serverResource.Exists()
	if err != nil {
		return false, err
	}

	clientResource := r.resourceForClientSalt(auth)
	clientExists, err := clientResource.Exists()
	if err != nil {
		return false, err
	}

	return serverExists && clientExists, nil
}

// Get loads an existing salt files
func (r *SaltRepository) Get(auth model.Auth) (UserSalts, error) {
	serverResource := r.resourceForServerSalt(auth)
	serverSalt, err := serverResource.Read()
	if err != nil {
		return UserSalts{}, err
	}

	clientResource := r.resourceForClientSalt(auth)
	clientSalt, err := clientResource.Read()
	if err != nil {
		return UserSalts{}, err
	}

	return UserSalts{ClientSalt: clientSalt, ServerSalt: serverSalt}, nil
}

// Delete removes an existing salt files
func (r *SaltRepository) Remove(auth model.Auth) error {
	serverResource := r.resourceForServerSalt(auth)
	err := serverResource.Remove()
	if err != nil {
		return err
	}

	clientResource := r.resourceForClientSalt(auth)
	err = clientResource.Remove()
	if err != nil {
		return err
	}

	return nil
}
