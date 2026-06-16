package worker

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/user"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/job"
)

// JobWorker processes background jobs from Redis.
type JobWorker struct {
	jobRepo     job.Repository
	userService user.Service
	queueName   string
	wg          *sync.WaitGroup
}

func NewJobWorker(jobRepo job.Repository, userService user.Service, queueName string) *JobWorker {
	return &JobWorker{
		jobRepo:     jobRepo,
		userService: userService,
		queueName:   queueName,
		wg:          &sync.WaitGroup{},
	}
}

// Start begins processing jobs in the background. It returns a function to gracefully stop.
func (w *JobWorker) Start(ctx context.Context) func() {
	stopCtx, cancel := context.WithCancel(ctx)
	w.wg.Add(1)

	go func() {
		defer w.wg.Done()
		log.Printf("Worker started for queue: %s", w.queueName)
		for {
			select {
			case <-stopCtx.Done():
				log.Printf("Worker stopped for queue: %s", w.queueName)
				return
			default:
				jobID, err := w.jobRepo.Dequeue(stopCtx, w.queueName)
				if err != nil {
					log.Printf("Error dequeueing job: %v", err)
					time.Sleep(1 * time.Second)
					continue
				}
				if jobID == "" {
					continue // Timeout, loop again
				}

				w.processJob(stopCtx, jobID)
			}
		}
	}()

	return func() {
		cancel()
		w.wg.Wait()
	}
}

func (w *JobWorker) processJob(ctx context.Context, jobID string) {
	j, err := w.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		log.Printf("Failed to fetch job %s: %v", jobID, err)
		return
	}

	j.Status = job.StatusProcessing
	_ = w.jobRepo.Save(ctx, j)

	var result map[string]interface{}

	switch j.Type {
	case "bulk_create_users":
		result = w.processBulkCreateUsers(ctx, j)
	default:
		j.Status = job.StatusFailed
		j.Result = "Unknown job type"
		_ = w.jobRepo.Save(ctx, j)
		return
	}

	resBytes, _ := json.Marshal(result)
	j.Result = string(resBytes)
	
	if result["failed_count"].(int) > 0 {
		// You might consider partial success as completed or failed depending on requirements.
		j.Status = job.StatusCompleted 
	} else {
		j.Status = job.StatusCompleted
	}
	
	_ = w.jobRepo.Save(ctx, j)
	log.Printf("Job %s processed with status %s", j.ID, j.Status)
}

func (w *JobWorker) processBulkCreateUsers(ctx context.Context, j *job.Job) map[string]interface{} {
	var req user.BulkCreateUsersRequest
	if err := json.Unmarshal([]byte(j.Payload), &req); err != nil {
		return map[string]interface{}{"error": "invalid payload", "failed_count": 1}
	}

	successCount := 0
	failedCount := 0
	var errors []string

	for _, uReq := range req.Users {
		uReq.TenantID = j.TenantID
		_, err := w.userService.RegisterUser(ctx, uReq)
		if err != nil {
			failedCount++
			errors = append(errors, "user "+uReq.Username+": "+err.Error())
		} else {
			successCount++
		}
	}

	return map[string]interface{}{
		"success_count": successCount,
		"failed_count":  failedCount,
		"errors":        errors,
	}
}
