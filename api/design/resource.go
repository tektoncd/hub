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

package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("resource", func() {
	Description("The resource service provides details about all kind of resources")

	Error("internal-error", ErrorResult, "Internal Server Error")
	Error("not-found", ErrorResult, "Resource Not Found Error")

	Method("Query", func() {
		Description("Find resources by a combination of name, kind and tags")
		Payload(func() {
			Attribute("name", String, "Name of resource", func() {
				Default("")
			})
			Attribute("kinds", ArrayOf(String), "Kinds of resource to filter by")
			Attribute("tags", ArrayOf(String), "Tags associated with a resource to filter by")
			Attribute("limit", UInt, "Maximum number of resources to be returned", func() {
				Default(100)
			})

			Attribute("match", String, "Strategy used to find matching resources", func() {
				Enum("exact", "contains")
				Default("contains")
			})
		})
		Result(CollectionOf(Resource), func() {
			View("withoutVersion")
		})

		HTTP(func() {
			GET("/query")
			Param("name")
			Param("kinds")
			Param("tags")
			Param("limit")
			Param("match")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
			Response("not-found", StatusNotFound)
		})
	})

	Method("List", func() {
		Description("List all resources sorted by rating and name")
		Payload(func() {
			Attribute("limit", UInt, "Maximum number of resources to be returned", func() {
				Default(100)
			})
		})
		Result(CollectionOf(Resource), func() {
			View("withoutVersion")
		})

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
			Attribute("id", UInt, "ID of a resource")
			Required("id")
		})
		Result(Versions)

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
			Attribute("catalog", String, "name of catalog")
			Attribute("kind", String, "kind of resource", func() {
				Enum("task", "pipeline")
			})
			Attribute("name", String, "name of resource")
			Attribute("version", String, "version of resource")

			Required("catalog", "kind", "name", "version")
		})
		Result(ResVersion, func() {
			View("default")
		})

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
			Attribute("versionID", UInt, "Version ID of a resource's version")
			Required("versionID")
		})
		Result(ResVersion, func() {
			View("default")
		})

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
			Attribute("catalog", String, "name of catalog")
			Attribute("kind", String, "kind of resource", func() {
				Enum("task", "pipeline")
			})
			Attribute("name", String, "Name of resource")
			Required("catalog", "kind", "name")
		})
		Result(CollectionOf(Resource), func() {
			View("withoutVersion")
		})

		HTTP(func() {
			GET("/resource/{catalog}/{kind}/{name}")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
			Response("not-found", StatusNotFound)
		})
	})

	Method("ById", func() {
		Description("Find a resource using it's id")
		Payload(func() {
			Attribute("id", UInt, "ID of a resource")
			Required("id")
		})
		Result(Resource, func() {
			View("default")
		})

		HTTP(func() {
			GET("/resource/{id}")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
			Response("not-found", StatusNotFound)
		})
	})

})
