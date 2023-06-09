package store

import "time"

type (
	Status string
	Type   string
)

const (
	StatusWaitingToStart Status = "Waiting to start"
	StatusJobQueued      Status = "Job queued"
	StatusProcessing     Status = "Processing"
	StatusComplete       Status = "Complete"
	StatusUnknown        Status = "Unknown"

	TypeImagine Type = "Imagine"
	TypeUpscale Type = "Upscale"
	TypeUnknown Type = "Unknown"

	Expired time.Duration = 3 * time.Hour
)
