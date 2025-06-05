package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("error loading .env file")
	}
}
