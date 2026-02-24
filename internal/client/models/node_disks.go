package models

// NodeDisk represents a physical disk on a node as returned by the disk list API.
type NodeDisk struct {
	Dev     string `json:"dev"`
	DevPath string `json:"devpath,omitempty"`
	Type    string `json:"type,omitempty"`
	Size    int64  `json:"size,omitempty"`
	Model   string `json:"model,omitempty"`
	Serial  string `json:"serial,omitempty"`
	Vendor  string `json:"vendor,omitempty"`
	WWN     string `json:"wwn,omitempty"`
	Health  string `json:"health,omitempty"`
	RPM     int    `json:"rpm,omitempty"`
	Used    string `json:"used,omitempty"`
	GPT     *int   `json:"gpt,omitempty"`
	Wearout int    `json:"wearout,omitempty"`
}

// NodeDiskSmart holds the SMART health data for a disk.
type NodeDiskSmart struct {
	Type       string                  `json:"type,omitempty"`
	Health     string                  `json:"health,omitempty"`
	Wearout    int                     `json:"wearout,omitempty"`
	Attributes []NodeDiskSmartAttribute `json:"attributes,omitempty"`
	Text       string                  `json:"text,omitempty"`
}

// NodeDiskSmartAttribute is a single SMART attribute value for a disk.
type NodeDiskSmartAttribute struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Value  int    `json:"value,omitempty"`
	Worst  int    `json:"worst,omitempty"`
	Thresh int    `json:"thresh,omitempty"`
	Raw    string `json:"raw,omitempty"`
	Flags  string `json:"flags,omitempty"`
}

// NodeDiskDirectoryCreateRequest is sent when initializing a disk as a directory storage.
type NodeDiskDirectoryCreateRequest struct {
	Dev        string `json:"dev"`
	Filesystem string `json:"filesystem,omitempty"`
	Name       string `json:"name"`
	AddStorage bool   `json:"add_storage,omitempty"`
}

// NodeDiskLVMCreateRequest is sent when creating an LVM volume group on a disk.
type NodeDiskLVMCreateRequest struct {
	Dev        string `json:"dev"`
	Name       string `json:"name"`
	AddStorage bool   `json:"add_storage,omitempty"`
}

// NodeDiskLVMThinCreateRequest is sent when creating an LVM-thin pool on a volume group.
type NodeDiskLVMThinCreateRequest struct {
	Dev         string `json:"dev"`
	Name        string `json:"name"`
	VolumeGroup string `json:"thinpool,omitempty"`
	AddStorage  bool   `json:"add_storage,omitempty"`
}

// NodeDiskZFSCreateRequest is sent when creating a new ZFS pool on a node.
type NodeDiskZFSCreateRequest struct {
	Devices     string `json:"devices"`
	Name        string `json:"name"`
	RAIDLevel   string `json:"raidlevel,omitempty"`
	Ashift      int    `json:"ashift,omitempty"`
	Compression string `json:"compression,omitempty"`
	AddStorage  bool   `json:"add_storage,omitempty"`
}
