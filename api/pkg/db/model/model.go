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
		Version             string
		Description         string
		URL                 string
		DisplayName         string
		MinPipelinesVersion string
		Resource            Resource
		ResourceID          uint
	}

	ResourceTag struct {
		ResourceID uint
		TagID      uint
	}
)
