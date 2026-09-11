package dlq

import (
	"sync"

	"go-taskqueue/internal/job"
)

type DLQ struct {
	mu sync.Mutex
	jobs []*job.Job
}

func NewDLQ() *DLQ {
	return &DLQ{}
}

func (d *DLQ) Add(j *job.Job) {
	d.mu.Lock() // 지금부터 나만 이 슬라이스를 건드릴 거야 다른 애들은 기다려
	defer d.mu.Unlock()  // 함수 끝나면 자동으로 "이제 다 썼어, 다음 사람 써도 돼
	d.jobs = append(d.jobs, j)
}

func (d *DLQ) List() []*job.Job {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.jobs
}