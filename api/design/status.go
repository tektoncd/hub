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

var _ = Service("status", func() {
	Description("Describes the status of each service")

	Method("Status", func() {
		Description("Return status of the services")
		Result(func() {
			Attribute("services", ArrayOf(types.HubService), "List of services and their status")
		})

		HTTP(func() {
			GET("/")
			GET("/v1")
			Response(StatusOK)
		})
	})
})
