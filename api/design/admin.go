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
	"github.com/tektoncd/hub/api/design/types"
	. "goa.design/goa/v3/dsl"
)

var _ = Service("admin", func() {
	Description("Admin service")

	Error("invalid-payload", ErrorResult, "Invalid request body")
	Error("invalid-token", ErrorResult, "Invalid User token")
	Error("invalid-scopes", ErrorResult, "Invalid Token scopes ")
	Error("internal-error", ErrorResult, "Internal server error")

	Method("UpdateAgent", func() {
		Description("Create or Update an agent user with required scopes")
		Security(types.JWTAuth, func() {
			Scope("agent:create")
		})
		Payload(func() {
			Token("token", String, "User JWT")
			Attribute("name", String, "Name of Agent")
			Attribute("scopes", ArrayOf(String), "Scopes required for Agent")
			Required("name", "scopes", "token")
		})
		Result(func() {
			Attribute("token", String, "Agent JWT")
			Required("token")
		})

		HTTP(func() {
			PUT("/system/user/agent")
			Header("token:Authorization")

			Response(StatusOK)
			Response("invalid-payload", StatusBadRequest)
			Response("invalid-token", StatusUnauthorized)
			Response("invalid-scopes", StatusForbidden)
			Response("internal-error", StatusInternalServerError)
		})
	})

	Method("RefreshConfig", func() {
		Description("Refresh the changes in config file")
		Security(types.JWTAuth, func() {
			Scope("config:refresh")
		})
		Payload(func() {
			Token("token", String, "User JWT")
			Required("token")
		})
		Result(func() {
			Attribute("checksum", String, "Config file checksum")
			Required("checksum")
		})

		HTTP(func() {
			POST("/system/config/refresh")
			Header("token:Authorization")

			Response(StatusOK)
			Response("invalid-token", StatusUnauthorized)
			Response("invalid-scopes", StatusForbidden)
			Response("internal-error", StatusInternalServerError)
		})
	})
})
