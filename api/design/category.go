package design
import (
	. "goa.design/goa/v3/dsl"
)
var _ = Service("category", func() {
	Description("The category service gives category details")

	Error("internal-error", ErrorResult, "Internal Server Error")

		//Method to get all categories with their tags
		Method("All", func() {
			Description("Get all Categories with their tags sorted by name")
			Result(ArrayOf(Category))
			HTTP(func() {
				GET("/categories")
				Response(StatusOK)
				Response("internal-error", StatusInternalServerError)
			})
		})
	})