package model

import (
	"github.com/jinzhu/gorm"
)

type (
	Category struct {
		gorm.Model
		Name string `gorm:"size:100;not null;unique"`
		Tags []Tag
	}

	Tag struct {
		gorm.Model
		Name       string `gorm:"size:100;not null;unique"`
		Category   Category
		CategoryID int
		Resources  []*Resource `gorm:"many2many:resource_tags;"`
	}

	Catalog struct {
		gorm.Model
		Name       string
		Type       string
		URL        string
		Owner      string
		ContextDir string
		Resources  []Resource
		Revision   string
	}

	Resource struct {
		gorm.Model
		Name      string
		Type      string
		Rating    float64
		Catalog   Catalog
		CatalogID uint
		Versions  []ResourceVersion
		Tags      []*Tag `gorm:"many2many:resource_tags;"`
	}

	ResourceVersion struct {
		gorm.Model
		Version     string
		Description string
		URL         string
		DisplayName string
		Resource    Resource
		ResourceID  uint
	}

	ResourceTag struct {
		ResourceID uint
		TagID      uint
	}
)
