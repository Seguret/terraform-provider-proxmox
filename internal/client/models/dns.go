package models

// NodeDNS holds the DNS resolver configuration for a node.
type NodeDNS struct {
	Search  string `json:"search,omitempty"`
	DNS1    string `json:"dns1,omitempty"`
	DNS2    string `json:"dns2,omitempty"`
	DNS3    string `json:"dns3,omitempty"`
}

// NodeDNSUpdateRequest is sent to change the DNS settings on a node.
type NodeDNSUpdateRequest struct {
	Search string `json:"search"`
	DNS1   string `json:"dns1,omitempty"`
	DNS2   string `json:"dns2,omitempty"`
	DNS3   string `json:"dns3,omitempty"`
}
