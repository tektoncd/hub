// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spannerdriver

import (
	"context"
	"database/sql/driver"

	"cloud.google.com/go/spanner"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type stmt struct {
	conn    *conn
	numArgs int
	query   string
}

func (s *stmt) Close() error {
	return nil
}

func (s *stmt) NumInput() int {
	return s.numArgs
}

func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, spanner.ToSpannerError(status.Errorf(codes.Unimplemented, "use ExecContext instead"))
}

func (s *stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	return s.conn.ExecContext(ctx, s.query, args)
}

func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, spanner.ToSpannerError(status.Errorf(codes.Unimplemented, "use QueryContext instead"))
}

func (s *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	ss, err := prepareSpannerStmt(s.query, args)
	if err != nil {
		return nil, err
	}

	var it rowIterator
	if s.conn.tx != nil {
		it = s.conn.tx.Query(ctx, ss)
	} else {
		it = &readOnlyRowIterator{s.conn.client.Single().WithTimestampBound(s.conn.readOnlyStaleness).Query(ctx, ss)}
	}
	return &rows{it: it}, nil
}

func prepareSpannerStmt(q string, args []driver.NamedValue) (spanner.Statement, error) {
	q, names, err := parseParameters(q)
	if err != nil {
		return spanner.Statement{}, err
	}
	ss := spanner.NewStatement(q)
	for i, v := range args {
		name := args[i].Name
		if name == "" && len(names) > i {
			name = names[i]
		}
		if name != "" {
			ss.Params[name] = convertParam(v.Value)
		}
	}
	// Verify that all parameters have a value.
	for _, name := range names {
		if _, ok := ss.Params[name]; !ok {
			return spanner.Statement{}, spanner.ToSpannerError(status.Errorf(codes.InvalidArgument, "missing value for query parameter %v", name))
		}
	}
	return ss, nil
}

func convertParam(v driver.Value) driver.Value {
	switch v := v.(type) {
	default:
		return v
	case int:
		return int64(v)
	case []int:
		if v == nil {
			return []int64(nil)
		}
		res := make([]int64, len(v))
		for i, val := range v {
			res[i] = int64(val)
		}
		return res
	case uint:
		return int64(v)
	case []uint:
		if v == nil {
			return []int64(nil)
		}
		res := make([]int64, len(v))
		for i, val := range v {
			res[i] = int64(val)
		}
		return res
	case uint64:
		return int64(v)
	case []uint64:
		if v == nil {
			return []int64(nil)
		}
		res := make([]int64, len(v))
		for i, val := range v {
			res[i] = int64(val)
		}
		return res
	case *uint64:
		if v == nil {
			return (*int64)(nil)
		}
		vi := int64(*v)
		return &vi
	case *int:
		if v == nil {
			return (*int64)(nil)
		}
		vi := int64(*v)
		return &vi
	case []*int:
		if v == nil {
			return []*int64(nil)
		}
		res := make([]*int64, len(v))
		for i, val := range v {
			if val == nil {
				res[i] = nil
			} else {
				z := int64(*val)
				res[i] = &z
			}
		}
		return res
	case *[]int:
		if v == nil {
			return []int64(nil)
		}
		res := make([]int64, len(*v))
		for i, val := range *v {
			res[i] = int64(val)
		}
		return res
	case *uint:
		if v == nil {
			return (*int64)(nil)
		}
		vi := int64(*v)
		return &vi
	case []*uint:
		if v == nil {
			return []*int64(nil)
		}
		res := make([]*int64, len(v))
		for i, val := range v {
			if val == nil {
				res[i] = nil
			} else {
				z := int64(*val)
				res[i] = &z
			}
		}
		return res
	case *[]uint:
		if v == nil {
			return []int64(nil)
		}
		res := make([]int64, len(*v))
		for i, val := range *v {
			res[i] = int64(val)
		}
		return res
	case []*uint64:
		if v == nil {
			return []*int64(nil)
		}
		res := make([]*int64, len(v))
		for i, val := range v {
			if val == nil {
				res[i] = nil
			} else {
				z := int64(*val)
				res[i] = &z
			}
		}
		return res
	case *[]uint64:
		if v == nil {
			return []int64(nil)
		}
		res := make([]int64, len(*v))
		for i, val := range *v {
			res[i] = int64(val)
		}
		return res
	}
}

type result struct {
	rowsAffected int64
}

func (r *result) LastInsertId() (int64, error) {
	return 0, spanner.ToSpannerError(status.Errorf(codes.Unimplemented, "Cloud Spanner does not support auto-generated ids"))
}

func (r *result) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}
