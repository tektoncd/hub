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

package model

import (
	"github.com/jinzhu/gorm"
)

// JobState defines the state of Sync Job
type JobState int

// Represents Job State
const (
	JobQueued JobState = iota
	JobRunning
	JobDone
	JobError
)

func (s JobState) String() string {
	return [...]string{"queued", "running", "done", "error"}[s]
}

type SyncJob struct {
	gorm.Model
	Catalog   Catalog
	CatalogID uint
	Status    string
	UserID    uint
	User      User
}

func (j *SyncJob) SetState(s JobState) {
	j.Status = s.String()
}

func (j *SyncJob) IsRunning() bool {
	return j.Status == JobRunning.String()
}
