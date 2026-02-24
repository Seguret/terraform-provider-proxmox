package models

// ReplicationJob is a ZFS replication job that syncs a VM to another cluster node.
type ReplicationJob struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"` // local
	Target   string  `json:"target"`
	Source   string  `json:"source,omitempty"`
	Guest    int     `json:"guest"`
	Schedule string  `json:"schedule,omitempty"`
	Rate     float64 `json:"rate,omitempty"`
	Comment  string  `json:"comment,omitempty"`
	Disable  *int    `json:"disable,omitempty"`
	Remove   *int    `json:"remove_job,omitempty"`
}

// ReplicationJobCreateRequest is sent when setting up a new replication job.
type ReplicationJobCreateRequest struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"` // local
	Target   string  `json:"target"`
	Guest    int     `json:"guest,omitempty"`
	Schedule string  `json:"schedule,omitempty"`
	Rate     float64 `json:"rate,omitempty"`
	Comment  string  `json:"comment,omitempty"`
	Disable  *int    `json:"disable,omitempty"`
}

// ReplicationJobUpdateRequest is sent to change schedule, rate or other settings on a job.
type ReplicationJobUpdateRequest struct {
	Schedule string  `json:"schedule,omitempty"`
	Rate     float64 `json:"rate,omitempty"`
	Comment  string  `json:"comment,omitempty"`
	Disable  *int    `json:"disable,omitempty"`
}
