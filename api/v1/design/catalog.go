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
	types "github.com/tektoncd/hub/api/design/types"
	. "goa.design/goa/v3/dsl"
)

var _ = Service("catalog", func() {
	Description("List of catalogs")

	HTTP(func() {
		Path("/v1")
	})

	Error("internal-error", ErrorResult, "Internal Server Error")

	Method("List", func() {
		Description("List all Catalogs")

		Result(func() {
			Attribute("data", ArrayOf(types.Catalog))
		})

		HTTP(func() {
			GET("/catalogs")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
		})
	})

})
