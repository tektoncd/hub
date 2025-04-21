// Copyright 2021 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
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
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// statementExecutor is an empty struct that is used to hold the execution methods
// of the different client side statements. This makes it possible to look up the
// methods using reflection, which is not possible if the methods do not belong to
// a struct. The methods all accept the same arguments and return the same types.
// This is to ensure that they can be assigned to a compiled clientSideStatement.
//
// The different methods of statementExecutor are invoked by a connection when one
// of the valid client side statements is executed on a connection. These methods
// are responsible for any argument parsing and translating that might be needed
// before the corresponding method on the connection can be called.
//
// The names of the methods are exactly equal to the naming in the
// client_side_statements.json file. This means that some methods do not adhere
// to the Go style guide, as these method names are equal for all languages that
// implement the Connection API.
type statementExecutor struct {
}

func (s *statementExecutor) ShowCommitTimestamp(_ context.Context, c *conn, _ string, _ []driver.NamedValue) (driver.Rows, error) {
	ts, err := c.CommitTimestamp()
	var commitTs *time.Time
	if err == nil {
		commitTs = &ts
	}
	it, err := createTimestampIterator("CommitTimestamp", commitTs)
	if err != nil {
		return nil, err
	}
	return &rows{it: it}, nil
}

func (s *statementExecutor) ShowRetryAbortsInternally(_ context.Context, c *conn, _ string, _ []driver.NamedValue) (driver.Rows, error) {
	it, err := createBooleanIterator("RetryAbortsInternally", c.RetryAbortsInternally())
	if err != nil {
		return nil, err
	}
	return &rows{it: it}, nil
}

func (s *statementExecutor) ShowAutocommitDmlMode(_ context.Context, c *conn, _ string, _ []driver.NamedValue) (driver.Rows, error) {
	it, err := createStringIterator("AutocommitDMLMode", c.AutocommitDMLMode().String())
	if err != nil {
		return nil, err
	}
	return &rows{it: it}, nil
}

func (s *statementExecutor) ShowReadOnlyStaleness(_ context.Context, c *conn, _ string, _ []driver.NamedValue) (driver.Rows, error) {
	it, err := createStringIterator("ReadOnlyStaleness", c.ReadOnlyStaleness().String())
	if err != nil {
		return nil, err
	}
	return &rows{it: it}, nil
}

func (s *statementExecutor) ShowExcludeTxnFromChangeStreams(_ context.Context, c *conn, _ string, _ []driver.NamedValue) (driver.Rows, error) {
	it, err := createBooleanIterator("ExcludeTxnFromChangeStreams", c.ExcludeTxnFromChangeStreams())
	if err != nil {
		return nil, err
	}
	return &rows{it: it}, nil
}

func (s *statementExecutor) StartBatchDdl(_ context.Context, c *conn, _ string, _ []driver.NamedValue) (driver.Result, error) {
	return c.startBatchDDL()
}

func (s *statementExecutor) StartBatchDml(_ context.Context, c *conn, _ string, _ []driver.NamedValue) (driver.Result, error) {
	return c.startBatchDML()
}

func (s *statementExecutor) RunBatch(ctx context.Context, c *conn, _ string, _ []driver.NamedValue) (driver.Result, error) {
	return c.runBatch(ctx)
}

func (s *statementExecutor) AbortBatch(_ context.Context, c *conn, _ string, _ []driver.NamedValue) (driver.Result, error) {
	return c.abortBatch()
}

func (s *statementExecutor) SetRetryAbortsInternally(_ context.Context, c *conn, params string, _ []driver.NamedValue) (driver.Result, error) {
	if params == "" {
		return nil, spanner.ToSpannerError(status.Error(codes.InvalidArgument, "no value given for RetryAbortsInternally"))
	}
	retry, err := strconv.ParseBool(params)
	if err != nil {
		return nil, spanner.ToSpannerError(status.Errorf(codes.InvalidArgument, "invalid boolean value: %s", params))
	}
	return c.setRetryAbortsInternally(retry)
}

func (s *statementExecutor) SetAutocommitDmlMode(_ context.Context, c *conn, params string, _ []driver.NamedValue) (driver.Result, error) {
	if params == "" {
		return nil, spanner.ToSpannerError(status.Error(codes.InvalidArgument, "no value given for AutocommitDMLMode"))
	}
	var mode AutocommitDMLMode
	switch strings.ToUpper(params) {
	case fmt.Sprintf("'%s'", strings.ToUpper(Transactional.String())):
		mode = Transactional
	case fmt.Sprintf("'%s'", strings.ToUpper(PartitionedNonAtomic.String())):
		mode = PartitionedNonAtomic
	default:
		return nil, spanner.ToSpannerError(status.Errorf(codes.InvalidArgument, "invalid AutocommitDMLMode value: %s", params))
	}
	return c.setAutocommitDMLMode(mode)
}

func (s *statementExecutor) SetExcludeTxnFromChangeStreams(_ context.Context, c *conn, params string, _ []driver.NamedValue) (driver.Result, error) {
	if params == "" {
		return nil, spanner.ToSpannerError(status.Error(codes.InvalidArgument, "no value given for ExcludeTxnFromChangeStreams"))
	}
	exclude, err := strconv.ParseBool(params)
	if err != nil {
		return nil, spanner.ToSpannerError(status.Errorf(codes.InvalidArgument, "invalid boolean value: %s", params))
	}
	return c.setExcludeTxnFromChangeStreams(exclude)
}

var strongRegexp = regexp.MustCompile("(?i)'STRONG'")
var exactStalenessRegexp = regexp.MustCompile(`(?i)'(?P<type>EXACT_STALENESS)[\t ]+(?P<duration>(\d{1,19})(s|ms|us|ns))'`)
var maxStalenessRegexp = regexp.MustCompile(`(?i)'(?P<type>MAX_STALENESS)[\t ]+(?P<duration>(\d{1,19})(s|ms|us|ns))'`)
var readTimestampRegexp = regexp.MustCompile(`(?i)'(?P<type>READ_TIMESTAMP)[\t ]+(?P<timestamp>(\d{4})-(\d{2})-(\d{2})([Tt](\d{2}):(\d{2}):(\d{2})(\.\d{1,9})?)([Zz]|([+-])(\d{2}):(\d{2})))'`)
var minReadTimestampRegexp = regexp.MustCompile(`(?i)'(?P<type>MIN_READ_TIMESTAMP)[\t ]+(?P<timestamp>(\d{4})-(\d{2})-(\d{2})([Tt](\d{2}):(\d{2}):(\d{2})(\.\d{1,9})?)([Zz]|([+-])(\d{2}):(\d{2})))'`)

func (s *statementExecutor) SetReadOnlyStaleness(_ context.Context, c *conn, params string, _ []driver.NamedValue) (driver.Result, error) {
	if params == "" {
		return nil, spanner.ToSpannerError(status.Error(codes.InvalidArgument, "no value given for ReadOnlyStaleness"))
	}
	invalidErr := spanner.ToSpannerError(status.Errorf(codes.InvalidArgument, "invalid ReadOnlyStaleness value: %s", params))

	var staleness spanner.TimestampBound

	if strongRegexp.MatchString(params) {
		staleness = spanner.StrongRead()
	} else if exactStalenessRegexp.MatchString(params) {
		d, err := parseDuration(exactStalenessRegexp, params)
		if err != nil {
			return nil, err
		}
		staleness = spanner.ExactStaleness(d)
	} else if maxStalenessRegexp.MatchString(params) {
		d, err := parseDuration(maxStalenessRegexp, params)
		if err != nil {
			return nil, err
		}
		staleness = spanner.MaxStaleness(d)
	} else if readTimestampRegexp.MatchString(params) {
		t, err := parseTimestamp(readTimestampRegexp, params)
		if err != nil {
			return nil, err
		}
		staleness = spanner.ReadTimestamp(t)
	} else if minReadTimestampRegexp.MatchString(params) {
		t, err := parseTimestamp(minReadTimestampRegexp, params)
		if err != nil {
			return nil, err
		}
		staleness = spanner.MinReadTimestamp(t)
	} else {
		return nil, invalidErr
	}
	return c.setReadOnlyStaleness(staleness)
}

func parseDuration(re *regexp.Regexp, params string) (time.Duration, error) {
	matches := matchesToMap(re, params)
	if matches["duration"] == "" {
		return 0, spanner.ToSpannerError(status.Error(codes.InvalidArgument, "No duration found in staleness string"))
	}
	d, err := time.ParseDuration(matches["duration"])
	if err != nil {
		return 0, spanner.ToSpannerError(status.Errorf(codes.InvalidArgument, "Invalid duration: %s", matches["duration"]))
	}
	return d, nil
}

func parseTimestamp(re *regexp.Regexp, params string) (time.Time, error) {
	matches := matchesToMap(re, params)
	if matches["timestamp"] == "" {
		return time.Time{}, spanner.ToSpannerError(status.Error(codes.InvalidArgument, "No timestamp found in staleness string"))
	}
	t, err := time.Parse(time.RFC3339Nano, matches["timestamp"])
	if err != nil {
		return time.Time{}, spanner.ToSpannerError(status.Errorf(codes.InvalidArgument, "Invalid timestamp: %s", matches["timestamp"]))
	}
	return t, nil
}

func matchesToMap(re *regexp.Regexp, s string) map[string]string {
	match := re.FindStringSubmatch(s)
	matches := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			matches[name] = match[i]
		}
	}
	return matches
}

// createBooleanIterator creates a row iterator with a single BOOL column with
// one row. This is used for client side statements that return a result set
// containing a BOOL value.
func createBooleanIterator(column string, value bool) (*clientSideIterator, error) {
	return createSingleValueIterator(column, value, sppb.TypeCode_BOOL)
}

// createStringIterator creates a row iterator with a single STRING column with
// one row. This is used for client side statements that return a result set
// containing a STRING value.
func createStringIterator(column string, value string) (*clientSideIterator, error) {
	return createSingleValueIterator(column, value, sppb.TypeCode_STRING)
}

// createTimestampIterator creates a row iterator with a single TIMESTAMP column with
// one row. This is used for client side statements that return a result set
// containing a TIMESTAMP value.
func createTimestampIterator(column string, value *time.Time) (*clientSideIterator, error) {
	return createSingleValueIterator(column, value, sppb.TypeCode_TIMESTAMP)
}

func createSingleValueIterator(column string, value interface{}, code sppb.TypeCode) (*clientSideIterator, error) {
	row, err := spanner.NewRow([]string{column}, []interface{}{value})
	if err != nil {
		return nil, err
	}
	return &clientSideIterator{
		metadata: &sppb.ResultSetMetadata{
			RowType: &sppb.StructType{
				Fields: []*sppb.StructType_Field{
					{Name: column, Type: &sppb.Type{Code: code}},
				},
			},
		},
		rows: []*spanner.Row{row},
	}, nil
}

// clientSideIterator implements the rowIterator interface for client side
// statements. All values are created and kept in memory, and this struct
// should only be used for small result sets.
type clientSideIterator struct {
	metadata *sppb.ResultSetMetadata
	rows     []*spanner.Row
	index    int
	stopped  bool
}

func (t *clientSideIterator) Next() (*spanner.Row, error) {
	if t.index == len(t.rows) {
		return nil, io.EOF
	}
	row := t.rows[t.index]
	t.index++
	return row, nil
}

func (t *clientSideIterator) Stop() {
	t.stopped = true
	t.rows = nil
	t.metadata = nil
}

func (t *clientSideIterator) Metadata() *sppb.ResultSetMetadata {
	return t.metadata
}
