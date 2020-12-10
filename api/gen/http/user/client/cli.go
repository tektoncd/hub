// Code generated by goa v3.2.2, DO NOT EDIT.
//
// user HTTP client CLI support package
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package client

import (
	user "github.com/tektoncd/hub/api/gen/user"
)

// BuildRefreshAccessTokenPayload builds the payload for the user
// RefreshAccessToken endpoint from CLI flags.
func BuildRefreshAccessTokenPayload(userRefreshAccessTokenRefreshToken string) (*user.RefreshAccessTokenPayload, error) {
	var refreshToken string
	{
		refreshToken = userRefreshAccessTokenRefreshToken
	}
	v := &user.RefreshAccessTokenPayload{}
	v.RefreshToken = refreshToken

	return v, nil
}

// BuildNewRefreshTokenPayload builds the payload for the user NewRefreshToken
// endpoint from CLI flags.
func BuildNewRefreshTokenPayload(userNewRefreshTokenRefreshToken string) (*user.NewRefreshTokenPayload, error) {
	var refreshToken string
	{
		refreshToken = userNewRefreshTokenRefreshToken
	}
	v := &user.NewRefreshTokenPayload{}
	v.RefreshToken = refreshToken

	return v, nil
}
