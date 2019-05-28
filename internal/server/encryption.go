package server

import (
"crypto/aes"
"crypto/cipher"
"crypto/sha256"
"crypto/rand"
"encoding/hex"
"io"
)

func makeSalt() ([]byte, error) {
	salt := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return salt, err
	}

	return salt, nil
}

func hash(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, hashedPassword []byte) ([]byte, error) {
	block, err := aes.NewCipher(hashedPassword[:32])
	if err != nil {
		panic(err)
		return make([]byte, 0), err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return make([]byte, 0), err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(data []byte, hashedPassword []byte) []byte {
	key := []byte(hashedPassword)[:32]
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}