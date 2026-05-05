package main

import (
	"database/sql"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var SqlDB *sql.DB

func ConnectDatabase() {
	// Fallback to default if no environment variable is provided
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(127.0.0.1:3306)/threat_visualizer?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database using GORM: %v. Please ensure MYSQL_DSN is set correctly.", err)
	} else {
		DB = db
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get raw db: %v", err)
		}
		SqlDB = sqlDB
		log.Println("Successfully connected to the database!")
	}
}
