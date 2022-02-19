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
	"strings"
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

	fetchSpec := git.FetchSpec{URL: catalog.URL, Revision: catalog.Revision, Path: clonePath, SSHUrl: catalog.SSHURL}
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

	rerr, err := s.updateResources(txn, log, &catalog, res)
	if err != nil {
		txn.Rollback()
		return err
	}

	result.Combine(rerr)

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
	catalog *model.Catalog, res []parser.Resource) (parser.Result, error) {

	if len(res) == 0 {
		return parser.Result{}, nil
	}

	var syncResourceID []uint
	rerr := parser.Result{}
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
			return parser.Result{}, err
		}
		syncResourceID = append(syncResourceID, dbRes.ID)

		log.Infof("Resource: %s  ID: %d stored", r.Name, dbRes.ID)

		rerr.Combine(s.updateResourceCategory(txn, log, &dbRes, r.Categories))
		s.updateResourceTags(txn, log, &dbRes, r.Tags)
		// platform ids on resource version level
		verPlatformIds := map[uint]bool{}
		s.updateResourceVersions(txn, log, catalog, dbRes.ID, r.Versions, &verPlatformIds)
		s.updateResourcePlatforms(txn, log, &dbRes, verPlatformIds)
	}

	var dltRes []model.Resource
	// Finds db resources which are deleted
	if err := txn.Model(&model.Resource{}).Where(&model.Resource{CatalogID: catalog.ID}).Not(map[string]interface{}{"id": syncResourceID}).Find(&dltRes).Error; err != nil {
		log.Error(err)
		return rerr, err
	}

	for _, r := range dltRes {
		var tags []model.ResourceTag
		var platforms []model.ResourcePlatform
		txn.Where(&model.ResourceTag{ResourceID: r.ID}).Find(&tags)
		if err := txn.Where(&model.ResourcePlatform{ResourceID: r.ID}).Find(&platforms).Error; err != nil {
			log.Error(err)
			return rerr, err
		}

		s.deleteResource(txn, log, r)
		s.deleteTag(txn, log, tags)
		s.deletePlatform(txn, log, platforms)
	}

	return rerr, nil
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

func (s *syncer) updateResourcePlatforms(
	txn *gorm.DB, log *zap.SugaredLogger,
	res *model.Resource, verPlatformIds map[uint]bool) {
	platformIds := []uint{}
	for verPlatformId := range verPlatformIds {
		platformIds = append(platformIds, verPlatformId)
		platform := model.Platform{}
		txn.First(&platform, verPlatformId)
		resPlatform := model.ResourcePlatform{ResourceID: res.ID, PlatformID: verPlatformId}
		txn.Model(&model.ResourcePlatform{}).Where(resPlatform).FirstOrCreate(&resPlatform)
		log.Infof("Resource: %d: %s | platform: %s (%d)", res.ID, res.Name, platform.Name, platform.ID)
	}
	// remove mapping for the platforms, which are not in the list
	txn.Unscoped().Where(&model.ResourcePlatform{ResourceID: res.ID}).
		Not(map[string]interface{}{"platform_id": platformIds}).
		Delete(&model.ResourcePlatform{})
}

func (s *syncer) updateResourceCategory(txn *gorm.DB, log *zap.SugaredLogger,
	res *model.Resource, categories []string) parser.Result {
	if len(categories) == 0 {
		return parser.Result{}
	}
	rerr := parser.Result{}
	for _, cat := range categories {
		c := model.Category{Name: cat}

		err := txn.Model(&model.Category{}).Where("LOWER(name) = ?", strings.ToLower(cat)).First(&c).Error
		if err != nil && err == gorm.ErrRecordNotFound {
			rerr.AddError(fmt.Errorf("Category `%s` for Resource `%s` is Invalid", strings.ReplaceAll(c.Name, " ", ""), res.Name))
			continue
		}
		resCategory := model.ResourceCategory{ResourceID: res.ID, CategoryID: c.ID}
		txn.Model(&model.ResourceCategory{}).Where(resCategory).FirstOrCreate(&resCategory)
		log.Infof("Resource: %d: %s | category: %s (%d)", res.ID, res.Name, c.Name, c.ID)
	}

	return rerr
}

func (s *syncer) updateResourceVersions(
	txn *gorm.DB, log *zap.SugaredLogger,
	catalog *model.Catalog,
	resourceID uint,
	versions []parser.VersionInfo,
	verPlatformIds *map[uint]bool) {

	for _, v := range versions {
		ver := model.ResourceVersion{
			Version:    v.Version,
			ResourceID: resourceID,
		}

		txn.Model(&model.ResourceVersion{}).
			Where(&model.ResourceVersion{ResourceID: resourceID, Version: v.Version}).FirstOrInit(&ver)

		ver.DisplayName = v.DisplayName
		ver.Deprecated = v.Deprecated
		ver.Description = v.Description
		ver.ModifiedAt = v.ModifiedAt
		ver.MinPipelinesVersion = v.MinPipelinesVersion
		switch catalog.Provider {
		case "github":
			ver.URL = fmt.Sprintf("%s/tree/%s/%s", catalog.URL, catalog.Revision, v.Path)
		case "bitbucket":
			ver.URL = fmt.Sprintf("%s/src/%s/%s", catalog.URL, catalog.Revision, v.Path)
		case "gitlab":
			ver.URL = fmt.Sprintf("%s/-/blob/%s/%s", catalog.URL, catalog.Revision, v.Path)
		}

		txn.Save(&ver)
		log.Infof(" Version: %d -> %s | Path: %s ", ver.ID, ver.Version, v.Path)

		platforms := v.Platforms
		if len(platforms) == 0 {
			return
		}
		platformIds := []uint{}
		for _, p := range platforms {

			platform := model.Platform{Name: p}
			txn.Model(&model.Platform{}).Where(&model.Platform{Name: p}).FirstOrCreate(&platform)
			platformIds = append(platformIds, platform.ID)
			verPlatform := model.VersionPlatform{ResourceVersionID: ver.ID, PlatformID: platform.ID}
			txn.Model(&model.VersionPlatform{}).Where(verPlatform).FirstOrCreate(&verPlatform)
			//save unique platform ids
			(*verPlatformIds)[platform.ID] = true
			log.Infof("Resource version: %d -> %s | platform: %s (%d)", ver.ID, ver.Version, platform.Name, platform.ID)
		}
		// remove mapping for the platforms, which are not in the list
		txn.Unscoped().Where(&model.VersionPlatform{ResourceVersionID: ver.ID}).
			Not(map[string]interface{}{"platform_id": platformIds}).
			Delete(&model.VersionPlatform{})
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

// Delete the resource by id
func (s *syncer) deleteResource(txn *gorm.DB, log *zap.SugaredLogger, res model.Resource) {

	txn.Unscoped().Where("id = ?", res.ID).Delete(&model.Resource{})

	log.Infof("Resource %s of kind %s has been deleted", res.Name, res.Kind)
}

// Delete the tags which doesn't belongs to any other existing resources
func (s *syncer) deleteTag(txn *gorm.DB, log *zap.SugaredLogger, tags []model.ResourceTag) {

	for _, t := range tags {
		var resTag []model.ResourceTag
		txn.Where(&model.ResourceTag{TagID: t.TagID}).Find(&resTag)

		if len(resTag) == 0 {
			txn.Unscoped().Where("id = ?", t.TagID).Delete(&model.Tag{})
			log.Infof("Tag with ID: %d has been deleted", t.TagID)
		}
	}
}

// Delete the platforms which doesn't belongs to any other existing resources
func (s *syncer) deletePlatform(txn *gorm.DB, log *zap.SugaredLogger, platforms []model.ResourcePlatform) {

	for _, t := range platforms {
		var resPlatform []model.ResourcePlatform
		txn.Model(&model.ResourcePlatform{}).Where("platform_id= ?", t.PlatformID).Find(&resPlatform)

		if len(resPlatform) == 0 {
			txn.Unscoped().Where("id = ?", t.PlatformID).Delete(&model.Platform{})
			log.Infof("Platform with ID: %d has been deleted", t.PlatformID)
		}
	}
}
