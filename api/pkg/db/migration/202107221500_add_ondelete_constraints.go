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
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"gorm.io/gorm"
)

// Drops existing foreign key constraints and create the same constraints
// in order to add another OnDelete constraints to foriegn key
func addOnDeleteConstraints(log *app.Logger) *gormigrate.Migration {

	return &gormigrate.Migration{
		ID: "202107221500_add_ondelete_constraints",
		Migrate: func(db *gorm.DB) error {
			txn := db.Begin()

			err := addOnDelete(txn, log)
			if err != nil {
				txn.Rollback()
				return err
			}
			txn.Commit()
			return nil
		},
	}
}

func addOnDelete(txn *gorm.DB, log *app.Logger) error {
	if err := txn.Migrator().DropConstraint(&model.ResourceCategory{}, "fk_resource_categories_resource"); err != nil {
		log.Error(err)
		return err
	}

	if err := txn.Exec("ALTER TABLE resource_categories ADD CONSTRAINT fk_resource_categories_resource  FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE;").Error; err != nil {
		log.Error(err)
		return err
	}

	if err := txn.Migrator().DropConstraint(&model.ResourceTag{}, "fk_resource_tags_resource"); err != nil {
		log.Error(err)
		return err
	}

	if err := txn.Exec("ALTER TABLE  resource_tags ADD CONSTRAINT fk_resource_tags_resource FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE;").Error; err != nil {
		log.Error(err)
		return err
	}

	if err := txn.Migrator().DropConstraint(&model.ResourceVersion{}, "fk_resource_versions_resource"); err != nil {
		log.Error(err)
		return err
	}

	if err := txn.Exec("ALTER TABLE resource_versions ADD CONSTRAINT fk_resource_versions_resource FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE;").Error; err != nil {
		log.Error(err)
		return err
	}

	if err := txn.Migrator().DropConstraint(&model.ResourceVersion{}, "fk_resources_versions"); err != nil {
		log.Error(err)
		return err
	}

	if err := txn.Exec("ALTER TABLE resource_versions ADD CONSTRAINT fk_resources_versions FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE;").Error; err != nil {
		log.Error(err)
		return err
	}

	if err := txn.Migrator().DropConstraint(&model.UserResourceRating{}, "fk_user_resource_ratings_resource"); err != nil {
		log.Error(err)
		return err
	}

	if err := txn.Exec("ALTER TABLE user_resource_ratings ADD CONSTRAINT fk_user_resource_ratings_resource FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE;").Error; err != nil {
		log.Error(err)
		return err
	}

	if txn.Migrator().HasConstraint(&model.ResourcePlatform{}, "fk_resource_platforms_platform") {
		err := txn.Migrator().DropConstraint(&model.ResourcePlatform{}, "fk_resource_platforms_platform")
		if err != nil {
			log.Error(err)
			return err
		}
	}

	if err := txn.Exec("ALTER TABLE resource_platforms ADD CONSTRAINT fk_resource_platforms_platform FOREIGN KEY (platform_id) REFERENCES platforms(id) ON DELETE CASCADE;").Error; err != nil {
		log.Error(err)
		return err
	}

	if txn.Migrator().HasConstraint(&model.ResourcePlatform{}, "fk_resource_platforms_resource") {
		err := txn.Migrator().DropConstraint(&model.ResourcePlatform{}, "fk_resource_platforms_resource")
		if err != nil {
			log.Error(err)
			return err
		}
	}

	if err := txn.Exec("ALTER TABLE resource_platforms ADD CONSTRAINT fk_resource_platforms_resource FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE;").Error; err != nil {
		log.Error(err)
		return err
	}

	if txn.Migrator().HasConstraint(&model.VersionPlatform{}, "fk_version_platforms_platform") {
		err := txn.Migrator().DropConstraint(&model.VersionPlatform{}, "fk_version_platforms_platform")
		if err != nil {
			log.Error(err)
			return err
		}
	}

	if err := txn.Exec("ALTER TABLE version_platforms ADD CONSTRAINT fk_version_platforms_platform FOREIGN KEY (platform_id) REFERENCES platforms(id) ON DELETE CASCADE;").Error; err != nil {
		log.Error(err)
		return err
	}

	if txn.Migrator().HasConstraint(&model.VersionPlatform{}, "fk_version_platforms_resource_version") {
		err := txn.Migrator().DropConstraint(&model.VersionPlatform{}, "fk_version_platforms_resource_version")
		if err != nil {
			log.Error(err)
			return err
		}
	}

	if err := txn.Exec("ALTER TABLE version_platforms ADD CONSTRAINT fk_version_platforms_resource_version FOREIGN KEY (resource_version_id) REFERENCES resource_versions(id) ON DELETE CASCADE;").Error; err != nil {
		log.Error(err)
		return err
	}

	return nil
}
