package catalog

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/tektoncd/hub/api/gen/catalog"
	"github.com/tektoncd/hub/api/pkg/app"
	"go.uber.org/zap"
)

type service struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

// New returns the catalog service implementation.
func New(api app.Config) catalog.Service {
	return &service{api.Logger(), api.DB()}
}

// refresh the catalog for new resources
func (s *service) Refresh(ctx context.Context) (err error) {
	s.logger.Info("catalog.Refresh")
	return
}
