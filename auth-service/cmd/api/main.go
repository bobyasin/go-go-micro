package main

import (
	"auth/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	webPort                   = "80"
	maxDbConnectionRetryCount = 10
)

var dbConnectionRetryCount int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Printf("Auth service started on port %s\n", webPort)

	conn := connectToDB()

	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	//TODO set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {

	dsn := os.Getenv("DSN")
	//dsn := "host=localhost port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5"

	for {
		connection, err := openDB(dsn)

		if err != nil {
			log.Println("Postgres not ready yet, retrying...")
			log.Println(err)
			dbConnectionRetryCount++
		} else {
			log.Println("Connected to Postgres!")

			if dbConnectionRetryCount > 0 {
				dbConnectionRetryCount = 0
			}

			return connection
		}

		if dbConnectionRetryCount > maxDbConnectionRetryCount {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
	}
}
