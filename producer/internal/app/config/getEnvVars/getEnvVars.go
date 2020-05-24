package getEnvVars

import (
	"github.com/joho/godotenv"
	"fmt"
)

func GetEnvVars() {
	err := godotenv.Load("./src/PubSub/credentials.env")
	if err != nil {
		// logger.SugarLogger.Error("Error laoding .env file")
		fmt.Println("Error laoding .env file")
	}
}