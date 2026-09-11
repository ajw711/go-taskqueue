package job

import (
	"encoding/json"
	"time"
)

// 작업 상태
type Status string
// 작업 우선순위
type Priority string
const (
	StatusPending Status = "pending"
	StatusProcessing Status = "processing"
    StatusCompleted  Status = "completed"
    StatusFailed     Status = "failed"
    StatusDead       Status = "dead"
)

const (
	PriorityLow Priority = "low"
	PriorityNormal Priority = "normal"
	PriorityHigh Priority = "high"
)

type Job struct {
    ID             string    `json:"id"`
    Type           string    `json:"type"`            // 작업 종류 (예: "email:send", "resize:image")
    Payload        string    `json:"payload"`         // 작업에 필요한 데이터
    Priority       Priority  `json:"priority"`
    Status         Status    `json:"status"`
    MaxAttempts    int       `json:"max_attempts"`    // 최대 재시도 횟수
    CurrentAttempt int       `json:"current_attempt"` // 현재 시도 횟수
    LastError      string    `json:"last_error,omitempty"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

func NewJob(id string, jobType string, payload string, priority Priority, maxAttempts int) (*Job, error) {

	now := time.Now()
	j := &Job{
        ID:             id,
        Type:           jobType,
        Payload:        payload,
        Priority:       priority,
        Status:         StatusPending,
        MaxAttempts:    maxAttempts,
        CurrentAttempt: 0,
        CreatedAt:      now,
        UpdatedAt:      now,
    }

	j.applyDefault()

 	if err := j.Validate(); err != nil {
            return nil, err
    }
	return j, nil
}

func (j *Job) applyDefault() {

	if j.Priority == "" {
		j.Priority = PriorityNormal
	}

	if j.MaxAttempts <= 0 {
		j.MaxAttempts = 3
	}
}

func FromJSON(data []byte) (*Job, error) {
	j := Job{}
	if err := json.Unmarshal(data, &j); err != nil {
		return nil, err
	}
	j.applyDefault()
	if err :=j.Validate(); err != nil {
		return nil, err
	}
	return &j, nil
}

func (j *Job) ToJSON() ([]byte, error) {
	data, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (j *Job) MarkProcessing() {
	j.Status = StatusProcessing
	j.UpdatedAt = time.Now()
}

func (j *Job) MarkCompleted() {
	j.Status = StatusCompleted
	j.UpdatedAt = time.Now()
}

func (j *Job) MarkFailed(err error) {
	j.Status = StatusFailed
	j.LastError = err.Error()
	j.UpdatedAt = time.Now()
	j.CurrentAttempt++
}

func (j *Job) MarkDead() {
	j.Status = StatusDead
	j.UpdatedAt = time.Now()
}

func (j *Job) IsExhausted() bool {
	return j.CurrentAttempt >= j.MaxAttempts
}