// Copyright © 2021 The Tekton Authors.
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

var _ = Service("user", func() {
	Description("The user service exposes endpoint to get user specific specs")

	Error("invalid-token", ErrorResult, "Invalid User token")
	Error("invalid-scopes", ErrorResult, "Invalid User scope")
	Error("internal-error", ErrorResult, "Internal Server Error")

	Method("RefreshAccessToken", func() {
		Description("Refresh the access token of User")
		Security(types.JWTAuth, func() {
			Scope("refresh:token")
		})
		Payload(func() {
			Token("refreshToken", String, "Refresh Token of User", func() {
				Example("refreshToken", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."+
					"eyJleHAiOjE1Nzc4ODM2MDAsImlhdCI6MTU3Nzg4MDAwMCwiaWQiOjExLCJpc3MiOiJUZWt0b24gSHViIiwic2NvcGVzIjpbInJlZnJlc2g6dG9rZW4iXSwidHlwZSI6InJlZnJlc2gtdG9rZW4ifQ."+
					"4RdUk5ttHdDiymurlZ_f7Uy5Pas3Lq9w04BjKQKRiCE")
			})
			Required("refreshToken")
		})
		Result(func() {
			Attribute("data", types.AccessToken, "User Access JWT")
			Required("data")
		})

		HTTP(func() {
			POST("/user/refresh/accesstoken")
			Header("refreshToken:Authorization")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
			Response("invalid-token", StatusUnauthorized)
			Response("invalid-scopes", StatusForbidden)
		})
	})
})
