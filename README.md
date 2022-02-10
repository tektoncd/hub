# Tekton Hub 

<p align="center">
<img width="250" height="175" src="https://github.com/cdfoundation/artwork/blob/main/tekton/additional-artwork/tekton-hub/color/TektonHub_color.svg" alt="Tekton Hub logo"></img>
</p>

The Tekton hub is a web based platform for developers to discover,
share and contribute tasks and pipelines for Tekton. This repository
contains the source code of this service.

> Hub is launched as [hub.tekton.dev](https://hub.tekton.dev) :mega:

## Key features

* Tekton Hub provides the following features:

1. Display of resources in a curated way:

     User gets detailed information such as description, YAML, tags and rating of resources.

  2. Filter by  categories :

     This allows user to filter resources based on categories which will help user to get multiple resources.

  3. Search resources on the basis of `name` or `displayName`

  4. Rating

     - User can rate a resource based on the experience
     - This can even help the other user as by knowing the rating they can directly compare and use the resource

  5. Sort resources on the basis of name or rating.

  6. Install resources on cluster :
     - The Tekton Hub provides cluster installation commands for tasks or pipelines.
    
## Hub APIs for Integration

The public hub APIs are exposed for Integration outside hub. You can find the API Policy [here](docs/API_POLICY.md).

For documentation of public APIs, you can use [Hub Swagger](https://swagger.hub.tekton.dev)

## Deploy your own instance

You can deploy your own instance of Tekton Hub. You can find the documentation [here](docs/DEPLOYMENT.md).

## Want to Contribute

We are so excited to have you!

- See [CONTRIBUTING.md](CONTRIBUTING.md) for an overview of our processes
- See [DEVELOPMENT.md](docs/DEVELOPMENT.md) for how to get started
- See [ROADMAP.md](ROADMAP.md) for the current roadmap
- Look at our
  [good first issues](https://github.com/tektoncd/hub/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
  and our
  [help wanted issues](https://github.com/tektoncd/hub/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22)
- If you are looking for support, enter an [issue][issue] or join our [Slack workspace][slack]


[swagger-def]:https://raw.githubusercontent.com/tektoncd/hub/main/api/v1/gen/http/openapi3.yaml
[issue]:https://github.com/tektoncd/hub/issues/new
[slack]:https://github.com/tektoncd/community/blob/main/contact.md#slack