package repository

import "github.com/ClintonMorrison/lorikeet/internal/model"

type UserRepository interface {
	IsUsernameAvailable(auth model.Auth) (bool, error)
	CreateUser(model.Auth, []byte) (*model.User, error)
	GetUser(model.Auth) (*model.User, error)
	UpdateUser(*model.User, model.UserUpdate) (*model.User, error)
	DeleteUser(*model.User) error
}
