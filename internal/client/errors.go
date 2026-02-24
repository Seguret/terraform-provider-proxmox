package client

import "fmt"

// APIError wraps an HTTP error response from the Proxmox API.
type APIError struct {
	StatusCode int
	Status     string
	Message    string
	Errors     map[string]string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("proxmox API error %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("proxmox API error %d: %s", e.StatusCode, e.Status)
}

// IsNotFound checks if this is a 404 - usefull for deciding whether to remove from state.
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == 404
}

// TaskError is returned when an async task finishes with a non-OK exit status.
type TaskError struct {
	UPID       string
	Node       string
	ExitStatus string
}

func (e *TaskError) Error() string {
	return fmt.Sprintf("proxmox task %s failed: %s", e.UPID, e.ExitStatus)
}
