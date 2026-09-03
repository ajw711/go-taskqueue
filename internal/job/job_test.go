package job

import (
	"errors"
	"testing"
)

func TestNewJob_Success(t *testing.T) {
	// given & when
	j, err := NewJob("job-1", "email:send", "{}", PriorityHigh, 5)

	// then
	if err != nil {
		t.Fatalf("에러 발생: %v", err)
	}

	if j.Status != StatusPending  {
		 t.Fatalf("상태 체크 pending: %v", j.Status)
	}
}

func TestNewJob_EmptyID(t *testing.T) {
	_, err := NewJob("", "email:send", "{}", PriorityHigh, 3)

	if err == nil {
		t.Fatalf("ID가 비어있으면 에러가 나야 하는데 안 남")
	}

	if !errors.Is(err, ErrEmptyID) {
    t.Fatalf("ErrEmptyID가 나야 하는데 다른 에러 남: %v", err)
	}
}

func TestNewJob_EmptyType(t *testing.T) {
    _, err := NewJob("job-1", "", "{}", PriorityHigh, 3)
    if err == nil {
        t.Fatalf("Type이 비어있으면 에러가 나야 하는데 안 남")
    }

	if !errors.Is(err, ErrEmptyType) {
    t.Fatalf("ErrEmptyType이가 나야 하는데 다른 에러 남: %v", err)
	}
}

func TestNewJob_InvalidPriority(t *testing.T) {
	// "" 는 applyDefault가 normal로 채워줌 → 에러 안 남
	// 진짜 잘못된 값으로 테스트해야 함
	_, err := NewJob("job-1", "email:send", "{}", Priority("invalid"), 3)

	if err == nil {
		t.Fatalf("잘못된 priority면 에러가 나야 하는데 안 남")
	}

	if !errors.Is(err, ErrInvalidPriority) {
		t.Fatalf("ErrInvalidPriority가 나야 하는데 다른 에러 남: %v", err)
	}
}

func TestNewJob_DefaultMaxAttempts(t *testing.T) {
    // MaxAttempts 0으로 넘기면 3으로 채워지는지
    j, _ := NewJob("job-1", "email:send", "{}", PriorityHigh, 0)
    if j.MaxAttempts != 3 {
        t.Fatalf("기본값 3이어야 하는데: %v", j.MaxAttempts)
    }
}

func TestFromJSON_Success(t *testing.T) {
    data := []byte(`{"id":"job-1","type":"email:send","payload":"{}","priority":"high","status":"pending","max_attempts":3}`)
    j, err := FromJSON(data)

    if err != nil {
        t.Fatalf("에러 발생: %v", err)
    }

    if j.ID != "job-1" {  
        t.Fatalf("ID가 job-1인지 확인: %v", j.ID)
    }

    if j.Type != "email:send" {
        t.Fatalf("Type이 email:send인지 확인: %v", j.Type)
    }
}

func TestFromJSON_InvalidJSON(t *testing.T) {
    data := []byte(`{"id":"job-1","type";"email:send","payload":"{}","priority":"high","status":"pending","max_attempts":3}`)

    _, err := FromJSON(data)  // j 안 씀, _ 로 버림
    if err == nil {
        t.Fatalf("잘못된 JSON 에러 발생")
    }
}

func TestToJSON(t *testing.T) {
    // given
    j, _ := NewJob("job-1", "email:send", "{}", PriorityHigh, 3)

    // when
    data, err := j.ToJSON()

    // then
    if err != nil {
        t.Fatalf("ToJSON 에러: %v", err)
    }

    // 다시 FromJSON으로 복원해서 원본이랑 비교
    restored, err := FromJSON(data)
    if err != nil {
        t.Fatalf("복원 실패: %v", err)
    }

    if restored.ID != j.ID {
        t.Fatalf("ID가 달라요: %v", restored.ID)
    }
}

