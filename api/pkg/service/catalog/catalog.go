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

package catalog

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/tektoncd/hub/api/gen/catalog"
	"github.com/tektoncd/hub/api/gen/log"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/git"
	"github.com/tektoncd/hub/api/pkg/parser"
)

type service struct {
	logger *log.Logger
	db     *gorm.DB
	wq     *catalogSyncer
}

var clonePath = "/tmp/catalog"

var (
	internalError = catalog.MakeInternalError(fmt.Errorf("failed to refresh catalog"))
	fetchError    = catalog.MakeInternalError(fmt.Errorf("failed to fetch catalog"))
	notFoundError = catalog.MakeNotFound(fmt.Errorf("resource not found"))
)

type catalogSyncer struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
	queued chan bool
	stop   chan bool
}

func newCatalogSyncer(api app.BaseConfig) *catalogSyncer {
	l := api.Logger("catalog-syncer")
	l.Info("create catalog syncer")
	return &catalogSyncer{
		db:     api.DB(),
		logger: api.Logger("catalog-syncer").SugaredLogger,
		queued: make(chan bool, 1),
		stop:   make(chan bool),
	}
}

func (cs *catalogSyncer) Enqueue() (*model.SyncJob, error) {
	catalog := model.Catalog{}
	if err := cs.db.Model(&model.Catalog{}).First(&catalog).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, notFoundError
		}
	}

	queued := &model.SyncJob{CatalogID: catalog.ID, Status: "queued"}
	running := &model.SyncJob{CatalogID: catalog.ID, Status: "running"}
	if err := cs.db.Where(queued).Or(running).FirstOrCreate(queued).Error; err != nil {
		return nil, internalError
	}
	cs.wakeUp()

	return queued, nil
}

func (cs *catalogSyncer) wakeUp() {
	// start a process if not already
	select {
	case cs.queued <- true:
		cs.logger.Info("signaled - start processing queue")
	default:
	}
}

func (cs *catalogSyncer) Run() {

	log := cs.logger.With("action", "run")
	log.Info("running catalog syncer ....")

	if err := cs.db.Model(model.SyncJob{}).
		Where(model.SyncJob{Status: "running"}).
		Update(model.SyncJob{Status: "queued"}).Error; ignoreNotFound(err) != nil {
		log.Error(err, "failed to update running -> queued")
	}

	go func() {
		defer log.Info("exiting job runner")
		for {
			select {
			case <-cs.stop:
				return
			case <-cs.queued:
				go func() {
					log.Info("processing the queue")
					cs.Process()
					cs.Next()
				}()
			}
		}
	}()

	cs.wakeUp()
}

func (cs *catalogSyncer) Stop() {
	close(cs.stop)
	close(cs.queued)
}

func (cs *catalogSyncer) Next() {
	log := cs.logger.With("action", "next")

	count := 0
	if err := cs.db.Model(&model.SyncJob{}).Where("status = ? ", "queued").Count(&count).Error; err != nil {
		log.Error(err)
		return
	}

	log.Info("queued job count: ", count)
	if count == 0 {
		return
	}
	cs.wakeUp()
}

func ignoreNotFound(err error) error {
	if gorm.IsRecordNotFoundError(err) {
		return nil
	}
	return err
}

func (cs *catalogSyncer) Process() error {
	log := cs.logger.With("action", "process")
	db := cs.db

	job := model.SyncJob{}

	// helper to update job state
	setJobState := func(s model.JobState) {
		job.SetState(s)
		db.Model(&job).Updates(job)
	}

	if err := db.
		Where("status = ?", model.Queued.String()).
		Order("created_at").
		First(&job).Error; err != nil {
		return ignoreNotFound(err)
	}

	job.SetState(model.Running)
	db.Model(&job).Updates(job)

	catalog := model.Catalog{}
	db.Model(job).Related(&catalog)

	fetchSpec := git.FetchSpec{
		URL:      catalog.URL,
		Revision: catalog.Revision,
		Path:     clonePath,
	}

	gitclient := git.New(cs.logger)

	repo, err := gitclient.Fetch(fetchSpec)
	if err != nil {
		log.Error(err, "clone failed")
		setJobState(model.Queued)
		return err
	}

	if repo.Head() == catalog.SHA {
		log.Infof("skipping already cloned catalog - %s | sha: %s", catalog.URL, catalog.SHA)
		setJobState(model.Done)
		return nil
	}
	// parse the catalog and fill the db

	parser := parser.ForCatalog(cs.logger, repo)

	res, result := parser.Parse()

	if len(res) == 0 {
		log.Warnf("parsing of resources failed err: %s", result.Error())
		// TODO(sthaha):  log errors against catalog
		// // TODO(sthaha): decide to requeue to retry N times
		setJobState(model.Queued)
		return result.Errors[0]

	}
	// Partial parsing of resources is allowed
	log.Warnf("Failed to parse some for the resources: %s found: %d ", err, len(res))

	if err := cs.updateResources(job, repo, res); err != nil {
		// TODO(sthaha): handle updation failure better
		log.Error(err, "updation of db failed")
		setJobState(model.Queued)
		return err
	}
	setJobState(model.Done)
	return nil
}

func (s *catalogSyncer) updateResources(job model.SyncJob, repo git.Repo, res []parser.Resource) error {
	log := s.logger.With("action", "updatedb")

	txn := s.db.Begin()

	catalog := model.Catalog{}
	txn.Model(&job).Related(&catalog)

	catalog.SHA = repo.Head()

	others := model.Category{}
	txn.Model(&model.Category{}).Where(&model.Category{Name: "Others"}).First(&others)

	for _, r := range res {

		s.logger.Infof("Res: %s | Name: %s ", r.Kind, r.Name)
		if len(r.Versions) == 0 {
			s.logger.Infof("      >>> Res: %s | Name: %s has no versions - skipping ", r.Kind, r.Name)
			continue
		}

		dbRes := model.Resource{
			Name:      r.Name,
			Kind:      r.Kind,
			CatalogID: catalog.ID,
		}

		txn.Model(&model.Resource{}).Where(&dbRes).FirstOrCreate(&dbRes)
		txn.Save(&dbRes)

		log.Info("Resource ID: ", dbRes.ID)

		for _, t := range r.Tags {
			tag := model.Tag{Name: t, CategoryID: others.ID}

			txn.Model(&model.Tag{}).Where(&model.Tag{Name: t}).FirstOrCreate(&tag)

			resTag := model.ResourceTag{ResourceID: dbRes.ID, TagID: tag.ID}
			txn.Model(&model.ResourceTag{}).Where(&resTag).FirstOrCreate(&resTag)
			s.logger.Infof("      >>> Resource: %d: %s | tag: %s (%d)", dbRes.ID, dbRes.Name, tag.Name, tag.ID)
		}

		for _, v := range r.Versions {
			ver := &model.ResourceVersion{
				Version:    v.Version,
				ResourceID: dbRes.ID,
				URL:        fmt.Sprintf("%s/tree/%s/%s", catalog.URL, catalog.Revision, v.Path),
			}

			txn.Model(&model.ResourceVersion{}).
				Where(&model.ResourceVersion{ResourceID: dbRes.ID, Version: v.Version}).FirstOrInit(&ver)

			ver.DisplayName = v.DisplayName
			ver.Description = v.Description
			ver.ModifiedAt = v.ModifiedAt
			ver.MinPipelinesVersion = v.MinPipelinesVersion

			txn.Save(&ver)
			s.logger.Infof("      >>> Version: %d -> %s | Path: %s ", ver.ID, ver.Version, v.Path)
		}

	}

	txn.Save(&catalog)
	txn.Commit()
	return nil
}

// New returns the catalog service implementation.
func New(api app.Config) catalog.Service {
	wq := newCatalogSyncer(api)

	// start running after some delay to allow for all services to mount
	time.AfterFunc(3*time.Second, wq.Run)

	s := &service{
		api.Logger("catalog"),
		api.DB(),
		wq,
	}
	return s
}

// refresh the catalog for new resources
func (s *service) Refresh(ctx context.Context, p *catalog.RefreshPayload) (*catalog.Job, error) {
	s.logger.Infof("going to enqueue")

	job, err := s.wq.Enqueue()
	if err != nil {
		return nil, err
	}

	res := &catalog.Job{ID: job.ID, Status: "queued"}
	s.logger.Infof("job %d queued for refresh", job.ID)

	return res, nil
}
