package crypt

import (
	"golang.org/x/crypto/bcrypt"
)

// Encrypts a string using bcrypt and returns the hash
func Encrypt(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(hash), err
}

// Compares a hash and a string to determine if the hashed string would match the hash.
// Returns 'matches' and possible errors.
func Compare(hash string, str string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	if err != nil {
		return false, err
	}
	return true, err
}


