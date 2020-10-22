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

package formatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/gen/http/resource/client"
)

func TestFormatName(t *testing.T) {
	name := FormatName("abc", "0.1")
	assert.Equal(t, name, "abc (0.1)")
}

func TestFormatDesc(t *testing.T) {

	// Description greater than 40 char
	desc := FormatDesc("Buildah task builds source into a container image and then pushes it to a container registry.")
	assert.Equal(t, "Buildah task builds source into a conta...", desc)

	// Description less than 40 char
	desc = FormatDesc("Buildah task builds images.")
	assert.Equal(t, "Buildah task builds images.", desc)

	// No Description
	desc = FormatDesc("")
	assert.Equal(t, "---", desc)
}

func TestFormatTags(t *testing.T) {

	tagName1 := "tag1"
	tagName2 := "tag2"

	res := []*client.TagResponseBody{
		&client.TagResponseBody{
			Name: &tagName1,
		},
		&client.TagResponseBody{
			Name: &tagName2,
		},
	}

	tags := FormatTags(res)
	assert.Equal(t, "tag1, tag2", tags)

	// No Tags
	tags = FormatTags(nil)
	assert.Equal(t, "---", tags)
}
