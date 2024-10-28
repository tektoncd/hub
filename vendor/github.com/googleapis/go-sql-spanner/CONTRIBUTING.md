# Contributing

1. [File an issue](https://github.com/googleapis/go-sql-spanner/issues/new/choose).
   The issue will be used to discuss the bug or feature and should be created
   before sending a PR.

1. [Install Go](https://golang.org/dl/).
    1. Ensure that your `GOBIN` directory (by default `$(go env GOPATH)/bin`)
       is in your `PATH`.
    1. Check it's working by running `go version`.
        * If it doesn't work, check the install location, usually
          `/usr/local/go`, is on your `PATH`.

1. Sign one of the
   [contributor license agreements](#contributor-license-agreements) below.

1. Clone the repo:
   `git clone https://github.com/googleapis/go-sql-spanner`

1. Change into the checked out source:
   `cd go-sql-spanner`

1. Fork the repo.

1. Set your fork as a remote:
   `git remote add fork git@github.com:GITHUB_USERNAME/go-sql-spanner.git`

1. Make changes, commit to your fork.

   Commit messages should follow the
   [Conventional Commits Style](https://www.conventionalcommits.org). The scope
   portion should always be filled with the name of the package affected by the
   changes being made. For example:
   ```
   feat(functions): add gophers codelab
   ```

1. Send a pull request with your changes.

   To minimize friction, consider setting `Allow edits from maintainers` on the
   PR, which will enable project committers and automation to update your PR.

1. A maintainer will review the pull request and make comments.

   Prefer adding additional commits over amending and force-pushing since it can
   be difficult to follow code reviews when the commit history changes.

   Commits will be squashed when they're merged.

## Testing

We test code against two versions of Go, the minimum and maximum versions
supported by our clients. To see which versions these are checkout our
[README](README.md#supported-versions).

### Integration Tests

In addition to the unit tests, you may run the integration test suite.

#### GCP Setup

To run the integrations tests, creation and configuration of a project in
the Google Developers Console is required.

After creating the project, you must [create a service account](https://developers.google.com/identity/protocols/OAuth2ServiceAccount#creatinganaccount)
for project. Ensure the project-level **Owner**
[IAM role](https://console.cloud.google.com/iam-admin/iam/project) role is added to
each service account. During the creation of the service account, you should
download the JSON credential file for use later.

#### Local Setup

Once the project is created and configured, set the following environment
variables:

- `SPANNER_TEST_PROJECT`: Developers Console project's ID (e.g.
  bamboo-shift-455).

Install the [gcloud command-line tool][gcloudcli] to your machine and use it to
create some resources used in integration tests.

From the project's root directory:

``` sh
# Sets the default project in your env.
$ gcloud config set project $SPANNER_TEST_PROJECT

# Authenticates the gcloud tool with your account.
$ gcloud auth login

# Creates a Spanner instance for the spanner integration tests.
$ gcloud beta spanner instances create go-integration-test --config regional-us-central1 --nodes 10 --description 'Instance for go client test'
# NOTE: Spanner instances are priced by the node-hour, so you may want to
# delete the instance after testing with 'gcloud beta spanner instances delete'.

$ export SPANNER_TEST_INSTANCE=go-integration-test
```

It may be useful to add exports to your shell initialization for future use.
For instance, in `.zshrc`:

```sh
#### START Test Variables
# Developers Console project's ID (e.g. bamboo-shift-455) for the general project.
export SPANNER_TEST_PROJECT=your-project

# Developers Console Spanner's instance ID (e.g. spanner-instance) for the running tests.
export SPANNER_TEST_INSTANCE=go-integration-test
#### END Test Variables
```

#### Running

Once you've done the necessary setup, you can run the integration tests by
running:

``` sh
$ go test -v ./...
```

## Contributor License Agreements

Before we can accept your pull requests you'll need to sign a Contributor
License Agreement (CLA):

- **If you are an individual writing original source code** and **you own the
  intellectual property**, then you'll need to sign an [individual CLA][indvcla].
- **If you work for a company that wants to allow you to contribute your
  work**, then you'll need to sign a [corporate CLA][corpcla].

You can sign these electronically (just scroll to the bottom). After that,
we'll be able to accept your pull requests.

[gcloudcli]: https://developers.google.com/cloud/sdk/gcloud/
[indvcla]: https://developers.google.com/open-source/cla/individual
[corpcla]: https://developers.google.com/open-source/cla/corporate