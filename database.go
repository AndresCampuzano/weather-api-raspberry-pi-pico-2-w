package main

import (
	"database/sql"
	"fmt"
	"os"
)

func NewPostgresStore() (*PostgresStore, error) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB_NAME")
	endpoint := os.Getenv("POSTGRES_DB_HOST")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=require", user, password, dbname, endpoint)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connectivity
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Install UUID on postgres
	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err != nil {
		fmt.Println("Error creating uuid-ossp extension:", err)
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}
