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

var jsonFile = `{
  "statements":
  [
	{
	  "name": "SHOW VARIABLE COMMIT_TIMESTAMP",
	  "executorName": "ClientSideStatementNoParamExecutor",
	  "resultType": "RESULT_SET",
	  "regex": "(?is)\\A\\s*show\\s+variable\\s+commit_timestamp\\s*\\z",
	  "method": "statementShowCommitTimestamp",
	  "exampleStatements": ["show variable commit_timestamp"],
	  "examplePrerequisiteStatements": ["update foo set bar=1"]
	},
	{
      "name": "SHOW VARIABLE RETRY_ABORTS_INTERNALLY",
      "executorName": "ClientSideStatementNoParamExecutor",
      "resultType": "RESULT_SET",
      "regex": "(?is)\\A\\s*show\\s+variable\\s+retry_aborts_internally\\s*\\z",
      "method": "statementShowRetryAbortsInternally",
      "exampleStatements": ["show variable retry_aborts_internally"],
      "examplePrerequisiteStatements": ["set readonly=false", "set autocommit=false"]
    },
    {
      "name": "SHOW VARIABLE AUTOCOMMIT_DML_MODE",
      "executorName": "ClientSideStatementNoParamExecutor",
      "resultType": "RESULT_SET",
      "regex": "(?is)\\A\\s*show\\s+variable\\s+autocommit_dml_mode\\s*\\z",
      "method": "statementShowAutocommitDmlMode",
      "exampleStatements": ["show variable autocommit_dml_mode"]
    },
    {
      "name": "SHOW VARIABLE READ_ONLY_STALENESS",
      "executorName": "ClientSideStatementNoParamExecutor",
      "resultType": "RESULT_SET",
      "regex": "(?is)\\A\\s*show\\s+variable\\s+read_only_staleness\\s*\\z",
      "method": "statementShowReadOnlyStaleness",
      "exampleStatements": ["show variable read_only_staleness"]
    },
	{
		"name": "SHOW VARIABLE EXCLUDE_TXN_FROM_CHANGE_STREAMS",
		"executorName": "ClientSideStatementNoParamExecutor",
		"resultType": "RESULT_SET",
		"statementType": "SHOW_EXCLUDE_TXN_FROM_CHANGE_STREAMS",
		"regex": "(?is)\\A\\s*show\\s+variable\\s+exclude_txn_from_change_streams\\s*\\z",
		"method": "statementShowExcludeTxnFromChangeStreams",
		"exampleStatements": ["show variable exclude_txn_from_change_streams"]
	},
	{
      "name": "START BATCH DDL",
      "executorName": "ClientSideStatementNoParamExecutor",
      "resultType": "NO_RESULT",
      "regex": "(?is)\\A\\s*(?:start)(?:\\s+batch)(?:\\s+ddl)\\s*\\z",
      "method": "statementStartBatchDdl",
      "exampleStatements": ["start batch ddl"]
    },
    {
      "name": "START BATCH DML",
      "executorName": "ClientSideStatementNoParamExecutor",
      "resultType": "NO_RESULT",
      "regex": "(?is)\\A\\s*(?:start)(?:\\s+batch)(?:\\s+dml)\\s*\\z",
      "method": "statementStartBatchDml",
      "exampleStatements": ["start batch dml"]
    },
    {
      "name": "RUN BATCH",
      "executorName": "ClientSideStatementNoParamExecutor",
      "resultType": "NO_RESULT",
      "regex": "(?is)\\A\\s*(?:run)(?:\\s+batch)\\s*\\z",
      "method": "statementRunBatch",
      "exampleStatements": ["run batch"],
      "examplePrerequisiteStatements": ["start batch ddl"]
    },
    {
      "name": "ABORT BATCH",
      "executorName": "ClientSideStatementNoParamExecutor",
      "resultType": "NO_RESULT",
      "regex": "(?is)\\A\\s*(?:abort)(?:\\s+batch)\\s*\\z",
      "method": "statementAbortBatch",
      "exampleStatements": ["abort batch"],
      "examplePrerequisiteStatements": ["start batch ddl"]
    },
    {
      "name": "SET RETRY_ABORTS_INTERNALLY = TRUE|FALSE",
      "executorName": "ClientSideStatementSetExecutor",
      "resultType": "NO_RESULT",
      "regex": "(?is)\\A\\s*set\\s+retry_aborts_internally\\s*(?:=)\\s*(.*)\\z",
      "method": "statementSetRetryAbortsInternally",
      "exampleStatements": ["set retry_aborts_internally = true", "set retry_aborts_internally = false"],
      "examplePrerequisiteStatements": ["set readonly = false", "set autocommit = false"],
      "setStatement": {
        "propertyName": "RETRY_ABORTS_INTERNALLY",
        "separator": "=",
        "allowedValues": "(TRUE|FALSE)",
        "converterName": "ClientSideStatementValueConverters$BooleanConverter"
      }
    },
    {
      "name": "SET AUTOCOMMIT_DML_MODE = 'PARTITIONED_NON_ATOMIC'|'TRANSACTIONAL'",
      "executorName": "ClientSideStatementSetExecutor",
      "resultType": "NO_RESULT",
      "regex": "(?is)\\A\\s*set\\s+autocommit_dml_mode\\s*(?:=)\\s*(.*)\\z",
      "method": "statementSetAutocommitDmlMode",
      "exampleStatements": ["set autocommit_dml_mode='PARTITIONED_NON_ATOMIC'", "set autocommit_dml_mode='TRANSACTIONAL'"],
      "setStatement": {
        "propertyName": "AUTOCOMMIT_DML_MODE",
        "separator": "=",
        "allowedValues": "'(PARTITIONED_NON_ATOMIC|TRANSACTIONAL)'",
        "converterName": "ClientSideStatementValueConverters$AutocommitDmlModeConverter"
      }
    },
    {
      "name": "SET READ_ONLY_STALENESS = 'STRONG' | 'MIN_READ_TIMESTAMP <timestamp>' | 'READ_TIMESTAMP <timestamp>' | 'MAX_STALENESS <int64>s|ms|us|ns' | 'EXACT_STALENESS (<int64>s|ms|us|ns)'",
      "executorName": "ClientSideStatementSetExecutor",
      "resultType": "NO_RESULT",
      "regex": "(?is)\\A\\s*set\\s+read_only_staleness\\s*(?:=)\\s*(.*)\\z",
      "method": "statementSetReadOnlyStaleness",
      "exampleStatements": ["set read_only_staleness='STRONG'",
        "set read_only_staleness='MIN_READ_TIMESTAMP 2018-01-02T03:04:05.123-08:00'",
        "set read_only_staleness='MIN_READ_TIMESTAMP 2018-01-02T03:04:05.123Z'",
        "set read_only_staleness='MIN_READ_TIMESTAMP 2018-01-02T03:04:05.123+07:45'",
        "set read_only_staleness='READ_TIMESTAMP 2018-01-02T03:04:05.54321-07:00'",
        "set read_only_staleness='READ_TIMESTAMP 2018-01-02T03:04:05.54321Z'",
        "set read_only_staleness='READ_TIMESTAMP 2018-01-02T03:04:05.54321+05:30'",
        "set read_only_staleness='MAX_STALENESS 12s'",
        "set read_only_staleness='MAX_STALENESS 100ms'",
        "set read_only_staleness='MAX_STALENESS 99999us'",
        "set read_only_staleness='MAX_STALENESS 10ns'",
        "set read_only_staleness='EXACT_STALENESS 15s'",
        "set read_only_staleness='EXACT_STALENESS 1500ms'",
        "set read_only_staleness='EXACT_STALENESS 15000000us'",
        "set read_only_staleness='EXACT_STALENESS 9999ns'"],
      "setStatement": {
        "propertyName": "READ_ONLY_STALENESS",
        "separator": "=",
        "allowedValues": "'((STRONG)|(MIN_READ_TIMESTAMP)[\\t ]+((\\d{4})-(\\d{2})-(\\d{2})([Tt](\\d{2}):(\\d{2}):(\\d{2})(\\.\\d{1,9})?)([Zz]|([+-])(\\d{2}):(\\d{2})))|(READ_TIMESTAMP)[\\t ]+((\\d{4})-(\\d{2})-(\\d{2})([Tt](\\d{2}):(\\d{2}):(\\d{2})(\\.\\d{1,9})?)([Zz]|([+-])(\\d{2}):(\\d{2})))|(MAX_STALENESS)[\\t ]+((\\d{1,19})(s|ms|us|ns))|(EXACT_STALENESS)[\\t ]+((\\d{1,19})(s|ms|us|ns)))'",
        "converterName": "ClientSideStatementValueConverters$ReadOnlyStalenessConverter"
      }
    },
	{
		"name": "SET EXCLUDE_TXN_FROM_CHANGE_STREAMS = TRUE|FALSE",
		"executorName": "ClientSideStatementSetExecutor",
		"resultType": "NO_RESULT",
		"statementType": "SET_EXCLUDE_TXN_FROM_CHANGE_STREAMS",
		"regex": "(?is)\\A\\s*set\\s+exclude_txn_from_change_streams\\s*(?:=)\\s*(.*)\\z",
		"method": "statementSetExcludeTxnFromChangeStreams",
		"exampleStatements": ["set exclude_txn_from_change_streams = true", "set exclude_txn_from_change_streams = false"],
		"setStatement": {
			"propertyName": "EXCLUDE_TXN_FROM_CHANGE_STREAMS",
			"separator": "=",
			"allowedValues": "(TRUE|FALSE)",
			"converterName": "ClientSideStatementValueConverters$BooleanConverter"
		}
	}
  ]
}
`
