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

package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Data struct {
	Catalogs   []Catalog
	Categories []Category
}

type Category struct {
	Name string
	Tags []string
}

type Catalog struct {
	Name       string
	Type       string
	URL        string
	Owner      string
	ContextDir string
	Revision   string
}

// dataFromURL reads data from file using URL or path
func dataFromURL(url string) ([]byte, error) {

	if strings.HasPrefix(url, "file://") {
		return readLocalFile(strings.TrimPrefix(url, "file://"))
	}
	return httpRead(url)
}

// httpRead reads data from a remote file using URL
func httpRead(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

// readLocalFile reads data from a local file using file path
func readLocalFile(path string) ([]byte, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
