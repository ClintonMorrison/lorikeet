package repository

import (
	"fmt"
	"os"
	"time"

	"github.com/ClintonMorrison/lorikeet/internal/errors"
	"github.com/ClintonMorrison/lorikeet/internal/model"
	"github.com/ClintonMorrison/lorikeet/internal/storage"
)

type V2 struct {
	dataPath               string
	userMetadataRepository UserMetadataRepository
	saltRepository         SaltRepository
	documentRepository     DocumentRepository
}

func NewRepositoryV2(baseDataPath string) *V2 {
	dataPath := fmt.Sprintf("%s/v2", baseDataPath)
	return &V2{
		dataPath:               dataPath,
		userMetadataRepository: UserMetadataRepository{dataPath: dataPath},
		saltRepository:         SaltRepository{dataPath: dataPath},
		documentRepository:     DocumentRepository{dataPath: dataPath},
	}
}

func (r *V2) CreateDataDirectory() {
	err := storage.CreateDirectory(r.dataPath)

	if err != nil {
		panic(err)
	}
}

func (r *V2) pathForUserFolder(auth model.Auth) string {
	return pathForUserFolder(r.dataPath, auth)
}

// NEW INTERFACE
func (r *V2) IsUsernameAvailable(auth model.Auth) (bool, error) {
	exists, err := r.saltRepository.Exists(auth)
	return !exists, err
}

func (r *V2) CreateUser(auth model.Auth, document []byte) (*model.User, error) {
	// Create user folder
	err := storage.CreateDirectory(r.pathForUserFolder(auth))
	if err != nil {
		return nil, fmt.Errorf("unable to create user folder: %s", err.Error())
	}

	// Create salt
	salts, err := r.saltRepository.Create(auth)
	if err != nil {
		return nil, fmt.Errorf("unable to create salt: %s", err.Error())
	}

	// Create metadata
	now := time.Now()
	metadata := model.UserMetadata{
		SignUpTime:     now,
		LastAccessTime: now,
		StorageVersion: 2,
	}
	err = r.userMetadataRepository.CreateOrUpdate(auth, metadata)
	if err != nil {
		return nil, err
	}

	// Create document
	err = r.documentRepository.CreateOrUpdate(document, auth, salts.ServerSalt)
	if err != nil {
		return nil, err
	}

	// Load user
	user, err := r.GetUser(auth)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *V2) GetUser(auth model.Auth) (*model.User, error) {
	// Load salt
	salts, err := r.saltRepository.Get(auth)
	if err != nil {
		return nil, err
	}

	// Load metadata
	metadata, err := r.userMetadataRepository.Get(auth)
	if err != nil {
		return nil, err
	}

	// Load document
	exists, err := r.documentRepository.Exists(auth, salts.ServerSalt)
	if !exists {
		return nil, errors.INVALID_CREDENTIALS
	}

	document, err := r.documentRepository.Get(auth, salts.ServerSalt)
	if err != nil {
		return nil, err
	}

	return &model.User{
		Username:   auth.Username,
		Auth:       auth,
		Metadata:   metadata,
		ClientSalt: salts.ClientSalt,
		ServerSalt: salts.ServerSalt,
		Document:   document,
	}, nil
}

func (r *V2) DeleteUser(user *model.User) error {
	// Remove document
	err := r.documentRepository.Remove(user.Auth, user.ServerSalt)
	if err != nil {
		return err
	}

	// Remove salt
	err = r.saltRepository.Remove(user.Auth)
	if err != nil {
		return err
	}

	// Remove metadata
	err = r.userMetadataRepository.Remove(user.Auth)
	if err != nil {
		return err
	}

	// Remove user folder
	err = os.Remove(r.pathForUserFolder(user.Auth))
	if err != nil {
		return err
	}

	return nil
}

func (r *V2) UpdateUser(user *model.User, update model.UserUpdate) (*model.User, error) {
	updatedUser := user

	// Update password, if present
	if len(update.Password) > 0 {
		newAuth := model.Auth{
			Ip:       user.Auth.Ip,
			Username: user.Auth.Username,
			Password: update.Password,
		}
		err := r.documentRepository.Move(user.Auth, user.ServerSalt, newAuth)
		if err != nil {
			return nil, err
		}

		updatedUser.Auth = newAuth
	}

	// Update document, if present
	if len(update.Document) > 0 {
		err := r.documentRepository.CreateOrUpdate(update.Document, updatedUser.Auth, updatedUser.ServerSalt)
		if err != nil {
			return nil, err
		}
	}

	// Update last access time, if present
	if !update.LastAccessTime.IsZero() {
		metadata := user.Metadata
		metadata.LastAccessTime = update.LastAccessTime
		err := r.userMetadataRepository.CreateOrUpdate(updatedUser.Auth, metadata)
		if err != nil {
			return nil, err
		}
	}

	// Reload updated user
	updatedUser, err := r.GetUser(updatedUser.Auth)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
