package models

// NodeConfig holds editable node-level settings like description and WoL config.
type NodeConfig struct {
	Description         string `json:"description,omitempty"`
	Wakeonlan           string `json:"wakeonlan,omitempty"`
	StartallOnbootDelay int    `json:"startall-onboot-delay,omitempty"`
}

// NodeConfigUpdateRequest is sent to update node settings.
type NodeConfigUpdateRequest struct {
	Description         string `json:"description,omitempty"`
	Wakeonlan           string `json:"wakeonlan,omitempty"`
	StartallOnbootDelay int    `json:"startall-onboot-delay,omitempty"`
	Delete              string `json:"delete,omitempty"`
}

// NodeSubscription holds the Proxmox VE subscription status and key info for a node.
type NodeSubscription struct {
	Status      string `json:"status"`
	ProductName string `json:"productname,omitempty"`
	RegDate     string `json:"regdate,omitempty"`
	NextDueDate string `json:"nextduedate,omitempty"`
	Key         string `json:"key,omitempty"`
	ServerID    string `json:"serverid,omitempty"`
	Sockets     int    `json:"sockets,omitempty"`
	Message     string `json:"message,omitempty"`
}
