package repository

import "github.com/ClintonMorrison/lorikeet/internal/model"

type HybridRepository struct {
	v1 V1
	v2 V2
}

const useV2 = true

func NewHybridRepository(dataPath string) HybridRepository {
	return HybridRepository{
		v1: *NewRepositoryV1(dataPath),
		v2: *NewRepositoryV2(dataPath),
	}
}

func (r *HybridRepository) IsUsernameAvailable(auth model.Auth) (bool, error) {
	if useV2 {
		return r.v2.IsUsernameAvailable(auth)
	}
	return r.v1.IsUsernameAvailable(auth)
}

func (r *HybridRepository) CreateUser(auth model.Auth, document []byte) (*model.User, error) {
	if useV2 {
		return r.v2.CreateUser(auth, document)
	}
	return r.v1.CreateUser(auth, document)

}
func (r *HybridRepository) GetUser(auth model.Auth) (*model.User, error) {
	if useV2 {
		return r.v2.GetUser(auth)
	}
	return r.v1.GetUser(auth)

}
func (r *HybridRepository) UpdateUser(user *model.User, update model.UserUpdate) (*model.User, error) {
	if useV2 {
		return r.v2.UpdateUser(user, update)
	}
	return r.v1.UpdateUser(user, update)

}
func (r *HybridRepository) DeleteUser(user *model.User) error {
	if useV2 {
		return r.v2.DeleteUser(user)
	}
	return r.v1.DeleteUser(user)

}
