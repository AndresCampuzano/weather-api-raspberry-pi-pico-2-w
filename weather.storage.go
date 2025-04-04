package main

func (s *PostgresStore) CreateWeatherTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS weather (
            id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
            temperature FLOAT NOT NULL,
            humidity FLOAT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
