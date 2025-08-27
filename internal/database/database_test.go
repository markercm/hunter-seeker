package database

import (
	"os"
	"testing"
	"time"

	"hunter-seeker/internal/models"
)

func TestDatabase(t *testing.T) {
	// Create a temporary database for testing
	dbPath := "./test_jobs.db"
	defer os.Remove(dbPath)

	// Initialize database
	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test creating a job application
	job := &models.JobApplication{
		DateApplied: time.Now(),
		JobTitle:    "Test Engineer",
		Company:     "Test Corp",
		Status:      models.StatusApplied,
		JobURL:      "https://test.com/job/123",
		Notes:       "This is a test job application",
	}

	err = db.CreateJobApplication(job)
	if err != nil {
		t.Fatalf("Failed to create job application: %v", err)
	}

	if job.ID == 0 {
		t.Error("Expected job ID to be set after creation")
	}

	// Test getting the job application
	retrievedJob, err := db.GetJobApplication(job.ID)
	if err != nil {
		t.Fatalf("Failed to get job application: %v", err)
	}

	if retrievedJob.JobTitle != job.JobTitle {
		t.Errorf("Expected job title %s, got %s", job.JobTitle, retrievedJob.JobTitle)
	}

	if retrievedJob.Company != job.Company {
		t.Errorf("Expected company %s, got %s", job.Company, retrievedJob.Company)
	}

	// Test updating the job application
	retrievedJob.Status = models.StatusInterview
	retrievedJob.Notes = "Updated notes"

	err = db.UpdateJobApplication(retrievedJob)
	if err != nil {
		t.Fatalf("Failed to update job application: %v", err)
	}

	// Verify the update
	updatedJob, err := db.GetJobApplication(job.ID)
	if err != nil {
		t.Fatalf("Failed to get updated job application: %v", err)
	}

	if updatedJob.Status != models.StatusInterview {
		t.Errorf("Expected status %s, got %s", models.StatusInterview, updatedJob.Status)
	}

	if updatedJob.Notes != "Updated notes" {
		t.Errorf("Expected notes to be updated")
	}

	// Test getting all job applications
	allJobs, err := db.GetAllJobApplications()
	if err != nil {
		t.Fatalf("Failed to get all job applications: %v", err)
	}

	if len(allJobs) != 1 {
		t.Errorf("Expected 1 job application, got %d", len(allJobs))
	}

	// Test getting job applications by status
	jobsByStatus, err := db.GetJobApplicationsByStatus(models.StatusInterview)
	if err != nil {
		t.Fatalf("Failed to get job applications by status: %v", err)
	}

	if len(jobsByStatus) != 1 {
		t.Errorf("Expected 1 job application with status %s, got %d", models.StatusInterview, len(jobsByStatus))
	}

	// Test status counts
	statusCounts, err := db.GetStatusCounts()
	if err != nil {
		t.Fatalf("Failed to get status counts: %v", err)
	}

	if count, exists := statusCounts[models.StatusInterview]; !exists || count != 1 {
		t.Errorf("Expected 1 job application with status %s", models.StatusInterview)
	}

	// Test deleting the job application
	err = db.DeleteJobApplication(job.ID)
	if err != nil {
		t.Fatalf("Failed to delete job application: %v", err)
	}

	// Verify deletion
	_, err = db.GetJobApplication(job.ID)
	if err == nil {
		t.Error("Expected error when getting deleted job application")
	}

	// Test getting all job applications after deletion
	allJobsAfterDeletion, err := db.GetAllJobApplications()
	if err != nil {
		t.Fatalf("Failed to get all job applications after deletion: %v", err)
	}

	if len(allJobsAfterDeletion) != 0 {
		t.Errorf("Expected 0 job applications after deletion, got %d", len(allJobsAfterDeletion))
	}
}

func TestDatabaseEdgeCases(t *testing.T) {
	// Create a temporary database for testing
	dbPath := "./test_edge_cases.db"
	defer os.Remove(dbPath)

	// Initialize database
	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test getting non-existent job application
	_, err = db.GetJobApplication(999)
	if err == nil {
		t.Error("Expected error when getting non-existent job application")
	}

	// Test updating non-existent job application
	nonExistentJob := &models.JobApplication{
		ID:          999,
		DateApplied: time.Now(),
		JobTitle:    "Non-existent",
		Company:     "Non-existent Corp",
		Status:      models.StatusApplied,
	}

	err = db.UpdateJobApplication(nonExistentJob)
	if err == nil {
		t.Error("Expected error when updating non-existent job application")
	}

	// Test deleting non-existent job application
	err = db.DeleteJobApplication(999)
	if err == nil {
		t.Error("Expected error when deleting non-existent job application")
	}

	// Test with empty status filter
	jobsByEmptyStatus, err := db.GetJobApplicationsByStatus("")
	if err != nil {
		t.Fatalf("Failed to get job applications by empty status: %v", err)
	}

	if len(jobsByEmptyStatus) != 0 {
		t.Errorf("Expected 0 job applications with empty status, got %d", len(jobsByEmptyStatus))
	}
}

func TestMultipleJobApplications(t *testing.T) {
	// Create a temporary database for testing
	dbPath := "./test_multiple.db"
	defer os.Remove(dbPath)

	// Initialize database
	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create multiple job applications
	jobs := []*models.JobApplication{
		{
			DateApplied: time.Now().AddDate(0, 0, -3),
			JobTitle:    "Software Engineer",
			Company:     "Company A",
			Status:      models.StatusApplied,
			JobURL:      "https://companya.com/job/1",
			Notes:       "Applied via website",
		},
		{
			DateApplied: time.Now().AddDate(0, 0, -2),
			JobTitle:    "Backend Developer",
			Company:     "Company B",
			Status:      models.StatusInReview,
			JobURL:      "https://companyb.com/job/2",
			Notes:       "Recruiter reached out",
		},
		{
			DateApplied: time.Now().AddDate(0, 0, -1),
			JobTitle:    "Full Stack Developer",
			Company:     "Company C",
			Status:      models.StatusInterview,
			JobURL:      "https://companyc.com/job/3",
			Notes:       "Interview scheduled",
		},
	}

	// Insert all job applications
	for _, job := range jobs {
		err = db.CreateJobApplication(job)
		if err != nil {
			t.Fatalf("Failed to create job application: %v", err)
		}
	}

	// Test getting all job applications (should be ordered by date applied desc)
	allJobs, err := db.GetAllJobApplications()
	if err != nil {
		t.Fatalf("Failed to get all job applications: %v", err)
	}

	if len(allJobs) != 3 {
		t.Errorf("Expected 3 job applications, got %d", len(allJobs))
	}

	// Verify ordering (newest first)
	if allJobs[0].Company != "Company C" {
		t.Errorf("Expected newest job application first, got %s", allJobs[0].Company)
	}

	// Test status counts
	statusCounts, err := db.GetStatusCounts()
	if err != nil {
		t.Fatalf("Failed to get status counts: %v", err)
	}

	expectedCounts := map[string]int{
		models.StatusApplied:   1,
		models.StatusInReview:  1,
		models.StatusInterview: 1,
	}

	for status, expectedCount := range expectedCounts {
		if count, exists := statusCounts[status]; !exists || count != expectedCount {
			t.Errorf("Expected %d applications with status %s, got %d", expectedCount, status, count)
		}
	}

	// Test filtering by status
	interviewJobs, err := db.GetJobApplicationsByStatus(models.StatusInterview)
	if err != nil {
		t.Fatalf("Failed to get job applications by status: %v", err)
	}

	if len(interviewJobs) != 1 {
		t.Errorf("Expected 1 job application with interview status, got %d", len(interviewJobs))
	}

	if interviewJobs[0].Company != "Company C" {
		t.Errorf("Expected Company C for interview status, got %s", interviewJobs[0].Company)
	}
}
