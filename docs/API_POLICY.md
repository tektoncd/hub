# Hub API Policy

This document proposes a policy regarding making updates to the APIs of this repo.

NOTE: Hub Release Version and API version are different and may not be in sync. Making a Hub major release does not necessarily increment API version too. API version will be incremented only when there is a breaking change.
eg. `/v1` -> `/v2` only when there is a breaking change.

Hub API service hosts two types of APIs:

## 1. Versioned APIs

- These are the APIs which are exposed for external users to use and integrate with their applications.
- These will have support provided by the Hub team.
- These are the APIs which return hub resource data.
- The version of the APIs will be prefixed to the API route. 
  Eg. resource API - /v1/resources
- Swagger Documentation for these APIs can be found at `https://api.hub.tekton.dev/<version>/schema/swagger.json`. Replace `<version>` with current version of APIs.
 

#### Additive changes

- Additive changes are changes that are added to the APIs in a backward compatible way and thus, do not cause problems for existing users of the API.
- If the current version of API is `/v1` then any new addition for eg. adding new field in response of HTTP Request which will not affect existing users of v1 version of that API.
- It is also possible to add new fields to request payload of an API as long as it is done in a backward compatible way, meaning it has a default value and can be omitted safely. 
- These changes must be approved by at least 2 OWNERS.


#### Breaking changes

- Breaking changes are those when users have to change their implementation as per the changes in API. for eg. removal of a field from response, removing an API, removing a query param of an API. 
- These changes will be introduced as a new version of the APIs.
- If current version of APIs is `/v1` and an API is going to be removed then all other APIs will be moved to `/v2`. 
  Users of v1 will be able to use old versioned APIs for at least 9 months then they will be removed.
- These changes must be approved by at least 2 OWNERS.

#### List of Versioned APIs
  - `/v1`
  - `/v1/query`
  - `/v1/resources`
  - `/v1/resource/<id>`
  - `/v1/resource/<catalog>/<kind>/<name>`
  - `/v1/resource/<catalog>/<kind>/<name>/<version>`
  - `/v1/resource/<id>/versions`
  - `/v1/resource/version/<versionID>`
  - `/v1/schema/swagger.json`


## 2. Internal APIs

- These are the APIs which are specifically used for Hub UI and hub services.
- These are not exposed for integration.
- These may change at any time as the Hub UI evolves.
- No support is provided for integration of this APIs.
- The Swagger documentation can be found [here](https://api.hub.tekton.dev/schema/swagger.json).
    Eg. Login, rating, Catalog refresh APIs, etc.

#### List of Internal APIs
- `/`
- `/categories`
- `/auth/login`
- `/system/user/agent`
- `/system/config/refresh`
- `/catalog/refresh`
- `/resource/<id>/rating`
- `/schema/swagger.json`
