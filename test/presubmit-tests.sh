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
declare -r UI_DIR="$SCRIPT_DIR/../ui"

source $(dirname $0)/../vendor/github.com/tektoncd/plumbing/scripts/presubmit-tests.sh

# TODO: enable this pleaseeeeee
#set -e -u -o pipefail

info() {
  echo "INFO: $@"
}

warn() {
  echo "WARN: $@"
}

install-node() {
  info Installing node

  curl -sL https://deb.nodesource.com/setup_14.x | bash -
  apt-get install -y nodejs

  node --version
}

ui-unittest(){
    install-node

    CI=true npm test
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
  set-pg-passwd "$POSTGRES_PASSWORD"
  pwd

  info Create test db - $POSTGRES_DB

  PGPASSWORD=$POSTGRES_PASSWORD \
    psql -h localhost -p 5432 \
    -U $POSTGRES_USER -c "create database $POSTGRES_DB;"

  info Running unittests

  go mod vendor
  go test -p 1 -v ./pkg/... && 
  go test -p 1 -v ./v1/service/...
}

api-build(){
  go mod vendor
  go build -mod=vendor ./cmd/...
}

ui-build(){
  cd $UI_DIR
  install-node
  npm install

  CI=true npm run build
}

### presubmit hooks ###

run_build_tests() {
  # run in a subshell so that path and shell options -eu -o pipefail will
  # will remain the same when it exits
  (
    set -eu -o pipefail
    api-build
    ui-build
  )
}

run_unit_tests() {
  # run in a subshell so that path and shell options -eu -o pipefail will
  # will remain the same when it exits
  (
    set -eu -o pipefail

    cd $API_DIR
    api-unittest
  )
  (
    set -eu -o pipefail

    cd $UI_DIR
    ui-unittest
  )
}

run_integration_tests() {
  warn "No integration tests to run"
  return 0
}

main $@
