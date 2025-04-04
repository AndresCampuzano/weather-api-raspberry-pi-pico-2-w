build:
	@go build -o bin/weather-api-raspberry-pi-pico-2-w

run: build
	@./bin/weather-api-raspberry-pi-pico-2-w

test:
	@go test -v ./...