package model

import "time"

type User struct {
	Username   string
	Auth       Auth
	Metadata   UserMetadata
	ClientSalt []byte // only used for client-side encryption
	ServerSalt []byte // only used for server-side encryption
	Document   []byte
}

type UserUpdate struct {
	Document       []byte
	Password       string
	LastAccessTime time.Time
}

type UserMetadata struct {
	SignUpTime           time.Time `json:"signUpTime"`
	LastAccessTime       time.Time `json:"lastAccessTime"`
	ClientStorageVersion int       `json:"clientStorageVersion"` // 1 = current
	ServerStorageVersion int       `json:"serverStorageVersion"` // 1 = legacy, 2 = new
}
