package main

import (
	"encoding/json"
	"net/http"
)

func (server *APIServer) handleCreateCity(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateCityRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	city, err := NewCity(
		req.Name,
	)
	if err != nil {
		return err
	}

	err = server.store.CreateCity(city)
	if err != nil {
		return err
	}

	// Recovering city from DB
	createdCity, err := server.store.GetCityByID(city.ID)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, createdCity)
}

func (server *APIServer) handleGetCityByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	city, err := server.store.GetCityByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, city)
}

func (server *APIServer) handleGetCities(w http.ResponseWriter, _ *http.Request) error {
	cities, err := server.store.GetCities()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, cities)
}

func (server *APIServer) handleUpdateCity(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	_, err = server.store.GetCityByID(id)
	if err != nil {
		return err
	}

	var city City
	if err := json.NewDecoder(r.Body).Decode(&city); err != nil {
		return err
	}

	city.ID = id

	if err := server.store.UpdateCity(&city); err != nil {
		return err
	}

	// Recovering data from DB to get the most up-to-date data
	updatedCity, err := server.store.GetCityByID(city.ID)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, updatedCity)
}

func (server *APIServer) handleDeleteCity(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	err = server.store.DeleteCity(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"deleted": id})
}
