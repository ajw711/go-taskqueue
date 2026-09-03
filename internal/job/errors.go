package job

import "fmt"

type ErrorCode string

const (
    ErrCodeEmptyID         ErrorCode = "JOB_EMPTY_ID"
    ErrCodeEmptyType       ErrorCode = "JOB_EMPTY_TYPE"
	ErrCodeInvalidPriority ErrorCode = "JOB_INVALID_PRIORITY"
)

type JobError struct {
	Code ErrorCode
	Message string
}

func (e *JobError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}


var (
	ErrEmptyID = &JobError {
		Code:    ErrCodeEmptyID,
        Message: "job id는 비어 있을 수 없습니다",
	}
	ErrEmptyType = &JobError {
		Code:    ErrCodeEmptyType,
        Message: "job type은 비어 있을 수 없습니다",
	}
	ErrInvalidPriority = &JobError {
		Code:    ErrCodeInvalidPriority,
        Message: "유효하지 않은 job priority입니다",
	}
)

func (j *Job) Validate() error {
	if j.ID == "" {
		return ErrEmptyID
	}

	if j.Type == "" {
		return ErrEmptyType
	}

	if j.Priority != PriorityLow && j.Priority != PriorityNormal && j.Priority != PriorityHigh {
        return ErrInvalidPriority
    }

	return nil
}
