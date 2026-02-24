package models

// HAResource is a resource managed by the Proxmox HA subsystem.
type HAResource struct {
	SID        string `json:"sid"`
	Type       string `json:"type,omitempty"`
	State      string `json:"state,omitempty"`
	Group      string `json:"group,omitempty"`
	MaxRestart int    `json:"max_restart,omitempty"`
	MaxRelocate int   `json:"max_relocate,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

// HAResourceCreateRequest is sent when registering a new HA-managed resource.
type HAResourceCreateRequest struct {
	SID        string `json:"sid"`
	Type       string `json:"type,omitempty"`
	State      string `json:"state,omitempty"`
	Group      string `json:"group,omitempty"`
	MaxRestart int    `json:"max_restart,omitempty"`
	MaxRelocate int   `json:"max_relocate,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

// HAResourceUpdateRequest is sent to update settings on an existing HA resource.
type HAResourceUpdateRequest struct {
	State      string `json:"state,omitempty"`
	Group      string `json:"group,omitempty"`
	MaxRestart *int   `json:"max_restart,omitempty"`
	MaxRelocate *int  `json:"max_relocate,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

// HAGroup defines which nodes a HA resource can run on and their priority.
type HAGroup struct {
	Group     string `json:"group"`
	Nodes     string `json:"nodes"`
	Type      string `json:"type,omitempty"`
	Restricted *int  `json:"restricted,omitempty"`
	NoFailback *int  `json:"nofailback,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

// HAGroupCreateRequest is sent when creating a new HA group.
type HAGroupCreateRequest struct {
	Group      string `json:"group"`
	Nodes      string `json:"nodes"`
	Restricted *int   `json:"restricted,omitempty"`
	NoFailback *int   `json:"nofailback,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

// HAGroupUpdateRequest is sent to update an existing HA group's node list or settings.
type HAGroupUpdateRequest struct {
	Nodes      string `json:"nodes,omitempty"`
	Restricted *int   `json:"restricted,omitempty"`
	NoFailback *int   `json:"nofailback,omitempty"`
	Comment    string `json:"comment,omitempty"`
}
