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
	cors "goa.design/plugins/v3/cors/dsl"

	// Enables the zaplogger plugin
	_ "goa.design/plugins/v3/zaplogger"
)

var _ = API("hub", func() {
	Title("Tekton Hub")
	Description("HTTP services for managing Tekton Hub")
	Version("1.0")
	Meta("swagger:example", "false")
	Server("hub", func() {
		Host("production", func() {
			URI("https://api.hub.tekton.dev")
		})

		Services(
			"admin",
			"catalog",
			"category",
			"rating",
			"status",
			"swagger",
		)
	})

	// TODO: restrict CORS origin | https://github.com/tektoncd/hub/issues/26
	cors.Origin("*", func() {
		cors.Headers("Content-Type", "Authorization")
		cors.Methods("GET", "POST", "PUT")
	})
})
