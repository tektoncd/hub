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
	"strings"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/tektoncd/hub/api/gen/catalog"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/git"
	gitclient "github.com/tektoncd/hub/api/pkg/git"
)

type service struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

var clonePath = "/tmp/catalog"

var (
	internalError = catalog.MakeInternalError(fmt.Errorf("failed to refresh catalog"))
)

// New returns the catalog service implementation.
func New(api app.Config) catalog.Service {
	return &service{api.Logger(), api.DB()}
}

// refresh the catalog for new resources
func (s *service) Refresh(ctx context.Context) (err error) {

	catalog := model.Catalog{}
	if err := s.db.Model(&model.Catalog{}).First(&catalog).Error; err != nil {
		return nil
	}

	job := model.SyncJob{Catalog: catalog, Status: "queued"}
	if err := s.db.Create(&job).Error; err != nil {
		return nil
	}

	if err := s.SyncCatalog(); err != nil {
		return nil
	}

	return nil
}

func (s *service) SyncCatalog() error {
	log := s.logger.With("action", "sync")
	db := s.db

	count := 0
	db.Model(&model.SyncJob{}).Count(&count)
	if count == 0 {
		log.Infof("skipping sync job count: %d", count)
		return nil
	}

	log.Info("job count: ", count)
	job := model.SyncJob{}

	// helper to update job state
	setJobState := func(s model.JobState) {
		job.SetState(s)
		db.Model(&job).Updates(job)
	}

	if err := db.Where("status = ?", model.Queued.String()).First(&job).Error; err != nil {
		return err
	}

	job.SetState(model.Running)
	db.Model(&job).Updates(job)
	// NOTE: only delete done jobs
	defer db.Unscoped().Where(&model.SyncJob{Status: model.Done.String()}).Delete(&job)

	catalog := model.Catalog{}
	db.Model(job).Related(&catalog)

	fetchSpec := gitclient.FetchSpec{
		URL:      catalog.URL,
		Revision: catalog.Revision,
		Path:     clonePath,
	}

	git := gitclient.New(s.logger)

	repo, err := git.Fetch(fetchSpec)
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

	res, err := repo.ParseTektonResources()
	if err != nil {
		if len(res) == 0 {
			log.Error(err, "parsing of resources failed")
			setJobState(model.Queued)
			return err

		}
		// Partial parsing of resources is allowed
		log.Warnf("Failed to parse some for the resources: %s found: %d ", err, len(res))
	}

	if err := s.updateResources(job, repo, res); err != nil {
		// TODO(sthaha): handle updation failure better
		log.Error(err, "updation of db failed")
		setJobState(model.Queued)
		return err
	}

	setJobState(model.Done)
	return nil
}

func (s *service) updateResources(job model.SyncJob, repo git.Repo, res []gitclient.TektonResource) error {
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

		for _, v := range r.Versions {
			ver := &model.ResourceVersion{
				Version:    v.Version,
				ResourceID: dbRes.ID,
			}

			txn.Model(&model.ResourceVersion{}).Where(&model.ResourceVersion{ResourceID: dbRes.ID, Version: v.Version}).FirstOrCreate(&ver)

			ver.DisplayName = v.DisplayName
			ver.Description = v.Description
			ver.ModifiedAt = v.ModifiedAt
			// TODO(sthaha): use gh client to get the path?
			// this heuristic works fine so far
			ver.URL = fmt.Sprintf("%s/tree/%s/%s", catalog.URL, catalog.Revision, v.Path)

			txn.Save(&ver)
			s.logger.Infof("      >>> Version: %d -> %s | Path: %s | tags: %s", ver.ID, ver.Version, v.Path, strings.Join(v.Tags, ", "))

			for _, t := range v.Tags {
				tag := model.Tag{Name: t, CategoryID: others.ID}

				txn.Model(&model.Tag{}).Where(&model.Tag{Name: t}).FirstOrCreate(&tag)

				resTag := model.ResourceTag{ResourceID: dbRes.ID, TagID: tag.ID}
				txn.Model(&model.ResourceTag{}).Where(&resTag).FirstOrCreate(&resTag)
				s.logger.Infof("      >>> Resource: %d: %s | tag: %s (%d)", dbRes.ID, dbRes.Name, tag.Name, tag.ID)
			}
		}

	}

	txn.Save(&catalog)
	txn.Commit()
	return nil
}
