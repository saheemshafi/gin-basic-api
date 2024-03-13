package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashBytes), err
}

func ComparePasswordHashes(password string, hashString string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashString), []byte(password)) == nil
}
