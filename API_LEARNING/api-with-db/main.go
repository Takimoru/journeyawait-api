package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/getproducts",
		reutrnAllProducts).Methods("GET")
	http.Handle("/", router)
	fmt.Println("Connected to port 8081")
	log.fatal(http.ListenAndServe(":8081", router))

	router := gin.Default()

}
