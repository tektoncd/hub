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
	"fmt"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"goa.design/goa/v3/security"

	"github.com/tektoncd/hub/api/gen/auth"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/token"
)

type contextKey string

var (
	userIDKey = contextKey("user-id")
)

var (
	tokenError  = auth.MakeInvalidToken(fmt.Errorf("invalid user token"))
	scopesError = auth.MakeInvalidScopes(fmt.Errorf("user not authorized"))
)

type Validator struct {
	DB     *gorm.DB
	Logger *zap.SugaredLogger
	JWTKey string
}

// JWTAuth implements the authorization logic for services for the "jwt" security scheme.
func (v *Validator) JWTAuth(ctx context.Context, jwt string, scheme *security.JWTScheme) (context.Context, error) {

	claims, err := token.Verify(jwt, v.JWTKey)
	if err != nil {
		return ctx, tokenError
	}

	err = token.ValidateScopes(claims, scheme)
	if err != nil {
		return ctx, scopesError
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		return ctx, tokenError
	}

	return WithUserID(ctx, uint(userID)), nil
}

// UserFromContext fetch user id from the passed context verfies if it exists in db
// returns the User object
func (v *Validator) UserFromContext(ctx context.Context) (*model.User, error) {
	userID := UserID(ctx)
	log := v.Logger.With("user-id", userID)

	var user model.User
	if err := v.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warnf("user not found for token %s", err.Error())
			return nil, tokenError
		}

		log.Errorf("error when looking up user. err: %s", err.Error())
		return nil, internalError
	}

	return &user, nil
}

// WithUserID adds user id in context passed to it
func WithUserID(ctx context.Context, id uint) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// UserID fetch the user id from passed context
func UserID(ctx context.Context) uint {
	return ctx.Value(userIDKey).(uint)
}
