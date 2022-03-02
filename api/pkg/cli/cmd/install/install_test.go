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

package install

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/pkg/cli/test"
	cb "github.com/tektoncd/hub/api/pkg/cli/test/builder"
	res "github.com/tektoncd/hub/api/v1/gen/resource"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	pipelinev1beta1test "github.com/tektoncd/pipeline/test"
	goa "goa.design/goa/v3/pkg"
	"gopkg.in/h2non/gock.v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var resVersion = &res.ResourceVersionData{
	ID:                  11,
	Version:             "0.3",
	DisplayName:         "foo-bar",
	Description:         "v0.3 Task to run foo",
	MinPipelinesVersion: "0.12",
	RawURL:              "http://raw.github.url/foo/0.3/foo.yaml",
	WebURL:              "http://web.github.com/foo/0.3/foo.yaml",
	UpdatedAt:           "2020-01-01 12:00:00 +0000 UTC",
	Resource: &res.ResourceData{
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
			{
				ID:   3,
				Name: "cli",
			},
		},
	},
}

var task1 = `---
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: foo
  labels:
    app.kubernetes.io/version: '0.3'
  annotations:
    tekton.dev/pipelines.minVersion: '0.13.1'
    tekton.dev/tags: cli
    tekton.dev/displayName: 'foo-bar'
spec:
  description: >-
    v0.3 Task to run foo
`

var taskWithNewVersionYaml = &res.ResourceContent{
	Yaml: &task1,
}

var task2 = `---
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: foo
  labels:
    app.kubernetes.io/version: '0.2'
  annotations:
    tekton.dev/pipelines.minVersion: '0.13.1'
    tekton.dev/tags: cli
    tekton.dev/displayName: 'foo-bar'
    tekton.dev/deprecated: 'true'
spec:
  description: >-
    v0.2 Task to run foo
`

var taskWithOldVersionYaml = &res.ResourceContent{
	Yaml: &task2,
}

func TestInstall_NewResource(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	rVer := &res.ResourceVersionYaml{Data: taskWithNewVersionYaml}
	resWithVersion := res.NewViewedResourceVersionYaml(rVer, "default")

	resInfo := fmt.Sprintf("%s/%s/%s/%s", "tekton", "task", "foo", "0.3")

	gock.New(test.API).
		Get("/resource/" + resInfo + "/yaml").
		Reply(200).
		JSON(&resWithVersion.Projected)

	resVersion := &res.ResourceVersion{Data: resVersion}
	resource := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&resource.Projected)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	version := "v1beta1"
	dynamic := test.DynamicClient()

	cs, _ := test.SeedV1beta1TestData(t, pipelinev1beta1test.Data{})
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})

	opts := &options{
		cs:      test.FakeClientSet(cs.Pipeline, dynamic, "hub"),
		cli:     cli,
		kind:    "task",
		args:    []string{"foo"},
		from:    "tekton",
		version: "0.3",
	}

	err := opts.run()
	assert.NoError(t, err)
	assert.Equal(t, "WARN: tekton pipelines version unknown, this resource is compatible with pipelines min version v0.12\nTask foo(0.3) installed in hub namespace\n", buf.String())
	assert.Equal(t, gock.IsDone(), true)
}

func TestInstall_ResourceNotFoundInHub(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3/yaml").
		Reply(404).
		JSON(&goa.ServiceError{
			ID:      "123456",
			Name:    "not-found",
			Message: "resource not found",
		})

	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(404).
		JSON(&goa.ServiceError{
			ID:      "123456",
			Name:    "not-found",
			Message: "resource not found",
		})

	opts := &options{
		cli:     cli,
		kind:    "task",
		args:    []string{"foo"},
		from:    "tekton",
		version: "0.3",
	}

	err := opts.run()
	assert.Error(t, err)
	assert.EqualError(t, err, "Task foo(0.3) from tekton catalog not found in Hub")
}

func TestInstall_ResourceAlreadyExistError(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	rVer := &res.ResourceVersionYaml{Data: taskWithNewVersionYaml}
	resWithVersion := res.NewViewedResourceVersionYaml(rVer, "default")

	resInfo := fmt.Sprintf("%s/%s/%s/%s", "tekton", "task", "foo", "0.3")

	gock.New(test.API).
		Get("/resource/" + resInfo + "/yaml").
		Reply(200).
		JSON(&resWithVersion.Projected)

	resVersion := &res.ResourceVersion{Data: resVersion}
	resource := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&resource.Projected)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	existingTask := &v1beta1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "hub",
		},
	}

	version := "v1beta1"
	dynamic := test.DynamicClient(cb.UnstructuredV1beta1T(existingTask, version))

	cs, _ := test.SeedV1beta1TestData(t, pipelinev1beta1test.Data{Tasks: []*v1beta1.Task{existingTask}})
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})

	opts := &options{
		cs:      test.FakeClientSet(cs.Pipeline, dynamic, "hub"),
		cli:     cli,
		kind:    "task",
		args:    []string{"foo"},
		from:    "tekton",
		version: "0.3",
	}

	err := opts.run()
	assert.Error(t, err)
	assert.EqualError(t, err, "Task foo already exists in hub namespace but seems to be missing version label. Use reinstall command to overwrite existing")
	assert.Equal(t, gock.IsDone(), true)
}

func TestInstall_UpgradeError(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	rVer := &res.ResourceVersionYaml{Data: taskWithNewVersionYaml}
	resWithVersion := res.NewViewedResourceVersionYaml(rVer, "default")

	resInfo := fmt.Sprintf("%s/%s/%s/%s", "tekton", "task", "foo", "0.3")

	gock.New(test.API).
		Get("/resource/" + resInfo + "/yaml").
		Reply(200).
		JSON(&resWithVersion.Projected)

	resVersion := &res.ResourceVersion{Data: resVersion}
	resource := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&resource.Projected)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	existingTask := &v1beta1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "hub",
			Labels:    map[string]string{"app.kubernetes.io/version": "0.1"},
		},
	}

	version := "v1beta1"
	dynamic := test.DynamicClient(cb.UnstructuredV1beta1T(existingTask, version))

	cs, _ := test.SeedV1beta1TestData(t, pipelinev1beta1test.Data{Tasks: []*v1beta1.Task{existingTask}})
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})

	opts := &options{
		cs:      test.FakeClientSet(cs.Pipeline, dynamic, "hub"),
		cli:     cli,
		kind:    "task",
		args:    []string{"foo"},
		from:    "tekton",
		version: "0.3",
	}

	err := opts.run()
	assert.Error(t, err)
	assert.EqualError(t, err, "Task foo(0.1) already exists in hub namespace. Use upgrade command to install v0.3")
	assert.Equal(t, gock.IsDone(), true)
}

func TestInstall_SameVersionError(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	rVer := &res.ResourceVersionYaml{Data: taskWithNewVersionYaml}
	resWithVersion := res.NewViewedResourceVersionYaml(rVer, "default")

	resInfo := fmt.Sprintf("%s/%s/%s/%s", "tekton", "task", "foo", "0.3")

	gock.New(test.API).
		Get("/resource/" + resInfo + "/yaml").
		Reply(200).
		JSON(&resWithVersion.Projected)

	resVersion := &res.ResourceVersion{Data: resVersion}
	resource := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&resource.Projected)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	existingTask := &v1beta1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "hub",
			Labels:    map[string]string{"app.kubernetes.io/version": "0.3"},
		},
	}

	version := "v1beta1"
	dynamic := test.DynamicClient(cb.UnstructuredV1beta1T(existingTask, version))

	cs, _ := test.SeedV1beta1TestData(t, pipelinev1beta1test.Data{Tasks: []*v1beta1.Task{existingTask}})
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})

	opts := &options{
		cs:      test.FakeClientSet(cs.Pipeline, dynamic, "hub"),
		cli:     cli,
		kind:    "task",
		args:    []string{"foo"},
		from:    "tekton",
		version: "0.3",
	}

	err := opts.run()
	assert.Error(t, err)
	assert.EqualError(t, err, "Task foo(0.3) already exists in hub namespace. Use reinstall command to overwrite existing")
	assert.Equal(t, gock.IsDone(), true)
}

func TestInstall_LowerVersionError(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	rVer := &res.ResourceVersionYaml{Data: taskWithNewVersionYaml}
	resWithVersion := res.NewViewedResourceVersionYaml(rVer, "default")

	resInfo := fmt.Sprintf("%s/%s/%s/%s", "tekton", "task", "foo", "0.3")

	gock.New(test.API).
		Get("/resource/" + resInfo + "/yaml").
		Reply(200).
		JSON(&resWithVersion.Projected)

	resVersion := &res.ResourceVersion{Data: resVersion}
	resource := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&resource.Projected)
	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	existingTask := &v1beta1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "hub",
			Labels:    map[string]string{"app.kubernetes.io/version": "0.7"},
		},
	}

	version := "v1beta1"
	dynamic := test.DynamicClient(cb.UnstructuredV1beta1T(existingTask, version))

	cs, _ := test.SeedV1beta1TestData(t, pipelinev1beta1test.Data{Tasks: []*v1beta1.Task{existingTask}})
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})

	opts := &options{
		cs:      test.FakeClientSet(cs.Pipeline, dynamic, "hub"),
		cli:     cli,
		kind:    "task",
		args:    []string{"foo"},
		from:    "tekton",
		version: "0.3",
	}

	err := opts.run()
	assert.Error(t, err)
	assert.EqualError(t, err, "Task foo(0.7) already exists in hub namespace. Use downgrade command to install v0.3")
	assert.Equal(t, gock.IsDone(), true)
}

func TestInstall_RespectingPipelinesVersion(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	rVer := &res.ResourceVersionYaml{Data: taskWithNewVersionYaml}
	resWithVersion := res.NewViewedResourceVersionYaml(rVer, "default")

	resInfo := fmt.Sprintf("%s/%s/%s/%s", "tekton", "task", "foo", "0.3")

	gock.New(test.API).
		Get("/resource/" + resInfo + "/yaml").
		Reply(200).
		JSON(&resWithVersion.Projected)

	resVersion := &res.ResourceVersion{Data: resVersion}
	resource := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&resource.Projected)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	version := "v1beta1"
	dynamic := test.DynamicClient()

	cs, _ := test.SeedV1beta1TestData(t, pipelinev1beta1test.Data{})
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})

	opts := &options{
		cs:      test.FakeClientSet(cs.Pipeline, dynamic, "hub"),
		cli:     cli,
		kind:    "task",
		args:    []string{"foo"},
		from:    "tekton",
		version: "0.3",
	}

	err := test.CreateTektonPipelineController(dynamic, "v0.14.0")
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	err = opts.run()
	assert.NoError(t, err)
	assert.Equal(t, "Task foo(0.3) installed in hub namespace\n", buf.String())
	assert.Equal(t, gock.IsDone(), true)
}

func TestInstall_RespectingPipelinesVersionFailure(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	rVer := &res.ResourceVersionYaml{Data: taskWithNewVersionYaml}
	resWithVersion := res.NewViewedResourceVersionYaml(rVer, "default")

	resInfo := fmt.Sprintf("%s/%s/%s/%s", "tekton", "task", "foo", "0.3")

	gock.New(test.API).
		Get("/resource/" + resInfo + "/yaml").
		Reply(200).
		JSON(&resWithVersion.Projected)

	resVersion := &res.ResourceVersion{Data: resVersion}
	resource := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&resource.Projected)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	version := "v1beta1"
	dynamic := test.DynamicClient()

	cs, _ := test.SeedV1beta1TestData(t, pipelinev1beta1test.Data{})
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})

	opts := &options{
		cs:      test.FakeClientSet(cs.Pipeline, dynamic, "hub"),
		cli:     cli,
		kind:    "task",
		args:    []string{"foo"},
		from:    "tekton",
		version: "0.3",
	}

	err := test.CreateTektonPipelineController(dynamic, "v0.11.0")
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	err = opts.run()
	assert.Error(t, err)
	assert.EqualError(t, err, "Task foo(0.3) requires Tekton Pipelines min version v0.12 but found v0.11.0")
	assert.Equal(t, gock.IsDone(), true)
}

func TestInstall_DeprecatedVersion(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	rVer := &res.ResourceVersionYaml{Data: taskWithOldVersionYaml}
	resWithVersion := res.NewViewedResourceVersionYaml(rVer, "default")

	resInfo := fmt.Sprintf("%s/%s/%s/%s", "tekton", "task", "foo", "0.2")

	gock.New(test.API).
		Get("/resource/" + resInfo + "/yaml").
		Reply(200).
		JSON(&resWithVersion.Projected)

	resVersion := &res.ResourceVersion{Data: resVersion}
	resource := res.NewViewedResourceVersion(resVersion, "0.2")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.2").
		Reply(200).
		JSON(&resource.Projected)

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	version := "v1beta1"
	dynamic := test.DynamicClient()

	cs, _ := test.SeedV1beta1TestData(t, pipelinev1beta1test.Data{})
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})

	opts := &options{
		cs:      test.FakeClientSet(cs.Pipeline, dynamic, "hub"),
		cli:     cli,
		kind:    "task",
		args:    []string{"foo"},
		from:    "tekton",
		version: "0.2",
	}

	err := opts.run()

	assert.NoError(t, err)
	assert.Equal(t, "WARN: tekton pipelines version unknown, this resource is compatible with pipelines min version v0.12\nWARN: This version has been deprecated\nTask foo(0.2) installed in hub namespace\n", buf.String())
	assert.Equal(t, gock.IsDone(), true)
}
