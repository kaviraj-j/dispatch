package queue

import (
	"errors"
	"sync"
)

// Manager owns and manages all queues.
type Manager struct {
	mu     sync.RWMutex
	queues map[string]*Queue
}

// NewManager initializes a queue manager.
func NewManager() *Manager {
	return &Manager{
		queues: make(map[string]*Queue),
	}
}

// CreateQueue explicitly creates a new queue.
// Returns an error if the queue already exists.
func (m *Manager) CreateQueue(name string) (*Queue, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.queues[name]; exists {
		return nil, errors.New("queue already exists")
	}

	q := NewQueue(name)
	m.queues[name] = q
	return q, nil
}

// GetQueue returns a queue if it exists.
func (m *Manager) GetQueue(name string) (*Queue, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	q, ok := m.queues[name]
	return q, ok
}

// EnsureQueue returns an existing queue or creates it if missing.
func (m *Manager) EnsureQueue(name string) *Queue {
	m.mu.Lock()
	defer m.mu.Unlock()

	if q, ok := m.queues[name]; ok {
		return q
	}

	q := NewQueue(name)
	m.queues[name] = q
	return q
}

// DeleteQueue removes a queue.
// Returns false if the queue does not exist.
func (m *Manager) DeleteQueue(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.queues[name]; !ok {
		return false
	}

	delete(m.queues, name)
	return true
}

// ListQueues returns all queue names.
func (m *Manager) ListQueues() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.queues))
	for name := range m.queues {
		names = append(names, name)
	}
	return names
}
