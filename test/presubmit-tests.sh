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

declare -r POSTGRES_CONTAINER=postgres

info() {
  echo "INFO: $@"
}

warn() {
  echo "WARN: $@"
}

container_exists() {
  local name=$1; shift
  [[ $(docker ps -a --filter "name=^/$name$" --format '{{.Names}}') == $name ]]
}

ensure_container_running() {
  local container=$1; shift
  container_exists "$container" || return 1
  [[ $(docker inspect --format '{{ .State.Status }}' $container) == running ]]
}

ensure_no_container_running(){
  local container=$1; shift
  container_exists "$container" || return 0

  info "removing container: $container"
  docker rm -f -v "$container"
}

##  override teardown_environment
teardown_environment() {
  header Teardown

  ensure_no_container_running "$POSTGRES_CONTAINER"
  rm -rf ${WORK_DIR}
}


api-unittest(){
  pwd
  go test -v ./pkg/...
}

run_db() {
  header "Running db"
  local container_id=""

  source "$API_DIR/test/config/env.test"
  container_id=$(docker run -d  \
    --name $POSTGRES_CONTAINER \
    -e POSTGRES_USER="$POSTGRESQL_USER"  \
    -e POSTGRES_PASSWORD="$POSTGRESQL_PASSWORD" \
    -e POSTGRES_DB="$POSTGRESQL_DATABASE" \
    -p "$POSTGRESQL_PORT:5432" \
    postgres
  )
  [[ "$container_id" != "" ]] || return 1
  return 0
}

api-test(){
  cd $API_DIR
  ensure_no_container_running "$POSTGRES_CONTAINER"

  info "running db"
  run_db || abort "Failed to run db container $POSTGRES_CONTAINER"
  sleep 5s
  docker ps -a
  ensure_container_running "$POSTGRES_CONTAINER"

  api-unittest
}

api-build(){
  cd $API_DIR
  go mod vendor
  go build ./cmd/...
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
    api-test
  )
}

run_integration_tests() {
  warn "No integration tests to run"
  return 0
}

main $@
