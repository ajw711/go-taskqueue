package queue

import (
	"errors"
	"go-taskqueue/internal/job"
	"testing"
	"time"
)
func TestQueue_RetryAndDLQ(t *testing.T)  {

	// given
	attemptCount := 0

	// 항상 실패하는 processor -> 결국 DLQ로 가야 함
	processor := func(j *job.Job) error {
		attemptCount++
		return errors.New("일부러 실패시킴")
	}

	q := NewQueue(1, processor)

	// when
	q.Start()
 	// 최대 3번 시도
	j, _:=job.NewJob("fail-job", "email:send", "{}", job.PriorityNormal, 3)
	q.Submit(j)
	// 재시도 텀이 2초니까, 3번 다 시도되려면 넉넉하게 대기 필요
	time.Sleep(15 * time.Second)

	// then
	// attemptCount가 3이어야 함
	if attemptCount != 3 {
		       t.Fatalf("3이여함")
	}

	dlqJobs := q.dlq.List()
	if len(dlqJobs) != 1 {
		 t.Fatalf("리스트 길이가 있어야함")
	}
	if q.dlq == nil {
		 t.Fatalf("dlq 있어야함")
	}

	if j.Status != job.StatusDead {
		 t.Fatalf("상태가 Dead이어야함")
	}

}