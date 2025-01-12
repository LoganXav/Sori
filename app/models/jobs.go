package models

import (
	"time"
)

type JobStatus string

const (
	JobPending   JobStatus = "pending"
	JobRunning   JobStatus = "running"
	JobCompleted JobStatus = "completed"
	JobFailed    JobStatus = "failed"
)

type JobType string

const (
	JobTypeQC         JobType = "quality_control"
	JobTypeAlignment  JobType = "alignment"
	JobTypeDownstream JobType = "downstream_analysis"
)

type Job struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"` 
	FileID    string    `gorm:"not null" json:"file_id"`           // Reference to the FASTQ file
	Type      JobType   `gorm:"not null" json:"type"`             
	Status    JobStatus `gorm:"not null" json:"status"`            
	ResultURL string    `json:"result_url"`                       // S3 URL for the job's output
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`  
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` 
	Error     string    `json:"error,omitempty"`                 
}
