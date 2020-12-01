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
	"errors"
	"fmt"

	kErr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	decoder "k8s.io/apimachinery/pkg/util/yaml"
)

const (
	installerLabel = "hub.tekton.dev/installer"
	catalogLabel   = "hub.tekton.dev/catalog"
	versionLabel   = "app.kubernetes.io/version"
)

// Errors
var (
	ErrNotFromCatalog = errors.New("resource already exists but doesn't seem to be from a catalog")
	ErrAlreadyExist   = errors.New("resource already exists")
)

// Install a resource in the namespace passed to it, will add a label tekton.dev/installer
// to the resource.
func (i *Installer) Install(data []byte, catalog, namespace string, overwrite bool) (*unstructured.Unstructured, error) {

	newRes, err := toUnstructured(data)
	if err != nil {
		return nil, err
	}

	// Check if resource already exists
	existingRes, err := i.get(newRes.GetName(), newRes.GetKind(), namespace, metav1.GetOptions{})
	if err != nil {
		// If error is notFoundError then create the resource
		if kErr.IsNotFound(err) {
			return i.createRes(newRes, catalog, namespace)
		}
		// otherwise return the error
		return nil, err
	}

	// Resource exists then check if it is a resource from a catalog
	// If version label is not available then it means the resource is not from a catalog. Currently,
	// we are checking this using version label later it can be replaced by catalog label once all the
	// resources in catalog have a label which has the name of catalog.
	if !isVersionLabelAvailable(existingRes) {
		// If overwrite is true then replace existing resource with new
		if overwrite {
			return i.updateRes(existingRes, newRes, catalog, namespace)
		}
		// If overwrite is false then return error resource not from a catalog
		return existingRes, ErrNotFromCatalog
	}

	return existingRes, ErrAlreadyExist
}

func (i *Installer) createRes(obj *unstructured.Unstructured, catalog, namespace string) (*unstructured.Unstructured, error) {
	addHubLabels(obj, catalog)
	res, err := i.create(obj, namespace, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *Installer) updateRes(existing, new *unstructured.Unstructured, catalog, namespace string) (*unstructured.Unstructured, error) {
	addHubLabels(new, catalog)
	// replace label, annotation and spec of old resource with new
	existing.SetLabels(new.GetLabels())
	existing.SetAnnotations(new.GetAnnotations())
	existing.Object["spec"] = new.Object["spec"]

	res, err := i.update(existing, namespace, metav1.UpdateOptions{})
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

// Checks if version label is available
func isVersionLabelAvailable(res *unstructured.Unstructured) bool {

	labels := res.GetLabels()
	if len(labels) == 0 {
		return false
	}

	_, ok := labels[versionLabel]
	return ok
}
