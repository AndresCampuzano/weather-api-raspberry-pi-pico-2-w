package main

import (
	"database/sql"
	"fmt"
	"log"
)

func (s *PostgresStore) CreateCityTable() error {
	// Create the table if it doesn't exist
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS cities (
            id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
            name TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP NULL
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
         	WHERE tgname = 'cities_updated_at_trigger' 
           	AND tgrelid = 'cities'::regclass)
           `).Scan(&triggerExists)
	if err != nil {
		return err
	}

	// Create the trigger only if it doesn't exist
	if !triggerExists {
		// Create the trigger
		_, err = s.db.Exec(`
            CREATE OR REPLACE FUNCTION update_city_timestamp()
            RETURNS TRIGGER AS $$
            BEGIN
                -- Only set updated_at when the record is actually modified
                -- (and not during the initial insert)
                IF OLD.name <> NEW.name THEN
                    NEW.updated_at = NOW();
                END IF;
                RETURN NEW;
            END;
            $$ LANGUAGE plpgsql;
            
            CREATE TRIGGER cities_updated_at_trigger
            BEFORE UPDATE ON cities
            FOR EACH ROW
            EXECUTE FUNCTION update_city_timestamp();
                `)
		if err != nil {
			return err
		}
	}

	// Check if the createdAt trigger already exists
	err = s.db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM pg_trigger 
            WHERE tgname = 'cities_created_at_trigger' 
            AND tgrelid = 'cities'::regclass)
    `).Scan(&triggerExists)
	if err != nil {
		return err
	}

	// Create the trigger only if it doesn't exist
	if !triggerExists {
		// Create the trigger for created_at
		_, err = s.db.Exec(`
            CREATE OR REPLACE FUNCTION set_city_created_at()
            RETURNS TRIGGER AS $$
            BEGIN
                NEW.created_at = NOW();
                RETURN NEW;
            END;
            $$ LANGUAGE plpgsql;
            
            CREATE TRIGGER cities_created_at_trigger
            BEFORE INSERT ON cities
            FOR EACH ROW
            EXECUTE FUNCTION set_city_created_at();
        `)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgresStore) CreateCity(city *City) error {
	query := `
		INSERT INTO cities (name, updated_at) 
		VALUES ($1, NULL)
		RETURNING id
	`

	var id string
	err := s.db.QueryRow(
		query,
		city.Name,
	).Scan(&id)
	if err != nil {
		return err
	}

	// Set the ID of the inserted city
	city.ID = id

	return nil
}

func (s *PostgresStore) GetCityByID(id string) (*City, error) {
	rows, err := s.db.Query("SELECT * FROM cities WHERE id = $1", id)
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
		return scanIntoCity(rows)
	}

	return nil, fmt.Errorf("city [%s] not found", id)
}

func scanIntoCity(rows *sql.Rows) (*City, error) {
	city := new(City)
	err := rows.Scan(
		&city.ID,
		&city.Name,
		&city.CreatedAt,
		&city.UpdatedAt,
	)

	return city, err
}

func (s *PostgresStore) GetCities() ([]*City, error) {
	rows, err := s.db.Query("SELECT * FROM cities")
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var cities []*City
	for rows.Next() {
		city, err := scanIntoCity(rows)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}

	return cities, nil
}

func (s *PostgresStore) UpdateCity(city *City) error {
	query := `
		UPDATE cities 
		SET name = $1, updated_at = NOW() 
		WHERE id = $2
	`

	_, err := s.db.Exec(
		query,
		city.Name,
		city.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteCity(id string) error {
	query := `
		DELETE FROM cities 
		WHERE id = $1
	`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
