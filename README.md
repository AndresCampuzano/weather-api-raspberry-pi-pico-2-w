# Weather API for Raspberry Pi Pico W

This project is a RESTful API designed to manage weather data and predictions for cities. It is built using Go and is intended to run on a Raspberry Pi Pico W. The API supports CRUD operations for weather data, cities, and predictions, and includes features like hourly weather averages and filtering by time.

## Features

- **Weather Management**: Create, retrieve, update, and delete weather data.
- **City Management**: Manage city information.
- **Predictions**: Add and retrieve weather predictions.
- **Hourly Averages**: Calculate hourly averages for weather data.
- **CORS Support**: Configurable allowed origins for cross-origin requests.

## Endpoints

- `/api/healthcheck`: Check API health.
- `/api/weather`: Manage weather data.
- `/api/weather/{id}`: Manage weather data by ID.
- `/api/cities`: Manage cities.
- `/api/cities/{id}`: Manage cities by ID.
- `/api/predictions`: Manage weather predictions.

## Setup

1. Clone the repository.
2. Install dependencies: `go mod tidy`.
3. Set environment variables:
   - `ALLOWED_ORIGINS`: Comma-separated list of allowed origins for CORS.
4. Build and run:
   ```bash
   make run
   ```

## Database

The project uses PostgreSQL for data storage. Ensure the database is set up with the required tables and triggers by calling the `Init` method in the `PostgresStore`.

## Testing

Run tests with:
```bash
make test
```

## License

This project is licensed under the MIT License.
