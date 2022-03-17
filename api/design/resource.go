// Copyright © 2022 The Tekton Authors.
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

// NOTE: APIs in the service are moved to v1. This APIs will be deprecated in the next release.

var _ = Service("resource", func() {
	Description("The resource service provides details about all kind of resources")

	Method("List", func() {
		Description("List all resources sorted by rating and name")
		Result(types.Resources)

		HTTP(func() {
			GET("/resources")
			Redirect("/v1/resources", StatusMovedPermanently)
		})
	})

})
