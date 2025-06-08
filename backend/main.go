package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"posteaze-backend/models"
	"posteaze-backend/pkg/config"
	"posteaze-backend/routes"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	port := os.Getenv("PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}

	if port == "" {
		log.Println("PORT is not found in the environment variable")
		port = "8080"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Kolkata", host, dbUser, dbPass, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Run extension for UUID generation
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`).Error; err != nil {
		log.Fatalf("Failed to enable pgcrypto extension: %v", err)
	}

	// Run DB migrations
	err = db.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.TeamMember{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	config.InitAppContext(db)
	router := gin.Default()

	// Register routes
	routes.RegisterRoutes(router)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "PostEaze backend running with Gin!",
		})
	})

	router.Run(":" + port)
}
