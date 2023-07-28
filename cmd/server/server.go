package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	HOST     = "localhost"
	PORT     = "5432"
	DBNAME   = "postgres"
	USER     = "postgres"
	PASSWORD = "postgres"
)

func main() {
	fmt.Println("hello")
	//ctx := context.Background()
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", HOST, PORT, DBNAME, USER, PASSWORD)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
