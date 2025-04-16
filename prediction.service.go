package main

import (
	"encoding/json"
	"net/http"
)

func (server *APIServer) handleCreatePrediction(w http.ResponseWriter, r *http.Request) error {
	var reqs []CreatePredictionRequest
	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		return err
	}

	var createdPredictions []*Prediction
	for _, req := range reqs {
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

		createdPredictions = append(createdPredictions, createdPrediction)
	}

	return WriteJSON(w, http.StatusOK, createdPredictions)
}
