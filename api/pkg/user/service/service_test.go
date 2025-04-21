package user

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestValidateScopes(t *testing.T) {

	claims := make(jwt.MapClaims)
	claims["scopes"] = []interface{}{"rating:read", "rating:write"}

	scheme := &JWTScheme{
		Name:           "",
		Scopes:         []string{"rating:read", "rating:write", "agent:create", "catalog:refresh", "config:refresh", "refresh:token"},
		RequiredScopes: []string{"rating:read", "rating:write"},
	}

	err := ValidateScopes(claims, scheme)

	assert.Equal(t, err, nil)
}

func TestInValidateScopes(t *testing.T) {
	claims := make(jwt.MapClaims)

	scheme := &JWTScheme{
		Name:           "",
		Scopes:         []string{"rating:read", "rating:write", "agent:create", "catalog:refresh", "config:refresh", "refresh:token"},
		RequiredScopes: []string{"rating:read", "rating:write"},
	}

	err := ValidateScopes(claims, scheme)

	assert.Equal(t, err.Error(), "invalid scopes")
}

func TestValidateMissingScopes(t *testing.T) {

	claims := make(jwt.MapClaims)
	claims["scopes"] = []interface{}{"abc:foo", "qwe:baar"}

	scheme := &JWTScheme{
		Name:           "",
		Scopes:         []string{"rating:read", "rating:write", "agent:create", "catalog:refresh", "config:refresh", "refresh:token"},
		RequiredScopes: []string{"rating:read", "rating:write"},
	}

	err := ValidateScopes(claims, scheme)

	assert.Equal(t, err.Error(), "missing scopes: rating:read, rating:write")
}
