package queue

import (
	"go-taskqueue/internal/job"
)

type Queue struct {
	jobs chan *job.Job
	processor func(*job.Job) error
	workerCount int
}

func NewQueue(workerCount int, processor func(*job.Job) error) *Queue {

	queue := &Queue{
		jobs: make(chan *job.Job, workerCount * 10),
		processor: processor,
		workerCount: workerCount,
	}
	return queue
}