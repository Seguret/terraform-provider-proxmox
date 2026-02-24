package models

// ClusterStatusEntry is a single entry from the cluster status (node or quorum info).
type ClusterStatusEntry struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	IP      string `json:"ip,omitempty"`
	Online  int    `json:"online,omitempty"`
	Local   int    `json:"local,omitempty"`
	NodeID  int    `json:"nodeid,omitempty"`
	Version int    `json:"version,omitempty"`
	Quorate int    `json:"quorate,omitempty"`
	Nodes   int    `json:"nodes,omitempty"`
}

// ClusterResource is a single resource in the cluster-wide resource list (VM, container, storage, node).
type ClusterResource struct {
	ID         string  `json:"id"`
	Type       string  `json:"type"`
	Node       string  `json:"node,omitempty"`
	Status     string  `json:"status,omitempty"`
	Name       string  `json:"name,omitempty"`
	VMID       int     `json:"vmid,omitempty"`
	Pool       string  `json:"pool,omitempty"`
	CPU        float64 `json:"cpu,omitempty"`
	MaxCPU     int     `json:"maxcpu,omitempty"`
	Mem        int64   `json:"mem,omitempty"`
	MaxMem     int64   `json:"maxmem,omitempty"`
	Disk       int64   `json:"disk,omitempty"`
	MaxDisk    int64   `json:"maxdisk,omitempty"`
	Uptime     int64   `json:"uptime,omitempty"`
	Storage    string  `json:"storage,omitempty"`
	PluginType string  `json:"plugintype,omitempty"`
	Content    string  `json:"content,omitempty"`
}

// ClusterTask is a single task entry from the cluster task history.
type ClusterTask struct {
	UPID      string `json:"upid"`
	Node      string `json:"node"`
	PID       int    `json:"pid"`
	PStart    int    `json:"pstart"`
	StartTime int64  `json:"starttime"`
	EndTime   int64  `json:"endtime,omitempty"`
	Type      string `json:"type"`
	ID        string `json:"id,omitempty"`
	User      string `json:"user"`
	Status    string `json:"status,omitempty"`
}

// HAStatusEntry is the current HA state for a managed resource.
type HAStatusEntry struct {
	SID         string `json:"sid"`
	State       string `json:"state"`
	Node        string `json:"node,omitempty"`
	MaxRestart  int    `json:"max_restart,omitempty"`
	MaxRelocate int    `json:"max_relocate,omitempty"`
	CRMState    string `json:"crm_state,omitempty"`
	Request     string `json:"request,omitempty"`
}

// NodeHardwarePCI is a PCI device available on a node.
type NodeHardwarePCI struct {
	ID          string `json:"id"`
	Class       string `json:"class,omitempty"`
	Device      string `json:"device,omitempty"`
	Vendor      string `json:"vendor,omitempty"`
	DeviceID    string `json:"device_id,omitempty"`
	VendorID    string `json:"vendor_id,omitempty"`
	SubVendorID string `json:"subsystem_vendor_id,omitempty"`
	SubDeviceID string `json:"subsystem_device_id,omitempty"`
	IOMMUGroup  int    `json:"iommugroup,omitempty"`
	MDevTypes   string `json:"mdev_types,omitempty"`
}

// NodeHardwareUSB is a USB device connected to a node.
type NodeHardwareUSB struct {
	BusNum       int    `json:"busnum"`
	Class        int    `json:"class,omitempty"`
	DevNum       int    `json:"devnum"`
	Level        int    `json:"level,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Port         int    `json:"port,omitempty"`
	Product      string `json:"product,omitempty"`
	ProdID       string `json:"prodid,omitempty"`
	Speed        string `json:"speed,omitempty"`
	UsbPath      string `json:"usbpath,omitempty"`
	VendorID     string `json:"vendid,omitempty"`
	Serialnumber string `json:"serialnumber,omitempty"`
}

// NodeTask is a single task in the per-node task history.
type NodeTask struct {
	UPID      string `json:"upid"`
	Node      string `json:"node"`
	PID       int    `json:"pid"`
	PStart    int    `json:"pstart"`
	StartTime int64  `json:"starttime"`
	EndTime   int64  `json:"endtime,omitempty"`
	Type      string `json:"type"`
	ID        string `json:"id,omitempty"`
	User      string `json:"user"`
	Status    string `json:"status,omitempty"`
}

// APTPackageVersion is a single installed package with version info from the APT database.
type APTPackageVersion struct {
	Package     string `json:"Package"`
	Version     string `json:"Version,omitempty"`
	OldVersion  string `json:"OldVersion,omitempty"`
	Priority    string `json:"Priority,omitempty"`
	Section     string `json:"Section,omitempty"`
	Title       string `json:"Title,omitempty"`
	Description string `json:"Description,omitempty"`
}

// ACMEDirectory is a known ACME directory URL (e.g. Let's Encrypt prod/staging).
type ACMEDirectory struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// ClusterLogEntry is a single event from the cluster-wide syslog.
type ClusterLogEntry struct {
	PID      int64  `json:"pid"`
	UID      int64  `json:"uid"`
	GID      int64  `json:"gid"`
	Node     string `json:"node"`
	UserID   string `json:"user"`
	Tag      string `json:"tag"`
	Severity string `json:"pri"`
	Msg      string `json:"msg"`
	Time     int64  `json:"time"`
}

// ClusterBackupInfoEntry is a VM that isnt covered by any backup schedule.
type ClusterBackupInfoEntry struct {
	VMID int    `json:"vmid"`
	Type string `json:"type"`
	Name string `json:"name"`
}
