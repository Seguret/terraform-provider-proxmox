package models

// StorageConfig is the cluster-wide definition of a storage backend.
type StorageConfig struct {
	Storage    string `json:"storage"`
	Type       string `json:"type"`
	Content    string `json:"content,omitempty"`
	Disable    *int   `json:"disable,omitempty"`
	Shared     *int   `json:"shared,omitempty"`
	Path       string `json:"path,omitempty"`
	Pool       string `json:"pool,omitempty"`
	VGName     string `json:"vgname,omitempty"`
	Nodes      string `json:"nodes,omitempty"`
	Server     string `json:"server,omitempty"`
	Export     string `json:"export,omitempty"`
	Share      string `json:"share,omitempty"`
	Username   string `json:"username,omitempty"`
	Domain     string `json:"domain,omitempty"`
	Datastore  string `json:"datastore,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	PruneBackups string `json:"prune-backups,omitempty"`
	MaxFiles   int    `json:"maxfiles,omitempty"`
	Preallocation string `json:"preallocation,omitempty"`
}

// StorageCreateRequest is sent when adding a new storage backend to the cluster.
type StorageCreateRequest struct {
	Storage       string `json:"storage"`
	Type          string `json:"type"`
	Content       string `json:"content,omitempty"`
	Disable       *int   `json:"disable,omitempty"`
	Shared        *int   `json:"shared,omitempty"`
	Path          string `json:"path,omitempty"`
	Pool          string `json:"pool,omitempty"`
	VGName        string `json:"vgname,omitempty"`
	Nodes         string `json:"nodes,omitempty"`
	Server        string `json:"server,omitempty"`
	Export        string `json:"export,omitempty"`
	Share         string `json:"share,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Domain        string `json:"domain,omitempty"`
	Datastore     string `json:"datastore,omitempty"`
	Namespace     string `json:"namespace,omitempty"`
	Fingerprint   string `json:"fingerprint,omitempty"`
	PruneBackups  string `json:"prune-backups,omitempty"`
	Preallocation string `json:"preallocation,omitempty"`
}

// StorageUpdateRequest is sent to modify a storage backend's settings.
type StorageUpdateRequest struct {
	Content    *string `json:"content,omitempty"`
	Disable    *int    `json:"disable,omitempty"`
	Shared     *int    `json:"shared,omitempty"`
	Nodes      *string `json:"nodes,omitempty"`
	PruneBackups *string `json:"prune-backups,omitempty"`
}
