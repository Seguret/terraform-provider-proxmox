package models

import (
	"encoding/json"
	"strconv"
)

// VMListEntry is a single VM as it appears in the node VM list (summary info only).
type VMListEntry struct {
	VMID      int     `json:"vmid"`
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	CPU       float64 `json:"cpu"`
	CPUs      int     `json:"cpus"`
	Mem       int64   `json:"mem"`
	MaxMem    int64   `json:"maxmem"`
	Disk      int64   `json:"disk"`
	MaxDisk   int64   `json:"maxdisk"`
	Uptime    int64   `json:"uptime"`
	Template  int     `json:"template"`
	QMPStatus string  `json:"qmpstatus"`
	PID       int     `json:"pid"`
	Tags      string  `json:"tags"`
	Lock      string  `json:"lock"`
}

// VMConfig holds the full configuration of a QEMU VM.
type VMConfig struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Tags        string `json:"tags,omitempty"`
	OnBoot      *int   `json:"onboot,omitempty"`
	Protection  *int   `json:"protection,omitempty"`
	Template    *int   `json:"template,omitempty"`
	Agent       string `json:"agent,omitempty"`
	OSType      string `json:"ostype,omitempty"`
	BIOS        string `json:"bios,omitempty"`
	Machine     string `json:"machine,omitempty"`
	Boot        string `json:"boot,omitempty"`

	// CPU
	Sockets int    `json:"sockets,omitempty"`
	Cores   int    `json:"cores,omitempty"`
	CPUType string `json:"cpu,omitempty"`
	NUMA    *int   `json:"numa,omitempty"`
	VCPUs   int    `json:"vcpus,omitempty"`

	// Memory
	Memory  int    `json:"memory,omitempty"`
	Balloon int    `json:"balloon,omitempty"`

	// Disks - returned as dynamic keys by proxmox: scsi0, virtio0, ide0, sata0, etc.
	SCSI0   string `json:"scsi0,omitempty"`
	SCSI1   string `json:"scsi1,omitempty"`
	SCSI2   string `json:"scsi2,omitempty"`
	SCSI3   string `json:"scsi3,omitempty"`
	VirtIO0 string `json:"virtio0,omitempty"`
	VirtIO1 string `json:"virtio1,omitempty"`
	VirtIO2 string `json:"virtio2,omitempty"`
	VirtIO3 string `json:"virtio3,omitempty"`
	IDE0    string `json:"ide0,omitempty"`
	IDE1    string `json:"ide1,omitempty"`
	IDE2    string `json:"ide2,omitempty"`
	EFIDisk0 string `json:"efidisk0,omitempty"`
	TPMState0 string `json:"tpmstate0,omitempty"`

	// Network - similar dynamic key pattern: net0, net1, etc.
	Net0 string `json:"net0,omitempty"`
	Net1 string `json:"net1,omitempty"`
	Net2 string `json:"net2,omitempty"`
	Net3 string `json:"net3,omitempty"`

	// Cloud-init
	CIUser     string `json:"ciuser,omitempty"`
	CIPassword string `json:"cipassword,omitempty"`
	CIType     string `json:"citype,omitempty"`
	IPConfig0  string `json:"ipconfig0,omitempty"`
	IPConfig1  string `json:"ipconfig1,omitempty"`
	IPConfig2  string `json:"ipconfig2,omitempty"`
	IPConfig3  string `json:"ipconfig3,omitempty"`
	Nameserver string `json:"nameserver,omitempty"`
	Searchdomain string `json:"searchdomain,omitempty"`
	SSHKeys    string `json:"sshkeys,omitempty"`

	// VGA
	VGA string `json:"vga,omitempty"`

	// Serial
	Serial0 string `json:"serial0,omitempty"`

	// SCSI controller
	SCSIHw string `json:"scsihw,omitempty"`

	// Raw config map - used to capture any dynamic keys we dont explicitly model
	RawConfig map[string]json.RawMessage `json:"-"`
}

// UnmarshalJSON handles the fact that Proxmox sometimes returns Memory
// as an int and sometimes as a string, depending on the version.
func (vc *VMConfig) UnmarshalJSON(data []byte) error {
	type Alias VMConfig
	aux := &struct {
		Memory  interface{} `json:"memory,omitempty"`
		Balloon interface{} `json:"balloon,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(vc),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// memory can come back as int or string
	if aux.Memory != nil {
		switch v := aux.Memory.(type) {
		case float64:
			vc.Memory = int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				vc.Memory = i
			}
		}
	}

	// same deal for balloon
	if aux.Balloon != nil {
		switch v := aux.Balloon.(type) {
		case float64:
			vc.Balloon = int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				vc.Balloon = i
			}
		}
	}

	return nil
}

// VMCreateRequest is sent when provisioning a new QEMU VM.
type VMCreateRequest struct {
	VMID        int    `json:"vmid,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Tags        string `json:"tags,omitempty"`
	OnBoot      *int   `json:"onboot,omitempty"`
	Protection  *int   `json:"protection,omitempty"`
	Agent       string `json:"agent,omitempty"`
	OSType      string `json:"ostype,omitempty"`
	BIOS        string `json:"bios,omitempty"`
	Machine     string `json:"machine,omitempty"`
	Boot        string `json:"boot,omitempty"`
	Start       *int   `json:"start,omitempty"`
	Pool        string `json:"pool,omitempty"`

	// CPU
	Sockets int    `json:"sockets,omitempty"`
	Cores   int    `json:"cores,omitempty"`
	CPUType string `json:"cpu,omitempty"`
	NUMA    *int   `json:"numa,omitempty"`

	// Memory
	Memory  int `json:"memory,omitempty"`
	Balloon int `json:"balloon,omitempty"`

	// SCSI controller
	SCSIHw string `json:"scsihw,omitempty"`

	// VGA
	VGA string `json:"vga,omitempty"`

	// Serial
	Serial0 string `json:"serial0,omitempty"`

	// Disks (dynamic)
	SCSI0   string `json:"scsi0,omitempty"`
	SCSI1   string `json:"scsi1,omitempty"`
	SCSI2   string `json:"scsi2,omitempty"`
	SCSI3   string `json:"scsi3,omitempty"`
	VirtIO0 string `json:"virtio0,omitempty"`
	VirtIO1 string `json:"virtio1,omitempty"`
	VirtIO2 string `json:"virtio2,omitempty"`
	VirtIO3 string `json:"virtio3,omitempty"`
	IDE0    string `json:"ide0,omitempty"`
	IDE1    string `json:"ide1,omitempty"`
	IDE2    string `json:"ide2,omitempty"`
	SATA0   string `json:"sata0,omitempty"`
	SATA1   string `json:"sata1,omitempty"`
	SATA2   string `json:"sata2,omitempty"`
	EFIDisk0 string `json:"efidisk0,omitempty"`
	TPMState0 string `json:"tpmstate0,omitempty"`

	// Network
	Net0 string `json:"net0,omitempty"`
	Net1 string `json:"net1,omitempty"`
	Net2 string `json:"net2,omitempty"`
	Net3 string `json:"net3,omitempty"`

	// Cloud-init
	CIUser       string `json:"ciuser,omitempty"`
	CIPassword   string `json:"cipassword,omitempty"`
	CIType       string `json:"citype,omitempty"`
	IPConfig0    string `json:"ipconfig0,omitempty"`
	IPConfig1    string `json:"ipconfig1,omitempty"`
	IPConfig2    string `json:"ipconfig2,omitempty"`
	IPConfig3    string `json:"ipconfig3,omitempty"`
	Nameserver   string `json:"nameserver,omitempty"`
	Searchdomain string `json:"searchdomain,omitempty"`
	SSHKeys      string `json:"sshkeys,omitempty"`
}

// VMCloneRequest is sent when cloning an existing VM.
type VMCloneRequest struct {
	NewID       int    `json:"newid"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Pool        string `json:"pool,omitempty"`
	Target      string `json:"target,omitempty"`
	Full        *int   `json:"full,omitempty"`
	Storage     string `json:"storage,omitempty"`
	Format      string `json:"format,omitempty"`
}

// VMResizeRequest is sent when resizing a VM disk.
type VMResizeRequest struct {
	Disk    string `json:"disk"`
	Size    string `json:"size"`
	Digest  string `json:"digest,omitempty"`
	SkipLock *int  `json:"skiplock,omitempty"`
}

// VMStatus holds the current runtime state of a VM.
type VMStatus struct {
	Status    string  `json:"status"`
	VMID      int     `json:"vmid"`
	Name      string  `json:"name"`
	QMPStatus string  `json:"qmpstatus"`
	CPU       float64 `json:"cpu"`
	CPUs      int     `json:"cpus"`
	Mem       int64   `json:"mem"`
	MaxMem    int64   `json:"maxmem"`
	Disk      int64   `json:"disk"`
	MaxDisk   int64   `json:"maxdisk"`
	Uptime    int64   `json:"uptime"`
	PID       int     `json:"pid"`
	Lock      string  `json:"lock"`
}

// VMProxyInfo holds the connection info returned when setting up a VNC, SPICE or terminal proxy.
type VMProxyInfo struct {
	Cert      string `json:"cert,omitempty"`
	Port      int    `json:"port"`
	Ticket    string `json:"ticket"`
	Upid      string `json:"upid,omitempty"`
	User      string `json:"user"`
}

// VMFeatureInfo indicates whether a feature is available for a VM and on which nodes.
type VMFeatureInfo struct {
	HasFeature bool   `json:"hasFeature"`
	Nodes      []string `json:"nodes,omitempty"`
}

