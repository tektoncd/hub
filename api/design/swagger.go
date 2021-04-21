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

import . "goa.design/goa/v3/dsl"

var _ = Service("swagger", func() {
	Description("The swagger service serves the API swagger definition.")

	HTTP(func() {
		Path("/schema")
	})

	// NOTE: The path is changed to docs/openapi3.json to make it work in container.
	// Copying the gen as it is doesn't seems to work properly, so in dockerfile, swagger will
	// file is copied to /docs. This will make the swagger api not work locally, as the file
	// generated is in gen directory. 
	Files("/swagger.json", "docs/openapi3.json", func() {
		Description("JSON document containing the API swagger definition")
	})
})
