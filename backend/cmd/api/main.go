package main

import (
	"backend/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB *sql.DB
}

func main() {
	log.Println("Starting backend service")

	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to postgres")
	}

	data.MigrateUp()

	// set up config
	app := Config{
		DB: conn,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.getRoutes(),
	}

	println("Starting web server on port", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
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

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}
		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing of for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
