package queue

import (
	"go-taskqueue/internal/job"
)

type Queue struct {
	jobs        chan *job.Job
	processor   func(*job.Job) error
	workerCount int
}

func NewQueue(workerCount int, processor func(*job.Job) error) *Queue {

	queue := &Queue{
		jobs:        make(chan *job.Job, workerCount*10),
		processor:   processor,
		workerCount: workerCount,
	}
	return queue
}

func (q *Queue) Submit(j *job.Job) error {
	if j == nil {
		return ErrNilJob
	}
	q.jobs <- j
	return nil
}

func (q *Queue) Start() {
	for i := 1; i <= q.workerCount; i++ {
		go q.worker(i) // 워커 수만큼 고루틴 실행
	}
}

func (q *Queue) worker(workerId int) {
	for j := range q.jobs {
		j.MarkProcessing()

		err := q.processor(j)
		if err != nil {
			j.MarkFailed(err)
		} else {
			j.MarkCompleted()
		}
	}
}
