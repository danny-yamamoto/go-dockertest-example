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
	dc "github.com/ory/dockertest/v3/docker"
)

const (
	DINDHOST     = "localhost"
	DINDPORT     = "0"
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

	ro := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.8",
		Env:        []string{"POSTGRES_USER=postgres", "POSTGRES_PASSWORD=postgres", "listen_addresses = '*'"},
	}

	// Delete dind's postgres so that there is no postgres left.
	resource, err := pool.RunWithOptions(ro,
		func(hc *dc.HostConfig) {
			hc.AutoRemove = true
			hc.RestartPolicy = dc.RestartPolicy{Name: "no"}
		})
	if err != nil {
		log.Fatalf("Failed to start docker: %s\n", err)
	}

	// The port to be opened is automatically waved.
	hostAndPort := resource.GetHostPort("5432/tcp")
	log.Printf("hostAndPort: %s\n", hostAndPort)
	databaseUrl := fmt.Sprintf("postgres://postgres:postgres@%s/postgres?sslmode=disable", hostAndPort)
	log.Printf("databaseUrl: %s\n", databaseUrl)
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
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}
