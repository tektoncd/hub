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

func createPlatformTables(log *app.Logger) *gormigrate.Migration {

	return &gormigrate.Migration{
		ID: "202106281500_create_platform_tables",
		Migrate: func(db *gorm.DB) error {
			if err := db.AutoMigrate(&model.Platform{}, &model.ResourceVersion{}, &model.Resource{}); err != nil {
				log.Error(err)
				return err
			}
			return nil
		},
	}
}
