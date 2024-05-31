package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type album struct {
	ID     string `json:"id"`
	Title  string `json:"name"`
	Artist string `json:"artist"`
	Price  int    `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Soire", Artist: "Hoshimachi Suisei", Price: 100},
	{ID: "2", Title: "Fanfare", Artist: "Sakura Miko", Price: 200},
	{ID: "3", Title: "RapGOD", Artist: "James Arthur", Price: 300},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
