package store

import (
	"sync"

	"go-taskqueue/internal/job"
)

type Store struct {
	mu sync.Mutex
	jobs map[string]*job.Job // key string : value job
}

func NewStore() *Store{
	return &Store{
		jobs: make(map[string]*job.Job),
	}
}
func (s *Store) Save(j*job.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[j.ID] = j
}

func (s *Store) Get(id string) (*job.Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	j :=s.jobs[id]
	if j == nil {
		return nil, false
	}
	return j, true
}