package handlers

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"hunter-seeker/internal/database"
	"hunter-seeker/internal/models"

	"github.com/gorilla/mux"
)

type Handler struct {
	db        *database.DB
	templates *template.Template
}

// New creates a new handler instance
func New(db *database.DB, templateDir string) (*Handler, error) {
	// Create template functions
	funcMap := template.FuncMap{
		"replace": func(s, old, new string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"lower": strings.ToLower,
		"formatDate": func(t time.Time) string {
			return t.Format("Jan 2, 2006")
		},
		"formatDateTime": func(t time.Time) string {
			return t.Format("Jan 2, 2006 at 3:04 PM")
		},
	}

	templates, err := template.New("").Funcs(funcMap).ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &Handler{
		db:        db,
		templates: templates,
	}, nil
}

// HomeHandler renders the main page with all job applications
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.db.GetAllJobApplications()
	if err != nil {
		log.Printf("Error getting job applications: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	statusCounts, err := h.db.GetStatusCounts()
	if err != nil {
		log.Printf("Error getting status counts: %v", err)
		statusCounts = make(map[string]int)
	}

	totalCount, err := h.db.GetTotalJobApplicationCount()
	if err != nil {
		log.Printf("Error getting total count: %v", err)
		totalCount = 0
	}

	// Handle status messages from delete operations
	var statusMessage string
	var statusType string

	if errorType := r.URL.Query().Get("error"); errorType != "" {
		statusType = "error"
		switch errorType {
		case "notfound":
			if id := r.URL.Query().Get("id"); id != "" {
				statusMessage = fmt.Sprintf("Job application with ID %s not found", id)
			} else {
				statusMessage = "Job application not found"
			}
		case "delete":
			statusMessage = "Failed to delete job application"
		}
	} else if success := r.URL.Query().Get("success"); success != "" {
		statusType = "success"
		switch success {
		case "deleted":
			statusMessage = "Job application deleted successfully"
		}
	}

	data := struct {
		Jobs          []*models.JobApplication
		StatusCounts  map[string]int
		TotalCount    int
		Statuses      []string
		CurrentFilter string
		StatusMessage string
		StatusType    string
	}{
		Jobs:          jobs,
		StatusCounts:  statusCounts,
		TotalCount:    totalCount,
		Statuses:      models.GetCommonStatuses(),
		CurrentFilter: "",
		StatusMessage: statusMessage,
		StatusType:    statusType,
	}

	if err := h.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

// AddJobHandler renders the add job form
func (h *Handler) AddJobHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Statuses []string
	}{
		Statuses: models.GetCommonStatuses(),
	}

	if err := h.templates.ExecuteTemplate(w, "add_job.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// CreateJobHandler creates a new job application
func (h *Handler) CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse date
	dateStr := r.FormValue("date_applied")
	dateApplied, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	job := &models.JobApplication{
		DateApplied: dateApplied,
		JobTitle:    r.FormValue("job_title"),
		Company:     r.FormValue("company"),
		Status:      r.FormValue("status"),
		JobURL:      r.FormValue("job_url"),
		Notes:       r.FormValue("notes"),
	}

	if err := h.db.CreateJobApplication(job); err != nil {
		log.Printf("Error creating job application: %v", err)
		http.Error(w, "Failed to create job application", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// EditJobHandler renders the edit job form
func (h *Handler) EditJobHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	job, err := h.db.GetJobApplication(id)
	if err != nil {
		log.Printf("Error getting job application: %v", err)
		http.Error(w, "Job application not found", http.StatusNotFound)
		return
	}

	data := struct {
		Job      *models.JobApplication
		Statuses []string
	}{
		Job:      job,
		Statuses: models.GetCommonStatuses(),
	}

	if err := h.templates.ExecuteTemplate(w, "edit_job.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// UpdateJobHandler updates an existing job application
func (h *Handler) UpdateJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse date
	dateStr := r.FormValue("date_applied")
	dateApplied, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	job := &models.JobApplication{
		ID:          id,
		DateApplied: dateApplied,
		JobTitle:    r.FormValue("job_title"),
		Company:     r.FormValue("company"),
		Status:      r.FormValue("status"),
		JobURL:      r.FormValue("job_url"),
		Notes:       r.FormValue("notes"),
	}

	if err := h.db.UpdateJobApplication(job); err != nil {
		log.Printf("Error updating job application: %v", err)
		http.Error(w, "Failed to update job application", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DeleteJobHandler deletes a job application
func (h *Handler) DeleteJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteJobApplication(id); err != nil {
		if errors.Is(err, database.ErrJobNotFound) {
			log.Printf("Job application not found: ID %d", id)
			http.Redirect(w, r, "/?error=notfound&id="+strconv.Itoa(id), http.StatusSeeOther)
			return
		}
		log.Printf("Error deleting job application: %v", err)
		http.Redirect(w, r, "/?error=delete", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/?success=deleted", http.StatusSeeOther)
}

// FilterHandler handles filtering by status
func (h *Handler) FilterHandler(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	var jobs []*models.JobApplication
	var err error

	if status != "" {
		jobs, err = h.db.GetJobApplicationsByStatus(status)
	} else {
		jobs, err = h.db.GetAllJobApplications()
	}

	if err != nil {
		log.Printf("Error getting job applications: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	statusCounts, err := h.db.GetStatusCounts()
	if err != nil {
		log.Printf("Error getting status counts: %v", err)
		statusCounts = make(map[string]int)
	}

	totalCount, err := h.db.GetTotalJobApplicationCount()
	if err != nil {
		log.Printf("Error getting total count: %v", err)
		totalCount = 0
	}

	data := struct {
		Jobs          []*models.JobApplication
		StatusCounts  map[string]int
		TotalCount    int
		Statuses      []string
		CurrentFilter string
		StatusMessage string
		StatusType    string
	}{
		Jobs:          jobs,
		StatusCounts:  statusCounts,
		TotalCount:    totalCount,
		Statuses:      models.GetCommonStatuses(),
		CurrentFilter: status,
		StatusMessage: "",
		StatusType:    "",
	}

	if err := h.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// DebugFilterHandler is a simple debug endpoint to test filter functionality
func (h *Handler) DebugFilterHandler(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	var jobs []*models.JobApplication
	var err error

	if status != "" {
		jobs, err = h.db.GetJobApplicationsByStatus(status)
	} else {
		jobs, err = h.db.GetAllJobApplications()
	}

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Database error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Filter: %s\n", status)
	fmt.Fprintf(w, "Found %d jobs:\n", len(jobs))
	for _, job := range jobs {
		fmt.Fprintf(w, "- ID: %d, Title: %s, Company: %s, Status: %s\n",
			job.ID, job.JobTitle, job.Company, job.Status)
	}
}

// StatsHandler returns job application statistics as JSON
func (h *Handler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	statusCounts, err := h.db.GetStatusCounts()
	if err != nil {
		log.Printf("Error getting status counts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statusCounts); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// ImportCSVHandler renders the CSV import form
func (h *Handler) ImportCSVHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Statuses []string
	}{
		Statuses: models.GetCommonStatuses(),
	}

	if err := h.templates.ExecuteTemplate(w, "import_csv.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// ProcessCSVHandler processes uploaded CSV file and creates job applications
func (h *Handler) ProcessCSVHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (10MB max)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from form
	file, _, err := r.FormFile("csv_file")
	if err != nil {
		http.Error(w, "Failed to get CSV file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Parse CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, "Failed to parse CSV file", http.StatusBadRequest)
		return
	}

	if len(records) == 0 {
		http.Error(w, "CSV file is empty", http.StatusBadRequest)
		return
	}

	var successCount, errorCount int
	var errors []string

	// Skip header row if present
	startIdx := 0
	if len(records) > 0 && isHeaderRow(records[0]) {
		startIdx = 1
	}

	// Process each row
	for i := startIdx; i < len(records); i++ {
		record := records[i]

		// Skip empty rows
		if len(record) == 0 || (len(record) == 1 && strings.TrimSpace(record[0]) == "") {
			continue
		}

		job, err := parseCSVRecord(record)
		if err != nil {
			errorCount++
			errors = append(errors, fmt.Sprintf("Row %d: %v", i+1, err))
			continue
		}

		if err := h.db.CreateJobApplication(job); err != nil {
			errorCount++
			errors = append(errors, fmt.Sprintf("Row %d: Failed to save %s at %s: %v", i+1, job.JobTitle, job.Company, err))
		} else {
			successCount++
		}
	}

	// Prepare response data
	data := struct {
		SuccessCount int
		ErrorCount   int
		Errors       []string
		TotalRows    int
	}{
		SuccessCount: successCount,
		ErrorCount:   errorCount,
		Errors:       errors,
		TotalRows:    len(records) - startIdx,
	}

	if err := h.templates.ExecuteTemplate(w, "import_result.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// isHeaderRow checks if the first row looks like a header
func isHeaderRow(record []string) bool {
	if len(record) == 0 {
		return false
	}

	// Check for common header keywords
	firstCol := strings.ToLower(strings.TrimSpace(record[0]))
	headerKeywords := []string{"date", "job", "title", "company", "position", "role"}

	for _, keyword := range headerKeywords {
		if strings.Contains(firstCol, keyword) {
			return true
		}
	}

	return false
}

// parseCSVRecord converts a CSV record to a JobApplication
func parseCSVRecord(record []string) (*models.JobApplication, error) {
	if len(record) < 3 {
		return nil, fmt.Errorf("insufficient columns (need at least: date_applied, job_title, company)")
	}

	// Parse date (required)
	dateStr := strings.TrimSpace(record[0])
	dateApplied, err := parseDate(dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format '%s': %v", dateStr, err)
	}

	// Job title (required)
	jobTitle := strings.TrimSpace(record[1])
	if jobTitle == "" {
		return nil, fmt.Errorf("job title is required")
	}

	// Company (required)
	company := strings.TrimSpace(record[2])
	if company == "" {
		return nil, fmt.Errorf("company is required")
	}

	job := &models.JobApplication{
		DateApplied: dateApplied,
		JobTitle:    jobTitle,
		Company:     company,
		Status:      models.StatusApplied, // Default status
	}

	// Status (optional, column 4)
	if len(record) > 3 {
		status := strings.TrimSpace(record[3])
		if status != "" {
			job.Status = status
		}
	}

	// Job URL (optional, column 5)
	if len(record) > 4 {
		job.JobURL = strings.TrimSpace(record[4])
	}

	// Notes (optional, column 6)
	if len(record) > 5 {
		job.Notes = strings.TrimSpace(record[5])
	}

	return job, nil
}

// parseDate attempts to parse various date formats
func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("date is required")
	}

	// Try common date formats
	formats := []string{
		"2006-01-02",      // ISO format (YYYY-MM-DD)
		"01/02/2006",      // US format (MM/DD/YYYY)
		"1/2/2006",        // US format without leading zeros
		"02/01/2006",      // Some other format (DD/MM/YYYY)
		"2/1/2006",        // Without leading zeros
		"2006/01/02",      // YYYY/MM/DD
		"2006/1/2",        // YYYY/M/D
		"Jan 2, 2006",     // Month name format
		"January 2, 2006", // Full month name
		"2 Jan 2006",      // European style
		"2006-1-2",        // ISO without leading zeros
		"Jan 2 2006",      // Month name without comma
		"January 2 2006",  // Full month name without comma
		"Feb 2 2006",      // Short month name without comma
		"Feb 20 2006",     // Short month name without comma (actual example)
		"March 2 2006",    // Month name variations
		"Apr 2 2006",
		"May 2 2006",
		"Jun 2 2006",
		"Jul 2 2006",
		"Aug 2 2006",
		"Sep 2 2006",
		"Oct 2 2006",
		"Nov 2 2006",
		"Dec 2 2006",
	}

	for _, format := range formats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("unrecognized date format")
}
