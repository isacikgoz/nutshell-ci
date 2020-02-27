package circleci

import (
	"time"
)

// Action is the definition and outcome of a Step
type Action struct {
	Name         string    `json:"name"`
	Failed       bool      `json:"failed"`
	Status       string    `json:"status"`
	ExitCode     int       `json:"exit_code"`
	StartTime    time.Time `json:"start_time,omitempty"`
	EndTime      time.Time `json:"end_time,omitempty"`
	AllocationID string    `json:"allocation_id"`
	Step         int       `json:"step"`
}

// Duration returns the duration of the action
func (a *Action) Duration() time.Duration {
	return a.EndTime.Sub(a.StartTime)
}
