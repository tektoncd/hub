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

package parser

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	decoder "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

type TektonParser struct {
	reader io.Reader
}

func (t *TektonParser) Parse() (*TektonResource, error) {

	// create a duplicate for  UniversalDeserializer and NewYAMLToJSONDecoder
	// to read from readers
	var dup bytes.Buffer
	r := io.TeeReader(t.reader, &dup)
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	object, gvk, err := scheme.Codecs.UniversalDeserializer().Decode(contents, nil, nil)
	if err != nil || !isTektonKind(gvk) {
		return nil, fmt.Errorf("parse error: invalid resource %+v:\n%s", err, contents)
	}

	decoder := decoder.NewYAMLToJSONDecoder(&dup)

	var res *unstructured.Unstructured
	for {
		res = &unstructured.Unstructured{}
		if err := decoder.Decode(res); err != nil {
			return nil, fmt.Errorf("failed to decode: %w", err)
		}

		if len(res.Object) == 0 {
			continue
		}

		if res.GetKind() != "" {
			break
		}
	}

	if _, err := convertToTyped(res); err != nil {
		return nil, err
	}

	return &TektonResource{
		Unstructured: res,
		Object:       object,
		Name:         res.GetName(),
		Kind:         gvk.Kind,
	}, nil
}
