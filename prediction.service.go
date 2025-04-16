package main

import (
	"encoding/json"
	"net/http"
)

func (server *APIServer) handleCreatePrediction(w http.ResponseWriter, r *http.Request) error {
	req := new(CreatePredictionRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	// Verify the city exists
	_, err := server.store.GetCityByID(req.CityID)
	if err != nil {
		return err
	}

	prediction, err := NewPrediction(
		req.CityID,
		req.Temperature,
		req.Humidity,
		req.ForecastFor,
	)
	if err != nil {
		return err
	}

	err = server.store.CreatePrediction(prediction)
	if err != nil {
		return err
	}

	// Recovering prediction from DB
	createdPrediction, err := server.store.GetPredictionByID(prediction.ID)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, createdPrediction)
}
