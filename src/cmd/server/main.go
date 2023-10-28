package main

import (
	"net/http"
	"os"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Book struct {
	Id        uint64 `json:"id"`
	Title			string `json:"title"`
	Isbn      string `json:"isbn"`
	Language  string `json:"language"`
	Publisher string `json:"publisher"`
	NumPages  int64  `json:"numPages"`
}

type Review struct {
	Id      uint64 `json:"id`
	BookId	uint64 `json:"bookId"`
	Rating  int64	`json:"rating"`
	Comment string `json:"comment"`
}

func main() {
	// Load environment variables
	dbHost := os.Getenv("DB_HOST") //
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSL_MODE")
	dbTimeZone := os.Getenv("DB_TIMEZONE")

	// Check if any of the environment variables are missing
	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" || dbSSLMode == "" {
			log.Fatal("Missing one or more required environment variables")
	}

	// Create the database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s" TimeZone=%s, dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode, dbTimeZone)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to the database", err)
	}

	route := gin.Default()

	route.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	router.GET("/books", func(c *gin.Context) {
    var books []Book
    result := db.Find(&books)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, books)
	})

	r.Run() // listen and serve on localhost:8080

}
