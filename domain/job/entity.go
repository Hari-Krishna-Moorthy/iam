package job

import (
	"context"
	"time"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

// Job represents an asynchronous background task.
type Job struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	Type       string    `json:"type"`       // e.g., "bulk_create_users"
	Payload    string    `json:"payload"`    // JSON encoded payload
	Status     Status    `json:"status"`
	Result     string    `json:"result"`     // JSON encoded result/errors
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Repository defines the persistence interface for Job.
type Repository interface {
	Save(ctx context.Context, job *Job) error
	GetByID(ctx context.Context, id string) (*Job, error)
	Enqueue(ctx context.Context, queueName string, jobID string) error
	Dequeue(ctx context.Context, queueName string) (string, error) // Returns jobID
}
