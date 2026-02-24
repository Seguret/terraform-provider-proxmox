package models

// NodeHosts holds the raw content of /etc/hosts on a node.
type NodeHosts struct {
	Data   string `json:"data"`
	Digest string `json:"digest,omitempty"`
}

// NodeHostsUpdateRequest is sent to overwrite the /etc/hosts file on a node.
type NodeHostsUpdateRequest struct {
	Data   string `json:"data"`
	Digest string `json:"digest,omitempty"`
}
