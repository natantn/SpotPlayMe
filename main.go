package main

import (
	"github.com/joho/godotenv"
	"github.com/natantn/SpotPlayMe/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	
	routes.HandleRequests()
}
