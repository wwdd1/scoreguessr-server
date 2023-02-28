package tools

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func GetPostgresSchema() string {
	return os.Getenv("POSTGRES_SCHEMA")
}

func CreatePostgresConnectionPool() *sql.DB {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASS")
	database := os.Getenv("POSTGRES_DB")
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, database,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	return db
}
