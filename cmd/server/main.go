package main

import (
	"fmt"
	"go-taskqueue/internal/job"
	"go-taskqueue/internal/queue"
	"time"
)

func main() {

	fakeProcessor := func(j *job.Job) error {
		fmt.Printf("--> [시작] Job %s 처리 시작 (Payload: %s)\n", j.ID, j.Payload)
		time.Sleep(1 * time.Second) // 1초 동안 일하는 척
		fmt.Printf("<-- [완료] Job %s 처리 완료\n", j.ID)
		return nil
	}

	q := queue.NewQueue(3, fakeProcessor)
	q.Start()
	fmt.Println("=== TaskQueue 워커 3명 가동 시작 ===")

	for i := 1; i <= 6; i++ {
		j, _ := job.NewJob(
			fmt.Sprintf("job-%d", i),
			"email:send",
			fmt.Sprintf("user%d@test.com", i),
			job.PriorityNormal,
			3,
		)
		err := q.Submit(j)
		if err != nil {
			return
		}
	}

	time.Sleep(3 * time.Second)
	fmt.Println("=== 모든 작업 처리 종료 ===")

}
