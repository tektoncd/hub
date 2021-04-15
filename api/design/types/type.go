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

package types

import (
	. "goa.design/goa/v3/dsl"
)

var Tag = Type("Tag", func() {
	Attribute("id", UInt, "ID is the unique id of tag", func() {
		Example("id", 1)
	})
	Attribute("name", String, "Name of tag", func() {
		Example("name", "image-build")
	})

	Required("id", "name")
})

var Tags = ArrayOf(Tag)

var Category = Type("Category", func() {
	Attribute("id", UInt, "ID is the unique id of the category", func() {
		Example("id", 1)
	})
	Attribute("name", String, "Name of category", func() {
		Example("name", "Image Builder")
	})
	Attribute("tags", Tags, "List of tags associated with the category", func() {
		Example("tags", func() {
			Value([]Val{
				{"id": 1, "name": "image-build"},
				{"id": 2, "name": "kaniko"},
			})
		})
	})

	Required("id", "name", "tags")
})

var Catalog = ResultType("application/vnd.hub.catalog", "Catalog", func() {
	Attribute("id", UInt, "ID is the unique id of the catalog", func() {
		Example("id", 1)
	})
	Attribute("name", String, "Name of catalog", func() {
		Example("name", "Tekton")
	})
	Attribute("type", String, "Type of catalog", func() {
		Enum("official", "community")
		Example("type", "community")
	})
	Attribute("url", String, "URL of catalog", func() {
		Example("url", "https://github.com/tektoncd/hub")
	})

	View("min", func() {
		Attribute("id")
		Attribute("name")
		Attribute("type")
	})

	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("type")
		Attribute("url")
	})

	Required("id", "name", "type", "url")
})

var ResourceVersionData = ResultType("application/vnd.hub.resource.version.data", "ResourceVersionData", func() {
	Description("The Version result type describes resource's version information.")

	Attribute("id", UInt, "ID is the unique id of resource's version", func() {
		Example("id", 1)
	})
	Attribute("version", String, "Version of resource", func() {
		Example("version", "0.1")
	})
	Attribute("displayName", String, "Display name of version", func() {
		Example("displayName", "Buildah")
	})
	Attribute("description", String, "Description of version", func() {
		Example("descripiton", "Buildah task builds source into a container image and then pushes it to a container registry.")
	})
	Attribute("minPipelinesVersion", String, "Minimum pipelines version the resource's version is compatible with", func() {
		Example("minPipelinesVersion", "0.12.1")
	})
	Attribute("rawURL", String, "Raw URL of resource's yaml file of the version", func() {
		Format(FormatURI)
		Example("rawURL", "https://raw.githubusercontent.com/tektoncd/catalog/main/task/buildah/0.1/buildah.yaml")
	})
	Attribute("webURL", String, "Web URL of resource's yaml file of the version", func() {
		Format(FormatURI)
		Example("webURL", "https://github.com/tektoncd/catalog/blob/main/task/buildah/0.1/buildah.yaml")
	})
	Attribute("updatedAt", String, "Timestamp when version was last updated", func() {
		Format(FormatDateTime)
		Example("updatedAt", "2020-01-01 12:00:00 +0000 UTC")
	})
	Attribute("resource", ResourceData, "Resource to which the version belongs", func() {
		View("info")
		Example("resource", func() {
			Value(Val{
				"id":      1,
				"name":    "buildah",
				"catalog": Val{"id": 1, "type": "community"},
				"kind":    "task",
				"tags":    []Val{{"id": 1, "name": "image-build"}},
				"rating":  4.3,
			})
		})
	})

	View("tiny", func() {
		Attribute("id")
		Attribute("version")
	})

	View("min", func() {
		Attribute("id")
		Attribute("version")
		Attribute("rawURL")
		Attribute("webURL")
	})

	View("withoutResource", func() {
		Attribute("id")
		Attribute("version")
		Attribute("displayName")
		Attribute("description")
		Attribute("minPipelinesVersion")
		Attribute("rawURL")
		Attribute("webURL")
		Attribute("updatedAt")
	})

	View("default", func() {
		Attribute("id")
		Attribute("version")
		Attribute("displayName")
		Attribute("description")
		Attribute("minPipelinesVersion")
		Attribute("rawURL")
		Attribute("webURL")
		Attribute("updatedAt")
		Attribute("resource")
	})

	Required("id", "version", "displayName", "description", "minPipelinesVersion", "rawURL", "webURL", "updatedAt", "resource")
})

var ResourceData = ResultType("application/vnd.hub.resource.data", "ResourceData", func() {
	Description("The resource type describes resource information.")

	Attribute("id", UInt, "ID is the unique id of the resource", func() {
		Example("id", 1)
	})
	Attribute("name", String, "Name of resource", func() {
		Example("name", "buildah")
	})
	Attribute("catalog", Catalog, "Type of catalog to which resource belongs", func() {
		View("min")
		Example("catalog", func() {
			Value(Val{"id": 1, "name": "tekton", "type": "community"})
		})
	})
	Attribute("kind", String, "Kind of resource", func() {
		Example("kind", "task")
	})
	Attribute("latestVersion", "ResourceVersionData", "Latest version of resource", func() {
		View("withoutResource")
		Example("latestVersion", func() {
			Value(Val{
				"id":                  1,
				"version":             "0.1",
				"description":         "Buildah task builds source into a container image and then pushes it to a container registry.",
				"displayName":         "Buildah",
				"minPipelinesVersion": "0.12.1",
				"rawURL":              "https://raw.githubusercontent.com/tektoncd/catalog/main/task/buildah/0.1/buildah.yaml",
				"webURL":              "https://github.com/tektoncd/catalog/blob/main/task/buildah/0.1/buildah.yaml",
				"updatedAt":           "2020-01-01 12:00:00 +0000 UTC",
			})
		})
	})
	Attribute("tags", Tags, "Tags related to resource", func() {
		Example("tags", func() {
			Value([]Val{
				{"id": 1, "name": "image-build"},
			})
		})
	})
	Attribute("rating", Float64, "Rating of resource", func() {
		Example("rating", 4.3)
	})
	Attribute("versions", ArrayOf("ResourceVersionData"), "List of all versions of a resource", func() {
		Example("versions", func() {
			Value([]Val{{
				"id":      1,
				"version": "0.1",
			}, {
				"id":      2,
				"version": "0.2",
			}})
		})
	})

	View("info", func() {
		Attribute("id")
		Attribute("name")
		Attribute("catalog")
		Attribute("kind")
		Attribute("tags")
		Attribute("rating")
	})

	View("withoutVersion", func() {
		Attribute("id")
		Attribute("name")
		Attribute("catalog", func() {
			View("min")
		})
		Attribute("kind")
		Attribute("latestVersion")
		Attribute("tags")
		Attribute("rating")
	})

	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("catalog")
		Attribute("kind")
		Attribute("latestVersion")
		Attribute("tags")
		Attribute("rating")
		Attribute("versions", func() {
			View("tiny")
		})
	})

	Required("id", "name", "catalog", "kind", "latestVersion", "tags", "rating", "versions")
})

var Versions = ResultType("application/vnd.hub.versions", "Versions", func() {
	Description("The Versions type describes response for versions by resource id API.")

	Attribute("latest", ResourceVersionData, "Latest Version of resource", func() {
		Example("latest", func() {
			Value(Val{
				"id":      2,
				"version": "0.2",
				"rawURL":  "https://raw.githubusercontent.com/tektoncd/catalog/main/task/buildah/0.2/buildah.yaml",
				"webURL":  "https://github.com/tektoncd/catalog/blob/main/task/buildah/0.2/buildah.yaml",
			})
		})
	})
	Attribute("versions", ArrayOf(ResourceVersionData), "List of all versions of resource", func() {
		Example("versions", func() {
			Value([]Val{{
				"id":      1,
				"version": "0.1",
				"rawURL":  "https://raw.githubusercontent.com/tektoncd/catalog/main/task/buildah/0.1/buildah.yaml",
				"webURL":  "https://github.com/tektoncd/catalog/blob/main/task/buildah/0.1/buildah.yaml",
			}, {
				"id":      2,
				"version": "0.2",
				"rawURL":  "https://raw.githubusercontent.com/tektoncd/catalog/main/task/buildah/0.2/buildah.yaml",
				"webURL":  "https://github.com/tektoncd/catalog/blob/main/task/buildah/0.2/buildah.yaml",
			}})
		})
	})

	View("default", func() {
		Attribute("latest", func() {
			View("min")
		})
		Attribute("versions", func() {
			View("min")
		})
	})

	Required("latest", "versions")
})

var JWTAuth = JWTSecurity("jwt", func() {
	Description("Secures endpoint by requiring a valid JWT retrieved via the /auth/login endpoint.")
	Scope("rating:read", "Read-only access to rating")
	Scope("rating:write", "Read and write access to rating")
	Scope("agent:create", "Access to create or update an agent")
	Scope("catalog:refresh", "Access to refresh catalog")
	Scope("config:refresh", "Access to refresh config file")
	Scope("refresh:token", "Access to refresh user access token")
})

var HubService = Type("HubService", func() {
	Description("Describes the services and their status")
	Attribute("name", String, "Name of the service", func() {
		Example("name", "api")
	})
	Attribute("status", String, "Status of the service", func() {
		Enum("ok", "error")
		Example("status", "ok")
	})
	Attribute("error", String, "Details of the error if any", func() {
		Example("error", "unable to reach db")
	})

	Required("name", "status")
})

var Job = ResultType("application/vnd.hub.job", "Job", func() {
	Description("The Job type describes a catalog refresh job that is run asynchronously")
	Attribute("id", UInt, "id of the job")
	Attribute("catalogName", String, "Name of the catalog")
	Attribute("status", String, "status of the job")
	Required("id", "catalogName", "status")
})

var Resources = ResultType("application/vnd.hub.resources", "Resources", func() {
	Attribute("data", CollectionOf(ResourceData), func() {
		View("withoutVersion")
	})
	Required("data")
})

var ResourceVersions = ResultType("application/vnd.hub.resource.versions", "ResourceVersions", func() {
	Attribute("data", Versions)
	Required("data")
})

var ResourceVersion = ResultType("application/vnd.hub.resource.version", "ResourceVersion", func() {
	Attribute("data", ResourceVersionData, func() {
		View("default")
	})
	Required("data")
})

var Resource = ResultType("application/vnd.hub.resource", "Resource", func() {
	Attribute("data", ResourceData, func() {
		View("default")
	})
	Required("data")
})

var Token = Type("Token", func() {
	Description("Token includes the JWT, Expire Duration & Time")
	Attribute("token", String, "JWT", func() {
		Example("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."+
			"eyJleHAiOjE1Nzc4ODAzMDAsImlhdCI6MTU3Nzg4MDAwMCwiaWQiOjExLCJpc3MiOiJUZWt0b24gSHViIiwic2NvcGVzIjpbInJhdGluZzpyZWFkIiwicmF0aW5nOndyaXRlIiwiYWdlbnQ6Y3JlYXRlIl0sInR5cGUiOiJhY2Nlc3MtdG9rZW4ifQ."+
			"6pDmziSKkoSqI1f0rc4-AqVdcfY0Q8wA-tSLzdTCLgM")
	})
	Attribute("refreshInterval", String, "Duration the token will Expire In", func() {
		Example("refreshInterval", "1h30m")
	})
	Attribute("expiresAt", Int64, "Time the token will expires at", func() {
		Example("expiresAt", 0)
	})

	Required("token", "refreshInterval", "expiresAt")
})

var AuthTokens = Type("AuthTokens", func() {
	Description("Auth tokens have access and refresh token for user")
	Attribute("access", Token, "Access Token")
	Attribute("refresh", Token, "Refresh Token")
})

var AccessToken = Type("AccessToken", func() {
	Description("Access Token for User")
	Attribute("access", Token, "Access Token for user")
})

var RefreshToken = Type("RefreshToken", func() {
	Description("Refresh Token for User")
	Attribute("refresh", Token, "Refresh Token for user")
})

var UserData = Type("UserData", func() {
	Description("Github user Information")
	Attribute("githubId", String, "Github id of User", func() {
		Example("githubId", "abc123")
	})
	Attribute("name", String, "Github user name", func() {
		Example("name", "abc")
	})
	Attribute("avatarUrl", String, "Github user's profile picture url", func() {
		Example("avatarUrl", "https://avatars.githubusercontent.com")
	})
	Required("githubId", "name", "avatarUrl")
})
