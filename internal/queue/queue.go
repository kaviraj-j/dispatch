package queue

import (
	"errors"
	"sync"
	"time"
)

type Queue struct {
	name       string
	mu         sync.Mutex
	jobs       []Job
	inProgress map[string]JobInProgress
}

type JobInProgress struct {
	job      Job
	jobTimer *time.Timer
}

type TimeoutPayload struct {
	queueName string
	jobID     string
}

// Errors
var (
	ErrQueueEmpty error = errors.New("no jobs in queue. queue is empty")
)

func NewQueue(name string) *Queue {
	return &Queue{
		name:       name,
		jobs:       make([]Job, 0),
		inProgress: make(map[string]JobInProgress),
	}
}

func (q *Queue) Name() string {
	return q.name
}

func (q *Queue) Enqueue(payload []byte) Job {
	q.mu.Lock()
	defer q.mu.Unlock()
	job := newJob(payload)
	q.jobs = append(q.jobs, job)
	return job
}

func (q *Queue) Dequeue() (*Job, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.jobs) == 0 {
		return nil, false
	}
	job := q.jobs[0]
	q.jobs = q.jobs[1:]
	jobTimer := time.AfterFunc(5*time.Second, q.JobTimeoutHanlder(job.ID))
	q.inProgress[job.ID] = JobInProgress{
		job:      job,
		jobTimer: jobTimer,
	}
	return &job, true
}

func (q *Queue) Ack(jobID string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	jobInProgress, ok := q.inProgress[jobID]
	if !ok {
		return false
	}
	jobInProgress.jobTimer.Stop()
	delete(q.inProgress, jobID)
	return true
}

// Size returns the number of queued (not in-progress) jobs.
func (q *Queue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.jobs)
}

// InProgressSize returns the number of in-progress jobs.
func (q *Queue) InProgressSize() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.inProgress)
}

func (q *Queue) JobTimeoutHanlder(jobID string) func() {
	return func() {
		q.mu.Lock()
		defer q.mu.Unlock()

		jobInProgress, ok := q.inProgress[jobID]
		if !ok {
			return
		}
		delete(q.inProgress, jobID)
		q.jobs = append(q.jobs, jobInProgress.job)
	}
}
