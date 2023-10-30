package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Book struct {
	Id        uint64   `json:"id" gorm:"primaryKey"`
	Title     string   `json:"title"`
	Isbn      string   `json:"isbn"`
	Language  string   `json:"language"`
	Publisher string   `json:"publisher"`
	NumPages  int64    `json:"numPages"`
	Reviews   []Review `json:"omitempty"`
}

type Review struct {
	Id      uint64 `json:"id" gorm:"primaryKey"`
	BookId  uint64 `json:"bookId"`
	Rating  int64  `json:"rating"`
	Comment string `json:"comment"`
	// Book    Book   `json:"-" gorm:"foreignKey:BookId"` // Foreign key reference
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

	router.GET("/api/books", getBooks)

	router.GET("/api/books/:id", getBookByID)

	router.GET("/api/books/:id/reviews", getBookReviews)

	router.POST("/api/books/:id/reviews", createReview)

	router.Run() // listen and serve on localhost:8080

}

func getBooks(c *gin.Context) {
	var books []Book
	result := db.Find(&books)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
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

func getBookReviews(c *gin.Context) {
	bookId := c.Param("id") // Get the user ID from the URL parameter
	var reviews []Review
	result := db.Where("book_id = ?", bookId).Find(&reviews)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, reviews)
}

func createReview(c *gin.Context) {

	id := c.Param("id") // Get the user ID from the URL parameter
	var review Review
	// Bind JSON request to the input struct
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check if the specified book exists
	var book Book
	result := db.First(&book, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	// Convert the string to a uint64
	bookId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Println("Conversion error:", err)
		return
	}
	review.BookId = bookId
	// Create a new review
	result = db.Create(&review)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, review)
}
