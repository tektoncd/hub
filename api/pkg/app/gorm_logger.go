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
	"fmt"
	"strings"
	"time"

	"github.com/tektoncd/hub/api/gen/log"
	glog "gorm.io/gorm/logger"
)

func newGormLogger(mode EnvMode, l *log.Logger) glog.Interface {
	if mode == Production {
		return glog.New(&prodWriter{l}, glog.Config{
			SlowThreshold: 100 * time.Millisecond,
			LogLevel:      glog.Info,
			Colorful:      false,
		})
	}

	return glog.New(&devWriter{l}, glog.Config{
		SlowThreshold: 100 * time.Millisecond,
		LogLevel:      glog.Info,
		Colorful:      true,
	})
}

// adaptor for gorm logger interface
type prodWriter struct {
	log *log.Logger
}

func (w *prodWriter) Printf(format string, data ...interface{}) {

	fields := strings.Fields(strings.Replace(format, "\n", " ", -1))
	log := w.log.With("file", data[0])

	data = data[1:]
	fields = fields[1:]

	msg := ""

	for i, d := range data {
		dt := d
		switch dt.(type) {
		case error:
			log = log.With("db-error", d)
		case float64:
			log = log.With("duration", fmt.Sprintf(fields[i], d))
		case int64:
			log = log.With("rows", d)
		case string:
			if i == len(data)-1 {
				msg = d.(string)
			} else {
				log = log.With("unknown", d)
			}
		default:
			log = log.With("unknown", d)
		}
	}
	log.Info(msg)
}

// adaptor for gorm logger interface
type devWriter struct {
	log *log.Logger
}

func (w *devWriter) Printf(msg string, data ...interface{}) {
	w.log.SugaredLogger.Infof(strings.Replace(msg, "%s ", "%s\n", 1), data...)
}
