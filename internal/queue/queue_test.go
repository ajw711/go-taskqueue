package queue

import (
	"fmt"
	"testing"
	"time"

	"go-taskqueue/internal/job"
)

func TestQueue_ProcessesJobsConcurrently(t *testing.T) {
	// given: 더미 processor + Queue 생성 + Start
	processor := func(j *job.Job) error {
		fmt.Print("처리 중:", j.ID, j.Type)
		return nil
	}
	q := NewQueue(3, processor)
	q.Start()

	// when: 작업 여러 개 등록
	j1, err := job.NewJob("job-1", "email:send", "{}", job.PriorityHigh, 5)
	if err != nil {
		t.Fatalf("NewJob 생성 실패: %v", err)
	}
	if err := q.Submit(j1); err != nil {
		t.Fatalf("Submit 실패: %v", err)
	}

	j2, err := job.NewJob("job-2", "email:send", "{}", job.PriorityHigh, 5)
	if err != nil {
		t.Fatalf("NewJob 생성 실패: %v", err)
	}
	if err := q.Submit(j2); err != nil {
		t.Fatalf("Submit 실패: %v", err)
	}

	j3, err := job.NewJob("job-3", "email:send", "{}", job.PriorityHigh, 5)
	if err != nil {
		t.Fatalf("NewJob 생성 실패: %v", err)
	}
	if err := q.Submit(j3); err != nil {
		t.Fatalf("Submit 실패: %v", err)
	}

	// then: 잠깐 대기 후 상태 확인
	time.Sleep(500 * time.Millisecond)
	if j1.Status != job.StatusCompleted {
		t.Errorf("기대한 상태는 %s인데 실제는 %s", job.StatusCompleted, j1.Status)
	}
	if j2.Status != job.StatusCompleted {
		t.Errorf("기대한 상태는 %s인데 실제는 %s", job.StatusCompleted, j2.Status)
	}
	if j3.Status != job.StatusCompleted {
		t.Errorf("기대한 상태는 %s인데 실제는 %s", job.StatusCompleted, j3.Status)
	}
}