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
	ID        string    `json:"id"`
	Payload   []byte    `json:"payload"`
	State     JobState  `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}

func newJob(payload []byte) Job {
	return Job{
		ID:        uuid.NewString(),
		Payload:   payload,
		State:     StateQueued,
		CreatedAt: time.Now(),
	}
}
