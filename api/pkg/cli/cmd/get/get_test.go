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

package get

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	res "github.com/tektoncd/hub/api/gen/resource"
	"github.com/tektoncd/hub/api/pkg/cli/test"
	goa "goa.design/goa/v3/pkg"
	"gopkg.in/h2non/gock.v1"
	"gotest.tools/v3/golden"
)

var want string = `
Get a Abc of name 'foo':

    tkn hub get abc foo

or

Get a Abc of name 'foo' of version '0.3':

    tkn hub get abc foo --version 0.3
`
var pipelineRes = &res.ResourceData{
	ID:   1,
	Name: "apple",
	Kind: "Pipeline",
	Catalog: &res.Catalog{
		ID:   1,
		Name: "fruit",
		Type: "community",
	},
	Rating: 2.3,
	LatestVersion: &res.ResourceVersionData{
		ID:                  03,
		Version:             "0.1",
		Description:         "v0.1 of Pipeline apple",
		DisplayName:         "Red apple",
		MinPipelinesVersion: "0.17.1",
		RawURL:              "http://raw.github.url/apple/0.1/apple.yaml",
		WebURL:              "http://web.github.com/apple/0.1/apple.yaml",
		UpdatedAt:           "2020-01-01 12:00:00 +0000 UTC",
	},
	Tags: []*res.Tag{
		&res.Tag{
			ID:   3,
			Name: "fruit",
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
	DisplayName:         "mango",
	Description:         "v0.3 of Pipeline mango",
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

func TestValidate(t *testing.T) {
	opt := options{
		version: "0.1",
	}
	err := opt.validate()
	assert.NoError(t, err)

	opt = options{
		version: "0.3.1",
	}
	err = opt.validate()
	assert.NoError(t, err)
}

func TestValidate_ErrorCase(t *testing.T) {
	opt := options{
		version: "abc",
	}
	err := opt.validate()
	assert.EqualError(t, err, "invalid value \"abc\" set for option version. valid eg. 0.1, 1.2.1")
}

func TestGet_WithoutVersion(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	resource := &res.Resource{Data: pipelineRes}
	res := res.NewViewedResource(resource, "default")
	gock.New(test.API).
		Get("/resource/fruit/pipeline/apple").
		Reply(200).
		JSON(&res.Projected)

	gock.New("http://raw.github.url").
		Get("/apple/0.1/apple.yaml").
		Reply(200).
		File("./testdata/pipeline-apple-v0.1.yaml")

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opt := &options{
		cli:  cli,
		kind: "pipeline",
		args: []string{"apple"},
		from: "fruit",
	}

	err := opt.run()
	assert.NoError(t, err)
	golden.Assert(t, buf.String(), fmt.Sprintf("%s.golden", t.Name()))
	assert.Equal(t, gock.IsDone(), true)
}

func TestGet_WithVersion(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	resVersion := &res.ResourceVersion{Data: pipelineResWithVersion}
	res := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/fruit/pipeline/mango/0.3").
		Reply(200).
		JSON(&res.Projected)

	gock.New("http://raw.github.url").
		Get("/mango/0.3/mango.yaml").
		Reply(200).
		File("./testdata/pipeline-mango-v0.3.yaml")

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opt := &options{
		cli:     cli,
		kind:    "pipeline",
		args:    []string{"mango"},
		from:    "fruit",
		version: "0.3",
	}

	err := opt.run()
	assert.NoError(t, err)
	golden.Assert(t, buf.String(), fmt.Sprintf("%s.golden", t.Name()))
	assert.Equal(t, gock.IsDone(), true)
}

func TestGet_ResourceNotFound(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	gock.New(test.API).
		Get("/resource/tekton/pipeline/xyz").
		Reply(404).
		JSON(&goa.ServiceError{
			ID:      "123456",
			Name:    "not-found",
			Message: "resource not found",
		})

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opt := &options{
		cli:  cli,
		kind: "pipeline",
		args: []string{"xyz"},
		from: "tekton",
	}

	err := opt.run()
	assert.Error(t, err)
	assert.EqualError(t, err, "No Resource Found")
	assert.Equal(t, gock.IsDone(), true)
}

func Test_examples(t *testing.T) {
	got := examples("abc")
	assert.Equal(t, want, got)
}
