package models

// ContainerListEntry is a single LXC container entry from the node container list.
type ContainerListEntry struct {
	VMID    int     `json:"vmid"`
	Name    string  `json:"name"`
	Status  string  `json:"status"`
	CPU     float64 `json:"cpu"`
	CPUs    int     `json:"cpus"`
	Mem     int64   `json:"mem"`
	MaxMem  int64   `json:"maxmem"`
	Disk    int64   `json:"disk"`
	MaxDisk int64   `json:"maxdisk"`
	Uptime  int64   `json:"uptime"`
	Tags    string  `json:"tags"`
	Lock    string  `json:"lock"`
	Type    string  `json:"type"`
}

// ContainerConfig holds the full configuration of an LXC container.
type ContainerConfig struct {
	Hostname    string `json:"hostname,omitempty"`
	Description string `json:"description,omitempty"`
	Tags        string `json:"tags,omitempty"`
	OSType      string `json:"ostype,omitempty"`
	OnBoot      *int   `json:"onboot,omitempty"`
	Protection  *int   `json:"protection,omitempty"`
	Unprivileged *int  `json:"unprivileged,omitempty"`
	Template    *int   `json:"template,omitempty"`

	// Networking
	Arch string `json:"arch,omitempty"`

	// CPU
	Cores      int    `json:"cores,omitempty"`
	CPULimit   int    `json:"cpulimit,omitempty"`
	CPUUnits   int    `json:"cpuunits,omitempty"`

	// Memory
	Memory  int `json:"memory,omitempty"`
	Swap    int `json:"swap,omitempty"`

	// Root filesystem
	RootFS  string `json:"rootfs,omitempty"`

	// Network interfaces (dynamic: net0, net1, ...)
	Net0 string `json:"net0,omitempty"`
	Net1 string `json:"net1,omitempty"`
	Net2 string `json:"net2,omitempty"`
	Net3 string `json:"net3,omitempty"`

	// Mount points (dynamic: mp0, mp1, ...)
	MP0 string `json:"mp0,omitempty"`
	MP1 string `json:"mp1,omitempty"`
	MP2 string `json:"mp2,omitempty"`

	// OS template
	OSTemplate string `json:"ostemplate,omitempty"`

	// DNS
	Nameserver   string `json:"nameserver,omitempty"`
	Searchdomain string `json:"searchdomain,omitempty"`

	// SSH Keys
	SSHKeys string `json:"ssh-public-keys,omitempty"`

	// Password
	Password string `json:"password,omitempty"`

	// Console
	Console  *int   `json:"console,omitempty"`
	TTY      int    `json:"tty,omitempty"`

	// Features
	Features string `json:"features,omitempty"`

	// Startup
	Startup string `json:"startup,omitempty"`

	// Pool
	Pool string `json:"pool,omitempty"`
}

// ContainerCreateRequest is sent when creating a new LXC container.
type ContainerCreateRequest struct {
	VMID         int    `json:"vmid,omitempty"`
	Hostname     string `json:"hostname,omitempty"`
	Description  string `json:"description,omitempty"`
	Tags         string `json:"tags,omitempty"`
	OSTemplate   string `json:"ostemplate"`
	OSType       string `json:"ostype,omitempty"`
	OnBoot       *int   `json:"onboot,omitempty"`
	Protection   *int   `json:"protection,omitempty"`
	Unprivileged *int   `json:"unprivileged,omitempty"`
	Start        *int   `json:"start,omitempty"`
	Pool         string `json:"pool,omitempty"`

	// CPU
	Cores    int `json:"cores,omitempty"`
	CPULimit int `json:"cpulimit,omitempty"`
	CPUUnits int `json:"cpuunits,omitempty"`

	// Memory
	Memory int `json:"memory,omitempty"`
	Swap   int `json:"swap,omitempty"`

	// Root filesystem
	RootFS string `json:"rootfs,omitempty"`

	// Network
	Net0 string `json:"net0,omitempty"`
	Net1 string `json:"net1,omitempty"`
	Net2 string `json:"net2,omitempty"`
	Net3 string `json:"net3,omitempty"`

	// Mount points
	MP0 string `json:"mp0,omitempty"`
	MP1 string `json:"mp1,omitempty"`
	MP2 string `json:"mp2,omitempty"`

	// DNS
	Nameserver   string `json:"nameserver,omitempty"`
	Searchdomain string `json:"searchdomain,omitempty"`

	// Auth
	Password string `json:"password,omitempty"`
	SSHKeys  string `json:"ssh-public-keys,omitempty"`

	// Console
	Console *int `json:"console,omitempty"`
	TTY     int  `json:"tty,omitempty"`

	// Features
	Features string `json:"features,omitempty"`

	// Clone
	Clone       *int   `json:"clone,omitempty"`
	CloneTarget string `json:"target,omitempty"`
	Full        *int   `json:"full,omitempty"`
}

// ContainerStatus holds the current runtime state of a container.
type ContainerStatus struct {
	Status  string  `json:"status"`
	VMID    int     `json:"vmid"`
	Name    string  `json:"name"`
	CPU     float64 `json:"cpu"`
	CPUs    int     `json:"cpus"`
	Mem     int64   `json:"mem"`
	MaxMem  int64   `json:"maxmem"`
	Disk    int64   `json:"disk"`
	MaxDisk int64   `json:"maxdisk"`
	Uptime  int64   `json:"uptime"`
	Lock    string  `json:"lock"`
}
