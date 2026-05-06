package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// UserData model is now defined in models.go

func ConnectDatabase() {
	user := os.Getenv("GO_DB_USER")
	if user == "" {
		user = "root"
	}
	pass := os.Getenv("GO_DB_PASSWORD")
	if pass == "" {
		pass = "password"
	}
	host := os.Getenv("GO_DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("GO_DB_PORT")
	if port == "" {
		port = "3306"
	}
	name := os.Getenv("GO_DB_NAME")
	if name == "" {
		name = "threat_visualizer"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")

	if err := DB.AutoMigrate(&UserData{}); err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
}
