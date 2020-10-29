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
	"bytes"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	decoder "k8s.io/apimachinery/pkg/util/yaml"
)

const (
	installerLabel = "hub.tekton.dev/installer"
	catalogLabel   = "hub.tekton.dev/catalog"
)

// Install a resource in the namespace passed to it, will add a label tekton.dev/installer
// to the resource.
func (i *Installer) Install(data []byte, catalog, namespace string) (*unstructured.Unstructured, error) {

	obj, err := toUnstructured(data)
	if err != nil {
		return nil, err
	}

	addHubLabels(obj, catalog)

	res, err := i.create(obj, namespace, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func toUnstructured(data []byte) (*unstructured.Unstructured, error) {

	r := bytes.NewReader(data)
	decoder := decoder.NewYAMLToJSONDecoder(r)

	res := &unstructured.Unstructured{}
	if err := decoder.Decode(res); err != nil {
		return nil, fmt.Errorf("failed to decode resource: %w", err)
	}

	return res, nil
}

func addHubLabels(obj *unstructured.Unstructured, catalog string) {
	labels := obj.GetLabels()
	if len(labels) == 0 {
		labels = make(map[string]string)
	}

	labels[installerLabel] = "hub"
	labels[catalogLabel] = catalog

	obj.SetLabels(labels)
}
