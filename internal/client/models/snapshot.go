package models

// VMSnapshot is a single snapshot entry for a QEMU VM.
type VMSnapshot struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Snaptime    int64   `json:"snaptime,omitempty"`
	VMSTATE     *int    `json:"vmstate,omitempty"`
	Parent      string  `json:"parent,omitempty"`
	Running     *int    `json:"running,omitempty"`
}

// VMSnapshotCreateRequest is sent when taking a new snapshot of a VM.
type VMSnapshotCreateRequest struct {
	Snapname    string `json:"snapname"`
	Description string `json:"description,omitempty"`
	VMSTATE     *int   `json:"vmstate,omitempty"`
}

// VMSnapshotUpdateRequest is sent to update the description of an existing snapshot.
type VMSnapshotUpdateRequest struct {
	Description string `json:"description,omitempty"`
}

// ContainerSnapshot is a single snapshot entry for an LXC container.
type ContainerSnapshot struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Snaptime    int64  `json:"snaptime,omitempty"`
	Parent      string `json:"parent,omitempty"`
}

// ContainerSnapshotCreateRequest is sent when snapshotting an LXC container.
type ContainerSnapshotCreateRequest struct {
	Snapname    string `json:"snapname"`
	Description string `json:"description,omitempty"`
}

// ContainerSnapshotUpdateRequest is sent to update the description of a container snapshot.
type ContainerSnapshotUpdateRequest struct {
	Description string `json:"description,omitempty"`
}
