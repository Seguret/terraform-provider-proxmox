package models

// BackupJob is a scheduled vzdump backup job at the cluster level.
type BackupJob struct {
	ID               string  `json:"id"`
	Enabled          *int    `json:"enabled,omitempty"`
	Storage          string  `json:"storage,omitempty"`
	Schedule         string  `json:"schedule,omitempty"`
	VMIDs            string  `json:"vmid,omitempty"`
	Nodes            string  `json:"node,omitempty"`
	All              *int    `json:"all,omitempty"`
	Compress         string  `json:"compress,omitempty"`
	Mode             string  `json:"mode,omitempty"`
	Comment          string  `json:"comment,omitempty"`
	Mailto           string  `json:"mailto,omitempty"`
	MailNotification string  `json:"mailnotification,omitempty"`
	MaxFiles         int     `json:"maxfiles,omitempty"`
	Remove           *int    `json:"remove,omitempty"`
	NotesTemplate    string  `json:"notes-template,omitempty"`
	BWLimit          float64 `json:"bwlimit,omitempty"`
}

// BackupJobCreateRequest is sent when creating a new backup schedule.
type BackupJobCreateRequest struct {
	Storage          string  `json:"storage,omitempty"`
	Schedule         string  `json:"schedule,omitempty"`
	VMIDs            string  `json:"vmid,omitempty"`
	Nodes            string  `json:"node,omitempty"`
	All              *int    `json:"all,omitempty"`
	Compress         string  `json:"compress,omitempty"`
	Mode             string  `json:"mode,omitempty"`
	Comment          string  `json:"comment,omitempty"`
	Mailto           string  `json:"mailto,omitempty"`
	MailNotification string  `json:"mailnotification,omitempty"`
	MaxFiles         int     `json:"maxfiles,omitempty"`
	Remove           *int    `json:"remove,omitempty"`
	NotesTemplate    string  `json:"notes-template,omitempty"`
	BWLimit          float64 `json:"bwlimit,omitempty"`
}

// BackupJobUpdateRequest is sent to modify an existing backup job.
type BackupJobUpdateRequest struct {
	Storage          string  `json:"storage,omitempty"`
	Schedule         string  `json:"schedule,omitempty"`
	VMIDs            string  `json:"vmid,omitempty"`
	Nodes            string  `json:"node,omitempty"`
	All              *int    `json:"all,omitempty"`
	Compress         string  `json:"compress,omitempty"`
	Mode             string  `json:"mode,omitempty"`
	Comment          string  `json:"comment,omitempty"`
	Mailto           string  `json:"mailto,omitempty"`
	MailNotification string  `json:"mailnotification,omitempty"`
	MaxFiles         int     `json:"maxfiles,omitempty"`
	Remove           *int    `json:"remove,omitempty"`
	NotesTemplate    string  `json:"notes-template,omitempty"`
	BWLimit          float64 `json:"bwlimit,omitempty"`
	Enabled          *int    `json:"enabled,omitempty"`
}
