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

# Define this variable if you want to run all tests and not just the modified one.
TEST_RUN_ALL_TESTS=${TEST_RUN_ALL_TESTS:-""}

source $(dirname $0)/../vendor/github.com/tektoncd/plumbing/scripts/presubmit-tests.sh

detect_changes() {
  local dir=${1}
  git --no-pager diff --name-only "${PULL_BASE_SHA}".."${PULL_PULL_SHA}"|grep "${dir}"
}

[[ -z ${TEST_RUN_ALL_TESTS} ]] && [[ ! -z $(detect_changes "test") ]] && TEST_RUN_ALL_TESTS=1

info() {
  echo "INFO: $@"
}

warn() {
  echo "WARN: $@"
}

err() {
  echo "ERROR: $@"
}

install-node() {
  info Installing node

  curl -sL https://deb.nodesource.com/setup_18.x | bash -
  apt-get install -y nodejs

  node --version

  apt-get install -y npm
  npm --version
}

ui-unittest() {
  install-node

  make ui-test || {
    err 'ui unit test failed'
    return 1
  }
}

install-postgres() {
  info Installing postgres 🛢🛢🛢
  # Update the system and install prerequisites
  apt-get update && apt-get install -y wget ca-certificates gnupg
  
  # Check if the system is on an unsupported version
  if grep -q "lunar" /etc/os-release; then
    echo "Detected Ubuntu Lunar (23.04), switching to old-releases repository."
    sed -i 's|http://archive.ubuntu.com/ubuntu%7Chttp://old-releases.ubuntu.com/ubuntu%7Cg' /etc/apt/sources.list
    sed -i 's|http://security.ubuntu.com/ubuntu%7Chttp://old-releases.ubuntu.com/ubuntu%7Cg' /etc/apt/sources.list
    apt-get update
  fi
  
  wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
  sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" >> /etc/apt/sources.list.d/pgdg.list'
  apt-get update
  apt-get -y install postgresql-15
  service postgresql start
  pg_ctlcluster 15 main start
}

set-pg-passwd() {
  local pass="$1"; shift
  su - postgres -c \
    "psql -c \"ALTER USER postgres PASSWORD '$pass';\""
}

api-unittest() {
  install-postgres
  source $API_DIR/test/config/env.test
  set-pg-passwd "$POSTGRES_PASSWORD"
  pwd

  info Create test db - $POSTGRES_DB

  PGPASSWORD=$POSTGRES_PASSWORD \
    psql -h localhost -p 5432 \
    -U $POSTGRES_USER -c "create database $POSTGRES_DB;"

  make api-test || {
    err 'api unit test failed'
    return 1
  }
}

api-golangci-lint() {

  export GOFLAGS="-mod=vendor"
  make api-lint || {
    err 'go lint failed'
    return 1
  }
}

yaml-lint() {

  make yaml-lint  || {
    err 'yaml lint failed'
    return 1
  }
}

goa-gen() {
  go version

  export GOFLAGS=""
  go install goa.design/goa/v3/cmd/goa@v3.19.1

  goa version

  make goa-gen  || {
      err 'goa gen failed'
      return 1
  }

  files=$(git diff api | wc -l)
	if [[ ${files} == 0 ]];then
    echo "    Git repo is clean."
  else
    echo "---------------------------------------"
    echo "  🔴 Files are changed!! Please run 'goa gen' command."
    return 1
  fi
}

api-build() {
  make api-build || {
    err 'Api build failed'
    return 1
  }
}

ui-build() {
  install-node

  make ui-build || {
    err 'UI build failed'
    return 1
  }
}

swagger-build() {
  install-node

  make swagger-build || {
    err 'Swagger UI build failed'
    return 1
  }
}

api-e2e() {
  info Runnning Hub CLI E2E tests

  go mod vendor
  go build -o tkn-hub github.com/tektoncd/hub/api/cmd/tkn-hub

  export TEST_CLIENT_BINARY="${PWD}/tkn-hub"

  go test -count=1 -tags=e2e ./test/... || {
    err 'api e2e test failed'
    return 1
  }

}

detect_api_related_changes() {
  [[ ! -z $(detect_changes "api") ]] || [[ ! -z $(detect_changes "go.*") ]] && echo "changes detected related to API"
}

### presubmit hooks ###

build_tests() {
  # run in a subshell so that path and shell options -eu -o pipefail will
  # will remain the same when it exits
  [[ ! -z ${TEST_RUN_ALL_TESTS} ]] || [[ ! -z $(detect_api_related_changes) ]] && {
    (
      set -eu -o pipefail

      api-build
    ) || exit 1
  } || echo "No changes detected related to API"

  [[ ! -z ${TEST_RUN_ALL_TESTS} ]] || [[ ! -z $(detect_changes "ui") ]] && {
    (
      set -eu -o pipefail

      ui-build
    ) || exit 1
  } || echo "No changes detected related to UI"

    [[ ! -z ${TEST_RUN_ALL_TESTS} ]] || [[ ! -z $(detect_changes "swagger") ]] && {
    (
      set -eu -o pipefail

      swagger-build
    ) || exit 1
  } || echo "No changes detected related to UI"
}

unit_tests() {
  # run in a subshell so that path and shell options -eu -o pipefail will
  # will remain the same when it exits
  [[ ! -z ${TEST_RUN_ALL_TESTS} ]] || [[ ! -z $(detect_api_related_changes) ]] && {
    (
      set -eu -o pipefail

      goa-gen || return 1
    ) || exit 1

    (
      set -eu -o pipefail

      api-unittest || return 1
      api-golangci-lint || return 1
    ) || exit 1
  } || echo "No changes detected related to API"

  [[ ! -z ${TEST_RUN_ALL_TESTS} ]] || [[ ! -z $(detect_changes "ui") ]] && {
    (
      set -eu -o pipefail

      ui-unittest || return 1
    ) || exit 1
  } || echo "No changes detected related to UI"
  (
    set -eu -o pipefail

    yaml-lint || return 1
  ) || exit 1

}

integration_tests() {
  (
    set -eu -o pipefail

    cd "$API_DIR"
    api-e2e || return 1
  ) || exit 1

}

main $@
