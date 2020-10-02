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

package test

import (
	"io"

	"github.com/tektoncd/hub/cli/pkg/app"
	"github.com/tektoncd/hub/cli/pkg/hub"
)

type cli struct {
	hub    hub.Client
	stream app.Stream
}

var _ app.CLI = (*cli)(nil)

func NewCLI() *cli {
	return &cli{
		stream: app.Stream{},
		hub:    hub.NewClient(),
	}
}

func (c *cli) Stream() app.Stream {
	return c.stream
}

func (c *cli) SetStream(out, err io.Writer) {
	c.stream = app.Stream{Out: out, Err: err}
}

func (c *cli) Hub() hub.Client {
	return c.hub
}
