# go-taskqueue

Go의 동시성 모델(goroutine, channel)을 사용해 분산 작업 큐 시스템을 설계했습니다. 외부 미들웨어(Kafka, Redis) 없이, 언어 자체 기능만으로 워커 풀, 우선순위 큐, 재시도, Dead Letter Queue를 구성했습니다.

## 왜 만들었는가

백그라운드 작업 처리(이메일 발송, 알림, 배치)는 실무에서 흔히 마주치는 문제이고, Go는 고루틴 기반 동시성으로 이런 문제를 다루는 데 강점이 있는 언어입니다. 외부 라이브러리에 기대지 않고, 워커 풀, 우선순위 큐, 재시도, DLQ 같은 분산 시스템의 핵심 패턴을 언어 레벨에서 구성했습니다.

## 아키텍처

```
[클라이언트 프로그램] --HTTP(JSON)--> [handler] --Submit()--> [queue]
  (언어 무관)                                                     |
                                                        highJobs / normalJobs (channel)
                                                                     |
                                                          워커 goroutine 풀 (select 기반)
                                                                     |
                                                    성공 -> completed / 실패 -> 재시도 or dlq
                                                                     |
                                             [store] (상태 조회용, mutex+map)   [dlq] (mutex+slice)
```

### 패키지 구성

| 패키지 | 역할 |
|---|---|
| `internal/job` | Job 데이터 모델, 상태/우선순위 타입, 검증, JSON 직렬화, 상태 전이 메서드 |
| `internal/queue` | 우선순위 채널, goroutine 워커 풀, 재시도(지수 백오프) 오케스트레이션 |
| `internal/dlq` | 재시도 소진된 작업 격리 보관 (동시성 안전) |
| `internal/store` | Job ID로 현재 상태 조회 (동시성 안전) |
| `internal/handler` | HTTP API 노출 |
| `cmd/server` | 전체 조립 및 서버 기동 |

## 설계 결정과 트레이드오프

**우선순위 큐 — 완벽한 보장이 아닌 soft priority**
채널 두 개(high/normal) + 이중 `select` 패턴을 사용했습니다. 1단계에서 `default`로 high를 non-blocking 확인하고, 2단계에서 두 채널을 동시에 blocking 대기합니다. 이 구조는 극히 짧은 타이밍(1단계 통과 직후 두 채널에 동시에 값이 도착하는 경우)에 normal이 먼저 처리될 가능성이 있습니다. 완벽한 우선순위가 필요하면 `container/heap` 기반 자료구조가 필요하지만, 이번 프로젝트 범위(동시성 패턴 구현)에는 이 정도의 soft priority로 설정했습니다.

**재시도 — 워커를 막지 않는 논블로킹 재시도**
실패 시 `time.Sleep`으로 워커 자신을 재우는 대신, `time.AfterFunc`으로 별도 타이머를 걸어 워커는 즉시 다음 작업으로 넘어가도록 했습니다.

**영속성 없음 — 인메모리 전용**
채널, store, dlq 모두 프로세스 메모리에만 존재합니다. 서버 재시작 시 처리 중이던 데이터는 유실됩니다. Core NATS와 유사한 포지션이며, Kafka/NATS JetStream처럼 디스크 기반 영속성을 가지려면 WAL(Write-Ahead Log) 또는 Redis 같은 외부 저장소 연동이 필요합니다. 동시성 모델 구현이 목적이라 영속성은 스코프에서 제외했습니다.

## 실행 방법

```bash
go run ./cmd/server
```

기본적으로 `:8080` 포트에서 서버가 기동됩니다.

## API

| Method | Path | 설명 |
|---|---|---|
| POST | `/jobs` | 작업 등록 |
| GET | `/jobs/{id}` | 작업 상태 조회 |
| GET | `/dead-letters` | DLQ 목록 조회 |

**POST /jobs 요청 예시**
```json
{
  "type": "email:send",
  "payload": "user@example.com",
  "priority": "high",
  "max_attempts": 3
}
```

## 테스트

**유닛 테스트** (`go test ./...`)
- `job`: 생성/검증/상태전이/직렬화
- `queue`: 워커 풀 동시 처리, 우선순위 처리 순서, 재시도 및 DLQ 격리

**별도 클라이언트를 통한 통합 테스트** ([go-taskqueue-client](https://github.com/<github계정>/go-taskqueue-client), 별도 저장소/별도 Go 모듈)
서버의 `internal` 패키지를 참조하지 않고, 순수 HTTP 통신만으로 서버와 통신하는 독립 프로그램으로 작성했습니다.

- 단일 작업 등록 → 상태 조회 (pending → processing → completed 전이 로그를 기록했습니다)
- `sync.WaitGroup` + goroutine으로 50개 작업을 동시에 등록했고, 워커 3개가 병렬로 처리하는 로그를 기록했습니다

**성능 측정(처리량/지연시간)은 이번 검증 범위 밖입니다.** 로컬 환경에서의 자체 스크립트 측정은 서버-클라이언트 자원 경쟁, 네트워크 지연 부재 등으로 신뢰할 수 있는 수치를 제공하지 못합니다. 정확한 처리량/지연시간/한계치를 측정하려면 실 배포 환경에서 `hey`, `k6` 같은 전문 부하 테스트 도구를 별도 머신에서 실행해야 합니다.

```bash
hey -n 500 -c 50 -m POST \
  -H "Content-Type: application/json" \
  -d '{"type":"email:send","payload":"test@test.com","priority":"normal","max_attempts":3}' \
  http://<서버주소>:8080/jobs
```
