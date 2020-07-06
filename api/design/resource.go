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
		Result(CollectionOf(Resource))

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
})
