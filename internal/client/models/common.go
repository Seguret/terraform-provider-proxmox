package models

// APIResponse is the generic envelope that Proxmox wraps all API responses in.
// The actual payload lives in the Data field.
type APIResponse[T any] struct {
	Data T `json:"data"`
}

// TaskStatus holds the current state of an async task (running, stopped, etc).
type TaskStatus struct {
	Status     string `json:"status"`     // "running" or "stopped"
	ExitStatus string `json:"exitstatus"` // "OK" or error message (only when stopped)
	Type       string `json:"type"`
	ID         string `json:"id"`
	Node       string `json:"node"`
	PID        int    `json:"pid"`
	StartTime  int64  `json:"starttime"`
	UPID       string `json:"upid"`
}

// Version is the response from the version endpoint.
type Version struct {
	Release string `json:"release"`
	RepoID  string `json:"repoid"`
	Version string `json:"version"`
}
