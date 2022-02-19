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

package status

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/ikawaha/goahttpcheck"
	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/gen/http/status/server"
	"github.com/tektoncd/hub/api/gen/log"
	"github.com/tektoncd/hub/api/gen/status"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/testutils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gotest.tools/v3/golden"
)

type statusTestConfig struct {
	*testutils.TestConfig
	service app.Service
	db      *gorm.DB
}

func newStatusTestConfig(t *testing.T) *statusTestConfig {

	tc := testutils.Setup(t)

	db, err := gorm.Open(postgres.Open(tc.Database().ConnectionString()), &gorm.Config{})
	assert.NoError(t, err)

	return &statusTestConfig{
		TestConfig: tc,
		db:         db,
		service:    &fakeService{db: db},
	}
}
func (tc statusTestConfig) Service(n string) app.Service {
	return tc.service
}

type fakeService struct {
	db *gorm.DB
}

var _ app.Service = (*fakeService)(nil)

func (fs *fakeService) Logger(ctx context.Context) *log.Logger {
	return nil
}

func (fs *fakeService) LoggerWith(ctx context.Context, args ...interface{}) *log.Logger {
	return nil
}

func (fs *fakeService) DB(ctx context.Context) *gorm.DB {
	return fs.db
}

func (fs *fakeService) CatalogClonePath() string {
	return os.Getenv("CLONE_BASE_PATH")
}

func TestOk_http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	checker := goahttpcheck.New()
	checker.Mount(
		server.NewStatusHandler,
		server.MountStatusHandler,
		status.NewStatusEndpoint(New(tc)),
	)

	checker.Test(t, http.MethodGet, "/").Check().
		HasStatus(http.StatusOK).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestDB_NotOK(t *testing.T) {

	tc := newStatusTestConfig(t)

	checker := goahttpcheck.New()
	checker.Mount(
		server.NewStatusHandler,
		server.MountStatusHandler,
		status.NewStatusEndpoint(New(tc)),
	)

	db, err := tc.db.DB()
	assert.NoError(t, err)
	db.Close()

	checker.Test(t, http.MethodGet, "/").Check().
		HasStatus(http.StatusOK).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})

	// ensure the db connnection is still intact for other tests to execute
	db, err = testutils.Config().DB().DB()
	assert.NoError(t, err)
	assert.NoError(t, db.Ping())
}
