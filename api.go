package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"strings"
)

// APIServer represents an HTTP server for handling API requests.
type APIServer struct {
	listenAddr string
	store      Storage
	Router     *mux.Router
}

// NewAPIServer creates a new instance of APIServer.
func NewAPIServer(listenAddr string, store Storage) *APIServer {
	router := mux.NewRouter()

	server := &APIServer{
		listenAddr: listenAddr,
		store:      store,
		Router:     router,
	}

	router.HandleFunc("/api/healthcheck", makeHTTPHandlerFunc(server.handleHealth))
	router.HandleFunc("/api/weather", makeHTTPHandlerFunc(server.handleWeather))

	return server
}

// Run starts the API server and listens for incoming requests.
func (server *APIServer) Run() {
	log.Println("JSON API server running on port: ", server.listenAddr)

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	origins := strings.Split(allowedOrigins, ",")

	c := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	})

	handler := c.Handler(server.Router)

	// Use the CORS-wrapped handler as your HTTP server's handler
	err := http.ListenAndServe(server.listenAddr, handler)
	if err != nil {
		log.Fatal(err)
	}
}

// handleHealth sends a 200 status code.
func (server *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return server.handleHealthCheck(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleWeather handles weather data retrieval.
func (server *APIServer) handleWeather(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodPost:
		return server.handleCreateWeather(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}
