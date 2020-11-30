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

var resource = &res.ResourceData{
	ID:   1,
	Name: "foo",
	Kind: "Task",
	Catalog: &res.Catalog{
		ID:   1,
		Name: "tekton",
		Type: "community",
	},
	Rating: 4.8,
	LatestVersion: &res.ResourceVersionData{
		ID:                  11,
		Version:             "0.1",
		Description:         "v0.1 Task to run foo",
		DisplayName:         "foo-bar",
		MinPipelinesVersion: "0.11",
		RawURL:              "http://raw.github.url/foo/0.1/foo.yaml",
		WebURL:              "http://web.github.com/foo/0.1/foo.yaml",
		UpdatedAt:           "2020-01-01 12:00:00 +0000 UTC",
	},
	Tags: []*res.Tag{
		&res.Tag{
			ID:   3,
			Name: "cli",
		},
	},
	Versions: []*res.ResourceVersionData{
		&res.ResourceVersionData{
			ID:      11,
			Version: "0.1",
		},
	},
}

var pipelineResWithVersion = &res.ResourceVersionData{
	ID:                  11,
	Version:             "0.3",
	Description:         "v0.3 of Pipeline mango.",
	MinPipelinesVersion: "0.12",
	RawURL:              "http://raw.github.url/mango/0.3/mango.yaml",
	WebURL:              "http://web.github.com/mango/0.3/mango.yaml",
	UpdatedAt:           "2020-01-01 12:00:00 +0000 UTC",
	Resource: &res.ResourceData{
		ID:   7,
		Name: "foo",
		Kind: "Pipeline",
		Catalog: &res.Catalog{
			ID:   1,
			Name: "fruit",
			Type: "community",
		},
		Rating: 4.3,
		Tags: []*res.Tag{
			&res.Tag{
				ID:   3,
				Name: "fruit",
			},
		},
	},
}

func TestInfoTask_WithoutVersion(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	resource := &res.Resource{Data: resource}
	res := res.NewViewedResource(resource, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo").
		Reply(200).
		JSON(&res.Projected)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opts := options{
		cli:  cli,
		kind: "task",
		args: []string{"foo"},
		from: "tekton",
	}

	err := opts.run()
	assert.NoError(t, err)
	golden.Assert(t, buf.String(), fmt.Sprintf("%s.golden", t.Name()))
	assert.Equal(t, gock.IsDone(), true)
}

func TestInfoPipeline_WithVersion(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	resVersion := &res.ResourceVersion{Data: pipelineResWithVersion}
	res := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/fruit/pipeline/mango/0.3").
		Reply(200).
		JSON(&res.Projected)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opts := options{
		cli:     cli,
		kind:    "pipeline",
		args:    []string{"mango"},
		from:    "fruit",
		version: "0.3",
	}

	err := opts.run()
	assert.NoError(t, err)
	golden.Assert(t, buf.String(), fmt.Sprintf("%s.golden", t.Name()))
	assert.Equal(t, gock.IsDone(), true)
}

func TestInfoPipeline_Description(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	pipelineResWithVersion.Description = "Fruits are important sources of vitamins and carbohydrates like fiber and sugar. They are low in calories and naturally sweet. Fruits and their juices are good sources of water, too."

	resVersion := &res.ResourceVersion{Data: pipelineResWithVersion}
	res := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/fruit/pipeline/mango/0.3").
		Reply(200).
		JSON(&res.Projected)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opts := options{
		cli:     cli,
		kind:    "pipeline",
		args:    []string{"mango"},
		from:    "fruit",
		version: "0.3",
	}

	err := opts.run()
	assert.NoError(t, err)
	golden.Assert(t, buf.String(), fmt.Sprintf("%s.golden", t.Name()))
	assert.Equal(t, gock.IsDone(), true)
}

func TestPipelineTask_MultiLineDescription(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	pipelineResWithVersion.Description = "Fruits are important sources of vitamins and carbohydrates like fiber and sugar. \nThey are low in calories and naturally sweet. \nFruits and their juices are good sources of water, too."

	resVersion := &res.ResourceVersion{Data: pipelineResWithVersion}
	res := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/fruit/pipeline/mango/0.3").
		Reply(200).
		JSON(&res.Projected)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opts := options{
		cli:     cli,
		kind:    "pipeline",
		args:    []string{"mango"},
		from:    "fruit",
		version: "0.3",
	}

	err := opts.run()
	assert.NoError(t, err)
	golden.Assert(t, buf.String(), fmt.Sprintf("%s.golden", t.Name()))
	assert.Equal(t, gock.IsDone(), true)
}
