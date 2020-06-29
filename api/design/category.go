package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("category", func() {
	Description("The category service provides details about category")

	HTTP(func() {
		Path("/categories")
	})

	Error("internal-error", ErrorResult, "Internal Server Error")

	Method("list", func() {
		Description("List all categories along with their tags sorted by name")
		Result(ArrayOf(Category))

		HTTP(func() {
			GET("/")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
		})
	})
})
