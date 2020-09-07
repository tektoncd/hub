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

package git

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"go.uber.org/zap"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type Repo struct {
	Path        string
	ContextPath string
	head        string
	Log         *zap.SugaredLogger
}

func (r Repo) Head() string {
	if r.head == "" {
		head, _ := rawGit("", "rev-parse", "HEAD")
		r.head = strings.TrimSuffix(head, "\n")
	}
	return r.head
}

type (
	TektonResource struct {
		Name     string
		Kind     string
		Versions []TekonResourceVersion
	}

	TekonResourceVersion struct {
		Version     string
		DisplayName string
		Path        string
		Description string
		ModifiedAt  time.Time
		Tags        []string
	}
)

func (r Repo) ParseTektonResources() ([]TektonResource, error) {
	// TODO(sthaha): may be in parallel
	// TODO(sthaha): replace it by channels and stream and write?
	// TODO(sthaha): get task kind from scheme ?
	kinds := []string{"Task", "Pipeline"}
	resources := []TektonResource{}

	var parseError error
	for _, k := range kinds {
		res, err := r.findResourcesByKind(k)
		if err != nil {
			parseError = err
			continue
		}
		resources = append(resources, res...)
	}
	return resources, parseError
}

func ignoreNotExists(err error) error {
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (r Repo) findResourcesByKind(kind string) ([]TektonResource, error) {
	log := r.Log.With("kind", kind)
	log.Info("looking for resources")

	kindPath := filepath.Join(r.Path, r.ContextPath, strings.ToLower(kind))
	resources, err := ioutil.ReadDir(kindPath)
	if err != nil {
		r.Log.Errorf("failed to find %s: %s", kind, err)
		// NOTE: returns empty task list; upto caller to check for error
		return []TektonResource{}, ignoreNotExists(err)
	}

	found := []TektonResource{}
	var parseError error
	for _, res := range resources {
		if !res.IsDir() {
			log.Warnf("ignoring %s  not a directory for %s", res.Name(), kind)
			continue
		}

		res, err := r.parseResource(kind, kindPath, res)
		if err != nil {
			// TODO(sthaha): do something about invalid tasks
			parseError = err
			r.Log.Error(err)
			continue
		}

		// NOTE: res can be nil if no files exists for resource (invalid resource)
		if res != nil {
			found = append(found, *res)
		}
	}

	r.Log.Infof("found %d resources of kind %s", len(found), kind)
	return found, parseError

}

var errInvalidResourceDir = errors.New("invalid resource dir")

func (r Repo) parseResource(kind, kindPath string, res os.FileInfo) (*TektonResource, error) {
	// TODO(sthaha): move this to a different package that can scan a Repo
	name := res.Name()
	r.Log.Info("checking path", kindPath, " resource: ", name)
	// path/<task>/<version>[>
	pattern := filepath.Join(kindPath, name, "*", name+".yaml")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		r.Log.Error(err, "failed to find tasks")
		return nil, errInvalidResourceDir
	}

	versions := []TekonResourceVersion{}
	var parseError error
	for _, m := range matches {
		r.Log.Info(" found file: ", m)

		version, err := r.parseResourceVersion(m, kind)
		if err != nil {
			parseError = err
			r.Log.Error(err)
			continue
		}
		versions = append(versions, *version)
	}

	if len(versions) == 0 {
		return nil, parseError
	}

	r.Log.Infof("found %d versions of resource %s/%s", len(versions), kind, name)
	ret := &TektonResource{
		Name:     res.Name(),
		Kind:     kind,
		Versions: versions,
	}
	return ret, parseError
}

// parseResourceVersion will read the contents of the file at filePath and use the K8s deserializer to attempt to marshal the textjj
// into a Tekton struct. This will fail if the resource is unparseable or not a Tekton resource.
func (r Repo) parseResourceVersion(filePath string, kind string) (*TekonResourceVersion, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	res, err := decodeResource(f, kind)
	if err != nil {
		r.Log.Error(err)
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	log := r.Log.With("kind", kind, "name", res.GetName())
	apiVersion := res.GetAPIVersion()
	log.Info("current kind: ", kind, apiVersion, res.GroupVersionKind())

	log.Info("apiVersion: ", apiVersion)
	log.Info("-----", v1beta1.SchemeGroupVersion.Identifier())
	if apiVersion != v1beta1.SchemeGroupVersion.Identifier() {
		log.Infof("Skipping unknown resource %s name: %s", res.GroupVersionKind(), res.GetName())
		return nil, errors.New("invalid resource" + apiVersion)
	}

	labels := res.GetLabels()
	version, ok := labels["app.kubernetes.io/version"]
	if !ok {
		log.Infof("Resource %s name: %s has no version information", res.GroupVersionKind(), res.GetName())
		return nil, fmt.Errorf("resource has no version info %s/%s", res.GroupVersionKind(), res.GetName())
	}

	annotations := res.GetAnnotations()
	displayName, ok := annotations["tekton.dev/displayName"]
	if !ok {
		log.With("action", "ignore").Infof(
			"Resource %s name: %s has no display name", res.GroupVersionKind(), res.GetName())
	}

	tags := annotations["tekton.dev/tags"]

	// first line
	description, found, err := unstructured.NestedString(res.Object, "spec", "description")
	if !found || err != nil {
		log.Infof("Resource %s name: %s has no description", res.GroupVersionKind(), res.GetName())
		return nil, fmt.Errorf("resource has no description %s/%s", res.GroupVersionKind(), res.GetName())
	}

	basePath := filepath.Join(r.Path, r.ContextPath)
	relPath, _ := filepath.Rel(basePath, filePath)

	modified, err := r.lastModifiedTime(filePath)
	if err != nil {
		log.Errorf("Failed to compute modified time for %s/%s: %s : err: %s", res.GetKind(), res.GetName(), relPath, err)
		return nil, fmt.Errorf("internal error computing modified time for %s/%s", res.GroupVersionKind(), res.GetName())
	}

	ret := &TekonResourceVersion{
		Version:     version,
		DisplayName: displayName,
		Tags:        strings.Split(tags, ","),
		Description: description,
		Path:        relPath,
		ModifiedAt:  modified,
	}
	return ret, nil
}

func (r *Repo) lastModifiedTime(path string) (time.Time, error) {
	gitPath, _ := filepath.Rel(r.Path, path)
	commitedAt, err := rawGit(r.Path, "log", "-1", "--pretty=format:%cI", gitPath)
	if err != nil {
		r.Log.Error(err, "git log failed")
		return time.Time{}, err
	}

	r.Log.Infof("%s last commited at %s", gitPath, commitedAt)
	return time.Parse(time.RFC3339, commitedAt)
}

func ignoreEOF(err error) error {
	if err == io.EOF {
		return nil
	}
	return err
}

// decode consumes the given reader and parses its contents as YAML.
func decodeResource(reader io.Reader, kind string) (*unstructured.Unstructured, error) {
	decoder := yaml.NewYAMLToJSONDecoder(reader)
	var res *unstructured.Unstructured

	for {
		res = &unstructured.Unstructured{}
		if err := decoder.Decode(res); err != nil {
			return nil, ignoreEOF(err)
		}

		if len(res.Object) == 0 {
			continue
		}
		if res.GetKind() == kind {
			break
		}
	}
	return res, nil
}
