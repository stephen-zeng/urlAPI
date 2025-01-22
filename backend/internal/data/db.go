package data

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sql.DB

func Connect() {
	fmt.Println("Connecting to database...")
	db, err = sql.Open("sqlite3", "backend/assets/database.db")
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error connecting to database")
	}
	//defer db.Close()
	fmt.Println("Connected to database")
}

func Disconnect() {
	fmt.Println("Disconnecting from database...")
	defer db.Close()
	fmt.Println("Disconnected from database")
}
