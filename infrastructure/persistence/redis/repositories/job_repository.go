package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/job"
	"github.com/redis/go-redis/v9"
)

type jobRepository struct {
	client *redis.Client
}

func NewJobRepository(client *redis.Client) job.Repository {
	return &jobRepository{client: client}
}

func (r *jobRepository) Save(ctx context.Context, j *job.Job) error {
	j.UpdatedAt = time.Now()
	data, err := json.Marshal(j)
	if err != nil {
		return err
	}
	// Jobs expire after 24 hours to keep Redis clean
	return r.client.Set(ctx, "job:"+j.ID, data, 24*time.Hour).Err()
}

func (r *jobRepository) GetByID(ctx context.Context, id string) (*job.Job, error) {
	data, err := r.client.Get(ctx, "job:"+id).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("job not found")
		}
		return nil, err
	}

	var j job.Job
	if err := json.Unmarshal(data, &j); err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *jobRepository) Enqueue(ctx context.Context, queueName string, jobID string) error {
	return r.client.LPush(ctx, "queue:"+queueName, jobID).Err()
}

func (r *jobRepository) Dequeue(ctx context.Context, queueName string) (string, error) {
	// Blocking pop, waits up to 5 seconds
	res, err := r.client.BRPop(ctx, 5*time.Second, "queue:"+queueName).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Timeout, no job
		}
		return "", err
	}
	if len(res) == 2 {
		return res[1], nil
	}
	return "", nil
}
