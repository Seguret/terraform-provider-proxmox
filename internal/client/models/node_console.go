package models

// ConsoleProxyInfo holds the port info returned when opening a console proxy session.
type ConsoleProxyInfo struct {
	Port string `json:"port"`
}

// QEMUMonitorInfo is the response from the QEMU monitor endpoint.
type QEMUMonitorInfo struct {
	Prompt string `json:"prompt,omitempty"`
	Status string `json:"status,omitempty"`
}
