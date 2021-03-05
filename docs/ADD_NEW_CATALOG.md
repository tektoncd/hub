# Adding New Catalog to Hub

This doc defines the steps to add a new catalog in Hub. The catalog **must** follow the structure defined in the [Catalog Organization TEP][tep].

Process to add a new catalog:

- Create a pull request to Hub repository adding your catalog details in [Hub Config file][config].
-  Make sure you give a unique name which is not used for other catalogs defined in Config file.
    This name will be used in identifying catalog and will be used in API to search resources.
 eg.`/resource/<catalog-name>/<resource-name>`
 
Once the pull request is reviewed and merged, the Next Steps are to be performed by the Hub maintainers.
    
- Invoke the `/system/config/refresh` API. This will add the new catalog details in db. To access the API, you need to have `config:refresh` scope. 
- Now, use the Catalog refresh API `catalog/<catalogName>/refresh` to add resources from catalog in hub database. To access the API, you need to have `catalog:refresh` scope.
- Setup a cronjob to refresh the catalog after a certain interval.

After the catalog refresh is done, UI will reflect the resources from the newly added catalog.

[tep]: https://github.com/tektoncd/community/blob/main/teps/0003-tekton-catalog-organization.md
[config]: https://github.com/tektoncd/hub/blob/d29cf3d2a522bc6d27357083aa0cf896ea22f242/config.yaml#L49
