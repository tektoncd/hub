// Copyright © 2022 The Tekton Authors.
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
	"github.com/tektoncd/hub/api/design/types"
	. "goa.design/goa/v3/dsl"
)

// NOTE: APIs in the service are moved to v1. This APIs will be deprecated in the next release.

var _ = Service("resource", func() {
	Description("The resource service provides details about all kind of resources")

	Error("internal-error", ErrorResult, "Internal Server Error")
	Error("not-found", ErrorResult, "Resource Not Found Error")

	Method("List", func() {
		Description("List all resources sorted by rating and name")
		Result(types.Resources)

		HTTP(func() {
			GET("/resources")
			Redirect("/v1/resources", StatusMovedPermanently)
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
		Result(func() {
			Attribute("location", String, "Redirect URL", func() {
				Format(FormatURI)
			})
			Required("location")
		})

		HTTP(func() {
			GET("/resource/{id}/versions")
			Response(StatusFound, func() {
				Header("location")
			})

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
		Result(func() {
			Attribute("location", String, "Redirect URL", func() {
				Format(FormatURI)
			})
			Required("location")
		})

		HTTP(func() {
			GET("/resource/{catalog}/{kind}/{name}/{version}")
			Response(StatusFound, func() {
				Header("location")
			})

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
		Result(func() {
			Attribute("location", String, "Redirect URL", func() {
				Format(FormatURI)
			})
			Required("location")
		})

		HTTP(func() {
			GET("/resource/version/{versionID}")
			Response(StatusFound, func() {
				Header("location")
			})
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
		Result(func() {
			Attribute("location", String, "Redirect URL", func() {
				Format(FormatURI)
			})
			Required("location")
		})

		HTTP(func() {
			GET("/resource/{catalog}/{kind}/{name}")

			Param("pipelinesversion")

			Response(StatusFound, func() {
				Header("location")
			})

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
		Result(func() {
			Attribute("location", String, "Redirect URL", func() {
				Format(FormatURI)
			})
			Required("location")
		})

		HTTP(func() {
			GET("/resource/{id}")
			Response(StatusFound, func() {
				Header("location")
			})
		})
	})

})
