// Copyright © 2022 The Tekton Authors.
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

package resource

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/gen/resource"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

func TestVersionsByID(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.VersionsByIDPayload{ID: 1}
	res, err := resourceSvc.VersionsByID(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(res.Data.Versions))
	assert.Equal(t, "0.2", res.Data.Latest.Version)
}

func TestVersionsByID_NotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.VersionsByIDPayload{ID: 11}
	_, err := resourceSvc.VersionsByID(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestCreationRawURLBitbucket(t *testing.T) {
	url := "https://bitbucket.org/org/catalog/src/main/task/name/0.1/name.yaml"
	replacer := getStringReplacer(url, "bitbucket")
	rawUrl := replacer.Replace(url)
	expected := "https://bitbucket.org/org/catalog/raw/main/task/name/0.1/name.yaml"
	assert.Equal(t, expected, rawUrl)
}

func TestCreationRawURLGitlab(t *testing.T) {
	url := "https://gitlab.com/org/catalog/-/blob/main/task/name/0.1/name.yaml"
	replacer := getStringReplacer(url, "gitlab")
	rawUrl := replacer.Replace(url)
	expected := "https://gitlab.com/org/catalog/-/raw/main/task/name/0.1/name.yaml"
	assert.Equal(t, expected, rawUrl)
}

func TestCreationRawURLGitlabEnterprise(t *testing.T) {
	url := "https://gitlab.myhost.com/org/catalog/-/blob/main/task/name/0.1/name.yaml"
	replacer := getStringReplacer(url, "gitlab")
	rawUrl := replacer.Replace(url)
	expected := "https://gitlab.myhost.com/org/catalog/-/raw/main/task/name/0.1/name.yaml"
	assert.Equal(t, expected, rawUrl)
}
