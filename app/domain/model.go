package domain

import "time"

const (
	StatusPending   string = "pending"
	StatusCompleted string = "completed"
	StatusFailed    string = "failed"
)

// ProcessedBlock is the model used to keep track of the processed blocks
type ProcessedBlock struct {
	ID          int64
	Height      int64
	ProcessedAt time.Time
	Status      string
}

// ProcessedBlockEvent is the model used to keep track of the processed block events
type ProcessedBlockEvent struct {
	ID          int64
	Height      int64
	ProcessedAt time.Time
	Status      string
}
