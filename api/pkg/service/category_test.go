package hub

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	category "github.com/tektoncd/hub/api/gen/category"
	app "github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/model"
)

var (
	categorySvc category.Service
	db          *gorm.DB
)

func TestMain(m *testing.M) {

	testConfig, err := app.TestConfigFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: failed to initialise: %s", err)
		os.Exit(1)
	}

	db = testConfig.DB()
	logger := testConfig.Logger()

	db.AutoMigrate(model.Category{}, model.Tag{})

	categorySvc = NewCategory(db, logger)

	defer os.Exit(m.Run())
	defer testConfig.Cleanup()
}

// LoadFixture ...
func LoadFixture(db *gorm.DB, fixtureDir string) error {
	fixtures, err := testfixtures.New(
		testfixtures.Database(db.DB()),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(fixtureDir),
	)
	if err != nil {
		return err
	}
	if err := fixtures.Load(); err != nil {
		return err
	}
	return nil
}

func Test_All(t *testing.T) {
	LoadFixture(db, "../../fixtures")
	all, err := categorySvc.All(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, len(all), 3)
	assert.Equal(t, all[0].Name, "abc") // categories are sorted by name
}
