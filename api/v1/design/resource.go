// Copyright © 2021 The Tekton Authors.
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

package design

import (
	types "github.com/tektoncd/hub/api/design/types"
	. "goa.design/goa/v3/dsl"
)

var _ = Service("resource", func() {
	Description("The resource service provides details about all kind of resources")

	HTTP(func() {
		Path("/v1")
	})

	Error("internal-error", ErrorResult, "Internal Server Error")
	Error("not-found", ErrorResult, "Resource Not Found Error")
	Error("invalid-kind", ErrorResult, "Invalid Resource Kind")

	// NOTE: Supported Tekton Resource kind by APIs are defined in /pkg/parser/kind.go

	Method("Query", func() {
		Description("Find resources by a combination of name, kind, catalog, categories, platforms and tags")
		Payload(func() {
			Attribute("name", String, "Name of resource", func() {
				Default("")
				Example("name", "buildah")
			})
			Attribute("catalogs", ArrayOf(String), "Catalogs of resource to filter by", func() {
				Example([]string{"tekton", "openshift"})
			})
			Attribute("kinds", ArrayOf(String), "Kinds of resource to filter by", func() {
				Example([]string{"task", "pipelines"})
			})
			Attribute("categories", ArrayOf(String), "Category associated with a resource to filter by", func() {
				Example([]string{"build", "tools"})
			})
			Attribute("tags", ArrayOf(String), "Tags associated with a resource to filter by", func() {
				Example([]string{"image", "build"})
			})
			Attribute("platforms", ArrayOf(String), "Platforms associated with a resource to filter by", func() {
				Example([]string{"linux/s390x", "linux/amd64"})
			})
			Attribute("limit", UInt, "Maximum number of resources to be returned", func() {
				Default(1000)
				Example("limit", 100)
			})

			Attribute("match", String, "Strategy used to find matching resources", func() {
				Enum("exact", "contains")
				Default("contains")
			})
		})
		Result(types.Resources)

		HTTP(func() {
			GET("/query")

			Param("name")
			Param("catalogs")
			Param("categories")
			Param("kinds")
			Param("tags")
			Param("platforms")
			Param("limit")
			Param("match")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
			Response("invalid-kind", StatusBadRequest)
			Response("not-found", StatusNotFound)
		})
	})

	Method("List", func() {
		Description("List all resources sorted by rating and name")
		Payload(func() {
			Attribute("limit", UInt, "Maximum number of resources to be returned", func() {
				Default(1000)
				Example("limit", 100)
			})
		})
		Result(types.Resources)

		HTTP(func() {
			GET("/resources")

			Param("limit")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
		})
	})

	Method("VersionsByID", func() {
		Description("Find all versions of a resource by its id")
		Payload(func() {
			Attribute("id", UInt, "ID of a resource", func() {
				Example("id", 1)
			})
			Required("id")
		})
		Result(types.ResourceVersions)

		HTTP(func() {
			GET("/resource/{id}/versions")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
			Response("not-found", StatusNotFound)
		})
	})

	Method("ByCatalogKindNameVersion", func() {
		Description("Find resource using name of catalog & name, kind and version of resource")
		Payload(func() {
			Attribute("catalog", String, "name of catalog", func() {
				Example("catalog", "tektoncd")
			})
			Attribute("kind", String, "kind of resource", func() {
				Enum("task", "pipeline")
			})
			Attribute("name", String, "name of resource", func() {
				Example("name", "buildah")
			})
			Attribute("version", String, "version of resource", func() {
				Example("version", "0.1")
			})

			Required("catalog", "kind", "name", "version")
		})
		Result(types.ResourceVersion)

		HTTP(func() {
			GET("/resource/{catalog}/{kind}/{name}/{version}")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
			Response("not-found", StatusNotFound)
		})
	})

	Method("ByVersionId", func() {
		Description("Find a resource using its version's id")
		Payload(func() {
			Attribute("versionID", UInt, "Version ID of a resource's version", func() {
				Example("versionID", 1)
			})
			Required("versionID")
		})
		Result(types.ResourceVersion)

		HTTP(func() {
			GET("/resource/version/{versionID}")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
			Response("not-found", StatusNotFound)
		})
	})

	Method("ByCatalogKindName", func() {
		Description("Find resources using name of catalog, resource name and kind of resource")
		Payload(func() {
			Attribute("catalog", String, "name of catalog", func() {
				Example("catalog", "tektoncd")
			})
			Attribute("kind", String, "kind of resource", func() {
				Enum("task", "pipeline")
			})
			Attribute("name", String, "Name of resource", func() {
				Example("name", "buildah")
			})
			Attribute("pipelinesversion", String, "To find resource compatible with a Tekton pipelines version, use this param", func() {
				Pattern(types.PipelinesVersionRegex)
				Example("pipelinesversion", "0.21.0")
			})
			Required("catalog", "kind", "name")
		})
		Result(types.Resource)

		HTTP(func() {
			GET("/resource/{catalog}/{kind}/{name}")

			Param("pipelinesversion")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
			Response("not-found", StatusNotFound)
		})
	})

	Method("ById", func() {
		Description("Find a resource using it's id")
		Payload(func() {
			Attribute("id", UInt, "ID of a resource", func() {
				Example("id", 1)
			})
			Required("id")
		})
		Result(types.Resource)

		HTTP(func() {
			GET("/resource/{id}")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
			Response("not-found", StatusNotFound)
		})
	})

})
