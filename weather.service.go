package main

import (
	"encoding/json"
	"net/http"
)

func (server *APIServer) handleCreateWeather(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateWeatherRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	weather, err := NewWeather(
		req.Temperature,
		req.Humidity,
	)
	if err != nil {
		return err
	}

	err = server.store.CreateWeather(weather)
	if err != nil {
		return err
	}

	// Recovering weather from DB
	createdWeather, err := server.store.GetWeatherByID(weather.ID)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, createdWeather)
	//return WriteJSON(w, http.StatusOK, weather)
}
