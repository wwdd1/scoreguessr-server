package tools

import (
	"database/sql"
)

var db *sql.DB

func Init() {
	db = CreatePostgresConnectionPool()
	if db == nil {
		panic("SQL db pool is not initiated.")
	}
	// set default schema for postgres
	//schema := os.Getenv("POSTGRES_SCHEMA")
	//db.Exec("set search_path=\"" + schema + "\"")
}

func GetDb() *sql.DB {
	// TODO
	err := db.Ping()
	if err != nil {
		Init()
	}
	return db
}
