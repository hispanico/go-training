package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type Book struct {
	Id        uint64 `json:"id" gorm:"primaryKey"`
	Title     string `json:"title"`
	Isbn      string `json:"isbn"`
	Language  string `json:"language"`
	Publisher string `json:"publisher"`
	NumPages  int64  `json:"numPages"`
}

type Review struct {
	Id      uint64 `json:"id" gorm:"primaryKey"`
	BookId  uint64 `json:"bookId"`
	Rating  int64  `json:"rating"`
	Comment string `json:"comment"`
	Book    Book   `gorm:"foreignKey:BookId"` // Foreign key reference
}

var db *gorm.DB

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
	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" || dbSSLMode == "" || dbTimeZone == "" {
		log.Fatal("Missing one or more required environment variables")
	}

	// Create the database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode, dbTimeZone)
	var err error
	// Connect to the database
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database", err)
	}

	// Migrate the Book and Review model to the database (create the "book" and review tables)
	db.AutoMigrate(&Book{}, &Review{})

	// Create a Gin router
	router := gin.Default()

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	router.GET("/api/books", func(c *gin.Context) {
		var books []Book
		result := db.Find(&books)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, books)
	})

	router.GET("/api/reviews", func(c *gin.Context) {
		var reviews []Review
		result := db.Find(&reviews)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, reviews)
	})

	router.GET("/api/books/:id", getBookByID)

	router.Run() // listen and serve on localhost:8080

}

func getBookByID(c *gin.Context) {
	id := c.Param("id") // Get the user ID from the URL parameter

	var book Book
	result := db.First(&book, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}
