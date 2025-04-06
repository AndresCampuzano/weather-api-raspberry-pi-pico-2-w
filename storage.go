package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateWeather(weather *Weather) error
	GetWeatherByID(id string) (*Weather, error)
	GetWeathers() ([]*Weather, error)
	UpdateWeather(weather *Weather) error
	DeleteWeather(id string) error
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
