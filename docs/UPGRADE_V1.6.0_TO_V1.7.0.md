# Upgrade Hub v1.6.0 to v1.7.0

This doc defines the steps to follow for upgrading from Hub v1.6.0 to Hub v1.7.0.

Navigate to `config` directoy in the root of the project.

## Get the backup of Hub v1.6.0 db

Follow these steps sequentially to take the backup of Hub v1.6.0 DB:-

Execute the below command to get the Hub v1.6.0 db pod

```
kubectl get pods -l app=db -n tekton-hub
```

Exec into the db pod

```
kubectl exec -it <db_pod_name> /bin/bash
```

Dump the Hub v1.6.0 db by running below commands inside the db pod

```
pg_dump -Ft -h localhost -U postgres hub -f /tmp/db.dump
```

Exit from the pod and copy the Hub v1.6.0 db dump file from the db pod to your system

```
kubectl cp <hub_v1.6.0_db_pod_name>:/tmp/db.dump ./backup_db.dump
```

## Deploy Databse

Execute the below command

```
kubectl apply -f 00-init/
```

This will create `tekton-hub` namespace, db deployment with `postgres` image, secret to save db credentials, pvc for the db and a ClusterIP service to expose the db internally.

All resources are created in `tekton-hub` namespace.

Wait till the pod comes in a running state

```
kubectl get pod -l app=tekton-hub-db -n tekton-hub -w
```

## Restore Hub v1.6.0 db backup into Hub v1.7.0 db

Access the Hub v1.7.0 db pod and exec into the pod and restore the db backup

```
kubectl get pods -l app=tekton-hub-db -n tekton-hub
```

Copy the db backup file from your system to the Hub v1.7.0 db pod

```

kubectl cp backup_db.dump <hub_v1.7.0_db_pod_name>:/tmp/backup_db.dump

```

Now restore the dump file into Hub v1.7.0 db pod by running below command inside the pod

```
pg_restore -d hub -h localhost -U postgres /tmp/backup_db.dump

```

To check succesful backup of hub v1.6.0 db, run the following commands inside the Hub v1.7.0 db pod and resources count shouldn't be zero

```

psql -U postgres

\c hub;

select count(*) from resources;

```

Exit from the db pod

## Run the DB Migration

Edit the `01-db/10-db-migration.yaml` and replace the image with Hub v1.7.0 DB migration image and apply the yaml.

```

kubectl apply -f 01-db/10-db-migration.yaml -n tekton-hub

```

## Deploy the API service

### Setup Route or Ingress

- If deploying on OpenShift:-

  ```bash
  kubectl apply -f 04-openshift/40-api-route.yaml -f 04-openshift/40-auth-route.yaml -n tekton-hub
  ```

````

- If deploying on Kubernetes:-

  - Create the secret containing tls cert and tls key both for API and Oauth server. Example as follows:

    ```bash
    kubectl create secret tls api-hub-tekton-dev-tls --cert=path/to/cert/file --key=path/to/key/file -n tekton-hub
    ```

  - Apply the Ingress

    ```bash
    kubectl apply -f 04-kubernetes/40-api-ingress.yaml -f 04-kubernetes/40-auth-ingress.yaml -n tekton-hub
    ```

### Update API Secret

Eexcute the below command to get the secret of Hub v1.6.0 api

```
kubectl get secret api -o yaml
```

Copy the secret values and update the `02-api/20-api-secret.yaml` file with same value

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: tekton-hub-api
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
  AUTH_BASE_URL: auth route of Hub v1.7.0
  GHE_URL: Add Github Enterprise URL in case of authenticating through Github Enterprise (Example (https|http)://myghe.com) --> Do not provide the catalog URL
  GLE_URL: Add Gitlab Enterprise URL in case of authenticating through Gitlab Enterprise (Example (https|http)://mygle.com) --> Do not provide the catalog URL
```

### Update API ConfigMap

Execute the below command to get the `CONFIG_FILE_URL` of Hub v1.6.0 api config

```
kubectl get cm/api -o=jsonpath="{.data['CONFIG_FILE_URL']}"
```

Copy the `CONFIG_FILE_URL` url path from Hub v1.6.0 configMap and update the `02-api/21-api-configmap.yaml` with same url path

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: tekton-hub-api
  labels:
    app: tekton-hub-api
data:
  CONFIG_FILE_URL: https://raw.githubusercontent.com/tektoncd/hub/master/config.yaml ## Change the file URL here with Hub v1.6.0 configMap's CONFIG_FILE_URL
```

### Update API Image

Edit the `02-api/22-api-deployment.yaml` and replace the image with Hub v1.7.0 API image

### Apply API configs

```bash
kubectl apply -f 02-api/ -n tekton-hub
```

This will create the pvc, deployment, secret, configmap and a NodePort service to expose the API server.

## Deploy the UI service

### Setup Route or Ingress

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

Edit `config/03-ui/30-ui-configmap.yaml` and by setting the API URL, API version auth URL and Redirect URL.

You can get the API URL,AUTH URL and UI URL for Hub v1.7.0 by using below command (OpenShift)

```
  kubectl get routes tekton-hub-api --template='https://{{ .spec.host }}' -n tekton-hub
```

```
  kubectl get routes tekton-hub-auth --template='https://{{ .spec.host }}' -n tekton-hub
```

```
  kubectl get routes tekton-hub-ui --template='https://{{ .spec.host }}' -n tekton-hub
```

```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: tekton-hub-ui
data:
  API_URL: API URL   <<< Update this by Hub v1.7.0 API url
  API_VERSION: API VERSION  <<< Update this by API version (For e.g-"v1")
  AUTH_BASE_URL: AUTH URL << Update this by hub v1.7.0 auth url
  REDIRECT_URI: UI URL << Update this by hub v1.7.0 ui url
```

### Update UI Image

Edit the `03-ui/31-ui-deployment.yaml` by replacing the image with Hub v1.7.0 UI image

### apply API configs

Execute follow command

```bash
kubectl apply -f 03-ui/ -n tekton-hub
```

### Update the Git OAuth Configuration with Hub v1.7.0 auth routes

- **Github** - Update GitHub OAuth with `Homepage URL` and `Authorization callback URL` as Huv v1.7.0 `<auth-route>`.
- **Gitlab** - Update Gitlab Oauth with `REDIRECT_URI` as `<auth-route>/auth/gitlab/callback`,use the Hub v1.7.0 auth route as `<auth_route>`.
- **BitBucket** - Update BitBucket Oauth with `Callback URL` as Hub v1.7.0 `<auth-route>`.

## Delete the db, api and ui resoruces for Hub v1.6.0

Verify the Hub v1.7.0 API, DB and UI if it's working fine then delete the all resources for Hub v1.6.0 by exucuting the following commands

```
kubectl delete secret,pvc,deployment,service -l app=db -n tekton-hub
```

```
kubectl delete secret,pvc,cm,deployment,service,route -l app=api -n tekton-hub
```

```
kubectl delete cm,deployment,service,route -l app=ui -n tekton-hub
```
````
