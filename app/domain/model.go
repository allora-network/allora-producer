package domain

import "time"

const (
	StatusPending   string = "pending"
	StatusCompleted string = "completed"
	StatusFailed    string = "failed"
)

type ProcessedBlockTransactions struct {
	ID          int64
	Height      int64
	ProcessedAt time.Time
	Status      string
}

type ProcessedBlockEvents struct {
	ID          int64
	Height      int64
	ProcessedAt time.Time
	Status      string
}
