package env

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func Get(key string) string {
	return os.Getenv(key)
}

func AsBool(key string) bool {
	return Get(key) == "true"
}
