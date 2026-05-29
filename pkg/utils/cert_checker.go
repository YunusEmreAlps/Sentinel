package utils

import (
	"context"
	"sync"

	"sentinel/internal/models"
	"sentinel/pkg/logger"
)

const (
	// DefaultWorkerPoolSize is the default number of concurrent workers for certificate checking
	DefaultWorkerPoolSize = 10
)

// CertificateCheckJob represents a single certificate check job
type CertificateCheckJob struct {
	Domain     string
	ExpireDays int
}

// CertificateCheckResult represents the result of a certificate check
type CertificateCheckResult struct {
	Domain    string
	IsExpired bool
	Log       *models.Log
	Error     error
}

// CertificateChecker provides concurrent certificate checking with a worker pool
type CertificateChecker struct {
	workerCount int
}

// NewCertificateChecker creates a new CertificateChecker with the specified worker count
func NewCertificateChecker(workerCount int) *CertificateChecker {
	if workerCount <= 0 {
		workerCount = DefaultWorkerPoolSize
	}
	return &CertificateChecker{
		workerCount: workerCount,
	}
}

// CheckCertificates checks multiple certificates concurrently
func (cc *CertificateChecker) CheckCertificates(ctx context.Context, domains []string, expireDays int) []models.Log {
	jobs := make(chan CertificateCheckJob, len(domains))
	results := make(chan CertificateCheckResult, len(domains))

	// Create worker pool
	var wg sync.WaitGroup
	for i := 0; i < cc.workerCount; i++ {
		wg.Add(1)
		go cc.worker(ctx, &wg, jobs, results)
	}

	// Send jobs
	for _, domain := range domains {
		jobs <- CertificateCheckJob{
			Domain:     domain,
			ExpireDays: expireDays,
		}
	}
	close(jobs)

	// Wait for all workers to complete in a separate goroutine
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var logs []models.Log
	for result := range results {
		if result.Error != nil {
			logger.CLogger.Errorf("ERROR: %s - %v", result.Domain, result.Error)
			continue
		}

		if result.IsExpired && result.Log != nil {
			logs = append(logs, *result.Log)
		} else if result.Log != nil {
			logger.CLogger.Infof("INFO: %s - %s", result.Domain, result.Log.Message)
		}
	}

	return logs
}

// worker processes certificate check jobs from the jobs channel
func (cc *CertificateChecker) worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan CertificateCheckJob, results chan<- CertificateCheckResult) {
	defer wg.Done()

	for job := range jobs {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			results <- CertificateCheckResult{
				Domain: job.Domain,
				Error:  ctx.Err(),
			}
			return
		default:
		}

		// Perform certificate check with retry logic
		var lastErr error
		var isExpired bool
		var log *models.Log

		config := DefaultRetryConfig()
		err := RetryWithContext(ctx, config, func(attemptCtx context.Context) error {
			var checkErr error
			isExpired, log = CheckDomainCertificateWithContext(attemptCtx, job.Domain, job.ExpireDays)
			if !isExpired && log == nil {
				checkErr = &CertificateCheckError{Domain: job.Domain}
			}
			return checkErr
		}, job.Domain)

		if err != nil {
			lastErr = err
		}

		results <- CertificateCheckResult{
			Domain:    job.Domain,
			IsExpired: isExpired,
			Log:       log,
			Error:     lastErr,
		}
	}
}

// CertificateCheckError represents an error during certificate checking
type CertificateCheckError struct {
	Domain string
}

func (e *CertificateCheckError) Error() string {
	return "failed to check certificate for " + e.Domain
}
