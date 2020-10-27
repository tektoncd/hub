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

package search

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

var res1 = &res.ResourceData{
	ID:   1,
	Name: "foo",
	Kind: "Task",
	Catalog: &res.Catalog{
		ID:   1,
		Type: "community",
	},
	Rating: 4.8,
	LatestVersion: &res.ResourceVersionData{
		ID:                  11,
		Version:             "0.1",
		Description:         "Description for task abc version 0.1",
		DisplayName:         "foo-0.1",
		MinPipelinesVersion: "0.11",
		RawURL:              "http://raw.github.url/foo/",
		WebURL:              "http://web.github.com/foo/",
		UpdatedAt:           "2020-01-01 12:00:00 +0000 UTC",
	},
	Tags: []*res.Tag{
		&res.Tag{ID: 3, Name: "tag3"},
		&res.Tag{ID: 1, Name: "tag1"},
	},
}

var res2 = &res.ResourceData{
	ID:   2,
	Name: "foo-bar",
	Kind: "Pipeline",
	Catalog: &res.Catalog{
		ID:   1,
		Type: "community",
	},
	Rating: 4,
	LatestVersion: &res.ResourceVersionData{
		ID:                  12,
		Version:             "0.2",
		Description:         "Description for pipeline foo-bar version 0.2",
		DisplayName:         "foo-bar-0.1",
		MinPipelinesVersion: "0.12",
		RawURL:              "http://raw.github.url/foo-bar/",
		WebURL:              "http://web.github.com/foo-bar/",
		UpdatedAt:           "2020-01-01 12:00:00 +0000 UTC",
	},
	Tags: []*res.Tag{},
}

func TestValidate(t *testing.T) {
	opt := options{
		kinds:  []string{"pipeline"},
		tags:   []string{"abc,def", "mno"},
		match:  "exact",
		output: "table",
	}

	err := opt.validate()
	assert.NoError(t, err)

	opt = options{
		args:   []string{"abc"},
		match:  "contains",
		output: "json",
	}

	err = opt.validate()
	assert.NoError(t, err)
}

func TestValidate_ErrorCases(t *testing.T) {

	opt := options{}
	err := opt.validate()
	assert.Error(t, err)
	assert.EqualError(t, err, "please specify a name, tag or a kind to search")

	opt = options{
		kinds:  []string{"abc"},
		match:  "exact",
		output: "table",
	}
	err = opt.validate()
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid value \"abc\" set for option kinds. supported kinds: [task, pipeline]")

	opt = options{
		kinds:  []string{"task"},
		match:  "abc",
		output: "table",
	}
	err = opt.validate()
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid value \"abc\" set for option match. Valid options: [contains, exact]")

	opt = options{
		kinds:  []string{"task"},
		match:  "exact",
		output: "abc",
	}
	err = opt.validate()
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid value \"abc\" set for option output. Valid options: [table, json]")
}

func TestSearch_TableFormat(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()
	rArr := &res.Resources{Data: res.ResourceDataCollection{res1, res2}}
	res := res.NewViewedResources(rArr, "withoutVersion")

	gock.New(test.API).
		Get("/query").
		Reply(200).
		JSON(&res.Projected)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opt := &options{
		cli:    cli,
		args:   []string{"foo"},
		match:  "contains",
		output: "table",
	}

	err := opt.run()
	assert.NoError(t, err)
	golden.Assert(t, buf.String(), fmt.Sprintf("%s.golden", t.Name()))
	assert.Equal(t, gock.IsDone(), true)
}

func TestSearch_JSONFormat(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()
	rArr := &res.Resources{Data: res.ResourceDataCollection{res2}}
	res := res.NewViewedResources(rArr, "withoutVersion")

	gock.New(test.API).
		Get("/query").
		Reply(200).
		JSON(&res.Projected)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opt := &options{
		cli:    cli,
		args:   []string{"foo-bar"},
		match:  "exact",
		output: "json",
	}

	err := opt.run()
	assert.NoError(t, err)
	golden.Assert(t, buf.String(), fmt.Sprintf("%s.golden", t.Name()))
	assert.Equal(t, gock.IsDone(), true)
}

func TestSearch_ResourceNotFound(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	gock.New(test.API).
		Get("/query").
		Reply(404).
		JSON(&goa.ServiceError{
			ID:      "123456",
			Name:    "not-found",
			Message: "resource not found",
		})

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opt := &options{
		cli:    cli,
		args:   []string{"xyz"},
		match:  "exact",
		output: "json",
	}

	err := opt.run()
	assert.NoError(t, err)
	assert.Equal(t, gock.IsDone(), true)
}

func TestSearch_InternalServerError(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	gock.New(test.API).
		Get("/query").
		Reply(500).
		JSON(&goa.ServiceError{
			ID:      "123456",
			Name:    "internal-error",
			Message: "failed to fetch resources",
		})

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	opt := &options{
		cli:    cli,
		args:   []string{"xyz"},
		match:  "exact",
		output: "json",
	}

	err := opt.run()
	assert.Error(t, err)
	assert.EqualError(t, err, "Internal server Error: consider filing a bug report")
	assert.Equal(t, gock.IsDone(), true)
}

func TestSearch_InvalidAPIServerURL(t *testing.T) {
	cli := test.NewCLI()

	err := cli.Hub().SetURL("api")
	assert.Error(t, err)
	assert.EqualError(t, err, "parse \"api\": invalid URI for request")
}
