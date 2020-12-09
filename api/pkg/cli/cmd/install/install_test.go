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
	"testing"

	"github.com/stretchr/testify/assert"
	res "github.com/tektoncd/hub/api/gen/resource"
	"github.com/tektoncd/hub/api/pkg/cli/test"
	cb "github.com/tektoncd/hub/api/pkg/cli/test/builder"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	pipelinev1beta1test "github.com/tektoncd/pipeline/test"
	goa "goa.design/goa/v3/pkg"
	"gopkg.in/h2non/gock.v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic/fake"
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
			&res.Tag{
				ID:   3,
				Name: "cli",
			},
		},
	},
}

func TestInstall_NewResource(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	resVersion := &res.ResourceVersion{Data: resVersion}
	res := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&res.Projected)

	gock.New("http://raw.github.url").
		Get("/foo/0.3/foo.yaml").
		Reply(200).
		File("./testdata/foo-v0.3.yaml")

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	version := "v1beta1"
	dynamic := fake.NewSimpleDynamicClient(runtime.NewScheme())

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
	assert.Equal(t, "Task foo(0.3) installed in hub namespace\n", buf.String())
	assert.Equal(t, gock.IsDone(), true)
}

func TestInstall_ResourceNotFoundInHub(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

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

	resVersion := &res.ResourceVersion{Data: resVersion}
	res := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&res.Projected)

	gock.New("http://raw.github.url").
		Get("/foo/0.3/foo.yaml").
		Reply(200).
		File("./testdata/foo-v0.3.yaml")

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	existingTask := &v1beta1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "hub",
		},
	}

	version := "v1beta1"
	dynamic := fake.NewSimpleDynamicClient(runtime.NewScheme(), cb.UnstructuredV1beta1T(existingTask, version))

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

	resVersion := &res.ResourceVersion{Data: resVersion}
	res := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&res.Projected)

	gock.New("http://raw.github.url").
		Get("/foo/0.3/foo.yaml").
		Reply(200).
		File("./testdata/foo-v0.3.yaml")

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
	dynamic := fake.NewSimpleDynamicClient(runtime.NewScheme(), cb.UnstructuredV1beta1T(existingTask, version))

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

	resVersion := &res.ResourceVersion{Data: resVersion}
	res := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&res.Projected)

	gock.New("http://raw.github.url").
		Get("/foo/0.3/foo.yaml").
		Reply(200).
		File("./testdata/foo-v0.3.yaml")

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
	dynamic := fake.NewSimpleDynamicClient(runtime.NewScheme(), cb.UnstructuredV1beta1T(existingTask, version))

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

	resVersion := &res.ResourceVersion{Data: resVersion}
	res := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&res.Projected)

	gock.New("http://raw.github.url").
		Get("/foo/0.3/foo.yaml").
		Reply(200).
		File("./testdata/foo-v0.3.yaml")

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
	dynamic := fake.NewSimpleDynamicClient(runtime.NewScheme(), cb.UnstructuredV1beta1T(existingTask, version))

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
	assert.EqualError(t, err, "Task foo(0.7) already exists in hub namespace. Use reinstall command to overwrite existing")
	assert.Equal(t, gock.IsDone(), true)
}
