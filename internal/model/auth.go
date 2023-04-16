package model

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ClintonMorrison/lorikeet/internal/utils"
)

type Auth struct {
	Username string
	Password string
	Ip       string
}

func (a Auth) SaltedPassword(salt []byte) (string, error) {
	data := make([]byte, 0)
	data = append(data, []byte(a.Password)...)
	data = append(data, salt...)

	return utils.Hash(data), nil
}

func (a Auth) Signature(salt []byte) (string, error) {
	data := make([]byte, 0)

	saltedPassword, err := a.SaltedPassword(salt)
	if err != nil {
		return "", err
	}

	data = append(data, []byte(a.Username)...)
	data = append(data, []byte(saltedPassword)...)

	return utils.Hash(data), nil
}

func AuthFromRequest(r *http.Request) (Auth, error) {
	username, password, ok := r.BasicAuth()

	username = strings.ToLower(username)
	ip := utils.GetIpFromRequest(r)

	if !ok {
		return Auth{}, errors.New("invalid authorization header")
	}

	return Auth{username, password, ip}, nil
}
