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

package testutils

import (
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/migration"
)

// LoadFixtures is called before executing each test, it clears db and
// loads data from fixtures so that each test is executed on new db
func LoadFixtures(t *testing.T, dir string) {
	tc := Config()
	fixtures, err := testfixtures.New(
		testfixtures.Database(tc.DB().DB()),
		testfixtures.Dialect(app.DBDialect),
		testfixtures.Directory(dir))
	assert.NoError(t, err)
	assert.NoError(t, fixtures.Load())
}

// applyMigration creates tables in test db
func applyMigration() error {
	tc := Config()
	logger := tc.Logger("test")
	if err := migration.Migrate(tc.APIBase); err != nil {
		logger.Errorf("DB initialisation failed !!")
		return err
	}
	logger.Info("DB initialisation successful !!")
	return nil
}
