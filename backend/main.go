package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	port := os.Getenv("PORT")
	db_user := os.Getenv("POSTGRES_USER")
	db_pass := os.Getenv("POSTGRES_PASSWORD")
    db_name := os.Getenv("POSTGRES_DB")

	if port == "" {
		log.Println("PORT is not found in the environment variable")
		port = "8080"
	}

	dsn := "host=localhost user=" + db_user + " password=" + db_pass + " dbname=" + db_name + " port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	fmt.Println("Connecting to DB...", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
	    log.Fatalf("Failed to connect to DB: %v", err)
	}
	log.Println("Connected to DB successfully!", db)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
	    c.JSON(200, gin.H{
	        "message": "PostEaze backend running with Gin!",
	    })
	})

	router.Run(":" + port)
}
