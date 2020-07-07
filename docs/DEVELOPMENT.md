# Development <!-- TOC omit:true -->

[toc]

## Getting Started

1. [Create a GitHub Account][join-github]
1. [Setup GitHub access via SSH][gh-ssh]

## Checkout your fork

To check out this repository:

1. Create your own [fork of this repository][fork-repo]
2. Clone it to your machine:

```shell
  git clone git@github.com:${YOUR_GITHUB_USERNAME}/hub.git
  cd hub

  git remote add upstream git@github.com:tektoncd/hub.git
  # prevent accidental push to upstream
  git remote set-url --push upstream no-push
  git fetch --all
```

Adding the upstream remote sets you up nicely for regularly [syncing your fork][sync-fork].

## Requirements

You must install these tools:

1. [`go`][install-go]: The language hub apis are built in.
2. [`goa`][install-goa]: Framework for api service
3. [`git`][install-git]: For source control

You may need to install more tools depending on the way you want to run the hub.


## Running On Local Machine

For running on local machine you need postgresql database.

Two ways to run postgresql database:
- [Install postgresql][install-pg]
- Run a postgresql container using [docker][install-docker] / [podman][install-podman]


### Running database

- If you have installed postgresql locally, you need to create a `hub` database.

  **NOTE:** Use the same configuration mentioned in `.env.dev` or
  update `.env.dev` with the configuration you used. The api service
  and db migration uses the db configuration from `.env.dev`.

- If you want to run a postgres container, source the `.env.dev` so that
  `docker` can use the same database configuration as in `.env.dev` to create a container.

  Ensure you are in `hub/api` directory.

  ```bash
  source .env.dev

  docker run -d  --name hub \
    -e POSTGRES_USER=$POSTGRESQL_USER \
    -e POSTGRES_PASSWORD=$POSTGRESQL_PASSWORD \
    -e POSTGRES_DB=$POSTGRESQL_DATABASE \
    -p $POSTGRESQL_PORT:5432 \
    postgres
  ```


### Populating Database

Once the database is up and running, you can run migration to create tables and populate the database.

Run the following command to run migration

```bash
go run ./cmd/db
```

Wait until the migration completes and logs to show

  > DB initialisation successful !!


### Running API Service

Once the database is setup and the migration has been run, you can run api service by

```bash
go run ./cmd/hub
```

### Running tests

To run the tests, we need a test db.

- If you have installed postgresql, create a `hub_test` database.
- If you are running a container, create `hub_test` database in the same container.

  ```bash
  source .env.dev

  docker exec -it hub bash -c \
    "PGPASSWORD=$POSTGRESQL_PASSWORD \
     psql -h localhost -p 5432 -U postgres -c 'create database hub_test;'"
  ```

Once the `hub_test` database is created, you can run the test using following command:

```bash
  go test  -count=1 -v ./...
```

**NOTE:** `tests` use the database configurations from [test/config/env.test][env-test-file]

---

## Running On Kubernetes Cluster

### Prerequisites:

- Kubernetes Cluster:
  - You can run a cluster on your local machine using [CRC][install-crc], [Minikube][install-minikube] or [kind][install-kind]

- Depending on the cluster you have, install the command line tool:

  - [`kubectl`][install-kubectl]: For interacting with cluster.
  - [`oc`][install-oc]: For interacting with OpenShift cluster.
  - [`ko`][install-ko]: tool for building and deploying Golang applications to Kubernetes.

### Deploying API and DB Service

Export `KO_DOCKER_REPO` for `ko` to publish image to. E.g.

```bash
export KO_DOCKER_REPO=quay.io/<username>
```

Log into the registry used in `KO_DOCKER_REPO` so that `ko` gets
access to push the image to the registry.

Make sure you are logged into the cluster and you are in the
`hub/api` directory before running the following command

```bash
ko apply -f config/
```

The command above will create an image and push it to registry pointed by  `KO_DOCKER_REPO` and then deploy the `api` and `db` container on the cluster.

**NOTE:** Ensure that the image is **publicly** available otherwise deployment will fail. The pod status will give `Error: ImagePullBackOff`

Watch the pods until both of them are running.
```
kubectl get pods -o wide -w
```

To create tables and to populate, we need to run the migration. Run the following command:

```
ko apply -f config/db-migration/14-db-migration.yaml
```

The command above will create a migration image, push to registry pointed by `KO_DOCKER_REPO` and run a `job` that performs the migration.

Check the logs using

```bash
kubectl logs job/db-migration
```

Wait until the migration log shows

> DB initialisation successful !!

### Expose API service

To expose the api service to be accessible outside cluster, run the following command based the `k8s` cluster you are using:

#### OpenShift

```bash
oc apply -f config/openshift/
```
This will expose the api service as a `NodePort` service and create a `route` that can be used to access the api.

you can get the route to api service using

```bash
oc get routes api --template='https://{{ .spec.host }} '
```

#### GKE

```bash
kubectl apply -f  config/gke/21-api-service.yaml
```

This will expose the deployment as `LoadBalancer` service. You can access the api service using `LoadBalancer Ingress` and `Port`.

```bash
kubectl describe svc api-load
```
Use the `LoadBalancerIngress:Port` eg. `http://35.223.169.99:8000`


#### Test API Service

Run `curl http://<ip>:port/categories` api to make sure
everything is working fine.


#### Redeploy API Service

**To redeploy api service after changes run the following command:**

```bash
ko apply -f config/20-api-deployment.yaml
```

Above command will create a new image from updated code, push it to registry and update the deployment.

**Any changes to database should be run as a migration, Once you write a migration apply the below command:**
```
ko apply -f config/db-migration/14-db-migration.yaml
```
Above command will create a new image of db migration, push it to registry and run a job.


[join-github]:https://github.com/join
[gh-ssh]:https://help.github.com/articles/connecting-to-github-with-ssh/
[fork-repo]:https://help.github.com/articles/fork-a-repo/
[sync-fork]:https://help.github.com/articles/syncing-a-fork/
[install-go]:https://golang.org/doc/install
[install-goa]:https://github.com/goadesign/goa
[install-git]:https://help.github.com/articles/set-up-git/
[install-pg]: https://www.postgresql.org/docs/12/tutorial-install.html
[install-docker]: https://docs.docker.com/engine/install/
[install-podman]: https://podman.io/getting-started/installation.html
[install-crc]:https://cloud.redhat.com/openshift/install/crc/installer-provisioned
[install-minikube]:https://kubernetes.io/docs/tasks/tools/install-minikube/
[install-kubectl]:https://kubernetes.io/docs/tasks/tools/install-kubectl/
[install-oc]:https://docs.openshift.com/container-platform/4.2/cli_reference/openshift_cli/getting-started-cli.html
[install-ko]:https://github.com/google/ko
[install-kind]:https://kind.sigs.k8s.io/docs/user/quick-start/
[env-test-file]: https://github.com/tektoncd/hub/blob/master/api/test/config/env.test