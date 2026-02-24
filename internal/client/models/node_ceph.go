package models

// NodeCephStatus holds Ceph cluster health info as reported from a specific node.
type NodeCephStatus struct {
	Health map[string]interface{} `json:"health,omitempty"`
	Mgr    map[string]interface{} `json:"mgr_map,omitempty"`
	Mon    map[string]interface{} `json:"mon_map,omitempty"`
	Osd    map[string]interface{} `json:"osd_map,omitempty"`
	PGInfo map[string]interface{} `json:"pg_info,omitempty"`
}

// CephOSD is a single Object Storage Daemon in the Ceph cluster.
type CephOSD struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Weight float64 `json:"weight"`
}

// CephMON is a Ceph monitor daemon (responsible for cluster state).
type CephMON struct {
	Addr string `json:"addr"`
	Host string `json:"host"`
	Name string `json:"name"`
}

// CephConfig holds the raw ceph.conf content.
type CephConfig struct {
	Content string `json:"content,omitempty"`
}

// CephPool is a named Ceph storage pool with replication and autoscale settings.
type CephPool struct {
	Name            string                 `json:"pool_name"`
	PoolID          int                    `json:"pool"`
	Size            int                    `json:"size"`
	MinSize         int                    `json:"min_size"`
	PGNum           int                    `json:"pg_num"`
	PGAutoscaleMode string                 `json:"pg_autoscale_mode"`
	CrushRule       string                 `json:"crush_rule"`
	ApplicationMeta map[string]interface{} `json:"application_metadata"`
}

// CephPoolCreateRequest is sent when creating a new Ceph pool.
type CephPoolCreateRequest struct {
	Name            string `json:"name"`
	Size            *int   `json:"size,omitempty"`
	MinSize         *int   `json:"min_size,omitempty"`
	PGNum           *int   `json:"pg_num,omitempty"`
	PGAutoscaleMode string `json:"pg_autoscale_mode,omitempty"`
	CrushRule       string `json:"crush_rule,omitempty"`
	Application     string `json:"application,omitempty"`
	AddStorages     *int   `json:"add_storages,omitempty"`
}

// CephPoolUpdateRequest is sent to modify an existing Ceph pool's settings.
type CephPoolUpdateRequest struct {
	Size            *int   `json:"size,omitempty"`
	MinSize         *int   `json:"min_size,omitempty"`
	PGNum           *int   `json:"pg_num,omitempty"`
	PGAutoscaleMode string `json:"pg_autoscale_mode,omitempty"`
	CrushRule       string `json:"crush_rule,omitempty"`
	Application     string `json:"application,omitempty"`
}

// CephOSDCreateRequest is sent when adding a new OSD to the Ceph cluster.
type CephOSDCreateRequest struct {
	Dev       string `json:"dev"`
	Encrypted *int   `json:"encrypted,omitempty"`
	DBDev     string `json:"db_dev,omitempty"`
	WALDev    string `json:"wal_dev,omitempty"`
}

// CephOSDListResponse is the top-level response when listing OSDs (contains a tree structure).
type CephOSDListResponse struct {
	Nodes []CephOSDTreeNode `json:"nodes"`
}

// CephOSDTreeNode is an entry in the CRUSH tree (host, OSD, root, etc).
type CephOSDTreeNode struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Status string  `json:"status,omitempty"`
	Weight float64 `json:"crush_weight,omitempty"`
}

// CephMDS is a Ceph Metadata Server daemon (required for CephFS).
type CephMDS struct {
	Name  string `json:"name"`
	State string `json:"state"`
	Addr  string `json:"addr,omitempty"`
}

// CephMGR is a Ceph Manager daemon that handles metrics and orchestration.
type CephMGR struct {
	ID    string `json:"id"`
	State string `json:"state"`
	Addr  string `json:"addr,omitempty"`
}

// CephMGRListResponse holds the active and standby MGR daemons.
type CephMGRListResponse struct {
	Active   *CephMGR  `json:"active"`
	Standbys []CephMGR `json:"standbys"`
}

// CephFS is a CephFS filesystem with its metadata and data pools.
type CephFS struct {
	Name         string `json:"name"`
	MetadataPool string `json:"metadata_pool"`
	DataPool     string `json:"data_pool"`
}
