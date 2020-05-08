package app

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2"

	// Blank for package side effect: loads postgres drivers
	_ "github.com/lib/pq"
)

type Base interface {
	Environment() EnvMode
	Logger() *zap.SugaredLogger
	Database() *Database
	DB() *gorm.DB
	Cleanup()
}

type Config interface {
	Base
	GitHub() *GitHub
	Addr() string
}

type EnvMode string

const (
	Production  EnvMode = "production"
	Development EnvMode = "development"
	Test        EnvMode = "test"
)

type Database struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func (db *Database) String() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=xxxxxx dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Name)
}

func (db *Database) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Name)
}

type GitHub struct {
	AccessToken   string
	OAuthClientID string
	OAuthSecret   string
	JWTSigningKey string
	Client        *github.Client
}

type BaseConfig struct {
	mode   EnvMode
	logger *zap.SugaredLogger
	dbConf *Database
	db     *gorm.DB
}

var _ Base = (*BaseConfig)(nil)

func (bc *BaseConfig) Environment() EnvMode {
	return bc.mode
}

func (bc *BaseConfig) Logger() *zap.SugaredLogger {
	return bc.logger
}

func (bc *BaseConfig) Database() *Database {
	return bc.dbConf
}

func (bc *BaseConfig) DB() *gorm.DB {
	return bc.db
}

func (bc *BaseConfig) Cleanup() {
	bc.db.Close()
	bc.logger.Sync()
}

type ApiConfig struct {
	*BaseConfig
	gh *GitHub
}

type TestConfig struct {
	*BaseConfig
}

var _ Config = (*ApiConfig)(nil)

func (e *ApiConfig) GitHub() *GitHub {
	return e.gh
}

func (e *ApiConfig) Addr() string {
	return ":5000"
}

func BaseConfigFromEnv() (*BaseConfig, error) {
	// load from .env file but skip if not found
	if err := godotenv.Load(); err != nil {
		fmt.Fprintf(os.Stdout, "SKIP: loading .ApiConfig failed: %s", err)
	}

	mode := Environment()
	var err error

	var log *zap.SugaredLogger
	if log, err = initLogger(mode); err != nil {
		return nil, err
	}

	log.With("name", "app").Infof("in %q mode ", mode)

	bc := &BaseConfig{mode: mode, logger: log}
	if bc.dbConf, err = initDB(mode); err != nil {
		log.Error(err, "failed to obtain database configuration")
		return nil, err
	}

	bc.db, err = gorm.Open("postgres", bc.dbConf.ConnectionString())
	if err != nil {
		log.Error(err, "failed to establish database connection")
		return nil, err
	}

	log.Infof("Successfully connected to db %s", bc.dbConf)

	return bc, nil
}

func TestConfigFromEnv() (*TestConfig, error) {
	bc, err := BaseConfigFromEnv()
	if err != nil {
		return nil, err
	}

	TestConfig := &TestConfig{BaseConfig: bc}

	return TestConfig, nil
}

func FromEnv() (*ApiConfig, error) {
	bc, err := BaseConfigFromEnv()
	if err != nil {
		return nil, err
	}

	ApiConfig := &ApiConfig{BaseConfig: bc}

	if ApiConfig.gh, err = initGithub(); err != nil {
		return nil, err
	}

	return ApiConfig, nil
}

func env(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("NO %q environment variable defined", key)
	}
	return val, nil
}

func Environment() EnvMode {
	mode := "production"
	if val, ok := os.LookupEnv("ENVIRONMENT"); ok {
		mode = val
	}

	switch strings.ToLower(mode) {
	case "development":
		return Development
	case "test":
		return Test
	default:
		return Production
	}
}

func initDB(mode EnvMode) (*Database, error) {
	var err error

	db := &Database{}

	if mode == Test {

		if db.Host, err = env("TEST_POSTGRESQL_HOST"); err != nil {
			return nil, err
		}
		if db.Port, err = env("TEST_POSTGRESQL_PORT"); err != nil {
			return nil, err
		}
		if db.Name, err = env("TEST_POSTGRESQL_DATABASE"); err != nil {
			return nil, err
		}
		if db.User, err = env("TEST_POSTGRESQL_USER"); err != nil {
			return nil, err
		}
		if db.Password, err = env("TEST_POSTGRESQL_PASSWORD"); err != nil {
			return nil, err
		}

		return db, nil
	}

	if db.Host, err = env("POSTGRESQL_HOST"); err != nil {
		return nil, err
	}
	if db.Port, err = env("POSTGRESQL_PORT"); err != nil {
		return nil, err
	}
	if db.Name, err = env("POSTGRESQL_DATABASE"); err != nil {
		return nil, err
	}
	if db.User, err = env("POSTGRESQL_USER"); err != nil {
		return nil, err
	}
	if db.Password, err = env("POSTGRESQL_PASSWORD"); err != nil {
		return nil, err
	}

	return db, nil
}

func initGithub() (*GitHub, error) {
	var err error
	gh := &GitHub{}
	if gh.AccessToken, err = env("GITHUB_TOKEN"); err != nil {
		return nil, err
	}
	if gh.OAuthClientID, err = env("CLIENT_ID"); err != nil {
		return nil, err
	}
	if gh.OAuthSecret, err = env("CLIENT_SECRET"); err != nil {
		return nil, err
	}
	if gh.JWTSigningKey, err = env("JWT_SIGNING_KEY"); err != nil {
		return nil, err
	}

	token := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: gh.AccessToken})
	client := oauth2.NewClient(context.Background(), token)
	gh.Client = github.NewClient(client)
	return gh, nil
}

func initLogger(mode EnvMode) (*zap.SugaredLogger, error) {

	var log *zap.Logger
	var err error

	switch mode {
	case Production:
		log, err = zap.NewProduction()

	default:
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		log, err = config.Build()
	}

	if err != nil {
		return nil, err
	}
	return log.Sugar(), nil
}
