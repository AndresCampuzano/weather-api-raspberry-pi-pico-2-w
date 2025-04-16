package main

import "time"

type Prediction struct {
	ID          string     `json:"id"`
	CityID      string     `json:"city_id"`
	Temperature float64    `json:"temperature"`
	Humidity    float64    `json:"humidity"`
	ForecastFor time.Time  `json:"forecast_for"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type CreatePredictionRequest struct {
	CityID      string    `json:"city_id"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	ForecastFor time.Time `json:"forecast_for"`
}

func NewPrediction(cityID string, temperature, humidity float64, forecastFor time.Time) (*Prediction, error) {
	return &Prediction{
		CityID:      cityID,
		Temperature: temperature,
		Humidity:    humidity,
		ForecastFor: forecastFor,
	}, nil
}
