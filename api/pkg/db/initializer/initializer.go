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

package initializer

import (
	"context"
	"strings"

	"github.com/tektoncd/hub/api/gen/log"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"gorm.io/gorm"
)

// Initializer defines the configuration required for initailizer
// to populate the tables
type Initializer struct {
	app.Service
	api app.BaseConfig
}

// New returns the Initializer implementation.
func New(api app.BaseConfig) *Initializer {
	return &Initializer{
		Service: api.Service("initializer"),
		api:     api,
	}
}

// Run executes the func which populate the tables
func (i *Initializer) Run(ctx context.Context) (*model.Config, error) {

	db := i.DB(ctx)
	logger := i.Logger(ctx)

	config := model.Config{}
	if err := db.Model(&model.Config{}).FirstOrInit(&config).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	data := i.api.Data()
	if config.Checksum == data.Checksum {
		logger.Info("SKIP: Config refresh as config file has not changed")
		return &config, nil
	}

	updateConfig := func(db *gorm.DB, log *log.Logger, data *app.Data) error {
		config.Checksum = data.Checksum
		if err := db.Save(&config).Error; err != nil {
			log.Error(err)
			return err
		}
		return nil
	}

	if err := withTransaction(db, logger, data,
		addCategories,
		addCatalogs,
		addUsers,
		updateConfig,
	); err != nil {
		return nil, err
	}

	return &config, nil
}

type initFn func(*gorm.DB, *log.Logger, *app.Data) error

func addCategories(db *gorm.DB, log *log.Logger, data *app.Data) error {

	var configCatID []uint
	for _, c := range data.Categories {
		cat := model.Category{Name: c.Name}
		if err := db.Where(cat).FirstOrCreate(&cat).Error; err != nil {
			log.Error(err)
			return err
		}
		configCatID = append(configCatID, cat.ID)
	}

	if len(configCatID) > 0 {
		// Deletes mapping of removed categories and resources from the database
		if err := db.Unscoped().Not(map[string]interface{}{"category_id": configCatID}).Delete(&model.ResourceCategory{}).Error; err != nil {
			log.Error(err)
			return err
		}
		// Deletes categories from database which are removed from the config
		if err := db.Unscoped().Not(map[string]interface{}{"id": configCatID}).Delete(&model.Category{}).Error; err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

func addCatalogs(db *gorm.DB, log *log.Logger, data *app.Data) error {

	for _, c := range data.Catalogs {
		cat := model.Catalog{
			Name:       c.Name,
			Org:        c.Org,
			Type:       c.Type,
			URL:        c.URL,
			Revision:   c.Revision,
			ContextDir: c.ContextDir,
		}
		if err := db.Where(&model.Catalog{Name: c.Name, Org: c.Org}).FirstOrCreate(&cat).Error; err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func addUsers(db *gorm.DB, log *log.Logger, data *app.Data) error {

	for _, s := range data.Scopes {

		// Check if scopes exist or create it
		q := db.Where(&model.Scope{Name: s.Name})

		scope := model.Scope{}
		if err := q.FirstOrCreate(&scope).Error; err != nil {
			log.Error(err)
			return err
		}

		for _, userID := range s.Users {

			// Checks if user exists
			q := db.Where("LOWER(github_login) = ?", strings.ToLower(userID))

			user := model.User{}
			if err := q.First(&user).Error; err != nil {
				// If user not found then create a new record
				if err != gorm.ErrRecordNotFound {
					log.Error(err)
					return err
				}

				log.Infof("user %s not found, create a new user", userID)
				user.GithubLogin = strings.ToLower(userID)
				if err = db.Create(&user).Error; err != nil {
					log.Error(err)
					return err
				}
			}

			// Add scopes for user if not added already
			us := model.UserScope{UserID: user.ID, ScopeID: scope.ID}
			q = db.Model(&model.UserScope{}).Where(&us)

			if err := q.FirstOrCreate(&us).Error; err != nil {
				log.Error(err)
				return err
			}
		}

	}
	return nil
}

func withTransaction(db *gorm.DB, log *log.Logger, data *app.Data, fns ...initFn) error {
	txn := db.Begin()
	for _, fn := range fns {
		if err := fn(txn, log, data); err != nil {
			txn.Rollback()
			return err
		}
	}

	txn.Commit()
	return nil
}
