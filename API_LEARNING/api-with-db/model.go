package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Users struct {
	ID        string `form:"id" json:"id" db:"id"`
	FirstName string `form:"first_name" json:"first_name"`
	LastName  string `form:"last_name" json:"last_name"`
}

type Response struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    []Users `json:"data"` // Exported field
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	var users Users
	var arr_users []Users
	var response Response

	db := connect() // Ensure connect function is defined elsewhere
	defer db.Close()

	rows, err := db.Query("SELECT id, first_name, last_name FROM person")
	if err != nil {
		log.Println(err)
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&users.ID, &users.FirstName, &users.LastName); err != nil {
			log.Println(err)
			http.Error(w, "Error scanning rows", http.StatusInternalServerError)
			return
		}
		arr_users = append(arr_users, users)
	}

	response.Status = 1
	response.Message = "success"
	response.Data = arr_users // Use the exported field

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response) // Use the previously declared err variable
	if err != nil {
		log.Println(err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func connect() *sql.DB {
	// Ensure that this function is properly defined to return a database connection
	// This is a placeholder function definition
	// You should replace this with actual implementation
}
