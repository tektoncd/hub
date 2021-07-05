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

package cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	goexpect "github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/icmd"
)

const (
	timeout = 5 * time.Minute
)

// TknHubRunner contains information about the current test execution tkn binary path
// under test
type TknHubRunner struct {
	path      string
	namespace string
}

// Prompt provides test utility for prompt test.
type Prompt struct {
	CmdArgs   []string
	Procedure func(*goexpect.Console) error
}

// NewTknHubRunner returns details about the tkn on particular namespace
func NewTknHubRunner() (TknHubRunner, error) {
	if os.Getenv("TEST_CLIENT_BINARY") != "" {
		return TknHubRunner{
			path: os.Getenv("TEST_CLIENT_BINARY"),
		}, nil
	}
	return TknHubRunner{
		path: os.Getenv("TEST_CLIENT_BINARY"),
	}, fmt.Errorf("Error: couldn't Create tknRunner, please do check tkn binary path: (%+v)", os.Getenv("TEST_CLIENT_BINARY"))
}

// Run will help you execute tkn command on a specific namespace, with a timeout
func (e TknHubRunner) Run(args ...string) *icmd.Result {
	if e.namespace != "" {
		args = append(args, "--namespace", e.namespace)
	}
	cmd := append([]string{e.path}, args...)
	return icmd.RunCmd(icmd.Cmd{Command: cmd, Timeout: timeout})
}

// RunInteractiveTests helps to run interactive tests.
func (e TknHubRunner) RunInteractiveTests(t *testing.T, ops *Prompt) *expect.Console {
	t.Helper()

	// Multiplex output to a buffer as well for the raw bytes.
	buf := new(bytes.Buffer)
	c, state, err := vt10x.NewVT10XConsole(goexpect.WithStdout(buf))
	assert.NilError(t, err)
	defer c.Close()

	if e.namespace != "" {
		ops.CmdArgs = append(ops.CmdArgs, "--namespace", e.namespace)
	}

	cmd := exec.Command(e.path, ops.CmdArgs[0:len(ops.CmdArgs)]...) //nolint:gosec
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	assert.NilError(t, cmd.Start())

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		if err := ops.Procedure(c); err != nil {
			t.Logf("procedure failed: %v", err)
		}
	}()

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	_ = c.Tty().Close()
	<-donec

	// Dump the terminal's screen.
	t.Logf("\n%s", goexpect.StripTrailingEmptyLines(state.String()))

	assert.NilError(t, cmd.Wait())

	return c
}

// MustSucceed asserts that the command ran with 0 exit code
func (e TknHubRunner) MustSucceed(t *testing.T, args ...string) *icmd.Result {
	return e.Assert(t, icmd.Success, args...)
}

// Assert runs a command and verifies exit code (0)
func (e TknHubRunner) Assert(t *testing.T, exp icmd.Expected, args ...string) *icmd.Result {
	res := e.Run(args...)
	res.Assert(t, exp)
	return res
}
