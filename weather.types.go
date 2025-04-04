package main

import "time"

type Weather struct {
	ID          string     `json:"id"`
	Temperature float64    `json:"temperature"`
	Humidity    float64    `json:"humidity"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type CreateWeatherRequest struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

func NewWeather(
	temperature float64,
	humidity float64,
) (*Weather, error) {
	return &Weather{
		Temperature: temperature,
		Humidity:    humidity,
	}, nil
}
