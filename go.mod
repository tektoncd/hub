module github.com/tektoncd/hub

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.2.12
	github.com/Netflix/go-expect v0.0.0-20200312175327-da48e75238e2
	github.com/fatih/color v1.13.0
	github.com/go-gormigrate/gormigrate/v2 v2.0.0
	github.com/go-testfixtures/testfixtures/v3 v3.2.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1
	github.com/hako/durafmt v0.0.0-20210608085754-5c1018a4e16b
	github.com/hinshun/vt10x v0.0.0-20180616224451-1954e6464174
	github.com/ikawaha/goahttpcheck v1.3.1
	github.com/joho/godotenv v1.3.0
	github.com/markbates/goth v1.68.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.0
	github.com/stretchr/testify v1.7.1
	github.com/tektoncd/pipeline v0.33.1
	github.com/tektoncd/plumbing v0.0.0-20211012143332-c7cc43d9bc0c
	go.uber.org/automaxprocs v1.4.0
	go.uber.org/zap v1.21.0
	goa.design/goa/v3 v3.4.0
	goa.design/plugins/v3 v3.1.3
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
	gopkg.in/h2non/gock.v1 v1.0.16
	gorm.io/driver/postgres v1.0.2
	gorm.io/gorm v1.20.7
	gotest.tools/v3 v3.1.0
	k8s.io/apimachinery v0.23.4
	k8s.io/client-go v1.5.2
	knative.dev/pkg v0.0.0-20220131144930-f4b57aef0006
)

replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.5
	k8s.io/client-go => k8s.io/client-go v0.22.5
)
