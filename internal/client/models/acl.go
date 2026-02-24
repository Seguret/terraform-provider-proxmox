package models

// ACLEntry is a single permission entry from the access control list.
type ACLEntry struct {
	Path      string `json:"path"`
	RoleID    string `json:"roleid"`
	Type      string `json:"type"`      // "user" or "group"
	UGid      string `json:"ugid"`      // user or group id
	Propagate *int   `json:"propagate"` // 1 or 0
}

// ACLUpdateRequest is the body for adding or removing ACL entries.
type ACLUpdateRequest struct {
	Path      string `json:"path"`
	Roles     string `json:"roles"`
	Users     string `json:"users,omitempty"`
	Groups    string `json:"groups,omitempty"`
	Propagate *int   `json:"propagate,omitempty"`
	Delete    *int   `json:"delete,omitempty"`
}
