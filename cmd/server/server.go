package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) authorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	queries := tutorial.New(h.db)
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(authors)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(authors); err != nil {
		http.Error(w, "Failed to encode authors to JSON", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func main() {
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", HOST, PORT, DBNAME, USER, PASSWORD)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	handler := NewHandler(db)
	http.HandleFunc("/", handler.authorHandler)
	http.ListenAndServe(":8080", nil)
}
