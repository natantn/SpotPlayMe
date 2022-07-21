package integrations

import (
	"fmt"
	"os"

	"github.com/natantn/SpotPlayMe/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func ConectDB() {
	dbConn := make(map[string]string)
	dbConn["name"] = os.Getenv("db_name")
	dbConn["host"] = os.Getenv("db_host")
	dbConn["port"] = os.Getenv("db_port")
	dbConn["user"] = os.Getenv("db_user")
	dbConn["password"] = os.Getenv("db_password")

	stringConn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConn["host"], dbConn["user"], dbConn["password"], dbConn["name"], dbConn["port"])
	DB, err = gorm.Open(postgres.Open(stringConn))
	if err != nil {
		panic(err.Error())
	}
	DB.AutoMigrate(&models.Playlist{})
	DB.AutoMigrate(&models.Music{})
}
