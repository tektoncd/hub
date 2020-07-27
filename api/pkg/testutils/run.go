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

package testutils

import (
	"sync"
	"testing"
)

var setup sync.Once

// Setup when called first time will connect to test db, run migration to create tables
// and return TestConfig configurations. After that each time it is called it just retuns
// TestConfig configurations
func Setup(t *testing.T) *TestConfig {
	tc := Config()
	setup.Do(func() {
		if !tc.IsValid() {
			t.Fatalf("Failed to create test configuration object: %s", tc.Error())
		}

		if err := applyMigration(); err != nil {
			t.Fatalf("Failed to apply migration: %s", tc.Error())
		}
	})

	// NOTE: not calling tc.Cleanup as the connection
	//       need to be resued by other tests
	// t.Cleanup(func() {
	// 	 tc.Cleanup()
	// })

	return tc
}
