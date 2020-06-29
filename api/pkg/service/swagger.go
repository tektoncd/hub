package hub

import (
	"log"

	swagger "github.com/tektoncd/hub/api/gen/swagger"
)

// swagger service example implementation.
// The example methods log the requests and return zero values.
type swaggersrvc struct {
	logger *log.Logger
}

// NewSwagger returns the swagger service implementation.
func NewSwagger(logger *log.Logger) swagger.Service {
	return &swaggersrvc{logger}
}
