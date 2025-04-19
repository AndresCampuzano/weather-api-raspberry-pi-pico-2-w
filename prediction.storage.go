package main

import (
	"database/sql"
	"fmt"
	"log"
)

func (s *PostgresStore) CreatePredictionTable() error {
	// Create the table if it doesn't exist
	_, err := s.db.Exec(`	
        CREATE TABLE IF NOT EXISTS predictions (
            id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
            city_id UUID NOT NULL,
            temperature FLOAT,
            humidity FLOAT,
            forecast_for TIMESTAMP NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP NULL,
            FOREIGN KEY (city_id) REFERENCES cities(id)
        )
    `)
	if err != nil {
		return err
	}

	// Check if the updated_at trigger already exists
	var triggerExists bool
	err = s.db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM pg_trigger 
            WHERE tgname = 'predictions_updated_at_trigger' 
            AND tgrelid = 'predictions'::regclass)
    `).Scan(&triggerExists)
	if err != nil {
		return err
	}

	// Create the updated_at trigger only if it doesn't exist
	if !triggerExists {
		_, err = s.db.Exec(`
            CREATE OR REPLACE FUNCTION update_prediction_timestamp()
            RETURNS TRIGGER AS $$
            BEGIN
                IF OLD.temperature <> NEW.temperature OR OLD.humidity <> NEW.humidity OR OLD.forecast_for <> NEW.forecast_for THEN
                    NEW.updated_at = NOW();
                END IF;
                RETURN NEW;
            END;
            $$ LANGUAGE plpgsql;
            
            CREATE TRIGGER predictions_updated_at_trigger
            BEFORE UPDATE ON predictions
            FOR EACH ROW
            EXECUTE FUNCTION update_prediction_timestamp();
        `)
		if err != nil {
			return err
		}
	}

	// Check if the created_at trigger already exists
	err = s.db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM pg_trigger 
            WHERE tgname = 'predictions_created_at_trigger' 
            AND tgrelid = 'predictions'::regclass)
    `).Scan(&triggerExists)
	if err != nil {
		return err
	}

	// Create the created_at trigger only if it doesn't exist
	if !triggerExists {
		_, err = s.db.Exec(`
            CREATE OR REPLACE FUNCTION set_prediction_created_at()
            RETURNS TRIGGER AS $$
            BEGIN
                NEW.created_at = NOW();
                RETURN NEW;
            END;
            $$ LANGUAGE plpgsql;
            
            CREATE TRIGGER predictions_created_at_trigger
            BEFORE INSERT ON predictions
            FOR EACH ROW
            EXECUTE FUNCTION set_prediction_created_at();
        `)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgresStore) CreatePrediction(prediction *Prediction) error {
	query := `
		INSERT INTO predictions (city_id, temperature, humidity, forecast_for) 
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id string
	err := s.db.QueryRow(
		query,
		prediction.CityID,
		prediction.Temperature,
		prediction.Humidity,
		prediction.ForecastFor,
	).Scan(&id)
	if err != nil {
		return err
	}

	prediction.ID = id
	return nil
}

func (s *PostgresStore) GetPredictionByID(id string) (*Prediction, error) {
	rows, err := s.db.Query("SELECT * FROM predictions WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	if rows.Next() {
		return scanIntoPrediction(rows)
	}

	return nil, fmt.Errorf("prediction [%s] not found", id)
}

func (s *PostgresStore) GetPredictionsByCityID(cityID string) ([]*Prediction, error) {
	rows, err := s.db.Query("SELECT * FROM predictions WHERE city_id = $1 ORDER BY forecast_for", cityID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var predictions []*Prediction
	for rows.Next() {
		prediction, err := scanIntoPrediction(rows)
		if err != nil {
			return nil, err
		}
		predictions = append(predictions, prediction)
	}

	return predictions, nil
}

func scanIntoPrediction(rows *sql.Rows) (*Prediction, error) {
	prediction := new(Prediction)
	err := rows.Scan(
		&prediction.ID,
		&prediction.CityID,
		&prediction.Temperature,
		&prediction.Humidity,
		&prediction.ForecastFor,
		&prediction.CreatedAt,
		&prediction.UpdatedAt,
	)

	return prediction, err
}
