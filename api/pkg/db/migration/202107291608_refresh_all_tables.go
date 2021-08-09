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

package migration

import (
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/tektoncd/hub/api/gen/log"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"gorm.io/gorm"
)

// This migration collects the existing data from db, deletes the existing tables,
// recreates the tables back and inserts the data to update the db and makes it
// compatible with `gorm 2.0`. Also removing `removeCatgoryIdColumnAndConstraintsFromTagtable`
// migration because refreshAllTables handles the case by dropping the tables and
// creating the tables back
func refreshAllTables(log *log.Logger) *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202107291608_refresh_all_tables",
		Migrate: func(db *gorm.DB) error {
			txn := db.Begin()
			err := migrateDB(txn, log)
			if err != nil {
				txn.Rollback()
				log.Error(err)
				return err
			}
			txn.Commit()
			return nil
		},
	}
}

func migrateDB(txn *gorm.DB, log *log.Logger) error {
	var user []model.User
	if err := txn.Find(&user).Error; err != nil {
		log.Error(err)
		return err
	}

	var users_last_value uint
	if err := txn.Table("users_id_seq").Pluck("last_value", &users_last_value).Error; err != nil {
		log.Error(err)
		return err
	}

	var user_scopes []model.UserScope
	if err := txn.Find(&user_scopes).Error; err != nil {
		log.Error(err)
		return err
	}

	var user_rating []model.UserResourceRating
	if err := txn.Find(&user_rating).Error; err != nil {
		log.Error(err)
		return err
	}

	var user_resource_ratings_last_value uint
	if err := txn.Table("user_resource_ratings_id_seq").Pluck("last_value", &user_resource_ratings_last_value).Error; err != nil {
		log.Error(err)
		return err
	}

	var categories []model.Category
	if err := txn.Find(&categories).Error; err != nil {
		log.Error(err)
		return err
	}

	var categories_last_value uint
	if err := txn.Table("categories_id_seq").Pluck("last_value", &categories_last_value).Error; err != nil {
		log.Error(err)
		return err
	}

	var tags []model.Tag
	if err := txn.Find(&tags).Error; err != nil {
		log.Error(err)
		return err
	}

	var tags_last_value uint
	if err := txn.Table("tags_id_seq").Pluck("last_value", &tags_last_value).Error; err != nil {
		log.Error(err)
		return err
	}

	var sync_job []model.SyncJob
	if err := txn.Find(&sync_job).Error; err != nil {
		log.Error(err)
		return err
	}

	var sync_jobs_last_value uint
	if err := txn.Table("sync_jobs_id_seq").Pluck("last_value", &sync_jobs_last_value).Error; err != nil {
		log.Error(err)
		return err
	}

	var scopes []model.Scope
	if err := txn.Find(&scopes).Error; err != nil {
		log.Error(err)
		return err
	}

	var scopes_last_value uint
	if err := txn.Table("scopes_id_seq").Pluck("last_value", &scopes_last_value).Error; err != nil {
		log.Error(err)
		return err
	}

	var resources []model.Resource
	if err := txn.Find(&resources).Error; err != nil {
		log.Error(err)
		return err
	}

	var resources_last_value uint
	if err := txn.Table("resources_id_seq").Pluck("last_value", &resources_last_value).Error; err != nil {
		log.Error(err)
		return err
	}

	var resource_versions []model.ResourceVersion
	if err := txn.Find(&resource_versions).Error; err != nil {
		log.Error(err)
		return err
	}

	var resource_versions_last_value uint
	if err := txn.Table("resource_versions_id_seq").Pluck("last_value", &resource_versions_last_value).Error; err != nil {
		log.Error(err)
		return err
	}

	var resource_tags []model.ResourceTag
	if err := txn.Find(&resource_tags).Error; err != nil {
		log.Error(err)
		return err
	}

	var configs []model.Config
	if err := txn.Find(&configs).Error; err != nil {
		log.Error(err)
		return err
	}

	var configs_last_value uint
	if err := txn.Table("configs_id_seq").Pluck("last_value", &configs_last_value).Error; err != nil {
		log.Error(err)
		return err
	}

	var catalog []model.Catalog
	if err := txn.Find(&catalog).Error; err != nil {
		log.Error(err)
		return err
	}

	var catalogs_last_value uint
	if err := txn.Table("catalogs_id_seq").Pluck("last_value", &catalogs_last_value).Error; err != nil {
		log.Error(err)
		return err
	}

	var catalog_error []model.CatalogError
	if err := txn.Find(&catalog_error).Error; err != nil {
		log.Error(err)
		return err
	}

	var catalog_errors_last_value uint
	if err := txn.Table("catalog_errors_id_seq").Pluck("last_value", &catalog_errors_last_value).Error; err != nil {
		log.Error(err)
		return err
	}

	if err := txn.Migrator().DropTable(
		&model.Category{},
		&model.Tag{},
		&model.Catalog{},
		&model.CatalogError{},
		&model.Resource{},
		&model.ResourceVersion{},
		&model.User{},
		&model.UserResourceRating{},
		&model.SyncJob{},
		&model.Scope{},
		&model.UserScope{},
		&model.Config{},
		&model.ResourceTag{},
	); err != nil {
		log.Error(err)
		return err
	}

	if err := txn.AutoMigrate(
		&model.Category{},
		&model.Tag{},
		&model.Catalog{},
		&model.CatalogError{},
		&model.Resource{},
		&model.ResourceVersion{},
		&model.User{},
		&model.UserResourceRating{},
		&model.SyncJob{},
		&model.Scope{},
		&model.UserScope{},
		&model.Config{},
		&model.ResourceTag{},
	); err != nil {
		log.Error(err)
		return err
	}

	if len(catalog) > 0 {
		if err := txn.Create(
			&catalog,
		).Error; err != nil {
			log.Error(err)
			return err
		}

		query := fmt.Sprintf("ALTER SEQUENCE catalogs_id_seq RESTART WITH %d", catalogs_last_value+1)
		if err := txn.Exec(query).Error; err != nil {
			log.Error(err)
			return err
		}
	}

	if len(catalog_error) > 0 {
		if err := txn.Create(
			&catalog_error,
		).Error; err != nil {
			log.Error(err)
			return err
		}

		query := fmt.Sprintf("ALTER SEQUENCE catalog_errors_id_seq RESTART WITH %d", catalog_errors_last_value+1)
		if err := txn.Exec(query).Error; err != nil {
			log.Error(err)
			return err
		}
	}

	if len(resources) > 0 {
		if err := txn.Create(
			&resources,
		).Error; err != nil {
			log.Error(err)
			return err
		}

		query := fmt.Sprintf("ALTER SEQUENCE resources_id_seq RESTART WITH %d", resources_last_value+1)
		if err := txn.Exec(query).Error; err != nil {
			log.Error(err)
			return err
		}
	}

	if len(resource_versions) > 0 {
		if err := txn.Create(
			&resource_versions,
		).Error; err != nil {
			log.Error(err)
			return err
		}

		query := fmt.Sprintf("ALTER SEQUENCE resource_versions_id_seq RESTART WITH %d", resource_versions_last_value+1)
		if err := txn.Exec(query).Error; err != nil {
			log.Error(err)
			return err
		}
	}

	if len(tags) > 0 {
		if err := txn.Create(
			&tags,
		).Error; err != nil {
			log.Error(err)
			return err
		}

		query := fmt.Sprintf("ALTER SEQUENCE tags_id_seq RESTART WITH %d", tags_last_value+1)
		if err := txn.Exec(query).Error; err != nil {
			log.Error(err)
			return err
		}
	}

	if len(categories) > 0 {
		if err := txn.Create(
			&categories,
		).Error; err != nil {
			log.Error(err)
			return err
		}

		query := fmt.Sprintf("ALTER SEQUENCE categories_id_seq RESTART WITH %d", categories_last_value+1)
		if err := txn.Exec(query).Error; err != nil {
			log.Error(err)
			return err
		}
	}

	if len(user) > 0 {
		if err := txn.Create(
			&user,
		).Error; err != nil {
			log.Error(err)
			return err
		}

		query := fmt.Sprintf("ALTER SEQUENCE users_id_seq RESTART WITH %d", users_last_value+1)
		if err := txn.Exec(query).Error; err != nil {
			log.Error(err)
			return err
		}
	}

	if len(user_rating) > 0 {
		if err := txn.Create(
			&user_rating,
		).Error; err != nil {
			log.Error(err)
			return err
		}

		query := fmt.Sprintf("ALTER SEQUENCE user_resource_ratings_id_seq RESTART WITH %d", user_resource_ratings_last_value+1)
		if err := txn.Exec(query).Error; err != nil {
			log.Error(err)
			return err
		}
	}

	if len(sync_job) > 0 {
		if err := txn.Create(
			&sync_job,
		).Error; err != nil {
			log.Error(err)
			return err
		}

		query := fmt.Sprintf("ALTER SEQUENCE sync_jobs_id_seq RESTART WITH %d", sync_jobs_last_value+1)
		if err := txn.Exec(query).Error; err != nil {
			log.Error(err)
			return err
		}
	}

	if len(scopes) > 0 {
		if err := txn.Create(
			&scopes,
		).Error; err != nil {
			log.Error(err)
			return err
		}

		query := fmt.Sprintf("ALTER SEQUENCE scopes_id_seq RESTART WITH %d", scopes_last_value+1)
		if err := txn.Exec(query).Error; err != nil {
			log.Error(err)
			return err
		}
	}

	if len(user_scopes) > 0 {
		if err := txn.Create(
			&user_scopes,
		).Error; err != nil {
			log.Error(err)
			return err
		}
	}

	if len(configs) > 0 {
		if err := txn.Create(
			&configs,
		).Error; err != nil {
			log.Error(err)
			return err
		}

		query := fmt.Sprintf("ALTER SEQUENCE configs_id_seq RESTART WITH %d", configs_last_value+1)
		if err := txn.Exec(query).Error; err != nil {
			log.Error(err)
			return err
		}
	}

	if len(resource_tags) > 0 {
		if err := txn.Create(
			&resource_tags,
		).Error; err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}
