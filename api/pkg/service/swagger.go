package hub

import (
	"go.uber.org/zap"

	swagger "github.com/tektoncd/hub/api/gen/swagger"
)

// swagger service example implementation.
// The example methods log the requests and return zero values.
type swaggersrvc struct {
	logger *zap.SugaredLogger
}

// NewSwagger returns the swagger service implementation.
func NewSwagger(logger *zap.SugaredLogger) swagger.Service {
	return &swaggersrvc{logger}
}
