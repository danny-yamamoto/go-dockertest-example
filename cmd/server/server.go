package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"

	"github.com/danny-yamamoto/go-dockertest-example/tutorial"
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
	ctx := context.Background()
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", HOST, PORT, DBNAME, USER, PASSWORD)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	queries := tutorial.New(db)
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(authors)
	insert, err := queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{Name: "danny", Bio: sql.NullString{String: "SRE", Valid: true}})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(insert)
	fetched, err := queries.GetAuthor(ctx, insert.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reflect.DeepEqual(insert, fetched))
}
