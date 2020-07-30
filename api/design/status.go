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
