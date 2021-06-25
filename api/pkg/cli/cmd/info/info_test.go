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

package info

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	res "github.com/tektoncd/hub/api/gen/resource"
	"github.com/tektoncd/hub/api/pkg/cli/test"
	"gopkg.in/h2non/gock.v1"
	"gotest.tools/v3/golden"
)

type InfoOptions struct {
	ResId   int
	Name    string
	Kind    string
	Catalog string
	Version string
}

var taskResWithLatestVersion = &res.ResourceVersionData{
	ID:                  12,
	Version:             "0.2",
	Description:         "Description for task foo-bar version 0.2",
	MinPipelinesVersion: "0.12",
	RawURL:              "http://raw.github.url/foo-bar/",
	WebURL:              "http://web.github.com/foo-bar/",
	UpdatedAt:           "2020-01-01 12:00:00 +0000 UTC",
	Resource: &res.ResourceData{
		ID:   2,
		Name: "foo-bar",
		Kind: "Task",
		Catalog: &res.Catalog{
			ID:   1,
			Name: "tekton",
			Type: "community",
		},
		Rating: 4,
		Tags: []*res.Tag{
			{
				ID:   3,
				Name: "foo",
			},
		},
	},
}

var taskResWithOldVersion = &res.ResourceVersionData{
	ID:                  12,
	Version:             "0.1",
	Description:         "Description for task foo-bar version 0.1",
	MinPipelinesVersion: "0.12",
	RawURL:              "http://raw.github.url/foo-bar/",
	WebURL:              "http://web.github.com/foo-bar/",
	UpdatedAt:           "2020-01-01 12:00:00 +0000 UTC",
	Resource: &res.ResourceData{
		ID:   2,
		Name: "foo-bar",
		Kind: "Task",
		Catalog: &res.Catalog{
			ID:   1,
			Name: "tekton",
			Type: "community",
		},
		Rating: 4,
		Tags: []*res.Tag{
			{
				ID:   3,
				Name: "foo",
			},
		},
	},
}

func mockApi(io InfoOptions, taskWithVersion *res.ResourceVersionData) {

	// Get ResourceId in order to get all versions of resource
	rVer := &res.ResourceVersion{Data: taskWithVersion}
	resWithVersion := res.NewViewedResourceVersion(rVer, "default")
	resInfo := fmt.Sprintf("%s/%s/%s", io.Catalog, io.Kind, io.Name)

	gock.New(test.API).
		Get("/resource/" + resInfo + "/" + io.Version).
		Reply(200).
		JSON(&resWithVersion.Projected)
}

func TestInfoTask_WithLatestVersion(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	mockApi(InfoOptions{
		ResId:   12,
		Name:    "foo-bar",
		Kind:    "task",
		Catalog: "tekton",
		Version: "0.2",
	}, taskResWithLatestVersion)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opts := options{
		cli:     cli,
		kind:    "task",
		args:    []string{"foo-bar"},
		from:    "tekton",
		version: "0.2",
	}

	err := opts.run()
	assert.NoError(t, err)
	golden.Assert(t, buf.String(), fmt.Sprintf("%s.golden", t.Name()))
	assert.Equal(t, gock.IsDone(), true)
}

func TestInfoTask_WithOldVersion(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	mockApi(InfoOptions{
		ResId:   12,
		Name:    "foo-bar",
		Kind:    "task",
		Catalog: "tekton",
		Version: "0.1",
	}, taskResWithOldVersion)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opts := options{
		cli:     cli,
		kind:    "task",
		args:    []string{"foo-bar"},
		from:    "tekton",
		version: "0.1",
	}

	err := opts.run()
	assert.NoError(t, err)
	golden.Assert(t, buf.String(), fmt.Sprintf("%s.golden", t.Name()))
	assert.Equal(t, gock.IsDone(), true)
}

func TestPipelineTask_MultiLineDescription(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	taskResWithLatestVersion.Description = "A Task is a collection of Steps that you define and arrange in a specific order of execution as part of your continuous integration flow. A Task executes as a Pod on your Kubernetes cluster. A Task is available within a specific namespace, while a ClusterTask is available across the entire cluster."

	mockApi(InfoOptions{
		ResId:   12,
		Name:    "foo-bar",
		Kind:    "task",
		Catalog: "tekton",
		Version: "0.2",
	}, taskResWithLatestVersion)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opts := options{
		cli:     cli,
		kind:    "task",
		args:    []string{"foo-bar"},
		from:    "tekton",
		version: "0.2",
	}

	err := opts.run()
	assert.NoError(t, err)
	golden.Assert(t, buf.String(), fmt.Sprintf("%s.golden", t.Name()))
	assert.Equal(t, gock.IsDone(), true)
}
