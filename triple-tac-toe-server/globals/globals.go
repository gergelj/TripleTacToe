package globals

import (
	"os"
)

const (
	portDefault           = "9095"
	dbHostDefault         = "localhost"
	dbPortDefault         = "5432"
	dbUserDefault         = "postgres"
	dbPasswordDefault     = "root"
	dbNameDefault         = "tripletactoe"
	jwtTokenSecretDefault = "12345"
	jwtTokenLifeLength    = "24"
	devEnvDefault         = "true"
)

var (
	Port               string = loadEnvValue("PORT", portDefault)
	DbPort             string = loadEnvValue("DB_PORT", dbPortDefault)
	DbHost             string = loadEnvValue("DB_HOST", dbHostDefault)
	DbUser             string = loadEnvValue("DB_USER", dbUserDefault)
	DbPassword         string = loadEnvValue("DB_PASSWORD", dbPasswordDefault)
	DbName             string = loadEnvValue("DB_NAME", dbNameDefault)
	JwtTokenSecret     string = loadEnvValue("JWT_SECRET", jwtTokenSecretDefault)
	JwtTokenLifeLength string = loadEnvValue("JWT_TOKEN_LIFE_LENGTH", jwtTokenLifeLength) //Hours
	DevEnvironment     string = loadEnvValue("DEV_ENV", devEnvDefault)
)

func loadEnvValue(envName string, defaultValue string) string {
	if val, present := os.LookupEnv(envName); present {
		return val
	} else {
		return defaultValue
	}
}

func IsDevEnvironment() bool {
	return DevEnvironment == "true" || DevEnvironment == "TRUE"
}
