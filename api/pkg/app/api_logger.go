// Copyright © 2024 The Tekton Authors.
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

import "go.uber.org/zap"

type Logger struct {
	*zap.SugaredLogger
}

// New creates a new zap logger
func New(serviceName string, production bool) *Logger {

	if production {
		l, _ := zap.NewProduction()
		return &Logger{l.Sugar().With(zap.String("service", serviceName))}
	} else {
		l, _ := zap.NewDevelopment()
		return &Logger{l.Sugar().With(zap.String("service", serviceName))}
	}
}

// Log is called by the log middleware to log HTTP requests key values
func (logger *Logger) Log(keyvals ...interface{}) error {
	logger.Infow("HTTP Request", keyvals...)
	return nil
}
