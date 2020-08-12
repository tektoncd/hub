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

var _ = Service("rating", func() {
	Description("The rating service exposes endpoints to read and write user's rating for resources")

	Error("not-found", ErrorResult, "Resource Not Found Error")
	Error("internal-error", ErrorResult, "Internal server error")
	Error("invalid-token", ErrorResult, "Invalid User token")
	Error("invalid-scopes", ErrorResult, "Invalid User scope")

	Method("Get", func() {
		Description("Find user's rating for a resource")
		Security(JWTAuth, func() {
			Scope("rating:read")
		})
		Payload(func() {
			Attribute("id", UInt, "ID of a resource")
			Token("token", String, "JWT")
			Required("id", "token")
		})
		Result(func() {
			Attribute("rating", Int, "User rating for resource")
			Required("rating")
		})

		HTTP(func() {
			GET("/resource/{id}/rating")
			Header("token:Authorization")

			Response(StatusOK)
			Response("not-found", StatusNotFound)
			Response("internal-error", StatusInternalServerError)
			Response("invalid-token", StatusUnauthorized)
			Response("invalid-scopes", StatusForbidden)
		})
	})

	Method("Update", func() {
		Description("Update user's rating for a resource")
		Security(JWTAuth, func() {
			Scope("rating:write")
		})
		Payload(func() {
			Attribute("id", UInt, "ID of a resource")
			Attribute("rating", UInt, "User rating for resource", func() {
				Minimum(0)
				Maximum(5)
			})
			Token("token", String, "JWT")
			Required("id", "token", "rating")
		})

		HTTP(func() {
			PUT("/resource/{id}/rating")
			Header("token:Authorization")

			Response(StatusOK)
			Response("not-found", StatusNotFound)
			Response("internal-error", StatusInternalServerError)
			Response("invalid-token", StatusUnauthorized)
			Response("invalid-scopes", StatusForbidden)
		})
	})
})
