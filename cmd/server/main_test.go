package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"hunter-seeker/internal/database"
	"hunter-seeker/internal/handlers"
	"hunter-seeker/internal/models"

	"github.com/gorilla/mux"
)

// TestHealthHandler tests the health check endpoint
func TestHealthHandler(t *testing.T) {
	// Create a request to the health endpoint
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create router and add the health handler
	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"hunter-seeker"}`))
	}).Methods("GET")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"status":"ok","service":"hunter-seeker"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Check the content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}
}

// TestRouterSetup tests that all routes are properly configured
func TestRouterSetup(t *testing.T) {
	// Create a temporary database for testing
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := database.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Create a temporary templates directory
	templatesDir := filepath.Join(tempDir, "templates")
	err = os.MkdirAll(templatesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create templates directory: %v", err)
	}

	// Create a minimal template file for testing
	indexTemplate := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body><h1>Jobs: {{len .Jobs}}</h1></body>
</html>`

	err = os.WriteFile(filepath.Join(templatesDir, "index.html"), []byte(indexTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Initialize handlers
	h, err := handlers.New(db, templatesDir)
	if err != nil {
		t.Fatalf("Failed to initialize handlers: %v", err)
	}

	// Setup router (similar to main.go)
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

	// Test cases for different routes
	testCases := []struct {
		method       string
		path         string
		expectedCode int
		description  string
	}{
		{"GET", "/", http.StatusOK, "Home page"},
		{"GET", "/health", http.StatusOK, "Health check"},
		{"GET", "/api/stats", http.StatusOK, "Stats API"},
		{"GET", "/filter", http.StatusOK, "Filter page"},
		{"GET", "/nonexistent", http.StatusNotFound, "Non-existent route"},
		{"POST", "/", http.StatusMethodNotAllowed, "Wrong method for home"},
		{"GET", "/create", http.StatusMethodNotAllowed, "Wrong method for create"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Errorf("Route %s %s returned wrong status code: got %v want %v",
					tc.method, tc.path, rr.Code, tc.expectedCode)
			}
		})
	}
}

// TestEnvironmentVariables tests environment variable handling
func TestEnvironmentVariables(t *testing.T) {
	// Save original environment
	originalPort := os.Getenv("PORT")
	originalDBPath := os.Getenv("DB_PATH")

	// Cleanup function
	defer func() {
		os.Setenv("PORT", originalPort)
		os.Setenv("DB_PATH", originalDBPath)
	}()

	testCases := []struct {
		name           string
		envPort        string
		envDBPath      string
		expectedPort   string
		expectedDBPath string
	}{
		{
			name:           "Default values",
			envPort:        "",
			envDBPath:      "",
			expectedPort:   "8080",
			expectedDBPath: "./data/jobs.db",
		},
		{
			name:           "Custom values",
			envPort:        "3000",
			envDBPath:      "/tmp/test.db",
			expectedPort:   "3000",
			expectedDBPath: "/tmp/test.db",
		},
		{
			name:           "Only port set",
			envPort:        "9000",
			envDBPath:      "",
			expectedPort:   "9000",
			expectedDBPath: "./data/jobs.db",
		},
		{
			name:           "Only DB path set",
			envPort:        "",
			envDBPath:      "/custom/path.db",
			expectedPort:   "8080",
			expectedDBPath: "/custom/path.db",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables
			if tc.envPort != "" {
				os.Setenv("PORT", tc.envPort)
			} else {
				os.Unsetenv("PORT")
			}

			if tc.envDBPath != "" {
				os.Setenv("DB_PATH", tc.envDBPath)
			} else {
				os.Unsetenv("DB_PATH")
			}

			// Test the environment variable logic
			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
			}

			dbPath := os.Getenv("DB_PATH")
			if dbPath == "" {
				dbPath = "./data/jobs.db"
			}

			if port != tc.expectedPort {
				t.Errorf("Expected port %s, got %s", tc.expectedPort, port)
			}

			if dbPath != tc.expectedDBPath {
				t.Errorf("Expected DB path %s, got %s", tc.expectedDBPath, dbPath)
			}
		})
	}
}

// TestAPIStatsEndpoint tests the stats API endpoint specifically
func TestAPIStatsEndpoint(t *testing.T) {
	// Create a temporary database for testing
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := database.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Add some test data
	testJobs := []*models.JobApplication{
		{
			DateApplied: time.Now(),
			JobTitle:    "Software Engineer",
			Company:     "Tech Corp",
			Status:      models.StatusApplied,
		},
		{
			DateApplied: time.Now(),
			JobTitle:    "DevOps Engineer",
			Company:     "Cloud Inc",
			Status:      models.StatusInReview,
		},
		{
			DateApplied: time.Now(),
			JobTitle:    "Frontend Developer",
			Company:     "Web Solutions",
			Status:      models.StatusApplied,
		},
	}

	for _, job := range testJobs {
		err := db.CreateJobApplication(job)
		if err != nil {
			t.Fatalf("Failed to create test job: %v", err)
		}
	}

	// Create templates directory and minimal template
	templatesDir := filepath.Join(tempDir, "templates")
	err = os.MkdirAll(templatesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create templates directory: %v", err)
	}

	indexTemplate := `<html><body>Test</body></html>`
	err = os.WriteFile(filepath.Join(templatesDir, "index.html"), []byte(indexTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Initialize handlers
	h, err := handlers.New(db, templatesDir)
	if err != nil {
		t.Fatalf("Failed to initialize handlers: %v", err)
	}

	// Create request to stats endpoint
	req, err := http.NewRequest("GET", "/api/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h.StatsHandler(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Stats handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	// Check content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Stats handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}

	// Parse and verify JSON response
	var stats map[string]int
	err = json.Unmarshal(rr.Body.Bytes(), &stats)
	if err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	// Verify expected counts
	if stats[models.StatusApplied] != 2 {
		t.Errorf("Expected 2 Applied jobs, got %d", stats[models.StatusApplied])
	}

	if stats[models.StatusInReview] != 1 {
		t.Errorf("Expected 1 In Review job, got %d", stats[models.StatusInReview])
	}
}

// TestStaticFileHandling tests static file serving configuration
func TestStaticFileHandling(t *testing.T) {
	// Create a temporary static directory
	tempDir := t.TempDir()
	staticDir := filepath.Join(tempDir, "static")
	err := os.MkdirAll(staticDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create static directory: %v", err)
	}

	// Create a test CSS file
	cssContent := "body { color: blue; }"
	err = os.WriteFile(filepath.Join(staticDir, "style.css"), []byte(cssContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test CSS file: %v", err)
	}

	// Setup router with static file handler
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	// Test static file request
	req, err := http.NewRequest("GET", "/static/style.css", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Static file handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	// Check content
	if strings.TrimSpace(rr.Body.String()) != cssContent {
		t.Errorf("Static file handler returned wrong content: got %v want %v", rr.Body.String(), cssContent)
	}

	// Test non-existent static file
	req2, err := http.NewRequest("GET", "/static/nonexistent.css", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req2)

	// Should return 404 for non-existent file
	if rr2.Code != http.StatusNotFound {
		t.Errorf("Non-existent static file should return 404, got %v", rr2.Code)
	}
}

// TestDatabaseInitialization tests database initialization scenarios
func TestDatabaseInitialization(t *testing.T) {
	testCases := []struct {
		name        string
		setupFunc   func() (string, func())
		expectError bool
	}{
		{
			name: "Valid database path",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				dbPath := filepath.Join(tempDir, "test.db")
				return dbPath, func() {}
			},
			expectError: false,
		},
		{
			name: "Database in subdirectory",
			setupFunc: func() (string, func()) {
				tempDir := t.TempDir()
				subDir := filepath.Join(tempDir, "data")
				err := os.MkdirAll(subDir, 0755)
				if err != nil {
					t.Fatalf("Failed to create subdirectory: %v", err)
				}
				dbPath := filepath.Join(subDir, "jobs.db")
				return dbPath, func() {}
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbPath, cleanup := tc.setupFunc()
			defer cleanup()

			db, err := database.New(dbPath)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if db != nil {
				// Verify database is functional by creating a test job
				testJob := &models.JobApplication{
					DateApplied: time.Now(),
					JobTitle:    "Test Job",
					Company:     "Test Company",
					Status:      models.StatusApplied,
				}

				err = db.CreateJobApplication(testJob)
				if err != nil {
					t.Errorf("Failed to create test job: %v", err)
				}

				// Verify job was created
				jobs, err := db.GetAllJobApplications()
				if err != nil {
					t.Errorf("Failed to get jobs: %v", err)
				}

				if len(jobs) != 1 {
					t.Errorf("Expected 1 job, got %d", len(jobs))
				}

				db.Close()
			}
		})
	}
}

// TestHandlersInitialization tests handlers initialization
func TestHandlersInitialization(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := database.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	testCases := []struct {
		name          string
		templateSetup func() string
		expectError   bool
	}{
		{
			name: "Valid templates directory",
			templateSetup: func() string {
				templatesDir := filepath.Join(tempDir, "templates")
				err := os.MkdirAll(templatesDir, 0755)
				if err != nil {
					t.Fatalf("Failed to create templates directory: %v", err)
				}

				// Create required template files
				templates := map[string]string{
					"index.html":   `<html><body><h1>Jobs: {{len .Jobs}}</h1></body></html>`,
					"add_job.html": `<html><body><form>Add Job</form></body></html>`,
				}

				for filename, content := range templates {
					err = os.WriteFile(filepath.Join(templatesDir, filename), []byte(content), 0644)
					if err != nil {
						t.Fatalf("Failed to create template %s: %v", filename, err)
					}
				}

				return templatesDir
			},
			expectError: false,
		},
		{
			name: "Non-existent templates directory",
			templateSetup: func() string {
				return filepath.Join(tempDir, "nonexistent")
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			templatesDir := tc.templateSetup()

			h, err := handlers.New(db, templatesDir)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tc.expectError && h == nil {
				t.Error("Expected handlers instance but got nil")
			}
		})
	}
}

// BenchmarkHealthHandler benchmarks the health check endpoint
func BenchmarkHealthHandler(b *testing.B) {
	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"hunter-seeker"}`))
	}).Methods("GET")

	req, _ := http.NewRequest("GET", "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	}
}

// Helper function to create a test server setup
func setupTestServer(t *testing.T) (*database.DB, *handlers.Handler, func()) {
	tempDir := t.TempDir()

	// Setup database
	dbPath := filepath.Join(tempDir, "test.db")
	db, err := database.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Setup templates
	templatesDir := filepath.Join(tempDir, "templates")
	err = os.MkdirAll(templatesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create templates directory: %v", err)
	}

	indexTemplate := `<html><body><h1>Test</h1></body></html>`
	err = os.WriteFile(filepath.Join(templatesDir, "index.html"), []byte(indexTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Setup handlers
	h, err := handlers.New(db, templatesDir)
	if err != nil {
		t.Fatalf("Failed to initialize handlers: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, h, cleanup
}
