package conf

import (
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

const (
	hostKey        = "WORDSHUB_HOST"
	portKey        = "WORDSHUB_PORT"
	dbHostKey      = "WORDSHUB_DB_HOST"
	dbPortKey      = "WORDSHUB_DB_PORT"
	dbNameKey      = "WORDSHUB_DB_NAME"
	dbUserKey      = "WORDSHUB_DB_USER"
	dbPasswordKey  = "WORDSHUB_DB_PASSWORD"
	jwtSecretKey   = "WORDSHUB_JWT_SECRET"
	appIdKey       = "APP_ID"
	appSecretKey   = "APPSECRET"
	appCloudEnvKey = "APP_CLOUD_ENV"
)

type Config struct {
	Host        string
	Port        string
	DbHost      string
	DbPort      string
	DbName      string
	DbUser      string
	DbPassword  string
	JwtSecret   string
	Env         string
	AppId       string
	AppSecret   string
	AppCloudEnv string
}

func NewConfig(env string) Config {
	host, ok := os.LookupEnv(hostKey)
	if !ok || host == "" {
		logAndPanic(hostKey)
	}

	port, ok := os.LookupEnv(portKey)
	if !ok || port == "" {
		if _, err := strconv.Atoi(port); err != nil {
			logAndPanic(portKey)
		}
	}

	dbHost, ok := os.LookupEnv(dbHostKey)
	if !ok || dbHost == "" {
		logAndPanic(dbHostKey)
	}

	dbPort, ok := os.LookupEnv(dbPortKey)
	if !ok || dbPort == "" {
		if _, err := strconv.Atoi(dbPort); err != nil {
			logAndPanic(dbPortKey)
		}
	}

	dbName, ok := os.LookupEnv(dbNameKey)
	if !ok || dbName == "" {
		logAndPanic(dbNameKey)
	}

	dbUser, ok := os.LookupEnv(dbUserKey)
	if !ok || dbUser == "" {
		logAndPanic(dbUserKey)
	}

	dbPassword, ok := os.LookupEnv(dbPasswordKey)
	if !ok || dbPassword == "" {
		logAndPanic(dbPasswordKey)
	}

	jwtSecret, ok := os.LookupEnv(jwtSecretKey)
	if !ok || jwtSecret == "" {
		logAndPanic(jwtSecretKey)
	}

	appId, ok := os.LookupEnv(appIdKey)
	if !ok || jwtSecret == "" {
		logAndPanic(appIdKey)
	}

	appSecret, ok := os.LookupEnv(appSecretKey)
	if !ok || jwtSecret == "" {
		logAndPanic(appSecretKey)
	}

	appCloudEnv, ok := os.LookupEnv(appCloudEnvKey)
	if !ok || appCloudEnv == "" {
		logAndPanic(appCloudEnvKey)
	}

	return Config{
		Host:        host,
		Port:        port,
		DbHost:      dbHost,
		DbPort:      dbPort,
		DbName:      dbName,
		DbUser:      dbUser,
		DbPassword:  dbPassword,
		JwtSecret:   jwtSecret,
		Env:         env,
		AppId:       appId,
		AppSecret:   appSecret,
		AppCloudEnv: appCloudEnv,
	}
}

func NewTestConfig() Config {
	testConfig := NewConfig("dev")
	testConfig.DbName = testConfig.DbName + "_test"
	return testConfig
}

func logAndPanic(envVar string) {
	log.Panic().Str("envVar", envVar).Msg("ENV variable not set or value not valid")
}
