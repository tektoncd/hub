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
	"go.uber.org/zap"
	"gopkg.in/gormigrate.v1"

	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
)

// Migrate create tables and populates master tables
func Migrate(api *app.APIConfig) error {

	log := api.Logger()

	// NOTE: when writing a migration for a new table, add the same in InitSchema
	migration := gormigrate.New(
		api.DB(),
		gormigrate.DefaultOptions,
		[]*gormigrate.Migration{
			addResourceAndResourceVersionTable(log),
			addMinimumPipelineVersionToResourceVersion(log),
			addOrgAndDropOwnerColumnInCatalog(log),
			addNotNullToResourceAndResourceVersion(log),
		},
	)

	migration.InitSchema(func(db *gorm.DB) error {
		if err := db.AutoMigrate(
			&model.Category{},
			&model.Tag{},
			&model.Catalog{},
			&model.Resource{},
			&model.ResourceVersion{},
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

		log.Info("Schema initialised successfully !!")

		return nil
	})

	if err := migration.Migrate(); err != nil {
		log.Error(err, " failed to migrate")
		return err
	}

	log.Info("Migration ran successfully !!")
	return nil
}

func fkey(log *zap.SugaredLogger, db *gorm.DB, model interface{}, args ...string) error {
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
