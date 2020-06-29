package hub

import (
	"context"
	"go.uber.org/zap"

	category "github.com/tektoncd/hub/api/gen/category"
)

// category service example implementation.
// The example methods log the requests and return zero values.
type categorysrvc struct {
	logger *zap.SugaredLogger
}

// NewCategory returns the category service implementation.
func NewCategory(logger *zap.SugaredLogger) category.Service {
	return &categorysrvc{logger}
}

// Get all Categories with their tags sorted by name
func (s *categorysrvc) All(ctx context.Context) (res []*category.Category, err error) {
	s.logger.Info("category.All")
	return
}
