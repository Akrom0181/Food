package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

const (
	// DebugMode indicates service mode is debug.
	DebugMode = "debug"
	// TestMode indicates service mode is test.
	TestMode = "test"
	// ReleaseMode indicates service mode is release.
	ReleaseMode = "release"

	LocalMode = "local"
)

type Config struct {
	ServiceName string
	Environment string
	Version     string

	HTTPPort   string
	HTTPScheme string

	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string

	PostgresMaxConnections int32
	DefaultOffset          string
	DefaultLimit           string

	RedisHost     string
	RedisPort     string
	RedisPassword string

	RedisURL string

	SecretKey string

	AuthServiceHost string
	AuthGRPCPort    string
}

// Load ...
func Load() Config {

	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file found")
	}

	config := Config{}

	config.DefaultOffset = cast.ToString(getOrReturnDefaultValue("DEFAULT_OFFSET", "0"))
	config.DefaultLimit = cast.ToString(getOrReturnDefaultValue("DEFAULT_LIMIT", "100"))

	config.ServiceName = cast.ToString(getOrReturnDefaultValue("SERVICE_NAME", "food"))
	config.Environment = cast.ToString(getOrReturnDefaultValue("ENVIRONMENT", DebugMode))
	config.Version = cast.ToString(getOrReturnDefaultValue("VERSION", "1.0"))

	config.HTTPPort = cast.ToString(getOrReturnDefaultValue("HTTP_PORT", ":8000"))
	config.HTTPScheme = cast.ToString(getOrReturnDefaultValue("HTTP_SCHEME", "http"))

	config.PostgresHost = cast.ToString(os.Getenv("POSTGRES_HOST"))
	config.PostgresPort = cast.ToInt(os.Getenv("POSTGRES_PORT"))
	config.PostgresUser = cast.ToString(os.Getenv("POSTGRES_USER"))
	config.PostgresPassword = cast.ToString(os.Getenv("POSTGRES_PASSWORD"))
	config.PostgresDatabase = cast.ToString(os.Getenv("POSTGRES_DATABASE"))
	config.PostgresMaxConnections = cast.ToInt32(getOrReturnDefaultValue("POSTGRES_MAX_CONNECTIONS", 30))

	config.AuthServiceHost = cast.ToString(getOrReturnDefaultValue("AUTH_SERVICE_HOST", "localhost"))
	config.AuthGRPCPort = cast.ToString(getOrReturnDefaultValue("AUTH_GRPC_PORT", ":9105"))

	config.RedisURL = cast.ToString(os.Getenv("REDIS_URL"))
	config.RedisHost = cast.ToString(os.Getenv("REDIS_HOST"))
	config.RedisPort = cast.ToString(os.Getenv("REDIS_PORT"))
	config.RedisPassword = cast.ToString(os.Getenv("REDIS_PASSWORD"))

	config.SecretKey = cast.ToString(getOrReturnDefaultValue("SECRET_KEY", "NVWmbbPGxh7gy1igr4irX3qaAYun9nxi"))

	return config
}

// func getValue(key string) interface{} {
// 	val, exists := os.LookupEnv(key)
// 	if exists {
// 		return val
// 	}
// 	return nil
// }

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
