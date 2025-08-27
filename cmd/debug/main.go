package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"hunter-seeker/internal/database"
	"hunter-seeker/internal/models"
)

func main() {
	dbPath := "./data/jobs.db"
	if len(os.Args) > 1 {
		if os.Args[1] == "add-test-data" {
			addTestData(dbPath)
			return
		}
		dbPath = os.Args[1]
	}

	// Initialize database
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Get all job applications
	jobs, err := db.GetAllJobApplications()
	if err != nil {
		log.Fatalf("Failed to get job applications: %v", err)
	}

	fmt.Printf("Database: %s\n", dbPath)
	fmt.Printf("Total job applications: %d\n\n", len(jobs))

	if len(jobs) == 0 {
		fmt.Println("No job applications found.")
		fmt.Println("Use 'go run cmd/debug/main.go add-test-data' to add sample data.")
		return
	}

	// Display all jobs
	for i, job := range jobs {
		fmt.Printf("--- Job Application %d ---\n", i+1)
		fmt.Printf("ID: %d\n", job.ID)
		fmt.Printf("Date Applied: %s\n", job.DateApplied.Format("2006-01-02"))
		fmt.Printf("Job Title: %s\n", job.JobTitle)
		fmt.Printf("Company: %s\n", job.Company)
		fmt.Printf("Status: %s\n", job.Status)
		if job.JobURL != "" {
			fmt.Printf("Job URL: %s\n", job.JobURL)
		}
		if job.Notes != "" {
			fmt.Printf("Notes: %s\n", job.Notes)
		}
		fmt.Printf("Created: %s\n", job.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", job.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	// Show status counts
	statusCounts, err := db.GetStatusCounts()
	if err != nil {
		log.Printf("Failed to get status counts: %v", err)
		return
	}

	fmt.Println("--- Status Summary ---")
	for status, count := range statusCounts {
		fmt.Printf("%s: %d\n", status, count)
	}
}

func addTestData(dbPath string) {
	fmt.Printf("Adding test data to database: %s\n", dbPath)

	// Initialize database
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Sample job applications
	testJobs := []*models.JobApplication{
		{
			DateApplied: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			JobTitle:    "Senior Software Engineer",
			Company:     "TechCorp",
			Status:      models.StatusApplied,
			JobURL:      "https://techcorp.com/jobs/123",
			Notes:       "Applied through company website. Looking for backend Go developers.",
		},
		{
			DateApplied: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			JobTitle:    "Full Stack Developer",
			Company:     "StartupCo",
			Status:      models.StatusInReview,
			JobURL:      "https://startupco.com/careers/dev",
			Notes:       "Interesting startup working on AI tools. Remote-first company.",
		},
		{
			DateApplied: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			JobTitle:    "Backend Engineer",
			Company:     "BigTech Inc",
			Status:      models.StatusInterview,
			JobURL:      "https://bigtech.com/positions/backend-eng",
			Notes:       "Phone screen scheduled for next week. They use microservices architecture.",
		},
	}

	for i, job := range testJobs {
		if err := db.CreateJobApplication(job); err != nil {
			log.Printf("Failed to create job application %d: %v", i+1, err)
		} else {
			fmt.Printf("Created job application: %s at %s\n", job.JobTitle, job.Company)
		}
	}

	fmt.Printf("Successfully added %d test job applications.\n", len(testJobs))
}
