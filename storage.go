package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateWeather(weather *Weather) error
	GetWeatherByID(id string) (*Weather, error)
}

type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) Init() error {
	err := s.CreateWeatherTable()
	if err != nil {
		return err
	}

	return nil
}
