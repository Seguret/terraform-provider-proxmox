package models

// StorageListEntry is one storage backend as seen from a specific node (includes usage stats).
type StorageListEntry struct {
	Storage    string  `json:"storage"`
	Type       string  `json:"type"`
	Content    string  `json:"content"`
	Status     string  `json:"status"`
	Active     int     `json:"active"`
	Enabled    int     `json:"enabled"`
	Shared     int     `json:"shared"`
	Total      int64   `json:"total"`
	Used       int64   `json:"used"`
	Avail      int64   `json:"avail"`
	UsedFrac   float64 `json:"used_fraction"`
}
