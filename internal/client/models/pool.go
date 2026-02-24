package models

// Pool is a resource pool that groups VMs and storage for access control purposes.
type Pool struct {
	PoolID  string       `json:"poolid"`
	Comment string       `json:"comment,omitempty"`
	Members []PoolMember `json:"members,omitempty"`
}

// PoolMember is a single resource (VM, container, or storage) that belongs to a pool.
type PoolMember struct {
	ID   string `json:"id"`
	Node string `json:"node"`
	Type string `json:"type"` // "qemu" or "lxc" or "storage"
	VMID int    `json:"vmid,omitempty"`
}

// PoolCreateRequest is sent when creating a new resource pool.
type PoolCreateRequest struct {
	PoolID  string `json:"poolid"`
	Comment string `json:"comment,omitempty"`
}

// PoolUpdateRequest is sent to add/remove members or update the comment on a pool.
type PoolUpdateRequest struct {
	Comment *string `json:"comment,omitempty"`
	VMs     string  `json:"vms,omitempty"`
	Storage string  `json:"storage,omitempty"`
	Delete  *int    `json:"delete,omitempty"`
}
