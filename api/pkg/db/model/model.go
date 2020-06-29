package model

import (
	"github.com/jinzhu/gorm"
)

type (

	//Category model represents categories which associated with group of tags
	Category struct {
		gorm.Model
		Name string `gorm:"size:100;not null;unique"`
		Tags []Tag
	}

	// Tag model represents tags associated with a resource
	Tag struct {
		gorm.Model
		Name       string `gorm:"size:100;not null;unique"`
		Category   Category
		CategoryID int
		Resources  []*Resource `gorm:"many2many:resource_tags;"`
	}

	//Catalog model represents origin repo to which resource belongs to
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

	// Resource represents model which describe resources from a catalog
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

	// ResourceVersion represents diffrent versions of a Resource
	ResourceVersion struct {
		gorm.Model
		Version     string
		Description string
		URL         string
		DisplayName string
		Resource    Resource
		ResourceID  uint
	}
	// ResourceTag represent struct for resource_tag, needed for creating foreign key
	ResourceTag struct {
		ResourceID uint
		TagID      uint
	}
)
