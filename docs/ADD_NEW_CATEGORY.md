# Adding New Categories to Hub

This doc defines the steps to add a new categories in Hub.

Process to add a new category:

- Create a pull request to Hub repository by adding your categories in [Hub Api ConfigMap][config].

Once the pull request is reviewed and merged, the Next Steps are to be performed by the Hub maintainers.

- Invoke the `/system/config/refresh` API. This will add the new categories in db. To access the API, you need to have `config:refresh` scope.

Once the config refresh is done you will be able to see the newly added categories on Hub UI

[config]: https://github.com/tektoncd/hub/blob/main/config/02-api/21-api-configmap.yaml
