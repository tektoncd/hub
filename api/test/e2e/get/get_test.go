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
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/test/cli"
	"gotest.tools/v3/icmd"
)

const expected = `---
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: tkn
  labels:
    app.kubernetes.io/version: "0.1"
  annotations:
    tekton.dev/pipelines.minVersion: "0.12.1"
    tekton.dev/categories: CLI
    tekton.dev/tags: cli
    tekton.dev/platforms: "linux/amd64,linux/s390x,linux/ppc64le"
spec:
  description: >-
    This task performs operations on Tekton resources using tkn

  params:
  - name: tkn-image
    description: tkn CLI container image to run this task
    default: gcr.io/tekton-releases/dogfooding/tkn@sha256:f79c0a56d561b1ae98afca545d0feaf2759f5c486ef891a79f23dc2451167dad #tag: latest
  - name: ARGS
    type: array
    description: tkn CLI arguments to run
  steps:
    - name: tkn
      image: "$(params.tkn-image)"
      command: ["/usr/local/bin/tkn"]
      args: ["$(params.ARGS)"]

`

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
	t.Logf("Running Get Command for task")

	t.Run("Interactive mode for get command", func(t *testing.T) {
		tknhub.RunInteractiveTests(t, &cli.Prompt{
			CmdArgs: []string{"get", "task"},
			Procedure: func(c *expect.Console) error {
				if _, err := c.ExpectString("Select task:"); err != nil {
					return err
				}
				if _, err := c.SendLine("az"); err != nil {
					return err
				}
				if _, err := c.ExpectString("Select version:"); err != nil {
					return err
				}
				if _, err := c.SendLine(string(terminal.KeyEnter)); err != nil {
					return err
				}
				if _, err := c.ExpectEOF(); err != nil {
					return err
				}

				c.Close()
				return nil
			}})
	})

	t.Run("Result for get command when resource name and version is passed", func(t *testing.T) {
		res := tknhub.MustSucceed(t, "get", "task", "tkn", "--version=0.1")
		assert.Equal(t, expected, res.Stdout())
	})
}
