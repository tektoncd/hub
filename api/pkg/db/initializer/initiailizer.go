package initializer

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
)

// Initializer defines the configuration required for initailizer
// to populate the tables
type Initializer struct {
	log  *zap.SugaredLogger
	db   *gorm.DB
	data *app.Data
}

// New returns the Initializer implementation.
func New(c app.Config) *Initializer {
	return &Initializer{
		log:  c.Logger().With("component", "initiailizer"),
		db:   c.DB(),
		data: c.Data(),
	}
}

// Run executes the func which populate the tables
func (i *Initializer) Run() error {

	if err := i.addCategories(); err != nil {
		return err
	}
	if err := i.addCatalogs(); err != nil {
		return err
	}
	return nil
}

func (i *Initializer) addCategories() error {

	db := i.db

	// Checks if tables exists
	if !db.HasTable(&model.Category{}) || !db.HasTable(model.Tag{}) {
		return fmt.Errorf("categories or tags table not found")
	}

	for _, c := range i.data.Categories {
		cat := &model.Category{Name: c.Name}
		if err := i.db.Where(cat).FirstOrCreate(cat).
			Error; err != nil {
			i.log.Error(err)
			return err
		}
		for _, t := range c.Tags {
			tag := &model.Tag{Name: t, CategoryID: cat.ID}
			if err := db.Where(tag).FirstOrCreate(tag).
				Error; err != nil {
				i.log.Error(err)
				return err
			}
		}
	}
	return nil
}

func (i *Initializer) addCatalogs() error {

	db := i.db

	// Checks if tables exists
	if !db.HasTable(&model.Catalog{}) {
		return fmt.Errorf("catalogs table not found")
	}

	for _, c := range i.data.Catalogs {
		cat := &model.Catalog{
			Name:       c.Name,
			Type:       c.Type,
			URL:        c.URL,
			Owner:      c.Owner,
			Revision:   c.Revision,
			ContextDir: c.ContextDir,
		}
		if err := db.Where(&model.Catalog{Name: c.Name, Type: c.Type}).
			FirstOrCreate(cat).
			Error; err != nil {
			i.log.Error(err)
			return err
		}
	}
	return nil
}
