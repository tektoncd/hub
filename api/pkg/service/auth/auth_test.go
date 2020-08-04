// Copyright © 2020 The Tekton Authors.
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

package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"

	"github.com/tektoncd/hub/api/gen/auth"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

func TestLogin(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	defer gock.Off()

	gock.New("https://github.com").
		Post("/login/oauth/access_token").
		Reply(200).
		JSON(map[string]string{
			"access_token": "test-token",
		})

	gock.New("https://api.github.com").
		Get("/user").
		Reply(200).
		JSON(map[string]string{
			"login": "test",
			"name":  "test-user",
		})

	authSvc := New(tc)
	payload := &auth.AuthenticatePayload{Code: "test-code"}
	res, err := authSvc.Authenticate(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, validToken, res.Token)
	assert.Equal(t, gock.IsDone(), true)
}

func TestLogin_again(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	defer gock.Off()

	gock.New("https://github.com").
		Post("/login/oauth/access_token").
		Reply(200).
		JSON(map[string]string{
			"access_token": "test-token",
		})

	gock.New("https://api.github.com").
		Get("/user").
		Times(2).
		Reply(200).
		JSON(map[string]string{
			"login": "test",
			"name":  "test-user",
		})

	authSvc := New(tc)
	payload := &auth.AuthenticatePayload{Code: "test-code"}
	res, err := authSvc.Authenticate(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, validToken, res.Token)

	gock.New("https://github.com").
		Post("/login/oauth/access_token").
		Reply(200).
		JSON(map[string]string{
			"access_token": "test-token-2",
		})
	payloadAgain := &auth.AuthenticatePayload{Code: "test-code-2"}
	resAgain, err := authSvc.Authenticate(context.Background(), payloadAgain)
	assert.NoError(t, err)
	assert.Equal(t, validToken, resAgain.Token)
	assert.Equal(t, gock.IsDone(), true)
}

func TestLogin_InvalidCode(t *testing.T) {
	tc := testutils.Setup(t)

	defer gock.Off()

	gock.New("https://github.com").
		Post("/login/oauth/access_token").
		Reply(200).
		SetError(errors.New("oauth2: server response missing access_token"))

	authSvc := New(tc)
	payload := &auth.AuthenticatePayload{Code: "test-code"}
	_, err := authSvc.Authenticate(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid authorization code")
	assert.Equal(t, gock.IsDone(), true)
}
