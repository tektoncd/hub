package token

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"goa.design/goa/v3/security"
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

// Verify takes jwt and key and verify if jwt is valid
func Verify(token string, jwtKey string) (jwt.MapClaims, error) {

	claims := make(jwt.MapClaims)

	// Parse JWT token
	_, err := jwt.ParseWithClaims(token, claims,
		func(_ *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// ValidateScopes takes user scopes and checks if it has the scope which
// is required for accessing the api
func ValidateScopes(claims jwt.MapClaims, scheme *security.JWTScheme) error {

	if claims["scopes"] == nil {
		return fmt.Errorf("invalid scopes")
	}

	scopes, ok := claims["scopes"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid scopes")
	}

	scopesInToken := make([]string, len(scopes))
	for _, scp := range scopes {
		scopesInToken = append(scopesInToken, scp.(string))
	}
	if err := scheme.Validate(scopesInToken); err != nil {
		return err
	}

	return nil
}
