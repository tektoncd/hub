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

	"gorm.io/gorm"
)

type (
	Category struct {
		gorm.Model
		Name      string      `gorm:"not null;unique"`
		Resources []*Resource `gorm:"many2many:resource_categories;"`
	}

	Tag struct {
		gorm.Model
		Name      string      `gorm:"not null;unique"`
		Resources []*Resource `gorm:"many2many:resource_tags;"`
	}

	Platform struct {
		gorm.Model
		Name             string             `gorm:"not null;unique"`
		ResourceVersions []*ResourceVersion `gorm:"many2many:version_platforms;"`
		Resource         []*Resource        `gorm:"many2many:resource_platforms;"`
	}

	Catalog struct {
		gorm.Model
		Name       string `gorm:"uniqueIndex:uix_name_org"`
		Org        string `gorm:"uniqueIndex:uix_name_org"`
		Type       string `gorm:"not null;default:null"`
		URL        string `gorm:"not null;default:null"`
		Revision   string `gorm:"not null;default:null"`
		ContextDir string
		SHA        string
		Resources  []Resource
		Errors     []CatalogError
	}

	CatalogError struct {
		gorm.Model
		Catalog   Catalog
		CatalogID uint
		Type      string
		Detail    string
	}

	Resource struct {
		gorm.Model
		Name       string `gorm:"not null;default:null"`
		Kind       string `gorm:"not null;default:null"`
		Rating     float64
		Catalog    Catalog
		Categories []*Category `gorm:"many2many:resource_categories;"`
		CatalogID  uint
		Versions   []ResourceVersion
		Tags       []*Tag      `gorm:"many2many:resource_tags;"`
		Platforms  []*Platform `gorm:"many2many:resource_platforms;constraint:OnDelete:CASCADE;"`
	}

	ResourceVersion struct {
		gorm.Model
		Version             string `gorm:"not null;default:null"`
		Description         string
		URL                 string `gorm:"not null;default:null"`
		DisplayName         string
		Deprecated          bool   `gorm:"default:false"`
		MinPipelinesVersion string `gorm:"not null;default:null"`
		Resource            Resource
		ResourceID          uint
		Platforms           []*Platform `gorm:"many2many:version_platforms;constraint:OnDelete:CASCADE;"`
		ModifiedAt          time.Time
	}

	ResourceTag struct {
		ResourceID uint
		TagID      uint
	}

	VersionPlatform struct {
		ResourceVersionID uint
		PlatformID        uint
	}

	ResourcePlatform struct {
		ResourceID uint
		PlatformID uint
	}

	ResourceCategory struct {
		ResourceID uint
		CategoryID uint
	}

	User struct {
		gorm.Model
		AgentName            string
		GithubLogin          string
		GithubName           string
		Type                 UserType
		Scopes               []*Scope `gorm:"many2many:user_scopes;"`
		RefreshTokenChecksum string
		AvatarURL            string
	}

	Scope struct {
		gorm.Model
		Name string `gorm:"not null;unique"`
	}

	UserResourceRating struct {
		gorm.Model
		UserID     uint
		User       User
		Resource   Resource
		ResourceID uint
		Rating     uint `gorm:"not null;default:null"`
	}

	UserScope struct {
		UserID  uint
		ScopeID uint
	}

	Config struct {
		gorm.Model
		Checksum string
	}
)

type UserType string

// Types of Users
const (
	NormalUserType UserType = "user"
	AgentUserType  UserType = "agent"
)
