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

package app

import (
	"context"

	"goa.design/goa/v3/middleware"
	"gorm.io/gorm"
)

// Service defines methods on BaseService
type Service interface {
	Logger(ctx context.Context) *Logger
	LoggerWith(ctx context.Context, args ...interface{}) *Logger
	DB(ctx context.Context) *gorm.DB
	CatalogClonePath() string
}

type environmenter interface {
	Environment() EnvMode
}

// BaseService defines configuraition for creating logger and
// db object with http request id
type BaseService struct {
	env      environmenter
	logger   *Logger
	db       *gorm.DB
	basePath string
}

// Logger looks for http request id in passed context and append it to
// logger. This would help in filtering logs with request id.
func (s *BaseService) Logger(ctx context.Context) *Logger {
	reqID := ctx.Value(middleware.RequestIDKey)
	if reqID == nil {
		return s.logger
	}
	return &Logger{SugaredLogger: s.logger.With("id", reqID.(string))}
}

// LoggerWith gets logger created with http request id from context
// then appends args to it
func (s *BaseService) LoggerWith(ctx context.Context, args ...interface{}) *Logger {
	return &Logger{SugaredLogger: s.Logger(ctx).With(args...)}
}

// Returns the base path where catalog is to be cloned and stored
func (s *BaseService) CatalogClonePath() string {
	return s.basePath
}

// DB gets logger initialized with http request id and creates a gorm db
// session replacing writer for gorm logger with Logger so that gorm log
// will have http request id.
func (s *BaseService) DB(ctx context.Context) *gorm.DB {
	return s.db.Session(&gorm.Session{
		Logger: newGormLogger(s.env.Environment(), s.Logger(ctx)),
	})
}
