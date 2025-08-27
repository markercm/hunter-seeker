package handlers

import (
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
		Statuses      []string
		CurrentFilter string
		StatusMessage string
		StatusType    string
	}{
		Jobs:          jobs,
		StatusCounts:  statusCounts,
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

	data := struct {
		Jobs          []*models.JobApplication
		StatusCounts  map[string]int
		Statuses      []string
		CurrentFilter string
		StatusMessage string
		StatusType    string
	}{
		Jobs:          jobs,
		StatusCounts:  statusCounts,
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
