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

	// Verify the city exists first
	_, err := server.store.GetCityByID(req.CityID)
	if err != nil {
		return err
	}

	weather, err := NewWeather(
		req.Temperature,
		req.Humidity,
		req.CityID,
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
}

func (server *APIServer) handleGetWeatherByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	weather, err := server.store.GetWeatherByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, weather)
}

func (server *APIServer) handleGetWeathers(w http.ResponseWriter, _ *http.Request) error {
	weathers, err := server.store.GetWeathers()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, weathers)
}

func (server *APIServer) handleUpdateWeather(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	_, err = server.store.GetWeatherByID(id)
	if err != nil {
		return err
	}

	var weather Weather
	if err := json.NewDecoder(r.Body).Decode(&weather); err != nil {
		return err
	}

	weather.ID = id

	// Verify the city exists
	if weather.CityID != "" {
		_, err = server.store.GetCityByID(weather.CityID)
		if err != nil {
			return err
		}
	}

	if err := server.store.UpdateWeather(&weather); err != nil {
		return err
	}

	// Recovering data from DB to get the most up-to-date data
	updatedWeather, err := server.store.GetWeatherByID(weather.ID)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, updatedWeather)
}

func (server *APIServer) handleDeleteWeather(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	err = server.store.DeleteWeather(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"deleted": id})
}
