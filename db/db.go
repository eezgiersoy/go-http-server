package db

import (
	"awesomeProject/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

type Instance struct {
	Db *gorm.DB
}

var Database Instance

func ConnectDb() {
	db, err := gorm.Open(sqlite.Open("api.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database. \n", err.Error())
		os.Exit(2)
	}

	log.Println("Connected to the database successfully")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")

	//TODO: Add migrations
	err = db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})
	if err != nil {
		return
	}

	Database = Instance{
		Db: db,
	}
}
