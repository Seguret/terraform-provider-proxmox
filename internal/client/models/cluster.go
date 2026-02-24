package models

// ClusterOptions holds cluster-wide settings like keyboard layout, proxy, and migration config.
type ClusterOptions struct {
	Keyboard        string `json:"keyboard,omitempty"`
	Language        string `json:"language,omitempty"`
	EmailFrom       string `json:"email_from,omitempty"`
	HTTPProxy       string `json:"http_proxy,omitempty"`
	MaxWorkers      *int   `json:"max_workers,omitempty"`
	MigrationUnsecure *int `json:"migration_unsecure,omitempty"`
	MigrationType   string `json:"migration_type,omitempty"`
	HAShutdownPolicy string `json:"ha_shutdown_policy,omitempty"`
	CRS             string `json:"crs,omitempty"`
	Notify          string `json:"notify,omitempty"`
	TagStyle        string `json:"tag-style,omitempty"`
	RegisteredTags  string `json:"registered-tags,omitempty"`
}

// ClusterOptionsUpdateRequest is sent to change cluster-wide options.
type ClusterOptionsUpdateRequest struct {
	Keyboard        string `json:"keyboard,omitempty"`
	Language        string `json:"language,omitempty"`
	EmailFrom       string `json:"email_from,omitempty"`
	HTTPProxy       string `json:"http_proxy,omitempty"`
	MaxWorkers      *int   `json:"max_workers,omitempty"`
	MigrationUnsecure *int `json:"migration_unsecure,omitempty"`
	MigrationType   string `json:"migration_type,omitempty"`
	HAShutdownPolicy string `json:"ha_shutdown_policy,omitempty"`
}

// ClusterConfig holds the corosync cluster topology information.
type ClusterConfig struct {
	Nodes []ClusterNode `json:"nodes,omitempty"`
	TotemInterface string `json:"totem_interface,omitempty"`
}

// ClusterNode is a single member node in the corosync cluster config.
type ClusterNode struct {
	Name   string `json:"name"`
	NodeID int    `json:"nodeid"`
	IP     string `json:"ip,omitempty"`
	Ring0Addr string `json:"ring0_addr,omitempty"`
	Ring1Addr string `json:"ring1_addr,omitempty"`
}

// HardwareMappingPCI is a cluster-level PCI device mapping entry.
type HardwareMappingPCI struct {
	ID          string   `json:"id"`
	Description string   `json:"description,omitempty"`
	Map         []string `json:"map"`
}

// HardwareMappingUSB is a cluster-level USB device mapping entry.
type HardwareMappingUSB struct {
	ID          string   `json:"id"`
	Description string   `json:"description,omitempty"`
	Map         []string `json:"map"`
}

// CephStatus holds the top-level health and map info for the Ceph cluster.
type CephStatus struct {
	Health        CephHealth            `json:"health,omitempty"`
	OSDMap        map[string]interface{} `json:"osdmap,omitempty"`
	PGMap         map[string]interface{} `json:"pgmap,omitempty"`
	MonMap        map[string]interface{} `json:"monmap,omitempty"`
}

// CephHealth holds the health status string and any active health checks.
type CephHealth struct {
	Status string `json:"status"`
	Checks map[string]interface{} `json:"checks,omitempty"`
}
