package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func EncodeJWT(id string, expires time.Time) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        id,
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func DecodeJWT(tokenString string) (*jwt.RegisteredClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, errors.New("failed to decode token")
	}

	claims, ok := token.Claims.(jwt.RegisteredClaims)

	if !token.Valid || !ok {
		return nil, errors.New("token is not valid")
	}

	return &claims, nil
}
