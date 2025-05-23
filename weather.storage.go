package main

import (
	"database/sql"
	"fmt"
	"log"
)

func (s *PostgresStore) CreateWeatherTable() error {
	// Create the table if it doesn't exist
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS weather (
            id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
            temperature FLOAT NOT NULL,
            humidity FLOAT NOT NULL,
            city_id UUID NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP NULL,
            FOREIGN KEY (city_id) REFERENCES cities(id)
        )
    `)
	if err != nil {
		return err
	}

	// Check if the trigger already exists
	var triggerExists bool
	err = s.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM pg_trigger 
         	WHERE tgname = 'weather_updated_at_trigger' 
           	AND tgrelid = 'weather'::regclass)
           `).Scan(&triggerExists)
	if err != nil {
		return err
	}

	// Create the trigger only if it doesn't exist
	if !triggerExists {
		// Create the trigger
		_, err = s.db.Exec(`
            CREATE OR REPLACE FUNCTION update_timestamp()
            RETURNS TRIGGER AS $$
            BEGIN
                -- Only set updated_at when the record is actually modified
                -- (and not during the initial insert)
                IF OLD.temperature <> NEW.temperature OR OLD.humidity <> NEW.humidity OR OLD.city_id <> NEW.city_id THEN
                    NEW.updated_at = NOW();
                END IF;
                RETURN NEW;
            END;
            $$ LANGUAGE plpgsql;
            
            CREATE TRIGGER weather_updated_at_trigger
            BEFORE UPDATE ON weather
            FOR EACH ROW
            EXECUTE FUNCTION update_timestamp();
                `)
		if err != nil {
			return err
		}
	}

	// Check if the createdAt trigger already exists
	err = s.db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM pg_trigger 
            WHERE tgname = 'weather_created_at_trigger' 
            AND tgrelid = 'weather'::regclass)
    `).Scan(&triggerExists)
	if err != nil {
		return err
	}

	// Create the trigger only if it doesn't exist
	if !triggerExists {
		// Create the trigger for created_at
		_, err = s.db.Exec(`
            CREATE OR REPLACE FUNCTION set_created_at()
            RETURNS TRIGGER AS $$
            BEGIN
                NEW.created_at = NOW();
                RETURN NEW;
            END;
            $$ LANGUAGE plpgsql;
            
            CREATE TRIGGER weather_created_at_trigger
            BEFORE INSERT ON weather
            FOR EACH ROW
            EXECUTE FUNCTION set_created_at();
        `)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgresStore) CreateWeather(weather *Weather) error {
	query := `
		INSERT INTO weather (temperature, humidity, city_id, updated_at) 
		VALUES ($1, $2, $3, NULL)
		RETURNING id
	`

	var id string
	err := s.db.QueryRow(
		query,
		weather.Temperature,
		weather.Humidity,
		weather.CityID,
	).Scan(&id)
	if err != nil {
		return err
	}

	// Set the ID of the inserted weather
	weather.ID = id

	return nil
}

func (s *PostgresStore) GetWeatherByID(id string) (*Weather, error) {
	rows, err := s.db.Query("SELECT * FROM weather WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	for rows.Next() {
		return scanIntoWeather(rows)
	}

	return nil, fmt.Errorf("weather [%s] not found", id)
}

func (s *PostgresStore) GetWeathersByCityID(cityID string) ([]*Weather, error) {
	rows, err := s.db.Query("SELECT * FROM weather WHERE city_id = $1", cityID)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var weathers []*Weather
	for rows.Next() {
		weather, err := scanIntoWeather(rows)
		if err != nil {
			return nil, err
		}
		weathers = append(weathers, weather)
	}

	return weathers, nil
}

func (s *PostgresStore) GetHourlyAveragesByCityID(cityID string) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			date_trunc('hour', created_at) AS hour,
			AVG(temperature) AS avg_temperature,
			AVG(humidity) AS avg_humidity
		FROM weather
		WHERE city_id = $1
		GROUP BY hour
		ORDER BY hour
	`

	rows, err := s.db.Query(query, cityID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var results []map[string]interface{}
	for rows.Next() {
		var hour string
		var avgTemperature, avgHumidity float64

		err := rows.Scan(&hour, &avgTemperature, &avgHumidity)
		if err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"hour":        hour,
			"temperature": avgTemperature,
			"humidity":    avgHumidity,
		})
	}

	return results, nil
}

func scanIntoWeather(rows *sql.Rows) (*Weather, error) {
	weather := new(Weather)
	err := rows.Scan(
		&weather.ID,
		&weather.Temperature,
		&weather.Humidity,
		&weather.CityID,
		&weather.CreatedAt,
		&weather.UpdatedAt,
	)

	return weather, err
}

func (s *PostgresStore) GetWeathers() ([]*Weather, error) {
	rows, err := s.db.Query("SELECT * FROM weather")
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var weathers []*Weather
	for rows.Next() {
		expense, err := scanIntoWeather(rows)
		if err != nil {
			return nil, err
		}
		weathers = append(weathers, expense)
	}

	return weathers, nil
}

func (s *PostgresStore) UpdateWeather(weather *Weather) error {
	query := `
		UPDATE weather 
		SET temperature = $1, humidity = $2, city_id = $3, updated_at = NOW() 
		WHERE id = $4
	`

	_, err := s.db.Exec(
		query,
		weather.Temperature,
		weather.Humidity,
		weather.CityID,
		weather.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteWeather(id string) error {
	query := `
		DELETE FROM weather 
		WHERE id = $1
	`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
