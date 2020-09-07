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

var _ = Service("catalog", func() {
	Description("The Catalog Service exposes endpoints to refresh catalog")

	Error("internal-error", ErrorResult, "Internal Server Error")

	Method("Refresh", func() {
		Description("Refresh the catalog for new resources")

		HTTP(func() {
			POST("/catalog/refresh")

			Response(StatusOK)
			Response("internal-error", StatusInternalServerError)
		})
	})

})
