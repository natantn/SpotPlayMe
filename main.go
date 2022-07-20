package main

import (
	"github.com/joho/godotenv"
	"github.com/natantn/SpotPlayMe/integrations/database"
	"github.com/natantn/SpotPlayMe/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	
	database.ConectDB()
	routes.HandleRequests()
}
