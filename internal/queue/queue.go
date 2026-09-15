package queue

import (
	"fmt"
	"go-taskqueue/internal/dlq"
	"go-taskqueue/internal/job"
	"go-taskqueue/internal/store"
	"time"
)

type Queue struct {
	highJobs chan *job.Job
	normalJobs chan *job.Job
	processor   func(*job.Job) error
	workerCount int
	dlq *dlq.DLQ
	store *store.Store
}

func NewQueue(workerCount int, processor func(*job.Job) error) *Queue {

	queue := &Queue{
		highJobs: make(chan *job.Job, workerCount*10),
		normalJobs:        make(chan *job.Job, workerCount*10),
		processor:   processor,
		workerCount: workerCount,
		dlq: dlq.NewDLQ(),
		store: store.NewStore(),
	}
	return queue
}

func (q *Queue) Submit(j *job.Job) error {
	if j == nil {
		return ErrNilJob
	}
	q.store.Save(j)
	if j.Priority == job.PriorityHigh {
		q.highJobs <- j
	} else {
		q.normalJobs <- j
	}
	return nil
}

func (q *Queue) Start() {
	for i := 1; i <= q.workerCount; i++ {
		go q.worker(i) // 워커 수만큼 고루틴 실행
	}
}

func (q *Queue) worker(workerId int) {
	for {
		var j *job.Job

		select {
		case j = <-q.highJobs:
			q.handle(j, workerId)
			continue
		default:
		}

		select {
		case j = <-q.highJobs:
			q.handle(j, workerId)
		case j = <-q.normalJobs:
			q.handle(j, workerId)
		}
	}
}

func (q *Queue) handle(j *job.Job, workerId int) {
	fmt.Printf("[worker-%d] 처리 시작: %s (priority=%s)\n", workerId, j.ID, j.Priority)
	j.MarkProcessing()
	err := q.processor(j)
	if err == nil {
		j.MarkCompleted()
		return
	}

	j.MarkFailed(err)

	if j.IsExhausted() {
		j.MarkDead()
		q.dlq.Add(j)
		return
	}
	backoff := time.Duration(j.CurrentAttempt) * 2 * time.Second
	time.AfterFunc(backoff,func() {
		q.Submit(j)
	})
}

func (q *Queue) GetJob(id string) (*job.Job, bool) {
	return q.store.Get(id)
}

func (q *Queue) ListDeadLetters() []*job.Job {
	return q.dlq.List()
}