/*
Copyright 2022 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.

You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/tektoncd/hub/api/gen/log"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"gorm.io/gorm"
)

// This migration backups user details from `users` table into `user_backups` table
// drops a few columns from `users` table and create `account` table which will
// be linked to `users` table.
func addUsersDetailsInAccountTable(log *log.Logger) *gormigrate.Migration {

	return &gormigrate.Migration{
		ID: "202111091037_backup_users_add_account_table_and_update_data",
		Migrate: func(db *gorm.DB) error {

			// Backup users into user_backups table so that if in case there is any error
			// while migrating then our data is already there.
			if err := db.Exec("CREATE TABLE user_backups AS SELECT * FROM users;").Error; err != nil {
				log.Error(err)
				return err
			}

			var users []model.UserBackup
			if err := db.Find(&users).Error; err != nil {
				log.Error(err)
				return err
			}

			// Update the user table based on the model
			if err := db.Migrator().DropColumn(model.User{}, "github_login"); err != nil {
				log.Error(err)
				return err
			}

			if err := db.Migrator().DropColumn(model.User{}, "github_name"); err != nil {
				log.Error(err)
				return err
			}

			if err := db.Migrator().DropColumn(model.User{}, "avatar_url"); err != nil {
				log.Error(err)
				return err
			}

			// Create the account table
			if err := db.AutoMigrate(&model.Account{}); err != nil {
				log.Error(err)
				return err
			}

			var accounts []model.Account
			for _, user := range users {
				account := model.Account{
					UserID:   user.ID,
					UserName: user.GithubLogin,
					Name:     user.GithubName,
					// Assumption taken over here is that all the existing users were from github
					Provider: "github",
				}
				accounts = append(accounts, account)
			}

			// Add user details in account table
			if err := db.Create(&accounts).Error; err != nil {
				log.Error(err)
				return err
			}

			return nil
		},
	}
}
