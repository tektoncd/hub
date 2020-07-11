#!/usr/bin/env bash

# Copyright 2020 Tekton Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script runs the presubmit tests; it is started by prow for each PR.
# For convenience, it can also be executed manually.
# Running the script without parameters, or with the --all-tests
# flag, causes all tests to be executed, in the right order.
# Use the flags --build-tests, --unit-tests and --integration-tests
# to run a specific set of tests.

# Markdown linting failures don't show up properly in Gubernator resulting
# in a net-negative contributor experience.


declare -r SCRIPT_PATH=$(readlink -f "$0")
declare -r SCRIPT_DIR=$(cd $(dirname "$SCRIPT_PATH") && pwd)
declare -r API_DIR="$SCRIPT_DIR/../api"

source $(dirname $0)/../vendor/github.com/tektoncd/plumbing/scripts/presubmit-tests.sh

# TODO: enable this pleaseeeeee
#set -e -u -o pipefail

info() {
  echo "INFO: $@"
}

warn() {
  echo "WARN: $@"
}

install-postgres() {
  info Installing postgres 🛢🛢🛢
  apt-get install -y postgresql postgresql-contrib
  pg_ctlcluster 11 main start
}

set-pg-passwd() {
  local pass="$1"; shift
  su - postgres -c \
    "psql -c \"ALTER USER postgres PASSWORD '$pass';\""
}

api-unittest(){
  install-postgres
  source $API_DIR/test/config/env.test
  set-pg-passwd "$POSTGRESQL_PASSWORD"
  pwd

  info Create test db - $POSTGRESQL_DATABASE

  PGPASSWORD=$POSTGRESQL_PASSWORD \
    psql -h localhost -p 5432 \
    -U $POSTGRESQL_USER -c "create database $POSTGRESQL_DATABASE;"

  info Running unittests

  go mod vendor
  go test -p 1 -v ./pkg/...
}

api-build(){
  go mod vendor
  go build -mod=vendor ./cmd/...
}

### presubmit hooks ###

run_build_tests() {
  # run in a subshell so that path and shell options -eu -o pipefail will
  # will remain the same when it exits
  (
    set -eu -o pipefail
    api-build
  )
}

run_unit_tests() {
  # run in a subshell so that path and shell options -eu -o pipefail will
  # will remain the same when it exits
  (
    set -eu -o pipefail


    cd $API_DIR
    api-build
    api-unittest
  )
}

run_integration_tests() {
  warn "No integration tests to run"
  return 0
}

main $@
