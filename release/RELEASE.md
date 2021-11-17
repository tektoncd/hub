# Release

- [Prerequisites](#prerequisites)
- [Implementation](implementation)

## Prerequisites

- [ko][ko]
- [docker][docker] or [podman][podman]
- [hub][hub]

## Implementation

1. Once the above tools are installed, run the `release.sh` script which generates the release yamls for `db-migration`, `api` and `ui`

   ```
   REGISTRY_BASE_URL=docker.io/<username> bash release.sh
   ```

   **NOTE**: By default the registry is `quay.io/tekton-hub` and docker is by default tool to build the images. If you have `podman` you can run the command as

   ```
   DOCKER_CMD=podman REGISTRY_BASE_URL=docker.io/<username> bash release.sh
   ```

   - This will generate the following release yamls using ko
     - db-migration
     - api-k8s
     - api-openshift
     - ui-k8s
     - ui-openshift
   - Build the images for db-migration, api and ui and pushes the images to the specified registry
   - Replaces the generated yamls with the newly build images
   - Creates a draft release and attaches the generated release yamls

2. Once the release yamls are attached to the draft release then edit the draft release by adding the release notes for the release and publish the release

[ko]: https://github.com/google/ko
[docker]: https://docs.docker.com/engine/install/
[podman]: https://podman.io/getting-started/installation
[hub]: https://github.com/github/hub#installation
