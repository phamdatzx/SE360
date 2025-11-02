package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Kiểm tra nếu thiếu thông tin nào đó
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		fmt.Println("Missing environment database information (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)")
		return
	}

	// DSN cho PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
	})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	DB = database
	fmt.Println("✅ Database connected")
}
