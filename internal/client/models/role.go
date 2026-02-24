package models

// Role is a named set of privileges used in access control rules.
type Role struct {
	RoleID string `json:"roleid"`
	Privs  string `json:"privs,omitempty"`
}

// RoleCreateRequest is sent when creating a new role.
type RoleCreateRequest struct {
	RoleID string `json:"roleid"`
	Privs  string `json:"privs,omitempty"`
}

// RoleUpdateRequest is sent to change which privileges a role has.
type RoleUpdateRequest struct {
	Privs string `json:"privs,omitempty"`
}
