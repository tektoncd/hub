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

type Provider struct {
	Name string `json:"name"`
}

type ProviderList struct {
	Data []Provider `json:"data"`
}

type AuthenticateResult struct {
	// User Tokens
	Data *AuthTokens `json:"data"`
}

// Auth tokens have access and refresh token for user
type AuthTokens struct {
	// Access Token
	Access *Token `json:"access"`
	// Refresh Token
	Refresh *Token `json:"refresh"`
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
