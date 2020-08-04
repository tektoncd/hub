package token

import (
	"github.com/dgrijalva/jwt-go"
)

// Create takes claim and jwtkey and returns a signed token
func Create(claim jwt.Claims, jwtKey string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	jwt, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return jwt, nil
}
