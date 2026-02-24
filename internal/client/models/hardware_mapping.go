package models

// PCIHardwareMapping is a named mapping of a PCI device across cluster nodes.
type PCIHardwareMapping struct {
	ID      string   `json:"id"`
	Comment string   `json:"description,omitempty"`
	Map     []string `json:"map,omitempty"`
	MDevs   string   `json:"mdev,omitempty"`
}

// PCIHardwareMappingCreateRequest is sent when creating a new PCI device mapping.
type PCIHardwareMappingCreateRequest struct {
	ID      string   `json:"id"`
	Comment string   `json:"description,omitempty"`
	Map     []string `json:"map"`
	MDevs   string   `json:"mdev,omitempty"`
}

// PCIHardwareMappingUpdateRequest is sent to update an existing PCI device mapping.
type PCIHardwareMappingUpdateRequest struct {
	Comment string   `json:"description,omitempty"`
	Map     []string `json:"map,omitempty"`
	MDevs   string   `json:"mdev,omitempty"`
}

// USBHardwareMapping is a named mapping of a USB device across cluster nodes.
type USBHardwareMapping struct {
	ID      string   `json:"id"`
	Comment string   `json:"description,omitempty"`
	Map     []string `json:"map,omitempty"`
}

// USBHardwareMappingCreateRequest is sent when creating a new USB device mapping.
type USBHardwareMappingCreateRequest struct {
	ID      string   `json:"id"`
	Comment string   `json:"description,omitempty"`
	Map     []string `json:"map"`
}

// USBHardwareMappingUpdateRequest is sent to update an existing USB device mapping.
type USBHardwareMappingUpdateRequest struct {
	Comment string   `json:"description,omitempty"`
	Map     []string `json:"map,omitempty"`
}

// DirHardwareMapping is a named directory path mapping, allowing per-node path overrides.
type DirHardwareMapping struct {
	ID      string                    `json:"id"`
	Comment string                    `json:"description,omitempty"`
	Map     []DirHardwareMappingEntry `json:"map,omitempty"`
}

// DirHardwareMappingEntry maps a single node name to a filesystem path.
type DirHardwareMappingEntry struct {
	Node string `json:"node"`
	Path string `json:"path"`
}

// DirHardwareMappingCreateRequest is sent when creating a new directory mapping.
type DirHardwareMappingCreateRequest struct {
	ID      string   `json:"id"`
	Comment string   `json:"description,omitempty"`
	Map     []string `json:"map"`
}

// DirHardwareMappingUpdateRequest is sent to update an existing directory mapping.
type DirHardwareMappingUpdateRequest struct {
	Comment string   `json:"description,omitempty"`
	Map     []string `json:"map,omitempty"`
}
