package main

import (
	"database/sql"
	"log"
)

const port = 80

type Config struct {
	DB *sql.DB
	// Models data.Models
}

func main() {
	log.Println("Starting authentication service")
}
