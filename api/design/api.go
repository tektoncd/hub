package design

import (
	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = API("hub", func() {
	Title("Tekton Hub")
	Description("HTTP services for managing Tekton Hub")
	Version("0.1")
	Meta("swagger:example", "false")
	Server("hub", func() {
		Host("production", func() {
			URI("http://api.hub.tekton.dev")
		})

		Services("category", "resource", "swagger")
	})

	// TODO: restrict CORS origin | https://github.com/tektoncd/hub/issues/26
	cors.Origin("*", func() {
		cors.Headers("Content-Type")
		cors.Methods("GET")
	})
})
