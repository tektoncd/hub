// Code generated by goa v3.3.1, DO NOT EDIT.
//
// hub HTTP client CLI support package
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	adminc "github.com/tektoncd/hub/api/gen/http/admin/client"
	authc "github.com/tektoncd/hub/api/gen/http/auth/client"
	catalogc "github.com/tektoncd/hub/api/gen/http/catalog/client"
	categoryc "github.com/tektoncd/hub/api/gen/http/category/client"
	ratingc "github.com/tektoncd/hub/api/gen/http/rating/client"
	resourcec "github.com/tektoncd/hub/api/gen/http/resource/client"
	statusc "github.com/tektoncd/hub/api/gen/http/status/client"
	userc "github.com/tektoncd/hub/api/gen/http/user/client"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//    command (subcommand1|subcommand2|...)
//
func UsageCommands() string {
	return `admin (update-agent|refresh-config)
auth authenticate
catalog (refresh|refresh-all)
category list
rating (get|update)
resource (query|list|versions-by-id|by-catalog-kind-name-version|by-version-id|by-catalog-kind-name|by-id)
status status
user (refresh-access-token|new-refresh-token|info)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` admin update-agent --body '{
      "name": "abc",
      "scopes": [
         "catalog-refresh",
         "agent:create"
      ]
   }' --token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODAzMDAsImlhdCI6MTU3Nzg4MDAwMCwiaWQiOjExLCJpc3MiOiJUZWt0b24gSHViIiwic2NvcGVzIjpbInJhdGluZzpyZWFkIiwicmF0aW5nOndyaXRlIiwiYWdlbnQ6Y3JlYXRlIl0sInR5cGUiOiJhY2Nlc3MtdG9rZW4ifQ.6pDmziSKkoSqI1f0rc4-AqVdcfY0Q8wA-tSLzdTCLgM"` + "\n" +
		os.Args[0] + ` auth authenticate --code "5628b69ec09c09512eef"` + "\n" +
		os.Args[0] + ` catalog refresh --catalog-name "tekton" --token "Occaecati officia inventore adipisci."` + "\n" +
		os.Args[0] + ` category list` + "\n" +
		os.Args[0] + ` rating get --id 11281538076509796713 --token "Placeat hic accusamus."` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	scheme, host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restore bool,
) (goa.Endpoint, interface{}, error) {
	var (
		adminFlags = flag.NewFlagSet("admin", flag.ContinueOnError)

		adminUpdateAgentFlags     = flag.NewFlagSet("update-agent", flag.ExitOnError)
		adminUpdateAgentBodyFlag  = adminUpdateAgentFlags.String("body", "REQUIRED", "")
		adminUpdateAgentTokenFlag = adminUpdateAgentFlags.String("token", "REQUIRED", "")

		adminRefreshConfigFlags     = flag.NewFlagSet("refresh-config", flag.ExitOnError)
		adminRefreshConfigTokenFlag = adminRefreshConfigFlags.String("token", "REQUIRED", "")

		authFlags = flag.NewFlagSet("auth", flag.ContinueOnError)

		authAuthenticateFlags    = flag.NewFlagSet("authenticate", flag.ExitOnError)
		authAuthenticateCodeFlag = authAuthenticateFlags.String("code", "REQUIRED", "")

		catalogFlags = flag.NewFlagSet("catalog", flag.ContinueOnError)

		catalogRefreshFlags           = flag.NewFlagSet("refresh", flag.ExitOnError)
		catalogRefreshCatalogNameFlag = catalogRefreshFlags.String("catalog-name", "REQUIRED", "Name of catalog")
		catalogRefreshTokenFlag       = catalogRefreshFlags.String("token", "REQUIRED", "")

		catalogRefreshAllFlags     = flag.NewFlagSet("refresh-all", flag.ExitOnError)
		catalogRefreshAllTokenFlag = catalogRefreshAllFlags.String("token", "REQUIRED", "")

		categoryFlags = flag.NewFlagSet("category", flag.ContinueOnError)

		categoryListFlags = flag.NewFlagSet("list", flag.ExitOnError)

		ratingFlags = flag.NewFlagSet("rating", flag.ContinueOnError)

		ratingGetFlags     = flag.NewFlagSet("get", flag.ExitOnError)
		ratingGetIDFlag    = ratingGetFlags.String("id", "REQUIRED", "ID of a resource")
		ratingGetTokenFlag = ratingGetFlags.String("token", "REQUIRED", "")

		ratingUpdateFlags     = flag.NewFlagSet("update", flag.ExitOnError)
		ratingUpdateBodyFlag  = ratingUpdateFlags.String("body", "REQUIRED", "")
		ratingUpdateIDFlag    = ratingUpdateFlags.String("id", "REQUIRED", "ID of a resource")
		ratingUpdateTokenFlag = ratingUpdateFlags.String("token", "REQUIRED", "")

		resourceFlags = flag.NewFlagSet("resource", flag.ContinueOnError)

		resourceQueryFlags        = flag.NewFlagSet("query", flag.ExitOnError)
		resourceQueryNameFlag     = resourceQueryFlags.String("name", "", "")
		resourceQueryCatalogsFlag = resourceQueryFlags.String("catalogs", "", "")
		resourceQueryKindsFlag    = resourceQueryFlags.String("kinds", "", "")
		resourceQueryTagsFlag     = resourceQueryFlags.String("tags", "", "")
		resourceQueryLimitFlag    = resourceQueryFlags.String("limit", "1000", "")
		resourceQueryMatchFlag    = resourceQueryFlags.String("match", "contains", "")

		resourceListFlags     = flag.NewFlagSet("list", flag.ExitOnError)
		resourceListLimitFlag = resourceListFlags.String("limit", "1000", "")

		resourceVersionsByIDFlags  = flag.NewFlagSet("versions-by-id", flag.ExitOnError)
		resourceVersionsByIDIDFlag = resourceVersionsByIDFlags.String("id", "REQUIRED", "ID of a resource")

		resourceByCatalogKindNameVersionFlags       = flag.NewFlagSet("by-catalog-kind-name-version", flag.ExitOnError)
		resourceByCatalogKindNameVersionCatalogFlag = resourceByCatalogKindNameVersionFlags.String("catalog", "REQUIRED", "name of catalog")
		resourceByCatalogKindNameVersionKindFlag    = resourceByCatalogKindNameVersionFlags.String("kind", "REQUIRED", "kind of resource")
		resourceByCatalogKindNameVersionNameFlag    = resourceByCatalogKindNameVersionFlags.String("name", "REQUIRED", "name of resource")
		resourceByCatalogKindNameVersionVersionFlag = resourceByCatalogKindNameVersionFlags.String("version", "REQUIRED", "version of resource")

		resourceByVersionIDFlags         = flag.NewFlagSet("by-version-id", flag.ExitOnError)
		resourceByVersionIDVersionIDFlag = resourceByVersionIDFlags.String("version-id", "REQUIRED", "Version ID of a resource's version")

		resourceByCatalogKindNameFlags                = flag.NewFlagSet("by-catalog-kind-name", flag.ExitOnError)
		resourceByCatalogKindNameCatalogFlag          = resourceByCatalogKindNameFlags.String("catalog", "REQUIRED", "name of catalog")
		resourceByCatalogKindNameKindFlag             = resourceByCatalogKindNameFlags.String("kind", "REQUIRED", "kind of resource")
		resourceByCatalogKindNameNameFlag             = resourceByCatalogKindNameFlags.String("name", "REQUIRED", "Name of resource")
		resourceByCatalogKindNamePipelinesversionFlag = resourceByCatalogKindNameFlags.String("pipelinesversion", "", "")

		resourceByIDFlags  = flag.NewFlagSet("by-id", flag.ExitOnError)
		resourceByIDIDFlag = resourceByIDFlags.String("id", "REQUIRED", "ID of a resource")

		statusFlags = flag.NewFlagSet("status", flag.ContinueOnError)

		statusStatusFlags = flag.NewFlagSet("status", flag.ExitOnError)

		userFlags = flag.NewFlagSet("user", flag.ContinueOnError)

		userRefreshAccessTokenFlags            = flag.NewFlagSet("refresh-access-token", flag.ExitOnError)
		userRefreshAccessTokenRefreshTokenFlag = userRefreshAccessTokenFlags.String("refresh-token", "REQUIRED", "")

		userNewRefreshTokenFlags            = flag.NewFlagSet("new-refresh-token", flag.ExitOnError)
		userNewRefreshTokenRefreshTokenFlag = userNewRefreshTokenFlags.String("refresh-token", "REQUIRED", "")

		userInfoFlags           = flag.NewFlagSet("info", flag.ExitOnError)
		userInfoAccessTokenFlag = userInfoFlags.String("access-token", "REQUIRED", "")
	)
	adminFlags.Usage = adminUsage
	adminUpdateAgentFlags.Usage = adminUpdateAgentUsage
	adminRefreshConfigFlags.Usage = adminRefreshConfigUsage

	authFlags.Usage = authUsage
	authAuthenticateFlags.Usage = authAuthenticateUsage

	catalogFlags.Usage = catalogUsage
	catalogRefreshFlags.Usage = catalogRefreshUsage
	catalogRefreshAllFlags.Usage = catalogRefreshAllUsage

	categoryFlags.Usage = categoryUsage
	categoryListFlags.Usage = categoryListUsage

	ratingFlags.Usage = ratingUsage
	ratingGetFlags.Usage = ratingGetUsage
	ratingUpdateFlags.Usage = ratingUpdateUsage

	resourceFlags.Usage = resourceUsage
	resourceQueryFlags.Usage = resourceQueryUsage
	resourceListFlags.Usage = resourceListUsage
	resourceVersionsByIDFlags.Usage = resourceVersionsByIDUsage
	resourceByCatalogKindNameVersionFlags.Usage = resourceByCatalogKindNameVersionUsage
	resourceByVersionIDFlags.Usage = resourceByVersionIDUsage
	resourceByCatalogKindNameFlags.Usage = resourceByCatalogKindNameUsage
	resourceByIDFlags.Usage = resourceByIDUsage

	statusFlags.Usage = statusUsage
	statusStatusFlags.Usage = statusStatusUsage

	userFlags.Usage = userUsage
	userRefreshAccessTokenFlags.Usage = userRefreshAccessTokenUsage
	userNewRefreshTokenFlags.Usage = userNewRefreshTokenUsage
	userInfoFlags.Usage = userInfoUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "admin":
			svcf = adminFlags
		case "auth":
			svcf = authFlags
		case "catalog":
			svcf = catalogFlags
		case "category":
			svcf = categoryFlags
		case "rating":
			svcf = ratingFlags
		case "resource":
			svcf = resourceFlags
		case "status":
			svcf = statusFlags
		case "user":
			svcf = userFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "admin":
			switch epn {
			case "update-agent":
				epf = adminUpdateAgentFlags

			case "refresh-config":
				epf = adminRefreshConfigFlags

			}

		case "auth":
			switch epn {
			case "authenticate":
				epf = authAuthenticateFlags

			}

		case "catalog":
			switch epn {
			case "refresh":
				epf = catalogRefreshFlags

			case "refresh-all":
				epf = catalogRefreshAllFlags

			}

		case "category":
			switch epn {
			case "list":
				epf = categoryListFlags

			}

		case "rating":
			switch epn {
			case "get":
				epf = ratingGetFlags

			case "update":
				epf = ratingUpdateFlags

			}

		case "resource":
			switch epn {
			case "query":
				epf = resourceQueryFlags

			case "list":
				epf = resourceListFlags

			case "versions-by-id":
				epf = resourceVersionsByIDFlags

			case "by-catalog-kind-name-version":
				epf = resourceByCatalogKindNameVersionFlags

			case "by-version-id":
				epf = resourceByVersionIDFlags

			case "by-catalog-kind-name":
				epf = resourceByCatalogKindNameFlags

			case "by-id":
				epf = resourceByIDFlags

			}

		case "status":
			switch epn {
			case "status":
				epf = statusStatusFlags

			}

		case "user":
			switch epn {
			case "refresh-access-token":
				epf = userRefreshAccessTokenFlags

			case "new-refresh-token":
				epf = userNewRefreshTokenFlags

			case "info":
				epf = userInfoFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     interface{}
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "admin":
			c := adminc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "update-agent":
				endpoint = c.UpdateAgent()
				data, err = adminc.BuildUpdateAgentPayload(*adminUpdateAgentBodyFlag, *adminUpdateAgentTokenFlag)
			case "refresh-config":
				endpoint = c.RefreshConfig()
				data, err = adminc.BuildRefreshConfigPayload(*adminRefreshConfigTokenFlag)
			}
		case "auth":
			c := authc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "authenticate":
				endpoint = c.Authenticate()
				data, err = authc.BuildAuthenticatePayload(*authAuthenticateCodeFlag)
			}
		case "catalog":
			c := catalogc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "refresh":
				endpoint = c.Refresh()
				data, err = catalogc.BuildRefreshPayload(*catalogRefreshCatalogNameFlag, *catalogRefreshTokenFlag)
			case "refresh-all":
				endpoint = c.RefreshAll()
				data, err = catalogc.BuildRefreshAllPayload(*catalogRefreshAllTokenFlag)
			}
		case "category":
			c := categoryc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "list":
				endpoint = c.List()
				data = nil
			}
		case "rating":
			c := ratingc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "get":
				endpoint = c.Get()
				data, err = ratingc.BuildGetPayload(*ratingGetIDFlag, *ratingGetTokenFlag)
			case "update":
				endpoint = c.Update()
				data, err = ratingc.BuildUpdatePayload(*ratingUpdateBodyFlag, *ratingUpdateIDFlag, *ratingUpdateTokenFlag)
			}
		case "resource":
			c := resourcec.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "query":
				endpoint = c.Query()
				data, err = resourcec.BuildQueryPayload(*resourceQueryNameFlag, *resourceQueryCatalogsFlag, *resourceQueryKindsFlag, *resourceQueryTagsFlag, *resourceQueryLimitFlag, *resourceQueryMatchFlag)
			case "list":
				endpoint = c.List()
				data, err = resourcec.BuildListPayload(*resourceListLimitFlag)
			case "versions-by-id":
				endpoint = c.VersionsByID()
				data, err = resourcec.BuildVersionsByIDPayload(*resourceVersionsByIDIDFlag)
			case "by-catalog-kind-name-version":
				endpoint = c.ByCatalogKindNameVersion()
				data, err = resourcec.BuildByCatalogKindNameVersionPayload(*resourceByCatalogKindNameVersionCatalogFlag, *resourceByCatalogKindNameVersionKindFlag, *resourceByCatalogKindNameVersionNameFlag, *resourceByCatalogKindNameVersionVersionFlag)
			case "by-version-id":
				endpoint = c.ByVersionID()
				data, err = resourcec.BuildByVersionIDPayload(*resourceByVersionIDVersionIDFlag)
			case "by-catalog-kind-name":
				endpoint = c.ByCatalogKindName()
				data, err = resourcec.BuildByCatalogKindNamePayload(*resourceByCatalogKindNameCatalogFlag, *resourceByCatalogKindNameKindFlag, *resourceByCatalogKindNameNameFlag, *resourceByCatalogKindNamePipelinesversionFlag)
			case "by-id":
				endpoint = c.ByID()
				data, err = resourcec.BuildByIDPayload(*resourceByIDIDFlag)
			}
		case "status":
			c := statusc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "status":
				endpoint = c.Status()
				data = nil
			}
		case "user":
			c := userc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "refresh-access-token":
				endpoint = c.RefreshAccessToken()
				data, err = userc.BuildRefreshAccessTokenPayload(*userRefreshAccessTokenRefreshTokenFlag)
			case "new-refresh-token":
				endpoint = c.NewRefreshToken()
				data, err = userc.BuildNewRefreshTokenPayload(*userNewRefreshTokenRefreshTokenFlag)
			case "info":
				endpoint = c.Info()
				data, err = userc.BuildInfoPayload(*userInfoAccessTokenFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// adminUsage displays the usage of the admin command and its subcommands.
func adminUsage() {
	fmt.Fprintf(os.Stderr, `Admin service
Usage:
    %s [globalflags] admin COMMAND [flags]

COMMAND:
    update-agent: Create or Update an agent user with required scopes
    refresh-config: Refresh the changes in config file

Additional help:
    %s admin COMMAND --help
`, os.Args[0], os.Args[0])
}
func adminUpdateAgentUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] admin update-agent -body JSON -token STRING

Create or Update an agent user with required scopes
    -body JSON: 
    -token STRING: 

Example:
    `+os.Args[0]+` admin update-agent --body '{
      "name": "abc",
      "scopes": [
         "catalog-refresh",
         "agent:create"
      ]
   }' --token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODAzMDAsImlhdCI6MTU3Nzg4MDAwMCwiaWQiOjExLCJpc3MiOiJUZWt0b24gSHViIiwic2NvcGVzIjpbInJhdGluZzpyZWFkIiwicmF0aW5nOndyaXRlIiwiYWdlbnQ6Y3JlYXRlIl0sInR5cGUiOiJhY2Nlc3MtdG9rZW4ifQ.6pDmziSKkoSqI1f0rc4-AqVdcfY0Q8wA-tSLzdTCLgM"
`, os.Args[0])
}

func adminRefreshConfigUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] admin refresh-config -token STRING

Refresh the changes in config file
    -token STRING: 

Example:
    `+os.Args[0]+` admin refresh-config --token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODAzMDAsImlhdCI6MTU3Nzg4MDAwMCwiaWQiOjExLCJpc3MiOiJUZWt0b24gSHViIiwic2NvcGVzIjpbInJhdGluZzpyZWFkIiwicmF0aW5nOndyaXRlIiwiYWdlbnQ6Y3JlYXRlIl0sInR5cGUiOiJhY2Nlc3MtdG9rZW4ifQ.6pDmziSKkoSqI1f0rc4-AqVdcfY0Q8wA-tSLzdTCLgM"
`, os.Args[0])
}

// authUsage displays the usage of the auth command and its subcommands.
func authUsage() {
	fmt.Fprintf(os.Stderr, `The auth service exposes endpoint to authenticate User against GitHub OAuth
Usage:
    %s [globalflags] auth COMMAND [flags]

COMMAND:
    authenticate: Authenticates users against GitHub OAuth

Additional help:
    %s auth COMMAND --help
`, os.Args[0], os.Args[0])
}
func authAuthenticateUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] auth authenticate -code STRING

Authenticates users against GitHub OAuth
    -code STRING: 

Example:
    `+os.Args[0]+` auth authenticate --code "5628b69ec09c09512eef"
`, os.Args[0])
}

// catalogUsage displays the usage of the catalog command and its subcommands.
func catalogUsage() {
	fmt.Fprintf(os.Stderr, `The Catalog Service exposes endpoints to interact with catalogs
Usage:
    %s [globalflags] catalog COMMAND [flags]

COMMAND:
    refresh: Refresh a Catalog by it's name
    refresh-all: Refresh all catalogs

Additional help:
    %s catalog COMMAND --help
`, os.Args[0], os.Args[0])
}
func catalogRefreshUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] catalog refresh -catalog-name STRING -token STRING

Refresh a Catalog by it's name
    -catalog-name STRING: Name of catalog
    -token STRING: 

Example:
    `+os.Args[0]+` catalog refresh --catalog-name "tekton" --token "Occaecati officia inventore adipisci."
`, os.Args[0])
}

func catalogRefreshAllUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] catalog refresh-all -token STRING

Refresh all catalogs
    -token STRING: 

Example:
    `+os.Args[0]+` catalog refresh-all --token "Quaerat animi qui."
`, os.Args[0])
}

// categoryUsage displays the usage of the category command and its subcommands.
func categoryUsage() {
	fmt.Fprintf(os.Stderr, `The category service provides details about category
Usage:
    %s [globalflags] category COMMAND [flags]

COMMAND:
    list: List all categories along with their tags sorted by name

Additional help:
    %s category COMMAND --help
`, os.Args[0], os.Args[0])
}
func categoryListUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] category list

List all categories along with their tags sorted by name

Example:
    `+os.Args[0]+` category list
`, os.Args[0])
}

// ratingUsage displays the usage of the rating command and its subcommands.
func ratingUsage() {
	fmt.Fprintf(os.Stderr, `The rating service exposes endpoints to read and write user's rating for resources
Usage:
    %s [globalflags] rating COMMAND [flags]

COMMAND:
    get: Find user's rating for a resource
    update: Update user's rating for a resource

Additional help:
    %s rating COMMAND --help
`, os.Args[0], os.Args[0])
}
func ratingGetUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] rating get -id UINT -token STRING

Find user's rating for a resource
    -id UINT: ID of a resource
    -token STRING: 

Example:
    `+os.Args[0]+` rating get --id 11281538076509796713 --token "Placeat hic accusamus."
`, os.Args[0])
}

func ratingUpdateUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] rating update -body JSON -id UINT -token STRING

Update user's rating for a resource
    -body JSON: 
    -id UINT: ID of a resource
    -token STRING: 

Example:
    `+os.Args[0]+` rating update --body '{
      "rating": 3
   }' --id 13327540982832606309 --token "Quaerat ea rerum qui hic fugit inventore."
`, os.Args[0])
}

// resourceUsage displays the usage of the resource command and its subcommands.
func resourceUsage() {
	fmt.Fprintf(os.Stderr, `The resource service provides details about all kind of resources
Usage:
    %s [globalflags] resource COMMAND [flags]

COMMAND:
    query: Find resources by a combination of name, kind , catalog and tags
    list: List all resources sorted by rating and name
    versions-by-id: Find all versions of a resource by its id
    by-catalog-kind-name-version: Find resource using name of catalog & name, kind and version of resource
    by-version-id: Find a resource using its version's id
    by-catalog-kind-name: Find resources using name of catalog, resource name and kind of resource
    by-id: Find a resource using it's id

Additional help:
    %s resource COMMAND --help
`, os.Args[0], os.Args[0])
}
func resourceQueryUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] resource query -name STRING -catalogs JSON -kinds JSON -tags JSON -limit UINT -match STRING

Find resources by a combination of name, kind , catalog and tags
    -name STRING: 
    -catalogs JSON: 
    -kinds JSON: 
    -tags JSON: 
    -limit UINT: 
    -match STRING: 

Example:
    `+os.Args[0]+` resource query --name "buildah" --catalogs '[
      "tekton",
      "openshift"
   ]' --kinds '[
      "task",
      "pipelines"
   ]' --tags '[
      "image",
      "build"
   ]' --limit 100 --match "exact"
`, os.Args[0])
}

func resourceListUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] resource list -limit UINT

List all resources sorted by rating and name
    -limit UINT: 

Example:
    `+os.Args[0]+` resource list --limit 100
`, os.Args[0])
}

func resourceVersionsByIDUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] resource versions-by-id -id UINT

Find all versions of a resource by its id
    -id UINT: ID of a resource

Example:
    `+os.Args[0]+` resource versions-by-id --id 1
`, os.Args[0])
}

func resourceByCatalogKindNameVersionUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] resource by-catalog-kind-name-version -catalog STRING -kind STRING -name STRING -version STRING

Find resource using name of catalog & name, kind and version of resource
    -catalog STRING: name of catalog
    -kind STRING: kind of resource
    -name STRING: name of resource
    -version STRING: version of resource

Example:
    `+os.Args[0]+` resource by-catalog-kind-name-version --catalog "tektoncd" --kind "task" --name "buildah" --version "0.1"
`, os.Args[0])
}

func resourceByVersionIDUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] resource by-version-id -version-id UINT

Find a resource using its version's id
    -version-id UINT: Version ID of a resource's version

Example:
    `+os.Args[0]+` resource by-version-id --version-id 1
`, os.Args[0])
}

func resourceByCatalogKindNameUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] resource by-catalog-kind-name -catalog STRING -kind STRING -name STRING -pipelinesversion STRING

Find resources using name of catalog, resource name and kind of resource
    -catalog STRING: name of catalog
    -kind STRING: kind of resource
    -name STRING: Name of resource
    -pipelinesversion STRING: 

Example:
    `+os.Args[0]+` resource by-catalog-kind-name --catalog "tektoncd" --kind "pipeline" --name "buildah" --pipelinesversion "0.21.0"
`, os.Args[0])
}

func resourceByIDUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] resource by-id -id UINT

Find a resource using it's id
    -id UINT: ID of a resource

Example:
    `+os.Args[0]+` resource by-id --id 1
`, os.Args[0])
}

// statusUsage displays the usage of the status command and its subcommands.
func statusUsage() {
	fmt.Fprintf(os.Stderr, `Describes the status of each service
Usage:
    %s [globalflags] status COMMAND [flags]

COMMAND:
    status: Return status of the services

Additional help:
    %s status COMMAND --help
`, os.Args[0], os.Args[0])
}
func statusStatusUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] status status

Return status of the services

Example:
    `+os.Args[0]+` status status
`, os.Args[0])
}

// userUsage displays the usage of the user command and its subcommands.
func userUsage() {
	fmt.Fprintf(os.Stderr, `The user service exposes endpoint to get user specific specs
Usage:
    %s [globalflags] user COMMAND [flags]

COMMAND:
    refresh-access-token: Refresh the access token of User
    new-refresh-token: Get a new refresh token of User
    info: Get the user Info

Additional help:
    %s user COMMAND --help
`, os.Args[0], os.Args[0])
}
func userRefreshAccessTokenUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] user refresh-access-token -refresh-token STRING

Refresh the access token of User
    -refresh-token STRING: 

Example:
    `+os.Args[0]+` user refresh-access-token --refresh-token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODM2MDAsImlhdCI6MTU3Nzg4MDAwMCwiaWQiOjExLCJpc3MiOiJUZWt0b24gSHViIiwic2NvcGVzIjpbInJlZnJlc2g6dG9rZW4iXSwidHlwZSI6InJlZnJlc2gtdG9rZW4ifQ.4RdUk5ttHdDiymurlZ_f7Uy5Pas3Lq9w04BjKQKRiCE"
`, os.Args[0])
}

func userNewRefreshTokenUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] user new-refresh-token -refresh-token STRING

Get a new refresh token of User
    -refresh-token STRING: 

Example:
    `+os.Args[0]+` user new-refresh-token --refresh-token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODM2MDAsImlhdCI6MTU3Nzg4MDAwMCwiaWQiOjExLCJpc3MiOiJUZWt0b24gSHViIiwic2NvcGVzIjpbInJlZnJlc2g6dG9rZW4iXSwidHlwZSI6InJlZnJlc2gtdG9rZW4ifQ.4RdUk5ttHdDiymurlZ_f7Uy5Pas3Lq9w04BjKQKRiCE"
`, os.Args[0])
}

func userInfoUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] user info -access-token STRING

Get the user Info
    -access-token STRING: 

Example:
    `+os.Args[0]+` user info --access-token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODM2MDAsImlhdCI6MTU3Nzg4MDAwMCwiaWQiOjExLCJpc3MiOiJUZWt0b24gSHViIiwic2NvcGVzIjpbInJlZnJlc2g6dG9rZW4iXSwidHlwZSI6InJlZnJlc2gtdG9rZW4ifQ.4RdUk5ttHdDiymurlZ_f7Uy5Pas3Lq9w04BjKQKRiCE"
`, os.Args[0])
}
