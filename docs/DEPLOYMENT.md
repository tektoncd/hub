# Deployment

- [Prerequisites](#prerequisites)
- [Build Images](#build-images)
- [Deploy Database](#deploy-database)
- [Run Db migration](run-db-migration)
- [Deploy API Service](#deploy-api-service)
  - [Update API Secret](#update-api-secret)
  - [Update API ConfigMap](#update-api-configmap)
  - [Update API Image](#update-api-image)
  - [Setup Ingress/Route](#setup-route-or-ingress)
- [Deploy UI](#deploy-ui)
  - [Update UI ConfigMap](#update-ui-configmap)
  - [Update UI Image](#update-ui-image)
  - [Setup Ingress/Route](#setup-route-or-ingress)
- [Update Catalog and Catalog Refresh](#update-catalog-and-catalog-refresh)
- [Schedule Catalog Refresh On an Interval](#schedule-catalog-refresh-on-an-interval)
- [Setup Catalog Refresh CronJob](#setup-catalog-refresh-cronjob)
- [Adding New Users in Config](#adding-new-users-in-config)
- [Deploying on Disconnected Cluster](#deploying-on-disconnected)
- [Troubleshooting](#troubleshooting)

## Prerequisites

- Kubernetes Cluster
  - You can also deploy on [Minkube][minikube], [kind][kind] or [OpenShift][openshift].
- [kubectl][kubectl]
- [docker][docker] or [podman][podman]

## Build Images

Lets build and push all images to a registry and then create deployments.

**Note** : You can use images from quay and skip the build image section.

- The image tag would be same as release tag. You can check the latest release [here](https://github.com/tektoncd/hub/releases)

- For example:

  - UI Image: quay.io/tekton-hub/ui:v1.3.0

  - DB Migration Image: quay.io/tekton-hub/db-migration:v1.3.0

  - API Image: quay.io/tekton-hub/api:v1.3.0

### UI Image

Ensure you are in the root of project and logged in you image registry.

```
docker build -f images/ui.Dockerfile -t <image> . && docker push <image>
```

Replace `<image>` with the registry and image name.
eg. `quay.io/<username>/ui`

### API & DB Migration Image

Ensure you are in the root of project and logged in you image registry.

Build the API image using below command

```
docker build -f images/api.Dockerfile -t <image> . && docker push <image>
```

Replace `<image>` with the registry and image name.
eg. `quay.io/<username>/api`

Now, Build the Db migration image using below command

```
docker build -f images/db.Dockerfile -t <image> . && docker push <image>
```

Replace `<image>` with the registry and image name.
eg. `quay.io/<username>/db-migration`

Make sure all images are public before creating deployments.

---

> ### NOTE: If You have deployed Hub v1.6.0 and want to upgrade to Hub v1.7.0 then please follow [this][upgrade-v1.6.0-to-v1.7.0] doc for next steps. As we have renamed config metadata name and labels in the configs of db, api and ui.

---

---

> ### NOTE: In case of using Github or Gitlab Enterprise make sure you deploy Hub in the same network where the Enterprise server is running. Example if your Github/Gitlab Enterprise is running behind VPN then your Kubernetes Cluster should also be behind VPN.

---

## Deploy Database

Now, we have all images pushed to registry. Lets deploy the database first.

Navigate to `config` directoy in the root of the project.

Execute below command

```
kubectl apply -f 00-init/ -n tekton-hub
```

This will create `tekton-hub` namespace, db deployment with `postgres` image, secret to save db credentials, pvc for the db and a ClusterIP service to expose the db internally.

All resources are created in `tekton-hub` namespace.

Wait till the pod comes in a running state

```
kubectl get pod -n tekton-hub -w
```

## Run DB migration

Once the pod is in running state now we can run db migration. This will create all the tables required in the database.

Edit the `01-db/10-db-migration.yaml` and replace the image you created previously and apply the yaml.

```
kubectl apply -f 01-db/10-db-migration.yaml -n tekton-hub
```

This will create a job which will read the db credentials from the secret created while deploying database and create all required tables.

```bash
$ kubectl get pods -n tekton-hub
NAME                   READY   STATUS      RESTARTS   AGE
tekton-hub-db-589d44fdd5-ksf8v    1/1     Running     0          50s
tekton-hub-db-migration-8vhpd     0/1     Completed   0          17s
```

The `Pod` will go into `Completed` state if the job succeeds. You can also check the logs using `kubectl logs -n tekton-hub tekton-hub-db-migration-8vhpd`. In case of any failure the `Pod` will go into `Error` state.

```bash
$ kubectl get pods -n tekton-hub
NAME                   READY   STATUS      RESTARTS   AGE
tekton-hub-db-589d44fdd5-ksf8v    1/1     Running     0          50s
tekton-hub-db-migration-8vhpd     0/1     Error       0          17s
```

## Deploy API Service

### Setup Route or Ingress

- If deploying on OpenShift:-

  ```bash
  kubectl apply -f 04-openshift/40-api-route.yaml -f 04-openshift/40-auth-route.yaml -n tekton-hub
  ```

- If deploying on Kubernetes:-

  - Create the secret containing tls cert and tls key both for API and Oauth server. Example as follows:

    ```bash
    kubectl create secret tls api-hub-tekton-dev-tls --cert=path/to/cert/file --key=path/to/key/file -n tekton-hub
    ```

  - Apply the Ingress

    ```bash
    kubectl apply -f 04-kubernetes/40-api-ingress.yaml -f 04-kubernetes/40-auth-ingress.yaml -n tekton-hub
    ```

### Create Git Oauth Applications

- **Github** - Create a GitHub OAuth with `Homepage URL` and `Authorization callback URL` as `<auth-route>`. Follow the steps given [here][github-oauth-steps] to create a GitHub OAuth.
- **Gitlab** - Create a Gitlab Oauth with `REDIRECT_URI` as `<auth-route>/auth/gitlab/callback`. Follow the steps given [here][gitlab-oauth-steps] to create a Gitlab Oauth
- **BitBucket** - Create a BitBucket Oauth with `Callback URL` as `<auth-route>`. Follow the steps given [here][bitbucket-oauth-steps] to create a BitBucket Oauth

### Update API Secret

Edit `02-api/20-api-secret.yaml` and update the configuration

- After creating the OAuth add the Client ID and Client Secret in the yaml file.
- For JWT_SIGNING_KEY, you can add any random string, this is used to sign the JWT created for users.
- For `ACCESS_JWT_EXPIRES_IN` and `REFRESH_JWT_EXPIRES_IN` add time you want the jwt to be expired in. Refresh time should be greater than Access time.
  eg. 1m = 1 minute, 1h = 1 hour, 1d = 1 day
  - **NOTE**: Supported formats for `ACCESS_JWT_EXPIRES_IN` and `REFRESH_JWT_EXPIRES_IN` are w(weeks), d(days), h(hours), m(min) and s(sec)

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: tekton-hub-api
  namespace: tekton-hub
type: Opaque
stringData:
  GH_CLIENT_ID: Github Oauth client id
  GH_CLIENT_SECRET: Github Oauth secret
  GL_CLIENT_ID: Gitlab Oauth client id
  GL_CLIENT_SECRET: Gitlab Oauth secret
  BB_CLIENT_ID: BitBucket Oauth client id
  BB_CLIENT_SECRET: BitBucket Oauth secret
  JWT_SIGNING_KEY: a-long-signing-key
  ACCESS_JWT_EXPIRES_IN: time such as 15m
  REFRESH_JWT_EXPIRES_IN: time such as 15m
  AUTH_BASE_URL: auth route
  GHE_URL: Add Github Enterprise URL in case of authenticating through Github Enterprise (Example (https|http)://myghe.com) --> Do not provide the catalog URL
  GLE_URL: Add Gitlab Enterprise URL in case of authenticating through Gitlab Enterprise (Example (https|http)://mygle.com) --> Do not provide the catalog URL
```

### Update API ConfigMap

- Update the [config.yaml][config-yaml] to add at least one user with all scopes such as

  - refresh a catalog,
  - refresh config file
  - create an agent token

- For example

```
scopes:
  - name: agent:create
    users: [foo]        <<< Where `foo` is your Github Handle
  - name: catalog:refresh
    users: [foo]
  - name: config:refresh
    users: [foo]
```

- Commit the changes and push the changes to your fork

- Edit the `02-api/21-api-configmap.yaml` and change the URL to point to your fork. By default it is pointing to [config.yaml][config-yaml] in Hub which has Application Data.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: tekton-hub-api
  labels:
    app: tekton-hub-api
data:
  CONFIG_FILE_URL: https://raw.githubusercontent.com/tektoncd/hub/master/config.yaml ## Change the file URL here
  CATALOG_REFRESH_INTERVAL: 30m ## change the catalog refresh interval here
```

All users have default scopes `rating:read` and `rating:write`. Once user get login to Hub, they get the appropriate scopes which allows them to rate the resource

### Schedule Catalog Refresh On an Interval

- Default Catalog refresh interval time is 30 minutes. After every 30 minutes catalog resources in the hub db would be refreshed. If you want to change catalog refresh interval, update the `CATALOG_REFRESH_INTERVAL` in `tekton-hub-api` configmap with any one of these supported time units
As(A seconds), Bm(B minutes) Ch(C hours), Dd(D days) and Ew(E weeks).

- If you want to skip automate catalog refresh, keep `CATALOG_REFRESH_INTERVAL` as empty or `0`. Negative value is also not supported.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: tekton-hub-api
  labels:
    app: tekton-hub-api
data:
  CONFIG_FILE_URL: https://raw.githubusercontent.com/tektoncd/hub/master/config.yaml
  CATALOG_REFRESH_INTERVAL: 30m ## change the catalog refresh interval here
```

**WARN** : Make sure you have updated Hub config before starting the api server

### Create SSH secrets (Optional)

In order to clone private repositories or repositories from private git instances,
we need to mount the ssh credentials so that cloning can happen. To make this happen
please follow the below steps:

1. Create the ssh keys if not already present in `~/.ssh` dir.
1. Create kubernetes secret which will contain the ssh keys using the following command
   ```sh
   kubectl create secret generic tekton-hub-api-ssh-crds --from-file=id_rsa="/path/to/id_rsa" --from-file=id_rsa.pub="/path/to/id_rsa.pub" --from-file=known_hosts="/path/to/known_hosts"
   ```
   Please make sure that secrets are created with the name `tekton-hub-api-ssh-crds`

### Update API deployment (Optional)

By default the catalog is cloned in `$HOME/catalog`. If in case you want to change the clone path, edit the `02-api/22-api-deployment.yaml`
with the following

```yaml
spec:
  ...
    spec:
      ...
      containers:
        - ...
          volumeMounts:
          - name: catalog-source
            mountPath: "/home/hub/catalog" <---- change here with path/where/catalog/to/be/cloned eg /tmp/catalog
          ...
          env:
          - name: CLONE_BASE_PATH
            value: <path/where/catalog/to/be/cloned> eg /tmp/catalog
```

### Update API Image

Edit the `02-api/22-api-deployment.yaml` and replace the image with the one created previously and executed below command

```bash
kubectl apply -f 02-api/ -n tekton-hub
```

This will create the pvc, deployment, secret, configmap and a NodePort service to expose the API server.

#### Verify if api route is accessible

For `OpenShift`:-

```bash
curl -k -X GET $(oc get -n tekton-hub routes api --template='https://{{ .spec.host }}/v1/categories')
```

## Deploy UI

### Setup Route or Ingress

After the deployment is done successfully, we need to expose the URL to access the UI.

- If deploying on OpenShift:-

  ```bash
  kubectl apply -f 04-openshift/41-ui-route.yaml -n tekton-hub
  ```

- If deploying on Kubernetes:-

  - Create the secret containing tls cert and tls key

    ```bash
    kubectl create secret tls ui-hub-tekton-dev-tls --cert=path/to/cert/file --key=path/to/key/file -n tekton-hub
    ```

  - Apply the Ingress

    ```bash
    kubectl apply -f 04-kubernetes/41-ui-ingress.yaml -n tekton-hub
    ```

### Update UI ConfigMap

Edit `config/03-ui/30-ui-configmap.yaml` and by setting the API URL, auth URL and Redirect URL.

You can get the API URL using below command (OpenShift)

```
  kubectl get -n tekton-hub routes tekton-hub-api --template='https://{{ .spec.host }}'
```

```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: tekton-hub-ui
data:
  API_URL: API URL   <<< Update this by API url
  API_VERSION: API VERSION  <<< Update this by API version (For e.g-"v1")
  AUTH_BASE_URL: AUTH URL << Update this by auth url
  REDIRECT_URI: UI URL << Update this by ui url
```

### Update UI Image

Edit the `03-ui/31-ui-deployment.yaml` by replacing the image with the one created previously and executed below command

```bash
kubectl apply -f 03-ui/ -n tekton-hub
```

This will create the deployment, configmap and a NodePort service to expose the UI.

Ensure pods are up and running

```bash
kubectl get pods -n tekton-hub
```

eg. wait till status of UI pod is running

```
NAME                   READY   STATUS      RESTARTS   AGE
tekton-hub-api-6dfc6f97d9-vk66r   1/1     Running     3          21m
tekton-hub-db-9bd4cdf99-zsq89     1/1     Running     0          21m
tekton-hub-db-migration-26ngs     0/1     Completed   0          21m
tekton-hub-ui-55fc66cc6b-69dsp    1/1     Running     0          58s
```

#### Verify if ui route is accessible

For `OpenShift`:-

```
    kubectl get routes -n tekton-hub tekton-hub-ui --template='https://{{ .spec.host }} '
```

Open the URL in a browser.

Note: On UI, you should be able to see the resources on homepage and Categories on the left side.

Update your GitHub OAuth created with the UI route in place of `Homepage URL` and `Authorization callback URL`.

## Update Catalog and Catalog Refresh

If you have added and modified catalog in the config file , then you will have to folllow following steps to see resources from catalog on Hub UI

1. Login through the Hub UI. Click on Login on right corner and then Sign in with GitHub/Github Enterprise.
2. Copy the Hub Token by clicking on the user profile which is at the right corner on the Home Page
3. Refresh the config file and you need to have `config:refresh` scope.

   - To Refresh config file

     ```
       curl -X POST -H "Authorization: <access-token>" \
       --header "Content-Type: application/json" \
       --data '{"force": true} \
       <api-route>/system/config/refresh
     ```

4. Call the Catalog Refresh API:

   - To refresh a catalog with name

     ```
     curl -X POST -H "Authorization: <jwt-token>" \
         <api-url>/catalog/<catalogName>/refresh
     ```

     Replace `<access-token>` with your Hub token copied from UI and replace `<api-url>` with your api pod url.

     This will give an output as below

     ```
     [{"id":1,"catalogName":"tekton","status":"queued"}]
     ```

   - To refresh all catalogs

     ```
     curl -X POST -H "Authorization: <jwt-token>" \
         <api-url>/catalog/refresh
     ```

     Replace `<access-token>` with your Hub token copied from UI and replace `<api-url>` with your api pod url.

     This will give an output as below

     ```
     [{"id":1,"catalogName":"tekton","status":"queued"}]
     ```

5. Refresh your UI, you will be able to see resources.

## Setup Catalog Refresh CronJob (Optional)

You can setup a cronjob which will refresh your db after an interval if there are any changes in your catalog.

NOTE: The catalog refresh will add new resources or update existing resource if it is updated in catalog but it doesn't delete a resource in db even if it is deleted in catalog.

1. You will need a JWT token with catalog refresh scope which can be used.

2. A User token is short lived so you can create an agent token which will not expire.

3. To create an agent, you can use the `/system/user/agent`. You will need `agent:create` scope in your JWT which has to be added in config.yaml

   ```
   curl -X PUT --header "Content-Type: application/json" \
       -H "Authorization: <access-token>" \
       --data '{"name":"catalog-refresh-agent","scopes": ["catalog:refresh"]}' \
       <api-route>/system/user/agent
   ```

   Replace `<access-token>` with your JWT, you can create an agent with any required scopes. In this case we need the agent with `catalog:refresh`. You can give any name to agent.

   This will gives an output as:

   ```
   {"token":"agent jwt token"}
   ```

   A token will be returned with the requested scopes. You can check using the jwt decoder. Use this token for the catalog refresh cron job.

4. Edit `05-catalog-refresh-cj/50-catalog-refresh-secret.yaml` and Set the Hub Token

   ```yaml
   apiVersion: v1
   kind: Secret
   metadata:
     name: catalog-refresh
   type: Opaque
   stringData:
     HUB_TOKEN: hub token
   ```

   Set the token created in last step as Hub token.

5. Apply the YAMLs

   ```
   kubectl apply -f 05-catalog-refresh-cj/ -n tekton-hub
   ```

   The cron job is configured to run every 30 min, you can change the interval in `05-catalog-refresh-cj/51-catalog-refresh-cronjob.yaml`.

## Adding New Users in Config

By default when any user login for the first time, they will have only default scope even if they are added in config.yaml. To get additinal scopes, make sure the user has logged in once.

Now, we need to refresh the config. To do that do a POST API call.

```
curl -X POST -H "Authorization: <access-token>" \
    --header "Content-Type: application/json" \
    --data '{"force": true}' \
  <api-route>/system/config/refresh
```

Replace `<access-token>` with your JWT token. you must have `config-refresh` scope to call this API.

## Deploying on Disconnected Cluster

- To deploy on a disconnected cluster you need to fork the [tektoncd/catalog][catalog] or you can use your own catalog but it must follow the [Catalog TEP][catalog-tep].
- Fork the tektoncd/hub and update the [hub config][hub-config] file for your catalog details.
  - Update the catalog details with your catalog
  ```
  catalogs:
    - name: tekton
      org: tektoncd
      type: community
      url: https://github.com/tektoncd/catalog
      revision: main
  ```
  You can find the template at bottom of [hub config][hub-config] file.
- Mirror the Hub Deployment Images to the registry disconnected cluster has access to.
  You can use release images or build your own. For hub release images, you can use latest release tag as image tag. for eg. `quay.io/tekton-hub/ui:v1.3.0`. You can check the latest release [here][hub-releases].
- List of images to be mirrored

  - UI - `quay.io/tekton-hub/ui`
  - API - `quay.io/tekton-hub/api`
  - DB-Migration - `quay.io/tekton-hub/db-migration`
  - DB - `postgres:13@sha256:260a98d976574b439712c35914fdcb840755233f79f3e27ea632543f78b7a21e`

  You can mirror the images using `oc` as below

  ```
  oc image mirror quay.io/tekton-hub/ui:v1.3.0 your.registry/project/ui:v1.3.0
  ```

- Now, follow the steps from [Deploy Database Section](#deploy-database) above. Make sure you update the images in deployment files before applying.

## Troubleshooting

#### UI is not showing resources but catalog refresh is successful

- If you have insecure connection the UI might not be able to hit the API.
- So, hit the API URL in your browser `https://<api-url>` once you get the response as `status: ok`, refresh you UI.

#### UI is not showing resources but catalog refresh is successful (Console Error)

- Verify catalog refresh is successful by hitting the `/v1/resources` API using the API URL.
- If is returning resources than check for console errors on UI.
- Right click -> Inspect or (Ctrl + Shift + I) (For Chrome). Then click on console tab.
- If there is cors error then check the UI config map and the URL you have added.
- The URL should be for example `https://api.hub.tekton.dev` and not `https://api.hub.tekton.dev/` there shouldn't be `/` at the end.

[ko]: https://github.com/google/ko
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[kind]: https://kind.sigs.k8s.io/docs/user/quick-start/
[minikube]: https://kubernetes.io/docs/tasks/tools/install-minikube/
[openshift]: https://www.openshift.com/try
[docker]: https://docs.docker.com/engine/install/
[podman]: https://podman.io/getting-started/installation
[config-yaml]: https://raw.githubusercontent.com/tektoncd/hub/master/config.yaml
[github-oauth-steps]: https://docs.github.com/en/developers/apps/creating-an-oauth-app
[gitlab-oauth-steps]: https://docs.gitlab.com/ee/integration/oauth_provider.html#user-owned-applications
[bitbucket-oauth-steps]: https://support.atlassian.com/bitbucket-cloud/docs/use-oauth-on-bitbucket-cloud
[catalog]: https://github.com/tektoncd/catalog
[catalog-tep]: https://github.com/tektoncd/community/blob/main/teps/0003-tekton-catalog-organization.md
[hub-config]: https://github.com/tektoncd/hub/blob/main/config.yaml
[upgrade-v1.6.0-to-v1.7.0]: https://github.com/tektoncd/hub/blob/main/docs/UPGRADE_V1.6.0_TO_V1.7.0.md
