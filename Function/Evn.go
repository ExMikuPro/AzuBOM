package Function

import (
	"github.com/joho/godotenv"
	"os"
)

func BeginEvn() error {
	err := godotenv.Load("Service.env")
	return err
}

func GetEvn(key string) string {
	return os.Getenv(key)
}
