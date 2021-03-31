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
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/tektoncd/pipeline v0.20.1-0.20210204110343-8c5a751b53ea
	go.uber.org/zap v1.16.0
	goa.design/goa/v3 v3.2.2
	goa.design/plugins/v3 v3.1.3
	golang.org/x/oauth2 v0.0.0-20210126194326-f9ce19ea3013
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf
	gopkg.in/h2non/gock.v1 v1.0.15
	gorm.io/driver/postgres v1.0.2
	gorm.io/gorm v1.20.7
	gotest.tools/v3 v3.0.2
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	knative.dev/pkg v0.0.0-20210203171706-6045ed499615
	maze.io/x/duration v0.0.0-20160924141736-faac084b6075
)
