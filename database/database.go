package database

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres" 

	"gorm.io/gorm"
)


var GlobalDB *gorm.DB


func InitDatabase() (err error) {

	config, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error reading .env file")
	}


	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config["DATABASE_HOST"],
		config["DB_USERNAME"],
		config["DB_PASSWORD"],
		config["DB_DATABASE"],
		config["DB_PORT"],
	)


	GlobalDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
