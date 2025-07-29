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

package info

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/test/cli"
	"gotest.tools/v3/icmd"
)

const expected = `📦 Name: tkn

📌 Version: 0.1 (Deprecated)

📖 Description: This task performs operations on Tekton resources using tkn

🗒  Minimum Pipeline Version: 0.12.1

⭐ ️Rating: 0

🏷️  ️Categories
  ∙ CLI

🏷 Tags
  ∙ cli

💻 Platforms
  ∙ linux/amd64
  ∙ linux/ppc64le
  ∙ linux/s390x

⚒ Install Command:
  tkn hub install task tkn --version 0.1
`

func TestGetIneractiveE2E(t *testing.T) {
	tknhub, err := cli.NewTknHubRunner()
	assert.Nil(t, err)

	t.Run("Get list of Tasks when none present", func(t *testing.T) {
		res := tknhub.Run("info", "task", "foo")
		expected := "Error: No Resource Found"
		res.Assert(t, icmd.Expected{
			ExitCode: 1,
			Err:      expected,
			Out:      icmd.Success.Out,
		})
	})

	t.Logf("Running Info Command for task")

	t.Run("Result for get command when resource name and version is passed", func(t *testing.T) {
		res := tknhub.MustSucceed(t, "info", "task", "tkn", "--version=0.1")
		fmt.Println(res.Stdout())
		assert.Equal(t, expected, res.Stdout())
	})
}
