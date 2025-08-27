package main

import (
	"log"
	"net/http"
	"os"

	"hunter-seeker/internal/database"
	"hunter-seeker/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Get environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/jobs.db"
	}

	// Initialize database
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize handlers
	h, err := handlers.New(db, "./web/templates")
	if err != nil {
		log.Fatalf("Failed to initialize handlers: %v", err)
	}

	// Setup router
	r := mux.NewRouter()

	// Web routes
	r.HandleFunc("/", h.HomeHandler).Methods("GET")
	r.HandleFunc("/add", h.AddJobHandler).Methods("GET")
	r.HandleFunc("/create", h.CreateJobHandler).Methods("POST")
	r.HandleFunc("/edit/{id}", h.EditJobHandler).Methods("GET")
	r.HandleFunc("/update/{id}", h.UpdateJobHandler).Methods("POST")
	r.HandleFunc("/delete/{id}", h.DeleteJobHandler).Methods("POST")
	r.HandleFunc("/filter", h.FilterHandler).Methods("GET")
	r.HandleFunc("/import-csv", h.ImportCSVHandler).Methods("GET")
	r.HandleFunc("/process-csv", h.ProcessCSVHandler).Methods("POST")

	// API routes
	r.HandleFunc("/api/stats", h.StatsHandler).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"hunter-seeker"}`))
	}).Methods("GET")

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	log.Printf("Server starting on port %s", port)
	log.Printf("Database: %s", dbPath)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
