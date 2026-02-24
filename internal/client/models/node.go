package models

// NodeListEntry is a single node as returned by the cluster nodes list.
type NodeListEntry struct {
	Node           string  `json:"node"`
	Status         string  `json:"status"`
	CPU            float64 `json:"cpu"`
	MaxCPU         int64   `json:"maxcpu"`
	Mem            int64   `json:"mem"`
	MaxMem         int64   `json:"maxmem"`
	Disk           int64   `json:"disk"`
	MaxDisk        int64   `json:"maxdisk"`
	Uptime         int64   `json:"uptime"`
	Level          string  `json:"level"`
	ID             string  `json:"id"`
	Type           string  `json:"type"`
	SSLFingerprint string  `json:"ssl_fingerprint"`
}

// NodeStatus holds comprehensive status info for a single Proxmox node.
type NodeStatus struct {
	// Boot info
	BootInfo *BootInfo `json:"boot-info,omitempty"`

	// CPU
	CPUInfo  *CPUInfo `json:"cpuinfo,omitempty"`
	CPU      float64  `json:"cpu"`
	LoadAvg  []string `json:"loadavg,omitempty"`
	Wait     float64  `json:"wait"`

	// Memory
	Memory   *MemoryInfo `json:"memory,omitempty"`
	Swap     *MemoryInfo `json:"swap,omitempty"`

	// Storage
	RootFS   *StorageInfo `json:"rootfs,omitempty"`

	// System
	Uptime     int64   `json:"uptime"`
	Idle       float64 `json:"idle"`
	KVersion   string  `json:"kversion"`
	PVEVersion string  `json:"pveversion"`
}

// BootInfo contains info about how the node booted (BIOS/UEFI, secure boot, etc).
type BootInfo struct {
	Mode    string `json:"mode"`
	SecBoot int    `json:"secureboot"`
}

// CPUInfo describes the physical CPU hardware on the node.
type CPUInfo struct {
	Cores   int    `json:"cores"`
	CPUs    int    `json:"cpus"`
	MHz     string `json:"mhz"`
	Model   string `json:"model"`
	Sockets int    `json:"sockets"`
	Threads int    `json:"threads"`
	HVM     string `json:"hvm"`
	Flags   string `json:"flags"`
}

// MemoryInfo holds total/used/free values for memory or swap.
type MemoryInfo struct {
	Total int64 `json:"total"`
	Used  int64 `json:"used"`
	Free  int64 `json:"free"`
}

// StorageInfo holds disk space usage for a filesystem.
type StorageInfo struct {
	Total int64 `json:"total"`
	Used  int64 `json:"used"`
	Free  int64 `json:"free"`
	Avail int64 `json:"avail"`
}
