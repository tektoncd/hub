module github.com/tektoncd/hub/api

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-testfixtures/testfixtures/v3 v3.2.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/ikawaha/goahttpcheck v1.3.1
	github.com/jinzhu/gorm v1.9.12
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.3.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.5.1
	github.com/tektoncd/pipeline v0.15.2
	go.uber.org/zap v1.15.0
	goa.design/goa/v3 v3.2.2
	goa.design/plugins/v3 v3.1.3
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200812155832-6a926be9bd1d // indirect
	golang.org/x/tools v0.0.0-20200811215021-48a8ffc5b207 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/gormigrate.v1 v1.6.0
	gopkg.in/h2non/gock.v1 v1.0.15
	gotest.tools/v3 v3.0.2
	k8s.io/apimachinery v0.17.6
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/pkg v0.0.0-20200702222342-ea4d6e985ba0
)

// Copied from Catlin
// Pin k8s deps to 1.17.6
replace (
	k8s.io/api => k8s.io/api v0.17.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/apiserver => k8s.io/apiserver v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
	k8s.io/code-generator => k8s.io/code-generator v0.17.6
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20200410145947-bcb3869e6f29
)

// Use the same Cobra as upstream tkn cli
replace github.com/spf13/cobra => github.com/chmouel/cobra v0.0.0-20200107083527-379e7a80af0c
