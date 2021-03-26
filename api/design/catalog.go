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

var _ = Service("catalog", func() {
	Description("The Catalog Service exposes endpoints to interact with catalogs")

	Error("internal-error", ErrorResult, "Internal Server Error")
	Error("not-found", ErrorResult, "Resource Not Found Error")

	Method("Refresh", func() {
		Description("Refresh a Catalog by it's name")
		Security(types.JWTAuth, func() {
			Scope("catalog:refresh")
		})
		Payload(func() {
			Attribute("catalogName", String, "Name of catalog", func() {
				Example("catalogName", "tekton")
			})
			Token("token", String, "JWT")
			Required("token", "catalogName")
		})
		Result(types.Job)

		HTTP(func() {
			POST("/catalog/{catalogName}/refresh")
			Header("token:Authorization")

			Response(StatusOK)
			Response("not-found", StatusNotFound)
			Response("internal-error", StatusInternalServerError)
		})
	})

	Method("RefreshAll", func() {
		Description("Refresh all catalogs")
		Security(types.JWTAuth, func() {
			Scope("catalog:refresh")
		})
		Payload(func() {
			Token("token", String, "JWT")
			Required("token")
		})
		Result(ArrayOf(types.Job))

		HTTP(func() {
			POST("/catalog/refresh")
			Header("token:Authorization")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
		})
	})

})
