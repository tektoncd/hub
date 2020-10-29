module github.com/tektoncd/hub/api

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-gormigrate/gormigrate/v2 v2.0.0
	github.com/go-testfixtures/testfixtures/v3 v3.2.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/ikawaha/goahttpcheck v1.3.1
	github.com/joho/godotenv v1.3.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.5.1
	github.com/tektoncd/pipeline v0.17.1-0.20201007165454-9611f3e4509e
	go.uber.org/zap v1.15.0
	goa.design/goa/v3 v3.2.2
	goa.design/plugins/v3 v3.1.3
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	gopkg.in/h2non/gock.v1 v1.0.15
	gorm.io/driver/postgres v1.0.2
	gorm.io/gorm v1.20.7
	gotest.tools/v3 v3.0.2
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/pkg v0.0.0-20200922164940-4bf40ad82aab
)

// Pin k8s deps to 0.18.9
replace (
	k8s.io/api => k8s.io/api v0.18.9
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.9
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.9
	k8s.io/apiserver => k8s.io/apiserver v0.18.9
	k8s.io/client-go => k8s.io/client-go v0.18.9
	k8s.io/code-generator => k8s.io/code-generator v0.18.9
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20200410145947-bcb3869e6f29
)

// Use the same Cobra as upstream tkn cli
replace github.com/spf13/cobra => github.com/chmouel/cobra v0.0.0-20200107083527-379e7a80af0c
