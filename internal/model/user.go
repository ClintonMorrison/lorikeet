package model

import "time"

type User struct {
	Username string
	Auth     Auth
	Metadata UserMetadata
	Salt     []byte
	Document []byte
}

type UserUpdate struct {
	Document       []byte
	Password       string
	LastAccessTime time.Time
}

type UserMetadata struct {
	SignUpTime     time.Time `json:"signUpTime"`
	LastAccessTime time.Time `json:"lastAccessTime"`
	StorageVersion int       `json:"storageVersion"` // 1 = legacy, 2 = new
}
