package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDatabase() {
	mongoURI := os.Getenv("MONGO_URI")
	dbname := os.Getenv("DB_NAME")

	// Kiểm tra nếu thiếu thông tin nào đó
	if mongoURI == "" || dbname == "" {
		fmt.Println("Missing environment database information (MONGO_URI, DB_NAME)")
		return
	}

	// Tạo context với timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Kết nối đến MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Ping database để kiểm tra kết nối
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping database: ", err)
	}

	DB = client.Database(dbname)
	fmt.Println("✅ Database connected")
}
