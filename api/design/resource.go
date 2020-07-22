package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("resource", func() {
	Description("The resource service provides details about all type of resources")

	Error("internal-error", ErrorResult, "Internal Server Error")
	Error("not-found", ErrorResult, "Resource Not Found Error")

	Method("Query", func() {
		Description("Find resources by a combination of name, type")
		Payload(func() {
			Attribute("name", String, "Name of resource", func() {
				Default("")
			})
			Attribute("type", String, "Type of resource", func() {
				Enum("task", "pipeline", "")
				Default("")
			})
			Attribute("limit", UInt, "Maximum number of resources to be returned", func() {
				Default(100)
			})
		})
		Result(CollectionOf(Resource), func() {
			View("default")
		})

		HTTP(func() {
			GET("/query")
			Param("name")
			Param("type")
			Param("limit")

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
			View("default")
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

	Method("ByTypeNameVersion", func() {
		Description("Find resource using name, type and version of resource")
		Payload(func() {
			Attribute("type", String, "type of resource", func() {
				Enum("task", "pipeline")
			})
			Attribute("name", String, "name of resource")
			Attribute("version", String, "version of resource")

			Required("type", "name", "version")
		})
		Result(ResVersion, func() {
			View("default")
		})

		HTTP(func() {
			GET("/resource/{type}/{name}/{version}")

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

	Method("ByTypeName", func() {
		Description("Find resources using name and type")
		Payload(func() {
			Attribute("type", String, "Type of resource", func() {
				Enum("task", "pipeline")
			})
			Attribute("name", String, "Name of resource")
			Required("type", "name")
		})
		Result(CollectionOf(Resource), func() {
			View("default")
		})

		HTTP(func() {
			GET("/resource/{type}/{name}")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
			Response("not-found", StatusNotFound)
		})
	})

})
