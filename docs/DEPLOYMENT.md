# Deployment

- [Prerequisites](#prerequisites)
- [Deploy API Service](#deploy-api-service)
  - [Update API Secret](#update-the-api-secret)
  - [Update API Config Map](#update-api-config-map)
  - [Prepare API and Db Migration release yaml](#Prepare-the-api-and-db-migration-release-yaml)
  - [Setup Ingress/Route](#setup-route-or-ingress)
- [Deploy UI](#deploy-ui)
  - [Build and Publish Image](#build-and-publish-image)
  - [Update the UI ConfigMap](#update-the-ui-configmap)
  - [Setup Ingress/Route](#setup-route-or-ingress)
- [Add resources in DB](#add-resources-in-db)
- [Setup Catalog Refresh CronJob](#setup-catalog-refresh-cronjob)
- [Adding New Users in Config](#adding-new-users-in-config)
- [Troubleshooting](#troubleshooting)

## Prerequisites

- Kubernetes Cluster
  - You can also deploy on [Minkube][minikube], [kind][kind] or [OpenShift][openshift].
- [kubectl][kubectl]
- [ko][ko]

  - You can find installation steps [here][ko] or if you have Go installed then you can execute the command :-

    ```bash
    go get github.com/google/ko/cmd/ko
    ```

- [docker][docker] or [podman][podman]

## Deploy API Service

Ensure you are in `api` directory

```bash
cd api
```

### Update the API Secret

Navigate to the `config/00-config/` and edit `31-api-secret.yaml` . Set GitHub `oauth` client id and client secret.

To create a GitHub OAuth follow the steps given [here][oauth-steps]. For now add any URL in place of `Homepage URL` and `Authorization callback URL` for eg. `http://localhost:5000`. We will update this URL once we deploy the UI pod.
 
- After creating the OAuth add the Client ID and Client Secret in the file.
- For JWT, you can add any random key, this is used to sign the JWT created for users.
- For `ACCESS_JWT_EXPIRES_IN` and `REFRESH_JWT_EXPIRES_IN` add time you want the jwt to be expired in. Refresh time should be greater than Access time.
eg. 1m = 1 minute, 1h = 1 hour, 1d = 1 day

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: api
  namespace: tekton-hub
type: Opaque
stringData:
  GH_CLIENT_ID: Oauth client id
  GH_CLIENT_SECRET: Oauth secret
  JWT_SIGNING_KEY: a-long-signing-key
  ACCESS_JWT_EXPIRES_IN: time such as 15m
  REFRESH_JWT_EXPIRES_IN: time such as 15m
```

### Update API Config Map

If you want to change the config file passed to hub, you can edit the `[30-api-configmap.yaml][hub-cm]` and change the URL to the file. By default it is pointing to [config.yaml][config-yaml] in Hub which has Application Data.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: api
  namespace: tekton-hub
  labels:
    app: api
data:
  CONFIG_FILE_URL: https://raw.githubusercontent.com/tektoncd/hub/master/config.yaml        ## Change the file URL here

```

In this file we add users who have additional scopes which can be used to refresh a catalog, refresh config file after changes, create an agent token which can be used in cron job to refresh the catalog after an interval.

All users have default scopes `rating:read` and `rating:write` they get after login which allows them to rate the resources.

#### Apply supporting Resources

```bash
kubectl apply -f config/00-config
```

This will create `tekton-hub` namespace and create the supporting resources.

### Prepare the API and DB Migration Release yaml

Export `KO_DOCKER_REPO` for ko to publish image to. E.g.

```bash
export KO_DOCKER_REPO=quay.io/<username>
```

```bash
ko resolve -f config > api.yaml
```

The command above will create a container image and push it to the registry pointed by `KO_DOCKER_REPO`. Ensure that the image is **publicly** available.

Apply the release yaml

```bash
kubectl apply -f api.yaml
```

The command above will create the deployment for Database and API and a job for DB-Migration.

The db-migration job will create all tables in the db.

```bash
$ kubectl get pods -n tekton-hub
NAME                   READY   STATUS      RESTARTS   AGE
api-86ccf7484f-qrz4k   1/1     Running     3          50s
db-589d44fdd5-ksf8v    1/1     Running     0          50s
db-migration-4mb75     0/1     Error       0          47s
db-migration-8vhpd     0/1     Completed   0          17s
db-migration-pmtw5     0/1     Error       0          50s
db-migration-tkrsn     0/1     Error       0          37s
```

One can also check the logs using `kubectl logs db-migration-8vhpd`

The migration logs at the end should show

```bash
2020-09-22T15:35:16.412Z INFO migration/migration.go:91 Migration ran successfully !! {"service": "migration"}
2020-09-22T15:35:16.412Z INFO db/main.go:39 DB initialisation successful !! {"service": "main"}
```

### Setup Route or Ingress

After the deployment is done successfully, we need to expose the URL to access the API.

- If deploying on OpenShift:-

  ```bash
  oc apply -f config/openshift/
  ```

- If deploying on Kubernetes:-

  - Create the secret containing tls cert and tls key

    ```bash
    kubectl create secret tls api-hub-tekton-dev-tls --cert=path/to/cert/file --key=path/to/key/file
    ```

  - Apply the Ingress

    ```bash
    kubectl apply -f config/99-post-deploy/33-api-ingress.yaml
    ```

#### Verify if api route is accessible

For `OpenShift`:-

```bash
curl -k -X GET -I $(oc get -n tekton-hub routes api --template='https://{{ .spec.host }}/v1/categories')
```

NOTE: At this moment, there are no resources in the db. Only the tags and categories from hub config are added in the db.

Let's deploy the UI first and then we will add the resources in db.

## Deploy UI

Ensure you are in `ui` directory and logged in a image registry.

```bash
cd ui
```

### Build and Publish Image

```bash
docker build -t <image> . && docker push <image>
```

eg. `quay.io/<username>/ui`

Make sure your image is public after pushing to a registry.

#### Update the deployment image

Update `config/11-deployment` to use the image built above

### Update the UI ConfigMap

Edit `config/10-config.yaml` and set your GitHub OAuth Client ID and the API URL.

You can get the API URL using below command
```
  oc get -n tekton-hub routes api --template='https://{{ .spec.host }}
```

```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: ui
  namespace: tekton-hub
data:
  API_URL: API URL   <<< update this by api route
  GH_CLIENT_ID: GH OAuth Client ID   <<< update this
```

#### Apply the manifests

```bash
kubectl apply -f config/
```

This will create the deployment and service for UI.

#### Ensure pods are up and running

```bash
kubectl get pods -n tekton-hub
```
eg. wait till status of UI pod is running
```
NAME                   READY   STATUS      RESTARTS   AGE
api-6dfc6f97d9-vk66r   1/1     Running     3          21m
db-9bd4cdf99-zsq89     1/1     Running     0          21m
db-migration-26ngs     0/1     Completed   0          21m
ui-55fc66cc6b-69dsp    1/1     Running     0          58s
```

### Setup Route or Ingress

After the deployment is done successfully, we need to expose the URL to access the UI.

- If deploying on OpenShift:-

  ```bash
  oc apply -f config/openshift/
  ```

- If deploying on Kubernetes:-

  - Create the secret containing tls cert and tls key

    ```bash
    kubectl create secret tls ui-hub-tekton-dev-tls --cert=path/to/cert/file --key=path/to/key/file -n tekton-hub
    ```

  - Apply the Ingress

    ```bash
    kubectl apply -f config/post-deploy
    ```

#### Verify if UI route is accessible

If on openshift verify if the ui route is accessible

Open: `oc get routes ui --template='https://{{ .spec.host }} '`

Note: In UI you won't be able see the resources on homepage as there are not any in the db but you should see Categories on the left side.

Update your GitHub OAuth created with the UI route in place of `Homepage URL` and `Authorization callback URL`.

## Add resources in DB

> Note: Make sure you have added your name in the scopes present [config.yaml][config-yaml]. You will need additional scope i.e. catalog-refresh to add resources.

Now, follow the below steps

1. Login through the UI. Click on `Login` on right corner and then `Sign in with GitHub`.

    This will add you in db, but yet you have only the default scopes i.e. `rating:read` and `rating:write` which you can check with your jwt. On UI, in right corner there is an option to `Copy Hub Token` and paste in any jwt decoder. for eg. https://jwt.io/
    
2. We save the checksum of config file in db to avoid reading the config file again and again if the api pod is deleted due to any reason. The config file is read by API pod before starting the server. 

3. To make API pod read config we need to delete the entry in db. Exec into the database pod
   ```bash
   kubectl exec -it <db-pod-name> bash
   ```
4. Open PostgreSQL terminal by executing the command

   ```bash
   psql -U postgres -d hub
   ```

   You'll see the terminal changed to

   ```bash
   hub=#
   ```
5. Execute the SQL query

   ```sql
   hub=# delete from configs;
   ```

6. Quit the database pod by first executing

   ```sql
   hub=# \q
   ```

   and then `exit`.

7. Delete the API pod by executing the following command
   ```bash
   oc delete pod <pod-name>
   ```
   This will delete the previously created pod and spin a new pod which will add the scopes in the DB.
   
8. Once the pod is running again, Logout from Hub UI and login again.


9. Now, if you copy your token from UI and check in jwt decoder, you will have additional scopes.

10. To add resources, you need to make a POST api call passing your jwt token with the addtional scopes in Header.

```
curl -X POST -H "Authorization: <access-token>" \
    <api-route>/catalog/refresh 
```
Replace `<access-token>`  with your Hub token copied from UI and replace `<api-route>` with your api url.

This will give an output as below

```
{"id":1,"status":"queued"}
```

11. Refresh your UI, you will be able to see resources.

## Setup Catalog Refresh CronJob

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

4. Navigate to `api/config/99-post-deploy` and edit `34-catalog-refresh-secret.yaml`. Set the Hub Token

   ```yaml
   apiVersion: v1
   kind: Secret
   metadata:
     name: catalog-refresh
     namespace: tekton-hub
   type: Opaque
   stringData:
     HUB_TOKEN: hub token
   ```
    Set the token created in last step as Hub token.

4. Make sure you are in `api` directory and run the below command. 

    ```
    kubectl apply -f config/99-post-deploy/34-catalog-refresh-secret.yaml \ 
                -f config/99-post-deploy/35-catalog-refresh-cronjob.yaml 
    ```
    The cron job is configured to run every 30 min, you can change the interval in `config/99-post-deploy/35-catalog-refresh-cronjob.yaml`.
    
## Adding New Users in Config

By default when any user login for the first time, they will have only default scope even if they are added in config.yaml. To get additinal scopes, make sure the user has logged in once.

Now, we need to refresh the config. To do that do a POST API call.

```
curl -X POST -H "Authorization: <access-token>" \
    --header "Content-Type: application/json" \
    --data '{"force": true} \
  <api-route>/system/config/refresh
```
Replace `<access-token>` with your JWT token. you must have `config-refresh` scope to call this API. 

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
[oauth-steps]:https://docs.github.com/en/developers/apps/creating-an-oauth-app 
[hub-cm]:https://github.com/tektoncd/hub/blob/master/api/config/00-config/30-api-configmap.yaml