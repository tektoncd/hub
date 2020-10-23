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

var resource = &res.Resource{
	ID:   1,
	Name: "foo",
	Kind: "Task",
	Catalog: &res.Catalog{
		ID:   1,
		Name: "tekton",
		Type: "community",
	},
	Rating: 4.8,
	LatestVersion: &res.Version{
		ID:                  11,
		Version:             "0.1",
		Description:         "Description for task abc version 0.1",
		DisplayName:         "foo-0.1",
		MinPipelinesVersion: "0.11",
		RawURL:              "http://raw.github.url/foo/0.1/foo.yaml",
		WebURL:              "http://web.github.com/foo/0.1/foo.yaml",
		UpdatedAt:           "2020-01-01 12:00:00 +0000 UTC",
	},
	Tags: []*res.Tag{
		&res.Tag{
			ID:   3,
			Name: "tag3",
		},
		&res.Tag{
			ID:   1,
			Name: "tag1",
		},
	},
	Versions: []*res.Version{
		&res.Version{
			ID:      11,
			Version: "0.1",
		},
	},
}

var resVersion = &res.Version{
	ID:                  11,
	Version:             "0.3",
	DisplayName:         "foo",
	Description:         "Description for task abc version 0.3",
	MinPipelinesVersion: "0.12",
	RawURL:              "http://raw.github.url/foo/0.3/foo.yaml",
	WebURL:              "http://web.github.com/foo/0.3/foo.yaml",
	UpdatedAt:           "2020-01-01 12:00:00 +0000 UTC",
	Resource: &res.Resource{
		ID:   1,
		Name: "foo",
		Kind: "Task",
		Catalog: &res.Catalog{
			ID:   1,
			Name: "tekton",
			Type: "community",
		},
		Rating: 4.8,
		Tags: []*res.Tag{
			&res.Tag{
				ID:   3,
				Name: "tag3",
			},
			&res.Tag{
				ID:   1,
				Name: "tag1",
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

	res := res.NewViewedResource(resource, "default")
	gock.New(test.API).
		Get("/resource").
		Reply(200).
		JSON(&res.Projected)

	gock.New("http://raw.github.url").
		Get("/foo").
		Reply(200).
		File("./testdata/foo-v0.1.yaml")

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opt := &options{
		cli:  cli,
		args: []string{"foo"},
		from: "tekton",
	}

	err := opt.run()
	assert.NoError(t, err)
	golden.Assert(t, buf.String(), fmt.Sprintf("%s.golden", t.Name()))
	assert.Equal(t, gock.IsDone(), true)
}

func TestGet_WithVersion(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	res := res.NewViewedVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource").
		Reply(200).
		JSON(&res.Projected)

	gock.New("http://raw.github.url").
		Get("/foo").
		Reply(200).
		File("./testdata/foo-v0.3.yaml")

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opt := &options{
		cli:     cli,
		args:    []string{"foo"},
		from:    "tekton",
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
		Get("/resource").
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
