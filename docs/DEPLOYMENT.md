# Deployment

### Prerequisites

- Kubernetes Cluster
  - You can also deploy on [Minkube][minikube], [kind][kind] or [OpenShift][openshift].
- [kubectl][kubectl]
- [ko][ko]

  - You can find installation steps [here][ko].

  ```
  go get github.com/google/ko/cmd/ko
  ```

- [docker][docker]

### Deploy API Service

Ensure you are in `api` directory

```
cd api
```

#### Deploy the database

```
kubectl apply -f config/00-config
```

#### Update the GitHub Api secret and client id

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

#### Prepare the API and DB Migration Release and apply the Yaml

Export `KO_DOCKER_REPO` for ko to publish image to. E.g.

```
export KO_DOCKER_REPO=quay.io/<username>
```

```
ko resolve -f config > api.yaml
```

Apply the release yaml

```
kubectl apply -f api.yaml
```

The command above will create a container image and push it to the registry pointed by `KO_DOCKER_REPO`. Ensure that the image is **publicly** available

```
kubectl get pods
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

```
2020-09-22T15:35:16.412Z INFO migration/migration.go:91 Migration ran successfully !! {"service": "migration"}
2020-09-22T15:35:16.412Z INFO db/main.go:39 DB initialisation successful !! {"service": "main"}
```

#### Verify if api route is accessible

```
curl -k -X GET -I $(oc get routes api --template='https://{{ .spec.host }}/categories')
```

#### In case of refreshing the catalog

1. Get the Hub Token
2. Navigate to `config/99-post-deploy` and edit `34-catalog-refresh-secret.yaml`. Set the Hub Token

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

3. `kubectl apply -f config/99-post-deploy/34-catalog-refresh-secret.yaml`

### Deploy UI

```
cd ui
```

#### Build and Publish Image

```
docker build -t <image> . && docker push <image>
```

#### Update the deployment image

Update `config/11-deployement` to use the image built above

#### Update the GitHub OAuth Client ID

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

#### Apply the manifests

```
kubectl apply -f config
```

#### Ensure pods are up and running

```
kubectl get pods
```

If on openshift verify if the ui route is accessible

Open: oc get routes ui --template='https://{{ .spec.host }} '

[ko]: https://github.com/google/ko
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[kind]: https://kind.sigs.k8s.io/docs/user/quick-start/
[minikube]: https://kubernetes.io/docs/tasks/tools/install-minikube/
[openshift]: https://www.openshift.com/try
[docker]: https://docs.docker.com/engine/install/
