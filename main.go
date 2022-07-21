package main

import (
	"github.com/joho/godotenv"
	integrations "github.com/natantn/SpotPlayMe/integrations/database"
	"github.com/natantn/SpotPlayMe/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	integrations.ConectDB()
	routes.HandleRequests()
}
