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

// The tests validates the constraints on the tables which are added through
// migrations.
package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

// Checks the Not Null Constraint
func TestCheckNotNull(t *testing.T) {
	tc := testutils.Setup(t)

	db := tc.DB()

	err := db.Create(&model.Catalog{Name: "tekton", Type: "", URL: "", Revision: "master"}).Error
	assert.Error(t, err)
	assert.Equal(t, "pq: null value in column \"type\" violates not-null constraint", err.Error())

	err = db.Create(&model.Resource{Name: "tekton", Rating: 4}).Error
	assert.Error(t, err)
	assert.Equal(t, "pq: null value in column \"type\" violates not-null constraint", err.Error())

	err = db.Create(&model.ResourceVersion{Version: "", Description: "task", URL: "", DisplayName: "Task", MinPipelinesVersion: ""}).Error
	assert.Error(t, err)
	assert.Equal(t, "pq: null value in column \"version\" violates not-null constraint", err.Error())
}

// Checks the Unique constraint
func TestCheckUnique(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	db := tc.DB()

	err := db.Create(&model.Catalog{Name: "catalog-official", Org: "tektoncd", Type: "tektoncd", URL: "url", Revision: "master"}).Error
	assert.Error(t, err)
	assert.Equal(t, "pq: duplicate key value violates unique constraint \"uix_name_org\"", err.Error())
}
