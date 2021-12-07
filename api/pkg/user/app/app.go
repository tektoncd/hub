// Copyright © 2021 The Tekton Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

type InfoResult struct {
	// User Information
	Data *UserData `json:"data"`
}

// Git user Information
type UserData struct {
	// Username of User
	UserName string `json:"user_name"`
	// Name of user
	Name string `json:"name"`
	// User's profile picture url
	AvatarURL string `json:"avatarUrl"`
}

// RefreshAccessTokenResult is the result type of the user service
// RefreshAccessToken method.
type RefreshAccessTokenResult struct {
	// User Access JWT
	Data *AccessToken `json:"data"`
}

// NewRefreshTokenResult is the result type of the user service NewRefreshToken
// method.
type NewRefreshTokenResult struct {
	// User Refresh JWT
	Data *RefreshToken `json:"data"`
}

// Token includes the JWT, Expire Duration & Time
type Token struct {
	// JWT
	Token string `json:"token"`
	// Duration the token will Expire In
	RefreshInterval string `json:"refreshInterval"`
	// Time the token will expires at
	ExpiresAt int64 `json:"expiresAt"`
}

// Access Token for user
type AccessToken struct {
	// Access Token for user
	Access *Token `json:"access"`
}

// Refresh Token for User
type RefreshToken struct {
	// Refresh Token for user
	Refresh *Token `json:"refresh"`
}
