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

package builder

import (
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func UnstructuredV1beta1T(task *v1beta1.Task, version string) *unstructured.Unstructured {
	task.APIVersion = "tekton.dev/" + version
	task.Kind = "task"
	object, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(task)
	return &unstructured.Unstructured{
		Object: object,
	}
}

func UnstructuredV1(task *v1.Task, version string) *unstructured.Unstructured {
	task.APIVersion = "tekton.dev/" + version
	task.Kind = "task"
	object, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(task)
	return &unstructured.Unstructured{
		Object: object,
	}
}
