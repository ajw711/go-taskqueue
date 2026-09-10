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

	j2, err := job.NewJob("job-2", "email:send", "{}", job.PriorityNormal, 5)
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

func TestQueue_PriorityJobsProcessedFirst(t *testing.T) {
	processedOrder := []string{}

	processor := func(j *job.Job) error {
		time.Sleep(100 * time.Millisecond)
		processedOrder = append(processedOrder, j.ID)
		fmt.Printf("처리 완료: %s (priority=%s)\n", j.ID, j.Priority)
		return nil
	}

	q := NewQueue(1, processor) // 워커 1개로 경쟁 상황 강제
	q.Start()

	for i :=1; i <= 3; i++ {
		j, _ := job.NewJob(fmt.Sprintf("normal-%d", i), "email:send", "{}", job.PriorityNormal, 5)
		q.Submit(j)
	}

	time.Sleep(100 * time.Millisecond)

	// noraml 다음으로 high 등록
	highJob, _:= job.NewJob("high-1", "email:send", "{}", job.PriorityHigh, 5)
	q.Submit(highJob)

	time.Sleep(1 * time.Second) // 모든 작업이 끝날 때까지 충분히 대기
	// [normal-1, high-1, normal-2, normal-3]처럼 나와야함
	fmt.Print("처리 순서: ",processedOrder)
}