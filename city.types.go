package main

import "time"

type City struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type CreateCityRequest struct {
	Name string `json:"name"`
}

func NewCity(
	name string,
) (*City, error) {
	return &City{
		Name: name,
	}, nil
}
