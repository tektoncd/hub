package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("status", func() {
	Description("Describes the status of the server")

	Method("Status", func() {
		Description("Return status 'ok' when the server has started successfully")
		Result(func() {
			Attribute("status", String, "Status of server", func() {
				Example("status", "ok")
			})
			Required("status")
		})

		HTTP(func() {
			GET("/")

			Response(StatusOK)
		})
	})
})
