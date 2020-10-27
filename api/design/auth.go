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

var _ = Service("auth", func() {
	Description("The auth service exposes endpoint to authenticate User against GitHub OAuth")

	Error("invalid-code", ErrorResult, "Invalid Authorization code")
	Error("invalid-token", ErrorResult, "Invalid User token")
	Error("invalid-scopes", ErrorResult, "Invalid User scope")
	Error("internal-error", ErrorResult, "Internal Server Error")

	Method("Authenticate", func() {
		Description("Authenticates users against GitHub OAuth")
		Payload(func() {
			Attribute("code", String, "OAuth Authorization code of User", func() {
				Example("code", "5628b69ec09c09512eef")
			})
			Required("code")
		})
		Result(func() {
			Attribute("data", AuthTokens, "User Tokens")
			Required("data")
		})

		HTTP(func() {
			POST("/auth/login")
			Param("code")

			Response(StatusOK)
			Response("invalid-code", StatusBadRequest)
			Response("internal-error", StatusInternalServerError)
			Response("invalid-token", StatusUnauthorized)
			Response("invalid-scopes", StatusForbidden)
		})
	})
})
