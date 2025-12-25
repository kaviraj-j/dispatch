package queue

import (
	"time"

	"github.com/google/uuid"
)

type JobState string

const (
	StateQueued     JobState = "QUEUED"
	StateInProgress JobState = "IN_PROGRESS"
	StateRetried    JobState = "RETRIED"
	StateDone       JobState = "DONE"
)

type Job struct {
	ID        string
	Payload   []byte
	State     JobState
	CreatedAt time.Time
}

func newJob(payload []byte) Job {
	return Job{
		ID:        uuid.NewString(),
		Payload:   payload,
		State:     StateQueued,
		CreatedAt: time.Now(),
	}
}
