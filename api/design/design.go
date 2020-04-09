package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("hub", func() {
	Title("Tekton Hub")
	Description("Service to get resource details")

	Server("hub", func() {
		Host("localhost", func() {
			URI("http://localhost:8000")
		})
	})
})

var _ = Service("api", func() {
	Description("The api service gives resource details")

	//Method to get all resources
	Method("list", func() {
		Description("Get all tasks and pipelines.")
		HTTP(func() {
			GET("/resources")
			Response(StatusOK)
		})
	})
	Files("/openapi.json", "./gen/http/openapi.json")
})
