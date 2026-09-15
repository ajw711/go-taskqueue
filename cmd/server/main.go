package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go-taskqueue/internal/handler"
	"go-taskqueue/internal/job"
	"go-taskqueue/internal/queue"
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

	h := handler.NewHandler(q)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", h.SubmitJob)
	mux.HandleFunc("GET /jobs/{id}", h.GetJob)
	mux.HandleFunc("GET /dead-letters", h.DeadLetters)

	fmt.Println("서버 시작: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
