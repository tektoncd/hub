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

package downgrade

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	res "github.com/tektoncd/hub/api/gen/resource"
	"github.com/tektoncd/hub/api/pkg/cli/test"
	cb "github.com/tektoncd/hub/api/pkg/cli/test/builder"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	pipelinev1beta1test "github.com/tektoncd/pipeline/test"
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

var ver03 = &res.ResourceVersionData{
	ID:      111,
	Version: "0.3",
	RawURL:  "http://raw.github.url/foo/0.3/foo.yaml",
	WebURL:  "http://web.github.com/foo/0.3/foo.yaml",
}
var ver02 = &res.ResourceVersionData{
	ID:      113,
	Version: "0.2",
	RawURL:  "http://raw.github.url/foo/0.2/foo.yaml",
	WebURL:  "http://web.github.com/foo/0.2/foo.yaml",
}

func TestDowngrade_ResourceNotExist(t *testing.T) {
	cli := test.NewCLI()

	version := "v1beta1"
	dynamic := fake.NewSimpleDynamicClient(runtime.NewScheme())

	cs, _ := test.SeedV1beta1TestData(t, pipelinev1beta1test.Data{})
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})

	opts := &options{
		cs:   test.FakeClientSet(cs.Pipeline, dynamic, "hub"),
		cli:  cli,
		kind: "task",
		args: []string{"foo"},
	}

	err := opts.run()
	assert.Error(t, err)
	assert.EqualError(t, err, "Task foo doesn't exist in hub namespace. Use install command to install the task")
}

func TestDowngrade_VersionCatalogMissing(t *testing.T) {
	cli := test.NewCLI()

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
		cs:   test.FakeClientSet(cs.Pipeline, dynamic, "hub"),
		cli:  cli,
		kind: "task",
		args: []string{"foo"},
	}

	err := opts.run()
	assert.Error(t, err)
	assert.EqualError(t, err, "Task foo seems to be missing version and catalog label. Use reinstall command to overwrite existing task")
}

func TestDowngrade_VersionMissing(t *testing.T) {
	cli := test.NewCLI()

	existingTask := &v1beta1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "hub",
			Labels:    map[string]string{"hub.tekton.dev/catalog": "abc"},
		},
	}

	version := "v1beta1"
	dynamic := fake.NewSimpleDynamicClient(runtime.NewScheme(), cb.UnstructuredV1beta1T(existingTask, version))

	cs, _ := test.SeedV1beta1TestData(t, pipelinev1beta1test.Data{Tasks: []*v1beta1.Task{existingTask}})
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})

	opts := &options{
		cs:   test.FakeClientSet(cs.Pipeline, dynamic, "hub"),
		cli:  cli,
		kind: "task",
		args: []string{"foo"},
	}

	err := opts.run()
	assert.Error(t, err)
	assert.EqualError(t, err, "Task foo seems to be missing version label. Use reinstall command to overwrite existing task")
}

func TestDowngrade(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	resVersion := &res.ResourceVersion{Data: resVersion}
	resource := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&resource.Projected)

	versions := &res.ResourceVersions{Data: &res.Versions{Latest: ver03, Versions: []*res.ResourceVersionData{ver02, ver03}}}
	ver := res.NewViewedResourceVersions(versions, "default")

	gock.New(test.API).
		Get("/resource/1/versions").
		Reply(200).
		JSON(&ver.Projected)

	gock.New("http://raw.github.url").
		Get("/foo/0.2/foo.yaml").
		Reply(200).
		File("./testdata/foo-v0.2.yaml")

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	existingTask := &v1beta1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "hub",
			Labels: map[string]string{
				"hub.tekton.dev/catalog":    "tekton",
				"app.kubernetes.io/version": "0.3",
			}},
	}

	version := "v1beta1"
	dynamic := fake.NewSimpleDynamicClient(runtime.NewScheme(), cb.UnstructuredV1beta1T(existingTask, version))

	cs, _ := test.SeedV1beta1TestData(t, pipelinev1beta1test.Data{Tasks: []*v1beta1.Task{existingTask}})
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})

	opts := &options{
		cs:   test.FakeClientSet(cs.Pipeline, dynamic, "hub"),
		cli:  cli,
		kind: "task",
		args: []string{"foo"},
	}

	err := opts.run()
	assert.NoError(t, err)
	assert.Equal(t, "Task foo downgraded to v0.2 in hub namespace\n", buf.String())
	assert.Equal(t, gock.IsDone(), true)
}

func TestDowngrade_ToSpecificVersion(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	resVersion := &res.ResourceVersion{Data: resVersion}
	resource := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&resource.Projected)

	versions := &res.ResourceVersions{Data: &res.Versions{Latest: ver03, Versions: []*res.ResourceVersionData{ver02, ver03}}}
	ver := res.NewViewedResourceVersions(versions, "default")

	gock.New(test.API).
		Get("/resource/1/versions").
		Reply(200).
		JSON(&ver.Projected)

	gock.New("http://raw.github.url").
		Get("/foo/0.2/foo.yaml").
		Reply(200).
		File("./testdata/foo-v0.2.yaml")

	buf := new(bytes.Buffer)
	cli.SetStream(buf, buf)

	existingTask := &v1beta1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "hub",
			Labels: map[string]string{
				"hub.tekton.dev/catalog":    "tekton",
				"app.kubernetes.io/version": "0.3",
			}},
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
		version: "0.2",
	}

	err := opts.run()
	assert.NoError(t, err)
	assert.Equal(t, "Task foo downgraded to v0.2 in hub namespace\n", buf.String())
	assert.Equal(t, gock.IsDone(), true)
}

func TestDowngrade_SameVersionError(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	resVersion := &res.ResourceVersion{Data: resVersion}
	resource := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&resource.Projected)

	versions := &res.ResourceVersions{Data: &res.Versions{Latest: ver03, Versions: []*res.ResourceVersionData{ver02, ver03}}}
	ver := res.NewViewedResourceVersions(versions, "default")

	gock.New(test.API).
		Get("/resource/1/versions").
		Reply(200).
		JSON(&ver.Projected)

	gock.New("http://raw.github.url").
		Get("/foo/0.3/foo.yaml").
		Reply(200).
		File("./testdata/foo-v0.3.yaml")

	existingTask := &v1beta1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "hub",
			Labels: map[string]string{
				"hub.tekton.dev/catalog":    "tekton",
				"app.kubernetes.io/version": "0.3",
			}},
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
		version: "0.3",
	}

	err := opts.run()
	assert.Error(t, err)
	assert.EqualError(t, err, "cannot downgrade task foo to v0.3. existing resource seems to be of same version. Use reinstall command to overwrite existing task")
	assert.Equal(t, gock.IsDone(), true)
}

func TestDowngrade_HigherVersionError(t *testing.T) {
	cli := test.NewCLI()

	defer gock.Off()

	resVersion := &res.ResourceVersion{Data: resVersion}
	resource := res.NewViewedResourceVersion(resVersion, "default")
	gock.New(test.API).
		Get("/resource/tekton/task/foo/0.3").
		Reply(200).
		JSON(&resource.Projected)

	versions := &res.ResourceVersions{Data: &res.Versions{Latest: ver03, Versions: []*res.ResourceVersionData{ver02, ver03}}}
	ver := res.NewViewedResourceVersions(versions, "default")

	gock.New(test.API).
		Get("/resource/1/versions").
		Reply(200).
		JSON(&ver.Projected)

	existingTask := &v1beta1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "hub",
			Labels: map[string]string{
				"hub.tekton.dev/catalog":    "tekton",
				"app.kubernetes.io/version": "0.3",
			}},
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
		version: "0.4",
	}

	err := opts.run()
	assert.Error(t, err)
	assert.EqualError(t, err, "cannot downgrade task foo to v0.4. existing resource seems to be of lower version(v0.3). Use upgrade command")
	assert.Equal(t, gock.IsDone(), true)
}
