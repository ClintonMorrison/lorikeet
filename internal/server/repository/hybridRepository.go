package repository

import (
	"fmt"

	"github.com/ClintonMorrison/lorikeet/internal/model"
)

type HybridRepository struct {
	v1 V1
	v2 V2
}

const useV2 = true

// Uses the v2 repository if user is on v2
// Falls back to v1 if user exists but is only on v1
//
// New users are created with v2

func NewHybridRepository(dataPath string) *HybridRepository {
	return &HybridRepository{
		v1: *NewRepositoryV1(dataPath),
		v2: *NewRepositoryV2(dataPath),
	}
}

func (r *HybridRepository) IsUsernameAvailable(auth model.Auth) (bool, error) {
	availableInV1, err := r.v1.IsUsernameAvailable(auth)
	if err != nil {
		return false, err
	}

	availableInV2, err := r.v2.IsUsernameAvailable(auth)
	if err != nil {
		return false, err
	}

	return availableInV1 && availableInV2, nil
}

// All new users are created on v2
func (r *HybridRepository) CreateUser(auth model.Auth, document []byte) (*model.User, error) {
	return r.v2.CreateUser(auth, document)
}

// Get from v2 if username is present there, otherwise use v1
func (r *HybridRepository) GetUser(auth model.Auth) (*model.User, error) {
	existsOnV2, err := r.userExistsOnV2(auth)
	if err != nil {
		return nil, err
	}

	if existsOnV2 {
		return r.v2.GetUser(auth)
	}

	return r.v1.GetUser(auth)
}

// Update user on v2 if they are already on v2
func (r *HybridRepository) UpdateUser(user *model.User, update model.UserUpdate) (*model.User, error) {
	existsOnV2, err := r.userExistsOnV2(user.Auth)
	if err != nil {
		return nil, err
	}

	if existsOnV2 {
		return r.v2.UpdateUser(user, update)
	}

	return r.v1.UpdateUser(user, update)

}

// Delete user from v2 if they are on v2 (all v1 data will be deleted soon)
func (r *HybridRepository) DeleteUser(user *model.User) error {
	existsOnV2, err := r.userExistsOnV2(user.Auth)
	if err != nil {
		return err
	}

	if existsOnV2 {
		return r.v2.DeleteUser(user)
	}

	return r.v1.DeleteUser(user)

}

func (r *HybridRepository) MigrateUser(user *model.User) (*model.User, error) {
	// Check if user exists on v2 already (do nothing if so)
	existsOnV2, err := r.userExistsOnV2(user.Auth)
	if err != nil {
		return nil, err
	}

	if existsOnV2 {
		return user, nil
	}

	// Verify user exists on v1
	existsOnV1, err := r.userExistsOnV1(user.Auth)
	if err != nil {
		return nil, err
	}

	if !existsOnV1 {
		return nil, fmt.Errorf("Cannot migrate because user not on v1")
	}

	// Recreate user on v2
	newUser, err := r.v2.CreateUser(user.Auth, user.Document)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *HybridRepository) InitialSetup() {
	r.v1.InitialSetup()
	r.v2.InitialSetup()
}

func (r *HybridRepository) userExistsOnV2(auth model.Auth) (bool, error) {
	available, err := r.v2.IsUsernameAvailable(auth)
	return !available, err
}

func (r *HybridRepository) userExistsOnV1(auth model.Auth) (bool, error) {
	available, err := r.v1.IsUsernameAvailable(auth)
	return !available, err
}
