package queue

import (
	"errors"
	"sync"
)

type Queue struct {
	name       string
	mu         sync.Mutex
	jobs       []Job
	inProgress map[string]Job
}

// Errors
var (
	ErrQueueEmpty error = errors.New("no jobs in queue. queue is empty")
)

func NewQueue(name string) *Queue {
	return &Queue{
		name:       name,
		jobs:       make([]Job, 0),
		inProgress: make(map[string]Job),
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
	q.inProgress[job.ID] = job
	return &job, true
}

func (q *Queue) Ack(jobID string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, ok := q.inProgress[jobID]; !ok {
		return false
	}

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
