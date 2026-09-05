package queue

import "fmt"

type ErrorCode string

const (
	ErrCodeNilJob    ErrorCode = "QUEUE_NIL_JOB"
	ErrCodeQueueFull ErrorCode = "QUEUE_FULL"
)

type QueueError struct {
	Code    ErrorCode
	Message string
}

func (e *QueueError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

var (
	ErrNilJob = &QueueError{
		Code:    ErrCodeNilJob,
		Message: "job은 nil일 수 없습니다",
	}
	ErrQueueFull = &QueueError{
		Code:    ErrCodeQueueFull,
		Message: "큐가 가득 찼습니다",
	}
)
