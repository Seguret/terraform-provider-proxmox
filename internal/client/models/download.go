package models

// DownloadURLRequest is sent to kick off a URL download to node storage.
type DownloadURLRequest struct {
	URL               string `json:"url"`
	Filename          string `json:"filename"`
	Content           string `json:"content"` // iso or vztmpl
	Checksum          string `json:"checksum,omitempty"`
	ChecksumAlgorithm string `json:"checksum-algorithm,omitempty"` // md5, sha1, sha224, sha256, sha384, sha512
	Verify            *int   `json:"verify-certificates,omitempty"`
}

// StorageContent is a single file entry in a storage volume (ISO, template, backup, etc).
type StorageContent struct {
	VolID   string `json:"volid"`
	Content string `json:"content"` // iso, vztmpl, backup, etc.
	Format  string `json:"format,omitempty"`
	Size    int64  `json:"size,omitempty"`
	Used    int64  `json:"used,omitempty"`
	CTime   int64  `json:"ctime,omitempty"`
	Notes   string `json:"notes,omitempty"`
}
