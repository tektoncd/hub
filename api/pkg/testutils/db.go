package testutils

import (
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
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
	db := tc.DB()
	db.AutoMigrate(model.Category{}, model.Tag{})
	if len(db.GetErrors()) > 0 {
		return db.GetErrors()[0]
	}
	return nil
}
