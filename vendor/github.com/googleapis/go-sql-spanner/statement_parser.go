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
	"encoding/json"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"cloud.google.com/go/spanner"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ddlStatements = map[string]bool{"CREATE": true, "DROP": true, "ALTER": true}
var selectStatements = map[string]bool{"SELECT": true, "WITH": true}
var dmlStatements = map[string]bool{"INSERT": true, "UPDATE": true, "DELETE": true}
var selectAndDmlStatements = union(selectStatements, dmlStatements)

func union(m1 map[string]bool, m2 map[string]bool) map[string]bool {
	res := make(map[string]bool, len(m1)+len(m2))
	for k, v := range m1 {
		res[k] = v
	}
	for k, v := range m2 {
		res[k] = v
	}
	return res
}

// parseParameters returns the parameters in the given sql string, if the input
// sql contains positional parameters it returns the converted sql string with
// all positional parameters replaced with named parameters.
// The sql string must be a valid Cloud Spanner sql statement. It may contain
// comments and (string) literals without any restrictions. That is, string
// literals containing for example an email address ('test@test.com') will be
// recognized as a string literal and not returned as a named parameter.
func parseParameters(sql string) (string, []string, error) {
	sql, err := removeCommentsAndTrim(sql)
	if err != nil {
		return sql, nil, err
	}
	return findParams('?', sql)
}

// RemoveCommentsAndTrim removes any comments in the query string and trims any
// spaces at the beginning and end of the query. This makes checking what type
// of query a string is a lot easier, as only the first word(s) need to be
// checked after this has been removed.
func removeCommentsAndTrim(sql string) (string, error) {
	const singleQuote = '\''
	const doubleQuote = '"'
	const backtick = '`'
	const hyphen = '-'
	const dash = '#'
	const slash = '/'
	const asterisk = '*'
	isInQuoted := false
	isInSingleLineComment := false
	isInMultiLineComment := false
	var startQuote rune
	lastCharWasEscapeChar := false
	isTripleQuoted := false
	res := strings.Builder{}
	res.Grow(len(sql))
	index := 0
	runes := []rune(sql)
	for index < len(runes) {
		c := runes[index]
		if isInQuoted {
			if (c == '\n' || c == '\r') && !isTripleQuoted {
				return "", spanner.ToSpannerError(status.Errorf(codes.InvalidArgument, "statement contains an unclosed literal: %s", sql))
			} else if c == startQuote {
				if lastCharWasEscapeChar {
					lastCharWasEscapeChar = false
				} else if isTripleQuoted {
					if len(runes) > index+2 && runes[index+1] == startQuote && runes[index+2] == startQuote {
						isInQuoted = false
						startQuote = 0
						isTripleQuoted = false
						res.WriteRune(c)
						res.WriteRune(c)
						index += 2
					}
				} else {
					isInQuoted = false
					startQuote = 0
				}
			} else if c == '\\' {
				lastCharWasEscapeChar = true
			} else {
				lastCharWasEscapeChar = false
			}
			res.WriteRune(c)
		} else {
			// We are not in a quoted string.
			if isInSingleLineComment {
				if c == '\n' {
					isInSingleLineComment = false
					// Include the line feed in the result.
					res.WriteRune(c)
				}
			} else if isInMultiLineComment {
				if len(runes) > index+1 && c == asterisk && runes[index+1] == slash {
					isInMultiLineComment = false
					index++
				}
			} else {
				if c == dash || (len(runes) > index+1 && c == hyphen && runes[index+1] == hyphen) {
					// This is a single line comment.
					isInSingleLineComment = true
				} else if len(runes) > index+1 && c == slash && runes[index+1] == asterisk {
					isInMultiLineComment = true
					index++
				} else {
					if c == singleQuote || c == doubleQuote || c == backtick {
						isInQuoted = true
						startQuote = c
						// Check whether it is a triple-quote.
						if len(runes) > index+2 && runes[index+1] == startQuote && runes[index+2] == startQuote {
							isTripleQuoted = true
							res.WriteRune(c)
							res.WriteRune(c)
							index += 2
						}
					}
					res.WriteRune(c)
				}
			}
		}
		index++
	}
	if isInQuoted {
		return "", spanner.ToSpannerError(status.Errorf(codes.InvalidArgument, "statement contains an unclosed literal: %s", sql))
	}
	trimmed := strings.TrimSpace(res.String())
	if len(trimmed) > 0 && trimmed[len(trimmed)-1] == ';' {
		return trimmed[:len(trimmed)-1], nil
	}
	return trimmed, nil
}

// Removes any statement hints at the beginning of the statement.
// It assumes that any comments have already been removed.
func removeStatementHint(sql string) string {
	// Valid statement hints at the beginning of a query statement can only contain a fixed set of
	// possible values. Although it is possible to add a @{FORCE_INDEX=...} as a statement hint, the
	// only allowed value is _BASE_TABLE. This means that we can safely assume that the statement
	// hint will not contain any special characters, for example a closing curly brace or one of the
	// keywords SELECT, UPDATE, DELETE, WITH, and that we can keep the check simple by just
	// searching for the first occurrence of a keyword that should be preceded by a closing curly
	// brace at the end of the statement hint.
	startStatementHintIndex := strings.Index(sql, "{")
	// Statement hints are allowed for both queries and DML statements.
	startQueryIndex := -1
	upperCaseSql := strings.ToUpper(sql)
	for keyword := range selectAndDmlStatements {
		if startQueryIndex = strings.Index(upperCaseSql, keyword); startQueryIndex > -1 {
			break
		}
	}
	// The startQueryIndex can theoretically be larger than the length of the SQL string,
	// as the length of the uppercase SQL string can be different from the length of the
	// lower/mixed case SQL string. This is however only the case for specific non-ASCII
	// characters that are not allowed in a statement hint, so in that case we can safely
	// assume the statement to be invalid.
	if startQueryIndex > -1 && startQueryIndex < len(sql) {
		endStatementHintIndex := strings.LastIndex(sql[:startQueryIndex], "}")
		if startStatementHintIndex == -1 || startStatementHintIndex > endStatementHintIndex || endStatementHintIndex >= len(sql)-1 {
			// Looks like an invalid statement hint. Just ignore at this point
			// and let the caller handle the invalid query.
			return sql
		}
		return strings.TrimSpace(sql[endStatementHintIndex+1:])
	}
	// Seems invalid, just return the original statement.
	return sql
}

// This function assumes that all comments have already
// been removed from the statement.
func findParams(positionalParamChar rune, sql string) (string, []string, error) {
	const paramPrefix = '@'
	const singleQuote = '\''
	const doubleQuote = '"'
	const backtick = '`'
	isInQuoted := false
	var startQuote rune
	lastCharWasEscapeChar := false
	isTripleQuoted := false
	hasNamedParameter := false
	hasPositionalParameter := false
	namedParams := make([]string, 0)
	parsedSQL := strings.Builder{}
	parsedSQL.Grow(len(sql))
	positionalParameterIndex := 1
	index := 0
	runes := []rune(sql)
	for index < len(runes) {
		c := runes[index]
		if isInQuoted {
			if (c == '\n' || c == '\r') && !isTripleQuoted {
				return sql, nil, spanner.ToSpannerError(status.Errorf(codes.InvalidArgument, "statement contains an unclosed literal: %s", sql))
			} else if c == startQuote {
				if lastCharWasEscapeChar {
					lastCharWasEscapeChar = false
				} else if isTripleQuoted {
					if len(runes) > index+2 && runes[index+1] == startQuote && runes[index+2] == startQuote {
						isInQuoted = false
						startQuote = 0
						isTripleQuoted = false
						parsedSQL.WriteRune(c)
						parsedSQL.WriteRune(c)
						index += 2
					}
				} else {
					isInQuoted = false
					startQuote = 0
				}
			} else if c == '\\' {
				lastCharWasEscapeChar = true
			} else {
				lastCharWasEscapeChar = false
			}
			parsedSQL.WriteRune(c)
		} else {
			// We are not in a quoted string. It's a parameter if it is an '@' followed by a letter or an underscore.
			// See https://cloud.google.com/spanner/docs/lexical#identifiers for identifier rules.
			if c == paramPrefix && len(runes) > index+1 && (unicode.IsLetter(runes[index+1]) || runes[index+1] == '_') {
				if hasPositionalParameter {
					return sql, nil, spanner.ToSpannerError(status.Errorf(codes.InvalidArgument, "statement must not contain both named and positional parameter: %s", sql))
				}
				parsedSQL.WriteRune(c)
				index++
				startIndex := index
				for index < len(runes) {
					if !(unicode.IsLetter(runes[index]) || unicode.IsDigit(runes[index]) || runes[index] == '_') {
						hasNamedParameter = true
						namedParams = append(namedParams, string(runes[startIndex:index]))
						parsedSQL.WriteRune(runes[index])
						break
					}
					if index == len(runes)-1 {
						hasNamedParameter = true
						namedParams = append(namedParams, string(runes[startIndex:]))
						parsedSQL.WriteRune(runes[index])
						break
					}
					parsedSQL.WriteRune(runes[index])
					index++
				}
			} else if c == positionalParamChar {
				if hasNamedParameter {
					return sql, nil, spanner.ToSpannerError(status.Errorf(codes.InvalidArgument, "statement must not contain both named and positional parameter: %s", sql))
				}
				hasPositionalParameter = true
				parsedSQL.WriteString("@p" + strconv.Itoa(positionalParameterIndex))
				namedParams = append(namedParams, "p"+strconv.Itoa(positionalParameterIndex))
				positionalParameterIndex++
			} else {
				if c == singleQuote || c == doubleQuote || c == backtick {
					isInQuoted = true
					startQuote = c
					// Check whether it is a triple-quote.
					if len(runes) > index+2 && runes[index+1] == startQuote && runes[index+2] == startQuote {
						isTripleQuoted = true
						parsedSQL.WriteRune(c)
						parsedSQL.WriteRune(c)
						index += 2
					}
				}
				parsedSQL.WriteRune(c)
			}
		}
		index++
	}
	if isInQuoted {
		return sql, nil, spanner.ToSpannerError(status.Errorf(codes.InvalidArgument, "statement contains an unclosed literal: %s", sql))
	}
	if hasNamedParameter {
		return sql, namedParams, nil
	}
	sql = strings.TrimSpace(parsedSQL.String())
	if len(sql) > 0 && sql[len(sql)-1] == ';' {
		sql = sql
	}
	return sql, namedParams, nil
}

// isDDL returns true if the given sql string is a DDL statement.
func isDDL(query string) (bool, error) {
	query, err := removeCommentsAndTrim(query)
	if err != nil {
		return false, err
	}
	// We can safely check if the string starts with a specific string, as we
	// have already removed all leading spaces, and there are no keywords that
	// start with the same substring as one of the DDL keywords.
	for ddl := range ddlStatements {
		if len(query) >= len(ddl) && strings.EqualFold(query[:len(ddl)], ddl) {
			return true, nil
		}
	}
	return false, nil
}

// clientSideStatements are loaded from the client_side_statements.json file.
type clientSideStatements struct {
	Statements []*clientSideStatement `json:"statements"`
	executor   *statementExecutor
}

// clientSideStatement is the definition of a statement that can be executed on
// a connection and that will be handled by the connection itself, instead of
// sending it to Spanner.
type clientSideStatement struct {
	Name                          string `json:"name"`
	ExecutorName                  string `json:"executorName"`
	execContext                   func(ctx context.Context, c *conn, params string, args []driver.NamedValue) (driver.Result, error)
	queryContext                  func(ctx context.Context, c *conn, params string, args []driver.NamedValue) (driver.Rows, error)
	ResultType                    string `json:"resultType"`
	Regex                         string `json:"regex"`
	regexp                        *regexp.Regexp
	MethodName                    string `json:"method"`
	method                        func(query string) error
	ExampleStatements             []string `json:"exampleStatements"`
	ExamplePrerequisiteStatements []string `json:"examplePrerequisiteStatements"`

	setStatement `json:"setStatement"`
}

type setStatement struct {
	PropertyName  string `json:"propertyName"`
	Separator     string `json:"separator"`
	AllowedValues string `json:"allowedValues"`
	ConverterName string `json:"converterName"`
}

var statementsInit sync.Once
var statements *clientSideStatements
var statementsCompileErr error

// compileStatements loads all client side statements from the json file and
// assigns the Go methods to the different statements that should be executed
// when on of the statements is executed on a connection.
func compileStatements() error {
	statements = new(clientSideStatements)
	err := json.Unmarshal([]byte(jsonFile), statements)
	if err != nil {
		return err
	}
	statements.executor = &statementExecutor{}
	for _, stmt := range statements.Statements {
		stmt.regexp, err = regexp.Compile(stmt.Regex)
		if err != nil {
			return err
		}
		i := reflect.ValueOf(statements.executor).MethodByName(strings.TrimPrefix(stmt.MethodName, "statement")).Interface()
		if execContext, ok := i.(func(ctx context.Context, c *conn, query string, args []driver.NamedValue) (driver.Result, error)); ok {
			stmt.execContext = execContext
		}
		if queryContext, ok := i.(func(ctx context.Context, c *conn, query string, args []driver.NamedValue) (driver.Rows, error)); ok {
			stmt.queryContext = queryContext
		}
	}
	return nil
}

// executableClientSideStatement is the combination of a pre-defined client-side
// statement, the connection it should be executed on and any additional
// parameters that were included in the statement.
type executableClientSideStatement struct {
	*clientSideStatement
	conn   *conn
	query  string
	params string
}

func (c *executableClientSideStatement) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if c.clientSideStatement.execContext == nil {
		return nil, spanner.ToSpannerError(status.Errorf(codes.InvalidArgument, "%q cannot be used with execContext", c.query))
	}
	return c.clientSideStatement.execContext(ctx, c.conn, c.params, args)
}

func (c *executableClientSideStatement) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if c.clientSideStatement.queryContext == nil {
		return nil, spanner.ToSpannerError(status.Errorf(codes.InvalidArgument, "%q cannot be used with queryContext", c.query))
	}
	return c.clientSideStatement.queryContext(ctx, c.conn, c.params, args)
}

// parseClientSideStatement returns the executableClientSideStatement that
// corresponds with the given query string, or nil if it is not a valid client
// side statement.
func parseClientSideStatement(c *conn, query string) (*executableClientSideStatement, error) {
	statementsInit.Do(func() {
		if err := compileStatements(); err != nil {
			statementsCompileErr = err
		}
	})
	if statementsCompileErr != nil {
		return nil, statementsCompileErr
	}
	for _, stmt := range statements.Statements {
		if stmt.regexp.MatchString(query) {
			var params string
			if stmt.setStatement.Separator != "" {
				p := strings.SplitN(query, stmt.setStatement.Separator, 2)
				if len(p) == 2 {
					params = strings.TrimSpace(p[1])
				}
			}
			return &executableClientSideStatement{stmt, c, query, params}, nil
		}
	}
	return nil, nil
}
