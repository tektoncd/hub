#!/usr/bin/env bash

# Copyright © 2022 The Tekton Authors.
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

RELEASE_VERSION=""

declare -r SCRIPT_PATH=$(readlink -f "$0")
declare -r SCRIPT_DIR=$(cd $(dirname "$SCRIPT_PATH") && pwd)
declare -r API_DIR="$SCRIPT_DIR/api"
declare -r UI_DIR="$SCRIPT_DIR/ui"
declare -r RELEASE_DIR="$SCRIPT_DIR/release"
DOCKER_CMD=${DOCKER_CMD:-docker}
REGISTRY_BASE_URL=${REGISTRY_BASE_URL:-quay.io/tekton-hub}

BINARIES="ko hub"

info() {
  echo "INFO: $@"
}

err() {
  echo "ERROR: $@"
}

getReleaseVersion() {
  [[ -z ${RELEASE_VERSION} ]] && {
    read -r -e -p "Enter a target release (i.e: v0.1.2): " RELEASE_VERSION
    [[ -z ${RELEASE_VERSION} ]] && {
      echo "no target release"
      exit 1
    }
  }
  [[ ${RELEASE_VERSION} =~ v[0-9]+\.[0-9]*\.[0-9]+ ]] || {
    echo "invalid version provided, need to match v\d+\.\d+\.\d+"
    exit 1
  }
}

buildDbMigrationImage() {
  info Building DB Migration Image
  echo -----------------------------------
  ${DOCKER_CMD} build -f images/db.Dockerfile -t ${REGISTRY_BASE_URL}/db-migration:${RELEASE_VERSION} . && ${DOCKER_CMD} push ${REGISTRY_BASE_URL}/db-migration:${RELEASE_VERSION}
  info DB Migration Image Build Successfully
  echo -----------------------------------
}

buildApiImage() {
  info Building API Image
  echo -----------------------------------
  ${DOCKER_CMD} build -f images/api.Dockerfile -t ${REGISTRY_BASE_URL}/api:${RELEASE_VERSION} . && ${DOCKER_CMD} push ${REGISTRY_BASE_URL}/api:${RELEASE_VERSION}
  info API Image Build Successfully
  echo -----------------------------------
}

buildUiImage() {
  info Building UI Image
  echo -----------------------------------
  ${DOCKER_CMD} build -f images/ui.Dockerfile -t ${REGISTRY_BASE_URL}/ui:${RELEASE_VERSION} . && ${DOCKER_CMD} push ${REGISTRY_BASE_URL}/ui:${RELEASE_VERSION}
  info UI Image Build Successfully
  echo -----------------------------------
}

db(){
	info Creating DB Release Yaml

  ko resolve -f 00-init  > "${RELEASE_DIR}"/db.yaml || {
    err 'db release build failed'
    return 1
  }
  echo "-----------------------------------------"
}

db-migration(){
	info Creating Db-Migration Release Yaml

  ko resolve -f 01-db  > "${RELEASE_DIR}"/db-migration.yaml || {
    err 'db-migration release build failed'
    return 1
  }
  echo "------------------------------------------"
}

api-k8s(){
	info Creating API Release Yaml

  ko resolve -f 02-api  > "${RELEASE_DIR}"/api-k8s.yaml || {
    err 'api release build failed'
    return 1
  }
  echo "------------------------------------------"
}

api-openshift(){
	info Creating API Release Yaml

  ko resolve -f 02-api -f 04-openshift/40-api-route.yaml -f 04-openshift/40-auth-route.yaml > "${RELEASE_DIR}"/api-openshift.yaml || {
    err 'api release build failed'
    return 1
  }
  echo "------------------------------------------"
}

ui-k8s(){
	info Creating UI Release Yaml

  ko resolve -f 03-ui > "${RELEASE_DIR}"/ui-k8s.yaml || {
    err 'ui release build failed'
    return 1
  }
  echo "------------------------------------------"
}

ui-openshift(){
	info Creating UI Release Yaml

  ko resolve -f 03-ui -f 04-openshift/41-ui-route.yaml > "${RELEASE_DIR}"/ui-openshift.yaml || {
    err 'ui release build failed'
    return 1
  }
  echo "------------------------------------------"
}

replaceImageName() {
  info Changing Image Name

  cd ${RELEASE_DIR}
  #  Replace the db-migration image name
  sed -i "s@image: quay.io/tekton-hub/db-migration@image: ${REGISTRY_BASE_URL}/db-migration:$RELEASE_VERSION@g" ${RELEASE_DIR}/db-migration.yaml

  # Replace the api image
  sed -i "s@image: quay.io/tekton-hub/api@image: ${REGISTRY_BASE_URL}/api:$RELEASE_VERSION@g" ${RELEASE_DIR}/api-k8s.yaml

  sed -i "s@image: quay.io/tekton-hub/api@image: ${REGISTRY_BASE_URL}/api:$RELEASE_VERSION@g" ${RELEASE_DIR}/api-openshift.yaml

  #Replace the ui image
  sed -i "s@image: quay.io/tekton-hub/ui@image: ${REGISTRY_BASE_URL}/ui:$RELEASE_VERSION@g" ${RELEASE_DIR}/ui-k8s.yaml

  sed -i "s@image: quay.io/tekton-hub/ui@image: ${REGISTRY_BASE_URL}/ui:$RELEASE_VERSION@g" ${RELEASE_DIR}/ui-openshift.yaml
}

createGitTag() {
  echo; echo 'Creating tag for new release:'

  hub release create --draft --prerelease -a db.yaml \
   -a db-migration.yaml \
   -a  api-k8s.yaml \
   -a api-openshift.yaml \
   -a  ui-k8s.yaml \
   -a ui-openshift.yaml \
   -m "${RELEASE_VERSION}" "${RELEASE_VERSION}"
}

main() {

  # Check if all required command exists
  for b in ${BINARIES};do
      type -p ${b} >/dev/null || { echo "'${b}' need to be avail"; exit 1 ;}
  done

  # Ask the release version to build images
  getReleaseVersion

  # Generate the release yamls for db, db-migration, api and ui
  echo "********************************************"
  info     Generate the Release Yamls for Hub
  echo "********************************************"
  cd config
  db
  db-migration
  api-k8s
  api-openshift
  ui-k8s
  ui-openshift

  # Build images for db-migration, api and ui
  echo "********************************************"
  info        Build the Images for Hub
  echo "********************************************"
  buildDbMigrationImage
  buildApiImage
  buildUiImage

  # Change the image name with the release version specified
  echo "********************************************"
  info      Replace the Images with New Version
  echo "********************************************"
  replaceImageName

  echo "********************************************"
  info            Create Git Tag
  echo "********************************************"
  createGitTag

  echo "********************************************"
  echo "***" Release Created for Hub successfully "***"
  echo "********************************************"
}

main $@
