# Release Script

The Release script does the following tasks:
- fetch the latest main branch
- create and push a git tag by version provided from the user
- create supporting configuration for hub on the cluster
- create secret and roles for the Tekton Pipelines
- install the pipeline which will deploy the hub on the clusters

## Prerequisites

- Kubernetes or OpenShift cluster with Tekton Pipelines installed
- kubectl CLI
- git CLI

### User Prerequisites

- Access to push images to `quay.io/tekton-hub`
- Access to push git tag to `tektoncd/hub`

## Before Running the Script

- Enable Tekton Bundle in Pipelines
- Check the last hub release version [here](https://github.com/tektoncd/hub/releases) and enter the next in the script to
- If you are deploying hub on a clean cluster
    - You will need GitHub OAuth. You can create one using the steps given [here](https://docs.github.com/en/developers/apps/creating-an-oauth-app). Use any url for creating OAuth, later you can change it once hub is deployed.


## Running the Script

This will deploy the hub in `tekton-hub` namespace and run the Tekton pipelines for it in the `tekton-hub-ci` namespace.

Input asked by the script:
- Target Release Version: This is the tag to be created and pushed to the repository. Check the last tag created and enter the next verion.
- Deploying on Openshift? Enter (Y or n) depending on your cluster. If OpenShift is selected, it will configure adm policy for the service account.
- If there are not existing hub configuration in `tekton-hub` namespace, the script will ask for 
    - Database Configuration: Enter Database name, user and password for the db to be created
    - GitHub OAuth Configuration: Enter the OAuth Configuration you have created.
    - JWT Signing key: Enter any random key which will be used to sign User JWTs.
    - Access and Refresh Expire Time: Enter the time the token should be expired in. Refresh Expire time must be greater than Access Expire time. You can input time as `1d` = 1 day, `15h` = 15 hours, `30m` = 30 minutes.
    - Hub Config Raw URL: No need to change it unless you are not deploying hub from tektoncd/hub.
NOTE: If you have already hub instance in tekton-hub namespace and the above configuration are already created using secrets and config maps, then the above step will be skipped.
- Quay registry credentials: Enter your credentials to push images to the registry

Once, that is done,the script will create all configuration and then install the pipeline and resources requires by it. And start the Pipeline.
