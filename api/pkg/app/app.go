// Copyright © 2020 The Tekton Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/tektoncd/hub/api/gen/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"maze.io/x/duration"
)

// BaseConfig defines methods on APIBase
type BaseConfig interface {
	Environment() EnvMode
	Service(name string) Service
	Logger(service string) *log.Logger
	DB() *gorm.DB
	Data() *Data
	ReloadData() error
	Cleanup()
}

// APIBase defines the base configuration every service requires
type APIBase struct {
	env    EnvMode
	dbConf *Database
	db     *gorm.DB
	logger *log.Logger
	data   Data
}

// Config defines methods on APIConfig includes BaseConfig
type Config interface {
	BaseConfig
	OAuthConfig() *oauth2.Config
	JWTConfig() *JWTConfig
	GhConfig() *GHConfig
}

// APIConfig defines struct on top of APIBase with GitHub Oauth,
// GHEConfig & JWT Configurations
type APIConfig struct {
	*APIBase
	conf      *oauth2.Config
	jwtConfig *JWTConfig
	ghConfig  *GHConfig
}

// GHConfig struct defines the github configuration
type GHConfig struct {
	IsGhe     bool
	Url       string
	ApiUrl    string
	UploadUrl string
}

// JWTConfig defines configuration requires to create token
type JWTConfig struct {
	SigningKey       string
	AccessExpiresIn  time.Duration
	RefreshExpiresIn time.Duration
}

var _ BaseConfig = (*APIBase)(nil)

var _ Config = (*APIConfig)(nil)

// EnvMode defines the mode the server is running in
type EnvMode string

// Types of EnvMode
const (
	Production  EnvMode = "production"
	Development EnvMode = "development"
	Test        EnvMode = "test"
)

// DBDialect defines dialect for db connection
const DBDialect = "postgres"

// Database Object defines db configuration fields
type Database struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func (db *Database) String() string {
	return fmt.Sprintf("database=%s user=%s host=%s:%s", db.Name, db.User, db.Host, db.Port)
}

// ConnectionString returns the db connection string
func (db Database) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Name)
}

// Environment returns the EnvMode server would be running
func (ab *APIBase) Environment() EnvMode {
	return ab.env
}

// DB returns gorm db object
func (ab *APIBase) DB() *gorm.DB {
	return ab.db
}

// DBWithLogger returns gorm db object initialised with logger
func DBWithLogger(env EnvMode, db *gorm.DB, logger *log.Logger) *gorm.DB {
	db = db.Session(&gorm.Session{Logger: newGormLogger(env, logger)})
	return db
}

// Database returns the database object used for initializing db connection
func (ab *APIBase) Database() Database {
	return *ab.dbConf
}

// Logger returns log.Logger appended with service name
func (ab *APIBase) Logger(service string) *log.Logger {
	return &log.Logger{
		SugaredLogger: ab.logger.With(zap.String("service", service)),
	}
}

// Service creates a base service object
func (ab *APIBase) Service(name string) Service {
	l := &log.Logger{
		SugaredLogger: ab.logger.With(zap.String("service", name)),
	}
	return &BaseService{
		logger: l,
		db:     ab.DB(),
		env:    ab,
	}
}

// Data returns Data object which consist app data from config file
func (ab *APIBase) Data() *Data {
	return &ab.data
}

// ReloadData reads config file and loads data in Data object
func (ab *APIBase) ReloadData() error {
	// Reads config file url from env
	url, err := configFileURL()
	if err != nil {
		return err
	}

	// Reads data from config file
	fileData, err := dataFromURL(url)
	if err != nil {
		ab.logger.Errorf("failed to read config file: %v", err)
		return err
	}

	// Viper unmarshals data from config file into Data Object
	var data Data
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer(fileData)); err != nil {
		ab.logger.Errorf("failed to read configuration file: %v", err)
		return err
	}
	if err := viper.Unmarshal(&data); err != nil {
		ab.logger.Errorf("failed to unmarshal config data: %v", err)
		return err
	}
	ab.data = data

	// computes checksum on config data
	hash := sha256.Sum256(fileData)
	ab.data.Checksum = hex.EncodeToString(hash[:])

	return nil
}

// Cleanup flushes any buffered log entries & closes the db connection
func (ab *APIBase) Cleanup() {
	// TODO --> Handle the error
	_ = ab.logger.Sync()
	db, _ := ab.db.DB()
	db.Close()
}

// OAuthConfig returns oauth2 config object
func (ac *APIConfig) OAuthConfig() *oauth2.Config {
	return ac.conf
}

// GheConfig returns Github Enterprise object which stores
// whether GHE url is present or not and some other urls which
// are generated on the basis of GHE url
func (ac *APIConfig) GhConfig() *GHConfig {
	return ac.ghConfig
}

// JWTConfig returns JWTConfig Object
func (ac *APIConfig) JWTConfig() *JWTConfig {
	return ac.jwtConfig
}

// FromEnv will initialise APIConfig Object. This is called while starting
// the api server. It passes .env.dev which contains configurations for
// development mode, if it doesn't find the file it skips it and will look
// for configration among env variable
func FromEnv() (*APIConfig, error) {
	// load from .env.dev file for development but skip if not found
	return FromEnvFile(".env.dev")
}

// FromEnvFile expects a file name containing configurations. This is called
// when for running test where test config file is passed to initialise a
// APIConfig Object.
func FromEnvFile(file string) (*APIConfig, error) {
	ab, err := APIBaseFromEnvFile(file)
	if err != nil {
		return nil, err
	}

	err = ab.ReloadData()
	if err != nil {
		return nil, err
	}

	ac := &APIConfig{APIBase: ab}

	if ac.ghConfig, err = initGh(); err != nil {
		return nil, err
	}

	if ac.conf, err = initOAuthConfig(ac.ghConfig.Url); err != nil {
		return nil, err
	}
	if ac.jwtConfig, err = jwtConfig(); err != nil {
		return nil, err
	}

	return ac, nil
}

// APIBaseFromEnv initialises APIBase Object passing .env.dev file to
// APIBaseFromEnvFile which will have configuration for development mode.
// This will initialise db connection and logger only. This is called while
// running db migration.
func APIBaseFromEnv() (*APIBase, error) {
	// load from .env.dev file for development but skip if not found
	return APIBaseFromEnvFile(".env.dev")
}

// APIBaseFromEnvFile expects a filepath to env file which has configurations
// It loads .env file, skips it if not found, initialises a db connection &
// logger depending on the EnvMode and returns a APIBase Object.
func APIBaseFromEnvFile(file string) (*APIBase, error) {
	if err := godotenv.Load(file); err != nil {
		fmt.Fprintf(os.Stderr, "SKIP: loading env file %s failed: %s\n", file, err)
	}

	// Enables viper to read Environment Variables
	// NOTE: DO NOT move this line; viper must be initialized before reading ENV variables
	viper.AutomaticEnv()

	env := Environment()

	var err error
	var l *log.Logger
	if l, err = initLogger(env); err != nil {
		return nil, err
	}

	ab := &APIBase{env: env, logger: l}
	log := ab.logger.With("app", "hub")

	log.Infof("in %q mode ", env)

	if ab.dbConf, err = initDB(); err != nil {
		log.Errorf("failed to obtain database configuration: %v", err)
		return nil, err
	}

	ab.db, err = gorm.Open(postgres.Open(ab.dbConf.ConnectionString()), &gorm.Config{
		Logger: newGormLogger(env, ab.logger),
	})
	if err != nil {
		log.Errorf("failed to establish database connection: [%s]: %s", ab.dbConf, err)
		return nil, err
	}
	log.Infof("Successfully connected to [%s]", ab.dbConf)

	return ab, nil
}

// Environment return EnvMode the Api server would be running in.
// It looks for 'ENVIRONMENT' to be defined as environment variable and
// if does not found it then set it as development mode
func Environment() EnvMode {
	mode := "production"
	if val := viper.GetString("ENVIRONMENT"); val != "" {
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

// initDB looks for db credentials in environment variables and returns as Database object
// if it does not find a field then returns error
func initDB() (*Database, error) {

	db := &Database{}
	if db.Host = viper.GetString("POSTGRES_HOST"); db.Host == "" {
		return nil, fmt.Errorf("no POSTGRES_HOST environment variable defined")
	}
	if db.Port = viper.GetString("POSTGRES_PORT"); db.Port == "" {
		return nil, fmt.Errorf("no POSTGRES_PORT environment variable defined")
	}
	if db.Name = viper.GetString("POSTGRES_DB"); db.Name == "" {
		return nil, fmt.Errorf("no POSTGRES_DB environment variable defined")
	}
	if db.User = viper.GetString("POSTGRES_USER"); db.User == "" {
		return nil, fmt.Errorf("no POSTGRES_USER environment variable defined")
	}
	if db.Password = viper.GetString("POSTGRES_PASSWORD"); db.Password == "" {
		return nil, fmt.Errorf("no POSTGRES_PASSWORD environment variable defined")
	}
	return db, nil
}

// initLogger returns a instance of log.Logger depending on the EnvMode
func initLogger(mode EnvMode) (*log.Logger, error) {

	var l *zap.Logger
	var err error

	switch mode {
	case Production:
		l, err = zap.NewProduction()
	default:
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		l, err = config.Build()
	}

	if err != nil {
		return nil, err
	}
	return &log.Logger{SugaredLogger: l.Sugar()}, nil
}

// configFileURL will look for CONFIG_FILE_URL to be defined among
// environment variables
func configFileURL() (string, error) {

	val := viper.GetString("CONFIG_FILE_URL")
	if val == "" {
		return "", fmt.Errorf("no CONFIG_FILE_URL environment variable defined")
	}
	return val, nil
}

// initGh looks for Github Enterprise url from the environment variables
// and initialises the GHEConfig
func initGh() (*GHConfig, error) {
	ghe := &GHConfig{}
	if ghe.Url = viper.GetString("GHE_URL"); ghe.Url == "" {
		ghe.Url = "https://github.com"
	}
	if !strings.HasPrefix(ghe.Url, "https://github.com") {
		parsedUrl, err := url.Parse(ghe.Url)
		if err != nil {
			return nil, fmt.Errorf("There was some problem while parsing the Github Enterprise URL")
		}

		if parsedUrl.Path != "" || parsedUrl.RawQuery != "" {
			return nil, fmt.Errorf("Invalid Github Enterprise URL")
		}

		ghe.IsGhe = true
		ghe.ApiUrl = fmt.Sprintf("%s://api.%s", parsedUrl.Scheme, parsedUrl.Host)        // https://api.myghe.com
		ghe.UploadUrl = fmt.Sprintf("%s://uploads.%s", parsedUrl.Scheme, parsedUrl.Host) // https://uploads.myghe.com

	}

	return ghe, nil
}

// initOAuthConfig looks for configuration among environment variables
// and intialises the GitHub Oauth Config on the basis of github url
func initOAuthConfig(ghUrl string) (*oauth2.Config, error) {

	var clientID, clientSecret string
	if clientID = viper.GetString("GH_CLIENT_ID"); clientID == "" {
		return nil, fmt.Errorf("no GH_CLIENT_ID environment variable defined")
	}
	if clientSecret = viper.GetString("GH_CLIENT_SECRET"); clientSecret == "" {
		return nil, fmt.Errorf("no GH_CLIENT_SECRET environment variable defined")
	}

	gheEndpoint := oauth2.Endpoint{
		AuthURL:  fmt.Sprintf("%s/oauth/authorize", ghUrl),
		TokenURL: fmt.Sprintf("%s/login/oauth/access_token", ghUrl),
	}

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     gheEndpoint,
	}
	return conf, nil
}

// jwtConfig will look for jwt configurations to be defined among
// environment variables
func jwtConfig() (*JWTConfig, error) {

	conf := &JWTConfig{}

	conf.SigningKey = viper.GetString("JWT_SIGNING_KEY")
	if conf.SigningKey == "" {
		return nil, fmt.Errorf("no JWT_SIGNING_KEY environment variable defined")
	}

	accessExpiresIn := viper.GetString("ACCESS_JWT_EXPIRES_IN")
	if accessExpiresIn == "" {
		return nil, fmt.Errorf("no ACCESS_JWT_EXPIRES_IN environment variable defined")
	}
	var err error
	if conf.AccessExpiresIn, err = computeDuration(accessExpiresIn); err != nil {
		return nil, fmt.Errorf("invalid time format specified for ACCESS_JWT_EXPIRES_IN: %v", err)
	}

	refreshExpiresIn := viper.GetString("REFRESH_JWT_EXPIRES_IN")
	if refreshExpiresIn == "" {
		return nil, fmt.Errorf("no REFRESH_JWT_EXPIRES_IN environment variable defined")
	}
	if conf.RefreshExpiresIn, err = computeDuration(refreshExpiresIn); err != nil {
		return nil, fmt.Errorf("invalid time format specified for REFRESH_JWT_EXPIRES_IN: %v", err)
	}

	return conf, nil
}

func computeDuration(dur string) (time.Duration, error) {
	d, err := duration.ParseDuration(dur)
	if err != nil {
		return 0, err
	}
	return time.Duration(d), nil
}
