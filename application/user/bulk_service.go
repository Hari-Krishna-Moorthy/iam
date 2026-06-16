package user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/job"
	"github.com/google/uuid"
)

// BulkCreateUsersRequest represents the payload for bulk user creation.
type BulkCreateUsersRequest struct {
	Users []RegistrationRequest `json:"users"`
}

type BulkService interface {
	SubmitBulkCreate(ctx context.Context, tenantID string, req BulkCreateUsersRequest) (string, error)
	GetJobStatus(ctx context.Context, jobID string) (*job.Job, error)
}

type bulkService struct {
	jobRepo job.Repository
}

func NewBulkService(jobRepo job.Repository) BulkService {
	return &bulkService{jobRepo: jobRepo}
}

func (s *bulkService) SubmitBulkCreate(ctx context.Context, tenantID string, req BulkCreateUsersRequest) (string, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	j := &job.Job{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		Type:      "bulk_create_users",
		Payload:   string(payload),
		Status:    job.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := s.jobRepo.Save(ctx, j); err != nil {
		return "", err
	}

	if err := s.jobRepo.Enqueue(ctx, "bulk_ops", j.ID); err != nil {
		return "", err
	}

	return j.ID, nil
}

func (s *bulkService) GetJobStatus(ctx context.Context, jobID string) (*job.Job, error) {
	return s.jobRepo.GetByID(ctx, jobID)
}
