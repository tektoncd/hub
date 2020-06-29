package model

import (
	"github.com/jinzhu/gorm"
	"github.com/tektoncd/hub/api/pkg/app"
	"gopkg.in/gormigrate.v1"
)

// Migrate create tables and populates master tables
func Migrate(api *app.ApiConfig) error {

	logger := api.Logger()

	// NOTE: If writing a migration for a new table then add the same in InitSchema
	migration := gormigrate.New(api.DB(), gormigrate.DefaultOptions, []*gormigrate.Migration{
		// NOTE: Add Migration Here
	})

	migration.InitSchema(func(db *gorm.DB) error {
		if err := db.AutoMigrate(
			&Category{},
			&Tag{},
			&Catalog{},
			&Resource{},
			&ResourceVersion{},
		).Error; err != nil {
			return err
		}

		logger.Info("Schema initialised successfully !!")

		fkey := func(model interface{}, args ...string) error {
			for i := 0; i < len(args); i += 2 {
				col := args[i]
				table := args[i+1]
				err := db.Model(model).AddForeignKey(col, table, "CASCADE", "CASCADE").Error
				if err != nil {
					return err
				}
			}
			return nil
		}

		if err := fkey(Tag{}, "category_id", "categories"); err != nil {
			return err
		}

		if err := fkey(Resource{}, "catalog_id", "catalogs"); err != nil {
			return err
		}

		if err := fkey(ResourceVersion{}, "resource_id", "resources"); err != nil {
			return err
		}
		if err := fkey(ResourceTag{}, "resource_id", "resources", "tag_id", "tags"); err != nil {
			return err
		}

		initialiseTables(db)

		logger.Info("Data added successfully !!")

		return nil
	})

	if err := migration.Migrate(); err != nil {
		logger.Error(err, "could not migrate")
		return err
	}

	logger.Info("Migration ran successfully !!")

	return nil
}

// Initialise category table with data and associate to tag table
func initialiseTables(db *gorm.DB) {
	var categories = map[string][]string{
		"Others":         []string{},
		"Build Tools":    []string{"build-tool"},
		"CLI":            []string{"cli"},
		"Cloud":          []string{"gcp", "aws", "azure", "cloud"},
		"Deploy":         []string{"deploy"},
		"Image Build":    []string{"image-build"},
		"Notification":   []string{"notification"},
		"Test Framework": []string{"test"},
	}

	for name, tags := range categories {
		cat := &Category{Name: name}
		db.Create(cat)

		for _, tag := range tags {
			db.Model(&cat).Association("Tags").Append(&Tag{Name: tag})
		}
	}
}
