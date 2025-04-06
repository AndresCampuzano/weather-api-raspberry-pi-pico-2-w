package main

import "time"

type Weather struct {
	ID          string     `json:"id"`
	Temperature float64    `json:"temperature"`
	Humidity    float64    `json:"humidity"`
	CityID      string     `json:"city_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type CreateWeatherRequest struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	CityID      string  `json:"city_id"`
}

func NewWeather(
	temperature float64,
	humidity float64,
	cityID string,
) (*Weather, error) {
	return &Weather{
		Temperature: temperature,
		Humidity:    humidity,
		CityID:      cityID,
	}, nil
}
