// +build e2e
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

package get

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/test/cli"
	"gotest.tools/v3/icmd"
)

var expected = `---
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: foo
  annotations:
    tekton.dev/pipelines.minVersion: '0.12.1'
    tekton.dev/tags: cli
    tekton.dev/displayName: 'foo-bar'
spec:
  description: >-
    v0.1 Task to run foo
`
var apiResponse = []byte(`{
	"data": {
"yaml": "---\napiVersion: tekton.dev/v1beta1\nkind: Task\nmetadata:\n  name: foo\n  annotations:\n    tekton.dev/pipelines.minVersion: '0.12.1'\n    tekton.dev/tags: cli\n    tekton.dev/displayName: 'foo-bar'\nspec:\n  description: >-\n    v0.1 Task to run foo"
}
}`)

func TestGetIneractiveE2E(t *testing.T) {
	tknhub, err := cli.NewTknHubRunner()
	assert.Nil(t, err)

	t.Run("Get list of Tasks when none present", func(t *testing.T) {
		res := tknhub.Run("get", "task", "foo")
		expected := "Error: No Resource Found"
		res.Assert(t, icmd.Expected{
			ExitCode: 1,
			Err:      expected,
			Out:      icmd.Success.Out,
		})

	})

	t.Run("Result for get command when resource name and version is passed", func(t *testing.T) {

		go func() {
			http.HandleFunc("/v1/resource/tekton/task/foo/0.1/yaml", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(apiResponse))
			})
			log.Fatal(http.ListenAndServe(":8080", nil))
		}()

		res := tknhub.MustSucceed(t, "get", "task", "foo", "--version=0.1", "--from=tekton", "--api-server=http://localhost:8080")
		assert.Equal(t, expected, res.Stdout())
	})
}
