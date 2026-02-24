package models

// Group is a Proxmox VE user group for managing access permissions.
type Group struct {
	GroupID string   `json:"groupid"`
	Comment string   `json:"comment,omitempty"`
	Members []string `json:"members,omitempty"`
	Users   []string `json:"users,omitempty"`
}

// GroupCreateRequest is sent when creating a new user group.
type GroupCreateRequest struct {
	GroupID string `json:"groupid"`
	Comment string `json:"comment,omitempty"`
}

// GroupUpdateRequest is sent to update a group's comment.
type GroupUpdateRequest struct {
	Comment *string `json:"comment,omitempty"`
}
