# Tekton Hub 2021 Roadmap

## Roadmap

### Multi Catalog Support

 - User should be able to provide multi catalogs in config
 - API's should be able to act for multiple catalogs
 - UI should show the resources from all catalogs mentioned and also able to filter based on catalog
 - CLI should be able to work with multiple catalogs

### Add support for Pipeline Kind

 - Hub should support the Pipelines Kind resources
 - APIs created to get Pipeline Kind data
 - CLI should be able to install, upgrade etc for Pipelines
 - Catlin should be able to validate based on Pipeline Structure Proposal

### Handle Deprecated Task in UI

 - UI should have a way to identify deprecated task
 - API should return whether the resource is deprecated or not

### CLI acts based on Minimum Pipeline Version of resource

 - CLI commands like install, upgrade, etc act based on Pipeline version installed on cluster
 - It should warn user in case of mismatch with supported Pipeline version of resource

### Automate Hub Deployment

 - Hub deployment should be automated through script.

### Versioning of the APIs

 - APIs should be versioned.
 - API Policy should be defined.

### Automation of Categories and Tags mapping

 - Comeup with a proposal to remove the manual categories and tags mapping

### Add supported resources of tasks/pipelines in HUB UI

 - Hub should show details about examples, samples, owners etc. of resources.

### Interactive CLI command

 - Hub CLI should ask inputs from users interactively if not provided as args

### Bundle support

 - Hub CLI and Hub UI uses bundles to fetch resource yaml's
 - Can store resource bundle URLs in database

### Support for other git providers

 - Hub should support catalogs hosted on GitHub Enterprise, Gitlab, BitBucket etc.

### UI Refactoring and Enhancements

 - Better User Details
 - Improved searching experience
 - Improvement in filtering and providing URL params
