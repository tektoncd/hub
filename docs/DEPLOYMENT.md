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

```
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

[ko]: https://github.com/google/ko
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[kind]: https://kind.sigs.k8s.io/docs/user/quick-start/
[minikube]: https://kubernetes.io/docs/tasks/tools/install-minikube/
[openshift]: https://www.openshift.com/try
