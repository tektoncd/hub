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
	"time"

	"github.com/jinzhu/gorm"
)

type (
	Category struct {
		gorm.Model
		Name string `gorm:"not null;unique"`
		Tags []Tag
	}

	Tag struct {
		gorm.Model
		Name       string `gorm:"not null;unique"`
		Category   Category
		CategoryID uint
		Resources  []*Resource `gorm:"many2many:resource_tags;"`
	}

	Catalog struct {
		gorm.Model
		Name       string `gorm:"unique_index:uix_name_org"`
		Org        string `gorm:"unique_index:uix_name_org"`
		Type       string `gorm:"not null;default:null"`
		URL        string `gorm:"not null;default:null"`
		Revision   string `gorm:"not null;default:null"`
		ContextDir string
		SHA        string
		Resources  []Resource
	}

	Resource struct {
		gorm.Model
		Name      string `gorm:"not null;default:null"`
		Kind      string `gorm:"not null;default:null"`
		Rating    float64
		Catalog   Catalog
		CatalogID uint
		Versions  []ResourceVersion
		Tags      []*Tag `gorm:"many2many:resource_tags;"`
	}

	ResourceVersion struct {
		gorm.Model
		Version             string `gorm:"not null;default:null"`
		Description         string
		URL                 string `gorm:"not null;default:null"`
		DisplayName         string
		MinPipelinesVersion string `gorm:"not null;default:null"`
		Resource            Resource
		ResourceID          uint
		ModifiedAt          time.Time
	}

	ResourceTag struct {
		ResourceID uint
		TagID      uint
	}

	User struct {
		gorm.Model
		GithubLogin string
		GithubName  string
	}

	UserResourceRating struct {
		gorm.Model
		UserID     uint
		User       User
		Resource   Resource
		ResourceID uint
		Rating     uint `gorm:"not null;default:null"`
	}
)
