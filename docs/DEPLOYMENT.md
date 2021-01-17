# Deployment

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

### Update the GitHub Api secret and client id

Navigate to the `config/00-config/` and edit `31-api-secret.yaml` . Set GitHub `oauth` client id and client secret.

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

**NOTE:** DO NOT commit and push

### Apply supporting Resources

```bash
kubectl apply -f config/00-config
```

### Prepare the API and DB Migration Release and apply the Yaml

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
ko apply -f api.yaml
```

The command above will create the deployment for Database and API and a job for DB-Migration.

```bash
$ kubectl get pods
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

### Setup Route/Ingress

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

### Verify if api route is accessible

For `OpenShift`:-

```bash
curl -k -X GET -I $(oc get routes api --template='https://{{ .spec.host }}/categories')
```

### Setup Catalog Refresh Secret

1. Get the Hub Token
2. Make sure you have your Github ID in Catalog Refresh scope present in [config.yaml][config.yaml]
3. Navigate to `config/99-post-deploy` and edit `34-catalog-refresh-secret.yaml`. Set the Hub Token

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

4. `kubectl apply -f config/99-post-deploy/34-catalog-refresh-secret.yaml`

## Deploy UI

```bash
cd ui
```

### Build and Publish Image

```bash
docker build -t <image> . && docker push <image>
```

### Update the deployment image

Update `config/11-deployement` to use the image built above

### Update the GitHub OAuth Client ID

Edit `config/10-config.yaml` and set your GitHub OAuth Client ID

```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: ui
  namespace: tekton-hub
data:
  API_URL: API URL   <<< update this by api routes
  GH_CLIENT_ID: GH OAuth Client ID   <<< update this
```

### Apply the manifests

```bash
kubectl apply -f config
```

### Ensure pods are up and running

```bash
kubectl get pods
```

### Setup Route/Ingress

After the deployment is done successfully, we need to expose the URL to access the UI.

- If deploying on OpenShift:-

  ```bash
  oc apply -f config/openshift/
  ```

- If deploying on Kubernetes:-

  - Create the secret containing tls cert and tls key

    ```bash
    kubectl create secret tls ui-hub-tekton-dev-tls --cert=path/to/cert/file --key=path/to/key/file
    ```

  - Apply the Ingress

    ```bash
    kubectl apply -f config/post-deploy
    ```

### Verify if UI route is accessible

If on openshift verify if the ui route is accessible

Open: `oc get routes ui --template='https://{{ .spec.host }} '`

## Deploying for the first time?

> Note: Make sure you have added your name in the scopes present [config.yaml][config.yaml]

1. Open currently deployed hub in the browser
2. Login into the Hub
3. After successful login, now open the terminal and follow below steps.
4. Exec into the database pod
   ```bash
   kubectl exec -it <db-pod-name> bash
   ```
5. Open PostgreSQL terminal by executing the command

   ```bash
   psql -U postgres
   ```

   You'll see the terminal changed to

   ```bash
   $ psql -U postgres
   psql (13beta3 (Debian 13~beta3-1.pgdg100+1))
   Type "help" for help.

   postgres=#
   ```

6. Connect to Hub database by executing
   ```sql
   postgres=# \c hub
   ```
7. Execute the SQL query

   ```sql
   postgres=# delete from configs;
   ```

8. Quit the database pod by first executing

   ```sql
   postgres=# \q
   ```

   and then `exit`.

9. Delete the API pod by executing the following command
   ```bash
   oc delete pod <pod-name>
   ```
   This will delete the previously created pod and spin a new pod which will add the scopes in the DB
10. Repeat the steps 1-2
11. Copy the Hub Token from the UI and add that in [catalog refresh step](#setup-catalog-refresh-secret)

[ko]: https://github.com/google/ko
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[kind]: https://kind.sigs.k8s.io/docs/user/quick-start/
[minikube]: https://kubernetes.io/docs/tasks/tools/install-minikube/
[openshift]: https://www.openshift.com/try
[docker]: https://docs.docker.com/engine/install/
[podman]: https://podman.io/getting-started/installation
[config.yaml]: https://raw.githubusercontent.com/tektoncd/hub/master/config.yaml
