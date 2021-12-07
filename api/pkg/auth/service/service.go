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

package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/markbates/goth"
	"github.com/tektoncd/hub/api/pkg/auth/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/token"
	"gorm.io/gorm"
)

func (r *request) userScopes(account *model.Account) ([]string, error) {

	var userScopes []string = r.defaultScopes

	scopes := []model.Scope{}

	if err := r.db.Model(&model.Scope{}).Joins("JOIN user_scopes as u on scopes.id=u.scope_id").
		Where("u.user_id = ?", account.UserID).Find(&scopes).Error; err != nil {
		r.log.Error(err)
		return nil, err
	}

	for _, s := range scopes {
		userScopes = append(userScopes, s.Name)
	}

	return userScopes, nil
}

func (r *request) createTokens(user *model.User, scopes []string, provider string) (*app.AuthenticateResult, error) {

	req := token.Request{
		User:      user,
		Scopes:    scopes,
		JWTConfig: r.jwtConfig,
		Provider:  provider,
	}

	accessToken, accessExpiresAt, err := req.AccessJWT()
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	refreshToken, refreshExpiresAt, err := req.RefreshJWT()
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	user.RefreshTokenChecksum = createChecksum(refreshToken)

	if err = r.db.Save(user).Error; err != nil {
		r.log.Error(err)
		return nil, err
	}

	data := &app.AuthTokens{
		Access: &app.Token{
			Token:           accessToken,
			RefreshInterval: r.jwtConfig.AccessExpiresIn.String(),
			ExpiresAt:       accessExpiresAt,
		},
		Refresh: &app.Token{
			Token:           refreshToken,
			RefreshInterval: r.jwtConfig.RefreshExpiresIn.String(),
			ExpiresAt:       refreshExpiresAt,
		},
	}

	return &app.AuthenticateResult{Data: data}, nil
}

func createChecksum(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

/*
(Keeping email as the primary key in users table)
We perform search query in users table on the basis of email:
If Email is not present(existing user or new user)
	query the accounts table on the basis of username and provider
		- If the above case matches -> update the email in users table and other details respectively in accounts table
		- If the above case returns empty record(new user)
			- create a new user in the database
If Email is already present
	fetch the user_id from users table and query the accounts table with `where` clause of `user_id` and `provider`
		- if the record doesn't exists(user trying to login with new provider)
			then create a new record in accounts table with new provider and the user_id already associated with the email
		- if record exists then update the record if there are any changes
*/
func (r *request) insertData(gitUser goth.User, code, provider string) error {

	var acc model.Account
	var user model.User

	userQuery := r.db.Model(&model.User{}).
		Where("email = ?", gitUser.Email)

	// Check if user exist
	err := userQuery.First(&user).Error

	// If email doesn't exists in users table
	if err != nil {
		// Check whether username and provider are matching
		accountQuery := r.db.Model(&model.Account{}).Where("user_name = ?", gitUser.NickName).Where("provider = ?", provider)
		err = accountQuery.First(&acc).Error

		// If user doesn't exist, create a new record
		if err == gorm.ErrRecordNotFound {
			if user, err = r.insertIntoUsersTable(gitUser, code); err != nil {
				r.log.Error(err)
			}

			if err = r.insertIntoAccountsTable(gitUser, provider, user.ID); err != nil {
				r.log.Error(err)
				return err
			}
		} else {
			// Account exists
			// Update the user table with the email
			if err := r.db.Model(&model.User{}).Where("id = ?", acc.UserID).
				Updates(model.User{Code: code, Email: gitUser.Email, Type: model.NormalUserType}).Error; err != nil {
				r.log.Error(err)
				return err
			}

			// Update the AvatarUrl and Name in Accounts table
			if err := updateAccountDetails(accountQuery, acc,
				model.Account{AvatarURL: gitUser.AvatarURL, Name: gitUser.Name}); err != nil {
				r.log.Error(err)
				return err
			}
		}
	} else { // when the email of user already exists
		// Update the users table with the auth code
		if err := userQuery.Update("code", code).Error; err != nil {
			r.log.Error(err)
			return err
		}

		// Check for the account on the basis of user_id and provider
		accountQuery := r.db.Model(&model.Account{}).Where("user_id = ?", user.ID).Where("provider = ?", provider)
		err = accountQuery.First(&acc).Error

		// If not found then create a new entry in accounts table
		if err == gorm.ErrRecordNotFound {
			if err = r.insertIntoAccountsTable(gitUser, provider, user.ID); err != nil {
				r.log.Error(err)
				return err
			}
			return nil
		} else if err != nil {
			r.log.Error(err)
			return err
		}

		// If account found then update the details of the user
		if err := updateAccountDetails(accountQuery, acc,
			model.Account{AvatarURL: gitUser.AvatarURL, Name: gitUser.Name, UserName: gitUser.NickName}); err != nil {
			r.log.Error(err)
			return err
		}

	}

	return nil
}

func updateAccountDetails(accountQuery *gorm.DB, existingAccountDetails model.Account, newAccountDetails model.Account) error {

	detailsToUpdate := model.Account{}
	if newAccountDetails.UserName != "" && newAccountDetails.UserName != existingAccountDetails.UserName {
		detailsToUpdate.UserName = newAccountDetails.UserName
	}

	if newAccountDetails.Name != "" && newAccountDetails.Name != existingAccountDetails.Name {
		detailsToUpdate.Name = newAccountDetails.Name
	}

	if newAccountDetails.AvatarURL != "" && newAccountDetails.AvatarURL != existingAccountDetails.AvatarURL {
		detailsToUpdate.AvatarURL = newAccountDetails.AvatarURL
	}

	if detailsToUpdate.UserName != "" || detailsToUpdate.Name != "" || detailsToUpdate.AvatarURL != "" {
		if err := accountQuery.Updates(detailsToUpdate).Error; err != nil {
			return err
		}
	}

	return nil
}

// Creates a new record in Users table
func (r *request) insertIntoUsersTable(gitUser goth.User, code string) (model.User, error) {
	user := model.User{
		Code:  code,
		Email: gitUser.Email,
		Type:  model.NormalUserType,
	}

	// User ID is by default set to zero so we need to update the value with (last inserted user_id+1)
	lastUser := model.User{}
	if err := r.db.Model(&model.User{}).Last(&lastUser).Error; err != nil {
		r.log.Error(err)
		return model.User{}, err
	}
	user.ID = lastUser.ID + 1

	err := r.db.Create(&user).Error
	if err != nil {
		r.log.Error(err)
		return model.User{}, err
	}
	return user, nil
}

// Creates a new record in Accounts table
func (r *request) insertIntoAccountsTable(gitUser goth.User, provider string, userId uint) error {
	acc := model.Account{
		Name:      gitUser.Name,
		UserName:  gitUser.NickName,
		AvatarURL: gitUser.AvatarURL,
		Provider:  provider,
		UserID:    userId,
	}

	if err := r.db.Create(&acc).Error; err != nil {
		r.log.Error(err)
		return err
	}
	return nil
}
