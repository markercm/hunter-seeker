package database

import (
	"database/sql"
	"errors"
	"fmt"

	"hunter-seeker/internal/models"

	_ "modernc.org/sqlite"
)

// Define specific errors
var (
	ErrJobNotFound = errors.New("job application not found")
)

type DB struct {
	conn *sql.DB
}

// New creates a new database connection and sets up tables
func New(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// createTables creates the necessary database tables
func (db *DB) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS job_applications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date_applied DATE NOT NULL,
		job_title TEXT NOT NULL,
		company TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'Applied',
		job_url TEXT,
		notes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TRIGGER IF NOT EXISTS update_job_applications_updated_at
	AFTER UPDATE ON job_applications
	BEGIN
		UPDATE job_applications SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;
	`

	_, err := db.conn.Exec(query)
	return err
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// CreateJobApplication creates a new job application
func (db *DB) CreateJobApplication(job *models.JobApplication) error {
	query := `
	INSERT INTO job_applications (date_applied, job_title, company, status, job_url, notes)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := db.conn.Exec(query, job.DateApplied, job.JobTitle, job.Company, job.Status, job.JobURL, job.Notes)
	if err != nil {
		return fmt.Errorf("failed to create job application: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	job.ID = int(id)
	return nil
}

// GetJobApplication retrieves a job application by ID
func (db *DB) GetJobApplication(id int) (*models.JobApplication, error) {
	query := `
	SELECT id, date_applied, job_title, company, status, job_url, notes, created_at, updated_at
	FROM job_applications
	WHERE id = ?
	`

	job := &models.JobApplication{}
	err := db.conn.QueryRow(query, id).Scan(
		&job.ID, &job.DateApplied, &job.JobTitle, &job.Company,
		&job.Status, &job.JobURL, &job.Notes, &job.CreatedAt, &job.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job application not found")
		}
		return nil, fmt.Errorf("failed to get job application: %w", err)
	}

	return job, nil
}

// GetAllJobApplications retrieves all job applications, ordered by date applied (newest first)
func (db *DB) GetAllJobApplications() ([]*models.JobApplication, error) {
	query := `
	SELECT id, date_applied, job_title, company, status, job_url, notes, created_at, updated_at
	FROM job_applications
	ORDER BY date_applied DESC, created_at DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query job applications: %w", err)
	}
	defer rows.Close()

	var jobs []*models.JobApplication
	for rows.Next() {
		job := &models.JobApplication{}
		err := rows.Scan(
			&job.ID, &job.DateApplied, &job.JobTitle, &job.Company,
			&job.Status, &job.JobURL, &job.Notes, &job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job application: %w", err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// UpdateJobApplication updates an existing job application
func (db *DB) UpdateJobApplication(job *models.JobApplication) error {
	query := `
	UPDATE job_applications
	SET date_applied = ?, job_title = ?, company = ?, status = ?, job_url = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	result, err := db.conn.Exec(query, job.DateApplied, job.JobTitle, job.Company, job.Status, job.JobURL, job.Notes, job.ID)
	if err != nil {
		return fmt.Errorf("failed to update job application: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("job application not found")
	}

	return nil
}

// DeleteJobApplication deletes a job application by ID
func (db *DB) DeleteJobApplication(id int) error {
	query := `DELETE FROM job_applications WHERE id = ?`

	result, err := db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete job application: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrJobNotFound
	}

	return nil
}

// GetJobApplicationsByStatus retrieves job applications filtered by status
func (db *DB) GetJobApplicationsByStatus(status string) ([]*models.JobApplication, error) {
	query := `
	SELECT id, date_applied, job_title, company, status, job_url, notes, created_at, updated_at
	FROM job_applications
	WHERE status = ?
	ORDER BY date_applied DESC, created_at DESC
	`

	rows, err := db.conn.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query job applications by status: %w", err)
	}
	defer rows.Close()

	var jobs []*models.JobApplication
	for rows.Next() {
		job := &models.JobApplication{}
		err := rows.Scan(
			&job.ID, &job.DateApplied, &job.JobTitle, &job.Company,
			&job.Status, &job.JobURL, &job.Notes, &job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job application: %w", err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// GetStatusCounts returns counts of job applications by status
func (db *DB) GetStatusCounts() (map[string]int, error) {
	query := `
	SELECT status, COUNT(*) as count
	FROM job_applications
	GROUP BY status
	ORDER BY count DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query status counts: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		err := rows.Scan(&status, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan status count: %w", err)
		}
		counts[status] = count
	}

	return counts, nil
}
