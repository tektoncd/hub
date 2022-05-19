module github.com/tektoncd/hub

go 1.16

require (
	github.com/ActiveState/vt10x v1.3.1
	github.com/AlecAivazis/survey/v2 v2.3.4
	github.com/Netflix/go-expect v0.0.0-20220104043353-73e0943537d2
	github.com/fatih/color v1.13.0
	github.com/go-gormigrate/gormigrate/v2 v2.0.1
	github.com/go-testfixtures/testfixtures/v3 v3.6.2
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1
	github.com/hako/durafmt v0.0.0-20210608085754-5c1018a4e16b
	github.com/ikawaha/goahttpcheck v1.10.0
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.4.0
	github.com/markbates/goth v1.72.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.11.0
	github.com/stretchr/testify v1.7.1
	github.com/tektoncd/pipeline v0.34.1
	github.com/tektoncd/plumbing v0.0.0-20211012143332-c7cc43d9bc0c
	go.uber.org/automaxprocs v1.5.1
	go.uber.org/zap v1.21.0
	goa.design/goa/v3 v3.7.5
	goa.design/plugins/v3 v3.7.5
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
	gopkg.in/h2non/gock.v1 v1.1.2
	gorm.io/driver/postgres v1.3.5
	gorm.io/gorm v1.23.4
	gotest.tools/v3 v3.2.0
	k8s.io/apimachinery v0.23.4
	k8s.io/client-go v1.5.2
	knative.dev/pkg v0.0.0-20220131144930-f4b57aef0006
)

replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.5
	k8s.io/client-go => k8s.io/client-go v0.22.5
)
