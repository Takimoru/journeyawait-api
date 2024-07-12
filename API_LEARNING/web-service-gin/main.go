package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func main() {
	var err error
	// Database connection
	db, err = sql.Open("mysql", "root:1234@tcp(localhost:3306)/golangdb")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// Verify connection to the database
	if err := db.Ping(); err != nil {
		log.Fatal("Database connection failed:", err)
	} else {
		fmt.Println("Database connection successful.")
	}

	// Set database connection settings
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Initialize Gin router
	router := gin.Default()

	// Set up CORS middleware
	router.Use(cors.Default())

	// Define route handlers
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.GET("/albums/:id", getAlbumById)

	// Start Gin router
	if err := router.Run("localhost:8081"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}

// getAlbums responds with the list of all albums as JSON, queried from the database.
func getAlbums(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, artist, price FROM album")
	if err != nil {
		log.Println("Failed to query albums:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query albums"})
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var albums []album
	for rows.Next() {
		var alb album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			log.Println("Failed to scan album:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan album"})
			return
		}
		albums = append(albums, alb)
	}
	if err = rows.Err(); err != nil {
		log.Println("Rows iteration error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rows iteration error"})
		return
	}

	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds a new album from JSON received in the request body and inserts it into the database.
func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		log.Println("Failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the new album into the database
	result, err := db.Exec("INSERT INTO album (id, title, artist, price) VALUES (?, ?, ?, ?)",
		newAlbum.ID, newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		log.Println("Failed to insert new album:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert new album"})
		return
	}

	// Log the number of rows affected and return the new album
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Failed to retrieve number of rows affected:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert new album"})
		return
	}
	log.Printf("Inserted new album with ID %s; %d rows affected.\n", newAlbum.ID, rowsAffected)

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumById locates the album whose ID value matches the id parameter sent by the client, then returns that album as JSON.
func getAlbumById(c *gin.Context) {
	id := c.Param("id")
	var alb album

	err := db.QueryRow("SELECT id, title, artist, price FROM album WHERE id = ?", id).Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Album with ID %s not found\n", id)
			c.JSON(http.StatusNotFound, gin.H{"message": "Album not found"})
		} else {
			log.Println("Failed to query album by ID:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query album by ID"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, alb)
}
