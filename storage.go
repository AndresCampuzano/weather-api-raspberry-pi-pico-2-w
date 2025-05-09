package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage interface {
	// Weather operations
	CreateWeather(weather *Weather) error
	GetWeatherByID(id string) (*Weather, error)
	GetWeathers() ([]*Weather, error)
	GetWeathersByCityID(cityID string) ([]*Weather, error)
	UpdateWeather(weather *Weather) error
	DeleteWeather(id string) error
	GetHourlyAveragesByCityID(cityID string) ([]map[string]interface{}, error)

	// City operations
	CreateCity(city *City) error
	GetCityByID(id string) (*City, error)
	GetCities() ([]*City, error)
	UpdateCity(city *City) error
	DeleteCity(id string) error

	// Prediction operations
	CreatePrediction(prediction *Prediction) error
	GetPredictionByID(id string) (*Prediction, error)
	GetPredictionsByCityID(cityID string) ([]*Prediction, error)
}

type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) Init() error {
	// First create the cities table since weather table has a foreign key to it
	err := s.CreateCityTable()
	if err != nil {
		return err
	}

	// Then create the weather table which references cities
	err = s.CreateWeatherTable()
	if err != nil {
		return err
	}

	// Finally create the predictions table
	err = s.CreatePredictionTable()
	if err != nil {
		return err
	}

	return nil
}
