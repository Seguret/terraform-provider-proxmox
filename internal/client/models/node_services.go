package models

// NodeService is a systemd service running on a Proxmox node.
type NodeService struct {
	Name          string `json:"name"`
	State         string `json:"state"`
	Desc          string `json:"desc,omitempty"`
	ActiveState   string `json:"active-state,omitempty"`
	SubState      string `json:"substate,omitempty"`
	UnitFileState string `json:"unit-file-state,omitempty"`
}
