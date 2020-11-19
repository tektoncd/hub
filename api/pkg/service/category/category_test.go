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

package category

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

func TestCategory_List(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	category := New(tc)
	all, err := category.List(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 3, len(all.Data))
	assert.Equal(t, 2, len(all.Data[0].Tags))
	assert.Equal(t, "abc", all.Data[0].Name)          // categories are sorted by name
	assert.Equal(t, "atag", all.Data[0].Tags[0].Name) // tags are sorted by name
}
