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

package installer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/pkg/cli/test"
	cb "github.com/tektoncd/hub/api/pkg/cli/test/builder"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	pipelinev1beta1test "github.com/tektoncd/pipeline/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic/fake"
)

const res = `---
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

func TestToUnstructuredAndAddLabel(t *testing.T) {
	obj, err := toUnstructured([]byte(res))
	assert.NoError(t, err)
	assert.Equal(t, "foo", obj.GetName())

	addCatalogLabel(obj, "tekton")
	assert.Equal(t, "tekton", obj.GetLabels()[catalogLabel])
}

func TestListInstalled(t *testing.T) {
	existingTask := &v1beta1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "hub",
			Labels: map[string]string{
				"hub.tekton.dev/catalog":    "tekton",
				"app.kubernetes.io/version": "0.1",
			}},
	}

	version := "v1beta1"
	dynamic := fake.NewSimpleDynamicClient(runtime.NewScheme(), cb.UnstructuredV1beta1T(existingTask, version))

	cs, _ := test.SeedV1beta1TestData(t, pipelinev1beta1test.Data{Tasks: []*v1beta1.Task{existingTask}})
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})

	clientSet := test.FakeClientSet(cs.Pipeline, dynamic, "hub")

	installer := New(clientSet)
	list, _ := installer.ListInstalled("task", "hub")

	assert.Equal(t, len(list), 1)
}
