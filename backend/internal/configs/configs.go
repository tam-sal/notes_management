package configs

import (
	"os"
	"strconv"
)

const (
	env            = "local"
	dbPORT         = 5400
	dbName         = "note_rev"
	httpPort       = 8025
	apiUrl         = "http://localhost:8025"
	allowedOrigins = ""
)

func New() *Config {
	dbUser := GetString("DB_USER", "")
	dbPassword := GetString("DB_PASSWORD", "")
	dbURI := "postgres://" + dbUser + ":" + dbPassword + "@localhost:" + strconv.Itoa(dbPORT) + "/" + dbName + "?sslmode=disable"
	return &Config{
		DB_URI:          GetString("DB_URI", dbURI),
		DB_NAME:         GetString("DB_NAME", dbName),
		DB_PORT:         GetInt("DB_PORT", dbPORT),
		ENV:             GetString("ENV", env),
		API_URL:         GetString("API_URL", apiUrl),
		HTTP_PORT:       GetInt("HTTP_PORT", httpPort),
		ALLOWED_ORIGINS: GetString("ALLOWED_ORIGINS", allowedOrigins),
	}
}

type Config struct {
	//DATABASE
	DB_URI  string
	DB_NAME string
	DB_PORT int

	//SERVER
	ENV             string
	API_URL         string
	HTTP_PORT       int
	ALLOWED_ORIGINS string
}

func GetString(key, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if !exists || val == "" {
		return defaultValue
	}
	return val
}

func GetInt(key string, defaultValue int) int {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return intVal
}

func GetBool(key string, defaultValue bool) bool {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		panic(err)
	}
	return boolVal
}
