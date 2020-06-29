package design
import (
	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl"
)
var _ = API("hub", func() {
	Title("Tekton Hub")
	Description("HTTP services for managing Tekton Hub")
	Version("1.0")
	Meta("swagger:example", "false")
	Server("hub", func() {
		Services("category", "swagger")
		Host("localhost", func() {
			URI("http://localhost:8000")
		})
	})
	cors.Origin("*", func() {
		cors.Headers("Content-Type")
		cors.Methods("GET")
	})
})