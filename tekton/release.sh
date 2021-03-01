#!/usr/bin/env bash
set -u -e -o pipefail

HUB_REPO="https://github.com/tektoncd/hub"
UPSTREAM_REMOTE="upstream"
BRANCH="main"
IMAGE_REGISTRY="quay.io/tekton-hub"

RELEASE_VERSION=""
HUB_NAMESPACE="tekton-hub"
HUB_CI_NAMESPACE="tekton-hub-ci"
CLUSTER=""

declare -a BINARIES=(
  kubectl
  git
)

checkPrerequisites() {
  for b in "${BINARIES[@]}"; do
    type -p "${b}" >/dev/null || {
      echo "'${b}' need to be available"
      exit 1
    }
  done

  kubectl version 2>/dev/null >/dev/null || {
    echo "you need to have access to a kubernetes cluster"
    exit 1
  }

  kubectl get pipelineresource 2>/dev/null >/dev/null || {
    echo "you need to have tekton install onto the cluster"
    exit 1
  }
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

whichCluster() {
  read -r -e -p "Are you deploying on OpenShift (Y/n): " YESORNO
  if [ "${YESORNO}" == "Y" ] || [ "${YESORNO}" == "y" ]; then
    CLUSTER='openshift'
    type -p oc >/dev/null || {
      echo "'oc' need to be available"
      exit 1
    }
  elif [ "${YESORNO}" == "N" ] || [ "${YESORNO}" == "n" ]; then
    CLUSTER='kubernetes'
  else
    echo 'invalid input'
    exit 1
  fi
}

createGitTag() {
  cd "${GOPATH}"/src/github.com/tektoncd/hub

  [[ -n $(git status --porcelain 2>&1) ]] && {
    echo "We have detected some changes in your repo"
    echo "Stash them before executing this script"
    exit 1
  }

  git checkout ${BRANCH}
  git reset --hard ${UPSTREAM_REMOTE}/${BRANCH}

  echo; echo 'Creating tag for new release:  '
  read -r -e -p "Enter tag message: " TAG_MESSAGE
  git tag -a "${RELEASE_VERSION}" -m "${TAG_MESSAGE}"
  git push ${UPSTREAM_REMOTE} --tags
}

createHubSecretAndCM() {
  kubectl create namespace ${HUB_NAMESPACE} 2>/dev/null || true

  kubectl -n ${HUB_NAMESPACE} get secret db 2>/dev/null >/dev/null || {

    echo; echo "Database Configurations:"
    read -r -e -p "Enter DB Name: " DB_NAME
    read -r -e -p "Enter DB Username: " DB_USERNAME
    read -r -e -p "Enter DB Password: " DB_PASSWORD

    kubectl -n ${HUB_NAMESPACE} create secret generic db \
      --from-literal=POSTGRES_DB="${DB_NAME}" \
      --from-literal=POSTGRES_USER="${DB_USERNAME}" \
      --from-literal=POSTGRES_PASSWORD="${DB_PASSWORD}" \
      --from-literal=POSTGRES_PORT="5432"

    kubectl -n ${HUB_NAMESPACE} label secret db app=db
    echo;
  }

  kubectl -n ${HUB_NAMESPACE} get secret api 2>/dev/null >/dev/null || {
    echo "API Configurations:"
    read -r -e -p "Enter GitHub OAuth Client ID: " GH_CLIENT_ID
    read -r -e -p "Enter GitHub OAuth Client Secret: " GH_CLIENT_SECRET
    read -r -e -p "Enter JWT Signing key: " JWT_SIGNING_KEY
    read -r -e -p "Enter the Access JWT expire time: (eg. 1d) " ACCESS_JWT_EXPIRES_IN
    read -r -e -p "Enter the Refresh JWT expire time: (eg. 1d) " REFRESH_JWT_EXPIRES_IN

    kubectl -n ${HUB_NAMESPACE} create secret generic api \
      --from-literal=GH_CLIENT_ID="${GH_CLIENT_ID}" \
      --from-literal=GH_CLIENT_SECRET="${GH_CLIENT_SECRET}" \
      --from-literal=JWT_SIGNING_KEY="${JWT_SIGNING_KEY}" \
      --from-literal=ACCESS_JWT_EXPIRES_IN="${ACCESS_JWT_EXPIRES_IN}" \
      --from-literal=REFRESH_JWT_EXPIRES_IN="${REFRESH_JWT_EXPIRES_IN}"

    kubectl -n ${HUB_NAMESPACE} label secret api app=api

    kubectl -n ${HUB_NAMESPACE} create cm ui \
      --from-literal=GH_CLIENT_ID="${GH_CLIENT_ID}" \
      --from-literal=API_URL="https://api.hub.tekton.dev" \
      --from-literal=API_VERSION="v1"

    kubectl -n ${HUB_NAMESPACE} label cm ui app=ui
    echo;
  }

  kubectl -n ${HUB_NAMESPACE} get cm api 2>/dev/null >/dev/null || {
    echo "Hub Config File:"
    read -r -e -p "Enter Raw URL of the hub config file (Default: https://raw.githubusercontent.com/tektoncd/hub/main/config.yaml): " HUB_CONFIG

    if [ -z "$HUB_CONFIG" ]; then
      HUB_CONFIG=https://raw.githubusercontent.com/tektoncd/hub/main/config.yaml
    fi

    kubectl -n ${HUB_NAMESPACE} create cm api \
      --from-literal=CONFIG_FILE_URL="${HUB_CONFIG}"

    kubectl -n ${HUB_NAMESPACE} label cm api app=api
    echo;
  }
}

createRegistrySecret() {
  kubectl create namespace ${HUB_CI_NAMESPACE} 2>/dev/null || true

  kubectl -n ${HUB_CI_NAMESPACE} delete secret registry-sec --ignore-not-found
  kubectl -n ${HUB_CI_NAMESPACE} get secret registry-sec 2>/dev/null >/dev/null || {

    echo; echo "Enter Quay registry credentials to push the images: (quay.io/tekton-hub) "
    read -r -e -p "Enter Username: " USERNAME
    read -r -e -sp "Enter Password: " PASSWORD

    kubectl -n ${HUB_CI_NAMESPACE} create secret generic registry-sec \
      --type="kubernetes.io/basic-auth" \
      --from-literal=username="${USERNAME}" \
      --from-literal=password="${PASSWORD}"

    kubectl -n ${HUB_CI_NAMESPACE} annotate secret registry-sec tekton.dev/docker-0=quay.io
  }
}

createNecessaryRoles() {

  echo; echo 'Creates service account and necessary role to create resources: '

  kubectl -n ${HUB_CI_NAMESPACE} delete serviceaccount registry-login --ignore-not-found
  cat <<EOF | kubectl -n ${HUB_CI_NAMESPACE} create -f-
apiVersion: v1
kind: ServiceAccount
metadata:
  name: registry-login
secrets:
  - name: registry-sec
EOF

  kubectl -n ${HUB_NAMESPACE} delete role hub-pipeline --ignore-not-found
  kubectl -n ${HUB_NAMESPACE} delete rolebinding hub-pipeline --ignore-not-found
  kubectl -n ${HUB_NAMESPACE} create role hub-pipeline \
    --resource=deployment,services,pvc,job \
    --verb=create,get,list,delete,patch
  kubectl -n ${HUB_NAMESPACE} create rolebinding hub-pipeline \
    --serviceaccount=${HUB_CI_NAMESPACE}:registry-login \
    --role=hub-pipeline

  if [ "${CLUSTER}" == "openshift" ]; then
    oc adm policy add-scc-to-user privileged system:serviceaccount:${HUB_CI_NAMESPACE}:registry-login

    kubectl -n ${HUB_NAMESPACE} delete role hub-pipeline-route --ignore-not-found
    kubectl -n ${HUB_NAMESPACE} delete rolebinding hub-pipeline-route --ignore-not-found
    kubectl -n ${HUB_NAMESPACE} create role hub-pipeline-route \
      --resource=route \
      --verb=create,get,list,delete,patch
    kubectl -n ${HUB_NAMESPACE} create rolebinding hub-pipeline-route \
    --serviceaccount=${HUB_CI_NAMESPACE}:registry-login \
    --role=hub-pipeline-route
  fi

  echo;
}

startPipelines() {
  echo 'Install Tasks: '
  kubectl -n ${HUB_CI_NAMESPACE} apply -f ./tekton/api/golang-db-test.yaml

  echo; echo 'Install Pipelines: '

  kubectl -n ${HUB_CI_NAMESPACE} apply -f ./tekton/api/pipeline.yaml
  kubectl -n ${HUB_CI_NAMESPACE} apply -f ./tekton/ui/pipeline.yaml

  echo; echo 'Start Pipelines: '

  cat <<EOF | kubectl -n ${HUB_CI_NAMESPACE} create -f-
apiVersion: tekton.dev/v1beta1
kind: PipelineRun
metadata:
  generateName: hub-api-
spec:
  serviceAccountName: registry-login
  pipelineRef:
    name: api-deploy
  params:
    - name: HUB_REPO
      value: ${HUB_REPO}
    - name: REVISION
      value: ${BRANCH}
    - name: API_IMAGE
      value: ${IMAGE_REGISTRY}/api
    - name: DB_MIGRATION_IMAGE
      value: ${IMAGE_REGISTRY}/db-migration
    - name: TAG
      value: ${RELEASE_VERSION}
    - name: HUB_NAMESPACE
      value: ${HUB_NAMESPACE}
    - name: K8S_VARIANT #it will accept either openshift or kubernetes
      value: ${CLUSTER}
  workspaces:
    - name: shared-workspace
      volumeClaimTemplate:
        spec:
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 10Gi
EOF

  cat <<EOF | kubectl -n ${HUB_CI_NAMESPACE} create -f-
apiVersion: tekton.dev/v1beta1
kind: PipelineRun
metadata:
  generateName: hub-ui-
spec:
  serviceAccountName: registry-login
  pipelineRef:
    name: ui-pipeline
  params:
    - name: HUB_REPO
      value: ${HUB_REPO}
    - name: REVISION
      value: ${BRANCH}
    - name: IMAGE
      value: ${IMAGE_REGISTRY}/ui
    - name: TAG
      value: ${RELEASE_VERSION}
    - name: HUB_NAMESPACE
      value: ${HUB_NAMESPACE}
    - name: K8S_VARIANT #it will accept either openshift or kubernetes
      value: ${CLUSTER}
  workspaces:
    - name: shared-workspace
      volumeClaimTemplate:
        spec:
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 10Gi
EOF
}

main() {

  checkPrerequisites

  echo 'Tekton Hub Release: '

  # ask for release version to create tag
  getReleaseVersion

  # OpenShift/Kubernetes
  whichCluster

  # creates a new git tag with release version and push to upstream
  createGitTag

  # creates hub supporting secrets and config map
  createHubSecretAndCM

  # create registry secret to push image to quay.io/tekton-hub
  createRegistrySecret

  # creates required role for pipeline service account to crud resources in HUB_NAMESPACE
  createNecessaryRoles

  # installs and starts pipelines
  startPipelines

  return 0
}

main "$@"
