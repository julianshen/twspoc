package data

import "time"

// Job represents an activity or unit of work in the system
type Job struct {
	JobID       string                 `json:"job_id"`
	Namespace   string                 `json:"namespace"`
	Status      string                 `json:"status"` // e.g., "started", "completed", "failed"
	Inputs      map[string]interface{} `json:"inputs,omitempty"`
	Result      map[string]interface{} `json:"result,omitempty"`
	CreatedAt   time.Time              `json:"created_at,omitempty"`
	StartedAt   time.Time              `json:"started_at,omitempty"`
	CompletedAt time.Time              `json:"completed_at,omitempty"`
	Retries     int                    `json:"retries,omitempty"`
	DependsOn   []string               `json:"depends_on,omitempty"`   // Array of task IDs for DAG-style orchestration
	TriggeredBy string                 `json:"triggered_by,omitempty"` // Event ID that caused the task to start
}
