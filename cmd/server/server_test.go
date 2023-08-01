package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
)

const (
	DINDHOST     = "localhost"
	DINDPORT     = "5432"
	DINDDBNAME   = "postgres"
	DINDUSER     = "postgres"
	DINDPASSWORD = "postgres"
)

func TestMain(m *testing.M) {
	log.Println("before")
	pool, err := dockertest.NewPool("")
	pool.MaxWait = time.Second * 10
	if err != nil {
		log.Fatalf("Unable to connect to docker: %s\n", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("docker does not accept communication: %s\n", err)
	}
	if err != nil {
		log.Fatalf("Failed to start docker: %s\n", err)
	}
	databaseUrl := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", DINDHOST, DINDPORT, DINDDBNAME, DINDUSER, DINDPASSWORD)
	fmt.Printf("databaseUrl: %s\n", databaseUrl)
	if err := pool.Retry(func() error {
		time.Sleep(time.Second * 10)
		db, err := sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Unable to connect to database: %s\n", err)
	}
	code := m.Run()
	os.Exit(code)
}
