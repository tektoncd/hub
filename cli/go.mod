module github.com/tektoncd/hub/cli

go 1.14

require (
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.5.1
	github.com/tektoncd/hub/api v0.0.0
	goa.design/goa/v3 v3.2.2
	gopkg.in/h2non/gock.v1 v1.0.15
	gotest.tools v2.2.0+incompatible
)

replace github.com/tektoncd/hub/api => ../api
