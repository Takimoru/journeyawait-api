package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

type user struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/traveldb")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	router.GET("/user/:id", getUser)
	router.POST("/register", postUser)
	router.PUT("/user/:id", updateUser)
	router.DELETE("/user/:id", deleteUser)
	router.POST("/login", loginUser)
	router.Run(":8080")
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	db := dbConn()
	defer db.Close()

	var u user
	err := db.QueryRow("SELECT id, name, password FROM users WHERE id=?", id).Scan(&u.ID, &u.Name, &u.Password)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

func postUser(c *gin.Context) {
	var newUser user
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := dbConn()
	defer db.Close()

	_, err := db.Exec("INSERT INTO users (id, name, password) VALUES (?, ?, ?)", newUser.ID, newUser.Name, newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newUser)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedUser user
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := dbConn()
	defer db.Close()

	_, err := db.Exec("UPDATE users SET name=?, password=? WHERE id=?", updatedUser.Name, updatedUser.Password, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")

	db := dbConn()
	defer db.Close()

	_, err := db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func loginUser(c *gin.Context) {
	var loginUser user
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := dbConn()
	defer db.Close()

	// Log the login attempt
	fmt.Printf("Login attempt: ID=%d, Password=%s\n", loginUser.ID, loginUser.Password)

	var u user
	err := db.QueryRow("SELECT id, name, password FROM users WHERE id=? AND password=?", loginUser.ID, loginUser.Password).Scan(&u.ID, &u.Name, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User not found or password incorrect")
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication failed: user not found or password incorrect"})
		} else {
			fmt.Printf("Database query error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database query error", "error": err.Error()})
		}
		return
	}

	fmt.Printf("User logged in successfully: ID=%d, Name=%s\n", u.ID, u.Name)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": u})
}
