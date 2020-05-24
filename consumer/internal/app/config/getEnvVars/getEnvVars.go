package getEnvVars

import (
	"fmt"

	"github.com/joho/godotenv"
)

func GetEnvVars() {
	err := godotenv.Load("./credentials.env")
	if err != nil {
		// logger.SugarLogger.Error("Error laoding .env file")
		fmt.Println("Error laoding .env file")
	}
}
