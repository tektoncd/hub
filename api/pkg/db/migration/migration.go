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

package migration

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"

	"github.com/tektoncd/hub/api/gen/log"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
)

// Migrate create tables and populates master tables
func Migrate(api *app.APIBase) error {

	log := api.Logger("migration")

	// NOTE: when writing a migration for a new table, add the same in InitSchema
	migration := gormigrate.New(
		api.DB(),
		gormigrate.DefaultOptions,
		[]*gormigrate.Migration{
			// NOTE: Checkout the migration template in migration_template.go.
			// Create a new file for a new migration and
			// add the migration function here.
		},
	)

	migration.InitSchema(func(db *gorm.DB) error {
		if err := db.AutoMigrate(
			&model.Category{},
			&model.Tag{},
			&model.Catalog{},
			&model.Resource{},
			&model.ResourceVersion{},
			&model.User{},
			&model.UserResourceRating{},
			&model.Scope{},
		).Error; err != nil {
			log.Error(err)
			return err
		}

		if err := fkey(log, db, model.Tag{}, "category_id", "categories"); err != nil {
			return err
		}

		if err := fkey(log, db, model.Resource{}, "catalog_id", "catalogs"); err != nil {
			return err
		}

		if err := fkey(log, db, model.ResourceVersion{}, "resource_id", "resources"); err != nil {
			return err
		}

		if err := fkey(log, db, model.ResourceTag{},
			"resource_id", "resources",
			"tag_id", "tags"); err != nil {
			return err
		}

		if err := fkey(log, db, model.UserResourceRating{},
			"resource_id", "resources",
			"user_id", "users"); err != nil {
			return err
		}

		if err := fkey(log, db, model.UserScope{},
			"user_id", "users",
			"scope_id", "scopes"); err != nil {
			return err
		}

		log.Info("Schema initialised successfully !!")

		if err := addScopes(log, db); err != nil {
			log.Error(err)
			return err
		}

		log.Info("Scopes added successfully !!")

		return nil
	})

	if err := migration.Migrate(); err != nil {
		log.Error(err, " failed to migrate")
		return err
	}

	log.Info("Migration ran successfully !!")
	return nil
}

func fkey(log *log.Logger, db *gorm.DB, model interface{}, args ...string) error {
	for i := 0; i < len(args); i += 2 {
		col := args[i]
		table := args[i+1]
		err := db.Model(model).AddForeignKey(col, table, "CASCADE", "CASCADE").Error
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func addScopes(log *log.Logger, db *gorm.DB) error {

	scopes := []string{"agent:create"}

	for _, s := range scopes {
		sc := &model.Scope{Name: s}
		if err := db.Create(sc).Error; err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}
