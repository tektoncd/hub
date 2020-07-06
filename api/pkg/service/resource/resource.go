	package resource

import (
	"context"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/tektoncd/hub/api/gen/resource"
	"github.com/tektoncd/hub/api/pkg/app"
)

type service struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

// New returns the resource service implementation.
func New(api app.Config) resource.Service {
	return &service{api.Logger(), api.DB()}
}

// Find resources based on name, type or both
func (s *service) Query(ctx context.Context, p *resource.QueryPayload) (res resource.ResourceCollection, err error) {
	s.logger.Info("Query")
	return
}
