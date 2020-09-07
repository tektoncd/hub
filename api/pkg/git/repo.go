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
	"path/filepath"
	"strings"
	"time"
)

type Repo struct {
	Path        string
	ContextPath string
	head        string
}

func (r Repo) Head() string {
	if r.head == "" {
		head, _ := rawGit("", "rev-parse", "HEAD")
		r.head = strings.TrimSuffix(head, "\n")
	}
	return r.head
}

func (r Repo) ModifiedTime(path string) (time.Time, error) {
	// gitPath should be relative to the repo and not the context
	gitPath, _ := filepath.Rel(r.Path, path)
	commitedAt, err := rawGit(r.Path, "log", "-1", "--pretty=format:%cI", gitPath)
	if err != nil {
		return time.Time{}, err
	}

	return time.Parse(time.RFC3339, commitedAt)
}

func (r Repo) RelPath(path string) (string, error) {
	basePath := filepath.Join(r.Path, r.ContextPath)
	return filepath.Rel(basePath, path)
}
