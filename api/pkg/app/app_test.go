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

package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeDuration(t *testing.T) {
	duration := "5d"
	dur, err := computeDuration(duration)
	assert.NoError(t, err)
	assert.Equal(t, dur.String(), "120h0m0s")
}

func TestComputeDurationError(t *testing.T) {
	duration := "5M"
	_, err := computeDuration(duration)
	assert.Equal(t, err.Error(), "JWT doesn't support the duration specified 5M. \nSupported formats are w(weeks), d(days), h(hours), m(min), s(sec)")
}
