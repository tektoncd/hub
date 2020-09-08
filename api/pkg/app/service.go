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

	"github.com/jinzhu/gorm"
	"goa.design/goa/v3/middleware"

	"github.com/tektoncd/hub/api/gen/log"
)

// Service defines methods on BaseService
type Service interface {
	Logger(ctx context.Context) *log.Logger
	LoggerWith(ctx context.Context, args ...interface{}) *log.Logger
	DB(ctx context.Context) *gorm.DB
}

// BaseService defines configuraition for creating logger and
// db object with http request id
type BaseService struct {
	logger *log.Logger
	db     *gorm.DB
}

// adaptor for gorm logger interface
type gormLogger struct {
	*log.Logger
}

func (l *gormLogger) Print(v ...interface{}) {
	l.Info(v...)
}

// Logger looks for http request id in passed context and append it to
// logger. This would help in filtering logs with request id.
func (s *BaseService) Logger(ctx context.Context) *log.Logger {
	reqID := ctx.Value(middleware.RequestIDKey)
	if reqID == nil {
		return s.logger
	}
	return &log.Logger{SugaredLogger: s.logger.With("id", reqID.(string))}
}

// LoggerWith gets logger created with http request id from context
// then appends args to it
func (s *BaseService) LoggerWith(ctx context.Context, args ...interface{}) *log.Logger {
	return &log.Logger{SugaredLogger: s.Logger(ctx).With(args...)}
}

// DB gets logger initialized with http request id and creates a gorm db
// object with replacing default logger with created one, so that each
// gorm log will have http request id.
func (s *BaseService) DB(ctx context.Context) *gorm.DB {
	logger := s.Logger(ctx)
	db := s.db.New()
	db.SetLogger(&gormLogger{logger})
	return db
}
