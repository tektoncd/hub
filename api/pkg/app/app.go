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
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	// Blank for package side effect: loads postgres drivers
	_ "github.com/lib/pq"
)

// Config defines methods on APIConfig
type Config interface {
	Environment() EnvMode
	Logger() *zap.SugaredLogger
	DB() *gorm.DB
	Cleanup()
}

// APIConfig defines the configuration a services requires
type APIConfig struct {
	mode   EnvMode
	dbConf *Database
	db     *gorm.DB
	logger *zap.SugaredLogger
}

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
func (db *Database) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Name)
}

// Environment returns the EnvMode server would be running
func (ac *APIConfig) Environment() EnvMode {
	return ac.mode
}

// Database returns Database object which consist of db configurations
func (ac *APIConfig) Database() *Database {
	return ac.dbConf
}

// DB returns gorm db object
func (ac *APIConfig) DB() *gorm.DB {
	return ac.db
}

// Logger returns suggared logger object
func (ac *APIConfig) Logger() *zap.SugaredLogger {
	return ac.logger
}

// Cleanup flushes any buffered log entries & closes the db connection
func (ac *APIConfig) Cleanup() {
	ac.logger.Sync()
	ac.db.Close()
}

// FromEnv is called while initailising the api service, it calls FromEnvFile
// passing .env.dev file which will have configuration while running in
// development mode
func FromEnv() (*APIConfig, error) {
	// load from .env.dev file for development but skip if not found
	return FromEnvFile(".env.dev")
}

// FromEnvFile expects a filepath to env file which has db configurations
// It loads .env file, initialises a db connection & logger depending on the EnvMode
// and returns a APIConfig Object
// If it doesn't finds a .env file, it looks for configuratin among environment variables
func FromEnvFile(file string) (*APIConfig, error) {
	if err := godotenv.Load(file); err != nil {
		fmt.Fprintf(os.Stderr, "SKIP: loading env file %s failed: %s\n", file, err)
	}
	mode := Environment()
	var err error

	var log *zap.SugaredLogger
	if log, err = initLogger(mode); err != nil {
		return nil, err
	}

	log.With("name", "app").Infof("in %q mode ", mode)

	ac := &APIConfig{mode: mode, logger: log}
	if ac.dbConf, err = initDB(); err != nil {
		log.Error(err, "failed to obtain database configuration")
		return nil, err
	}
	ac.db, err = gorm.Open(DBDialect, ac.dbConf.ConnectionString())
	if err != nil {
		log.Errorf("failed to establish database connection: [%s]: %s", ac.dbConf, err)
		return nil, err
	}
	log.Infof("Successfully connected to [%s]", ac.dbConf)

	return ac, nil
}

// env look for the input key to be defined as a environment variable
// returns error if the key is not found
func env(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("NO %q environment variable defined", key)
	}
	return val, nil
}

// Environment return EnvMode the Api server would be running in.
// It looks for 'ENVIRONMENT' to be defined as environment variable and
// if does not found it then set it as development
func Environment() EnvMode {
	mode := "development"
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

// initDB looks for db credentials in environment variables and returns as Database object
// if it does not find a field then returns error
func initDB() (*Database, error) {
	var err error
	db := &Database{}
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

// initLogger returns a instance of SugaredLogger depending on the EnvMode
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
