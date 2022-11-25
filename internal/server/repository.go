package server

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ClintonMorrison/lorikeet/internal/storage"
)

type Repository struct {
	dataPath string
}

func (r *Repository) createDataDirectory() {
	err := storage.CreateDirectory(r.dataPath)

	if err != nil {
		panic(err)
	}
}

func (r *Repository) pathForDocument(auth Auth, salt []byte) (string, error) {
	signature, err := auth.Signature(salt)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s.txt", r.dataPath, signature), nil
}

func (r *Repository) pathForSalt(auth Auth) string {
	return fmt.Sprintf("%s/%s.salt.txt", r.dataPath, auth.username)
}

//
// Salt Files
//
func (r *Repository) saltFileExists(auth Auth) (bool, error) {
	filename := r.pathForSalt(auth)
	return storage.FileExists(filename)
}

func (r *Repository) writeSaltFile(auth Auth) ([]byte, error) {
	fileName := r.pathForSalt(auth)

	salt, err := makeSalt()
	if err != nil {
		return salt, err
	}

	err = ioutil.WriteFile(fileName, salt, 0644)
	if err != nil {
		return salt, err
	}

	return salt, nil
}

func (r *Repository) readSaltFile(auth Auth) ([]byte, error) {
	fileName := r.pathForSalt(auth)
	salt, err := ioutil.ReadFile(fileName)
	if err != nil {
		return salt, err
	}

	return salt, nil
}

func (r *Repository) deleteSaltFile(auth Auth) error {
	fileName := r.pathForSalt(auth)

	err := os.Remove(fileName)
	if err != nil {
		return err
	}

	return nil
}

//
// Document Files
//
func (r *Repository) documentExists(auth Auth, salt []byte) (bool, error) {
	filename, err := r.pathForDocument(auth, salt)
	if err != nil {
		return false, err
	}

	return storage.FileExists(filename)
}

func (r *Repository) writeDocument(data []byte, auth Auth, salt []byte) error {
	filename, err := r.pathForDocument(auth, salt)
	if err != nil {
		return err
	}

	saltedPassword, err := auth.SaltedPassword(salt)
	if err != nil {
		return err
	}

	encrypted, err := encrypt(data, []byte(saltedPassword))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, encrypted, 0644)
	if err != nil {
		return nil
	}

	return nil
}

func (r *Repository) readDocument(auth Auth, salt []byte) ([]byte, error) {
	data := make([]byte, 0)
	filename, err := r.pathForDocument(auth, salt)
	if err != nil {
		return data, err
	}

	data, err = ioutil.ReadFile(filename)
	if err != nil {
		return data, err
	}

	saltedPassword, err := auth.SaltedPassword(salt)

	data, err = ioutil.ReadFile(filename)
	if err != nil {
		return data, err
	}

	decrypted := decrypt(data, []byte(saltedPassword))

	return decrypted, nil
}

func (r *Repository) deleteDocument(auth Auth, salt []byte) error {
	fileName, err := r.pathForDocument(auth, salt)
	if err != nil {
		return err
	}

	err = os.Remove(fileName)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) moveDocument(currentAuth Auth, salt []byte, newAuth Auth) error {
	currentFilename, err := r.pathForDocument(currentAuth, salt)
	if err != nil {
		return err
	}

	newFilename, err := r.pathForDocument(newAuth, salt)
	if err != nil {
		return err
	}

	err = os.Rename(currentFilename, newFilename)
	if err != nil {
		return err
	}

	return nil
}
