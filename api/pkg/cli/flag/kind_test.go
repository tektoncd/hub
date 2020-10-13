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

package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {

	// Valid Tekton Resource Kind
	k := "task"
	err := Kind(k).IsValid()
	assert.NoError(t, err)

	// Invalid Case
	k = "abc"
	err = Kind(k).IsValid()
	assert.Error(t, err, "invalid value \"abc\" set for option kinds. Valid options: [task, pipeline]")
}
