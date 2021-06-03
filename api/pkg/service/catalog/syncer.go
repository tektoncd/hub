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
	"fmt"
	"time"

	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/git"
	"github.com/tektoncd/hub/api/pkg/parser"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var clonePath = "/tmp/catalog"

type syncer struct {
	db      *gorm.DB
	logger  *zap.SugaredLogger
	running bool
	limit   chan bool
	stop    chan bool
	git     git.Client
}

var (
	queued  = &model.SyncJob{Status: model.JobQueued.String()}
	running = &model.SyncJob{Status: model.JobRunning.String()}
)

func newSyncer(api app.BaseConfig) *syncer {
	logger := api.Logger("syncer")
	return &syncer{
		db:     app.DBWithLogger(api.Environment(), api.DB(), logger),
		logger: logger.SugaredLogger,
		limit:  make(chan bool, 1),
		stop:   make(chan bool),
		git:    git.New(api.Logger("git").SugaredLogger),
	}
}

func (s *syncer) Enqueue(userID, catalogID uint) (*model.SyncJob, error) {

	s.logger.Infof("Enqueueing User: %d catalogID %d", userID, catalogID)
	queued := &model.SyncJob{CatalogID: catalogID, Status: "queued"}
	running := &model.SyncJob{CatalogID: catalogID, Status: "running"}
	newJob := model.SyncJob{CatalogID: catalogID, Status: "queued", UserID: userID}

	if err := s.db.Where(queued).Or(running).FirstOrCreate(&newJob).Error; err != nil {
		return nil, internalError
	}
	s.wakeUp()

	return &newJob, nil
}

func (s *syncer) wakeUp() {
	// start by freeing up limit if it occupied already
	select {
	case <-s.limit:
	default:
	}
}

func (s *syncer) Run() {

	if s.running {
		return
	}

	log := s.logger.With("action", "run")
	log.Info("running catalog syncer ....")

	// all running jobs should be queued so that they can be retried
	if err := s.db.Model(model.SyncJob{}).Where(running).Updates(queued).Error; ignoreNotFound(err) != nil {
		log.Error(err, "failed to update running -> queued")
	}

	go func() {
		defer log.Info("exiting job runner")
		for {
			select {
			case <-s.stop:
				return
			case s.limit <- true:
				log.Info("processing the queue")
				if err := s.Process(); err != nil {
					time.AfterFunc(30*time.Second, s.Next)
					return
				}
				s.Next()
			}
		}
	}()

	s.wakeUp()
	s.running = true
}

func (s *syncer) Stop() {
	close(s.stop)
	close(s.limit)
	s.running = false
}

func (s *syncer) Next() {
	log := s.logger.With("action", "next")

	var count int64
	if err := s.db.Model(&model.SyncJob{}).Where(queued).Count(&count).Error; err != nil {
		log.Error(err)
		return
	}

	log.Info("queued job count: ", count)
	if count == 0 {
		return
	}
	s.wakeUp()
}

func ignoreNotFound(err error) error {
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func (s *syncer) Process() error {
	log := s.logger.With("action", "process")
	db := s.db

	syncJob := model.SyncJob{}

	// helper to update job state
	setJobState := func(s model.JobState) {
		syncJob.SetState(s)
		db.Model(&syncJob).Update("status", syncJob.Status)
	}

	if err := db.Model(&model.SyncJob{}).Where(queued).Order("created_at").First(&syncJob).Error; err != nil {
		if ignoreNotFound(err) != nil {
			log.Error(err)
			return err
		}
		log.Info("nothing to sync")
		return nil
	}

	setJobState(model.JobRunning)

	catalog := model.Catalog{}
	if err := db.Model(&syncJob).Association("Catalog").Find(&catalog); err != nil {
		log.Error(err)
		return err
	}

	fetchSpec := git.FetchSpec{URL: catalog.URL, Revision: catalog.Revision, Path: clonePath}
	repo, err := s.git.Fetch(fetchSpec)
	if err != nil {
		log.Error(err, "clone failed")
		setJobState(model.JobError)
		return nil
	}

	if repo.Head() == catalog.SHA {
		log.Infof("skipping already cloned catalog - %s | sha: %s", catalog.URL, catalog.SHA)
		setJobState(model.JobDone)
		return nil
	}

	// parse the catalog and fill the db
	parser := parser.ForCatalog(s.logger, repo, catalog.ContextDir)

	res, result := parser.Parse()
	if err = s.updateJob(syncJob, repo.Head(), res, result); err != nil {
		log.Error(err, "updation of db failed")
		setJobState(model.JobQueued)
		return err
	}
	setJobState(model.JobDone)
	return nil
}

func (s *syncer) updateJob(syncJob model.SyncJob, sha string, res []parser.Resource, result parser.Result) error {
	log := s.logger.With("action", "update-job", "job-id", syncJob.ID)

	txn := s.db.Begin()

	catalog := model.Catalog{}
	if err := txn.Model(&syncJob).Association("Catalog").Find(&catalog); err != nil {
		return err
	}
	catalog.SHA = sha

	if err := s.updateResources(txn, log, &catalog, res); err != nil {
		txn.Rollback()
		return err
	}

	if err := s.updateCatalogResult(txn, log, &catalog, result); err != nil {
		txn.Rollback()
		return err
	}

	if err := txn.Save(catalog).Error; err != nil {
		txn.Rollback()
		return err
	}

	txn.Commit()
	return nil
}

func (s *syncer) updateResources(
	txn *gorm.DB, log *zap.SugaredLogger,
	catalog *model.Catalog, res []parser.Resource) error {

	if len(res) == 0 {
		return nil
	}

	for _, r := range res {

		log.Infof("Res: %s | Name: %s ", r.Kind, r.Name)
		if len(r.Versions) == 0 {
			log.Warnf("Res: %s | Name: %s has no versions - skipping ", r.Kind, r.Name)
			continue
		}

		dbRes := model.Resource{
			Name:      r.Name,
			Kind:      r.Kind,
			CatalogID: catalog.ID,
		}

		txn.Model(&model.Resource{}).Where(dbRes).FirstOrCreate(&dbRes)
		if err := txn.Save(&dbRes).Error; err != nil {
			return err
		}

		log.Infof("Resource: %s  ID: %d stored", r.Name, dbRes.ID)

		s.updateResourceCategory(txn, log, &dbRes, r.Categories)
		s.updateResourceTags(txn, log, &dbRes, r.Tags)
		s.updateResourceVersions(txn, log, catalog, dbRes.ID, r.Versions)

	}
	return nil
}

func (s *syncer) updateResourceTags(
	txn *gorm.DB, log *zap.SugaredLogger,
	res *model.Resource, tags []string) {

	if len(tags) == 0 {
		return
	}

	for _, t := range tags {

		tag := model.Tag{Name: t}

		txn.Model(&model.Tag{}).Where(&model.Tag{Name: t}).FirstOrCreate(&tag)

		resTag := model.ResourceTag{ResourceID: res.ID, TagID: tag.ID}
		txn.Model(&model.ResourceTag{}).Where(resTag).FirstOrCreate(&resTag)
		log.Infof("Resource: %d: %s | tag: %s (%d)", res.ID, res.Name, tag.Name, tag.ID)
	}
}

func (s *syncer) updateResourceCategory(txn *gorm.DB, log *zap.SugaredLogger,
	res *model.Resource, categories []string) {
	if len(categories) == 0 {
		return
	}

	for _, t := range categories {

		c := model.Category{}
		txn.Model(&model.Category{}).Where(&model.Category{Name: t}).FirstOrCreate(&c)

		resCategory := model.ResourceCategory{ResourceID: res.ID, CategoryID: c.ID}
		txn.Model(&model.ResourceCategory{}).Where(resCategory).FirstOrCreate(&resCategory)

	}
}

func (s *syncer) updateResourceVersions(
	txn *gorm.DB, log *zap.SugaredLogger,
	catalog *model.Catalog,
	resourceID uint,
	versions []parser.VersionInfo) {

	for _, v := range versions {
		ver := model.ResourceVersion{
			Version:    v.Version,
			ResourceID: resourceID,
		}

		txn.Model(&model.ResourceVersion{}).
			Where(&model.ResourceVersion{ResourceID: resourceID, Version: v.Version}).FirstOrInit(&ver)

		ver.DisplayName = v.DisplayName
		ver.Description = v.Description
		ver.ModifiedAt = v.ModifiedAt
		ver.MinPipelinesVersion = v.MinPipelinesVersion
		ver.URL = fmt.Sprintf("%s/tree/%s/%s", catalog.URL, catalog.Revision, v.Path)

		txn.Save(&ver)
		log.Infof(" Version: %d -> %s | Path: %s ", ver.ID, ver.Version, v.Path)
	}
}

func (s *syncer) updateCatalogResult(
	txn *gorm.DB, log *zap.SugaredLogger,
	catalog *model.Catalog, result parser.Result) error {

	// delete all old records
	txn.Unscoped().Where(&model.CatalogError{CatalogID: catalog.ID}).Delete(&model.CatalogError{})
	for _, err := range result.Errors {
		if err := txn.Create(&model.CatalogError{
			CatalogID: catalog.ID, Type: "error",
			Detail: err.Error(),
		}).Error; err != nil {
			return err
		}
	}

	for _, issue := range result.Issues {
		if err := txn.Create(&model.CatalogError{
			CatalogID: catalog.ID,
			Type:      issue.Type.String(),
			Detail:    issue.Message,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}
