package models

import "time"

// JobApplication represents a job application entry
type JobApplication struct {
	ID          int       `json:"id" db:"id"`
	DateApplied time.Time `json:"date_applied" db:"date_applied"`
	JobTitle    string    `json:"job_title" db:"job_title"`
	Company     string    `json:"company" db:"company"`
	Status      string    `json:"status" db:"status"`
	JobURL      string    `json:"job_url" db:"job_url"`
	Notes       string    `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// JobStatus constants for common statuses
const (
	StatusApplied     = "Applied"
	StatusInReview    = "In Review"
	StatusPhoneScreen = "Phone Screen"
	StatusInterview   = "Interview"
	StatusTechnical   = "Technical Test"
	StatusOffer       = "Offer"
	StatusRejected    = "Rejected"
	StatusWithdrawn   = "Withdrawn"
	StatusNoResponse  = "No Response"
)

// GetCommonStatuses returns a list of common job application statuses
func GetCommonStatuses() []string {
	return []string{
		StatusApplied,
		StatusInReview,
		StatusPhoneScreen,
		StatusInterview,
		StatusTechnical,
		StatusOffer,
		StatusRejected,
		StatusWithdrawn,
		StatusNoResponse,
	}
}
