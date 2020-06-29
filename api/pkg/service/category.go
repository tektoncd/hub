package hub

import (
	"context"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	category "github.com/tektoncd/hub/api/gen/category"
	app "github.com/tektoncd/hub/api/pkg/app"
)

// category service example implementation.
// The example methods log the requests and return zero values.
type categorysrvc struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

// NewCategory returns the category service implementation.
func NewCategory(api *app.ApiConfig) category.Service {
	return &categorysrvc{api.Logger(), api.DB()}
}

// Get all Categories with their tags sorted by name
func (s *categorysrvc) All(ctx context.Context) (res []*category.Category, err error) {
	s.logger.Info("category.All")
	return
}
