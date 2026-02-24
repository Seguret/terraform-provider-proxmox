package models

// NodeCapabilities describes what operations a node supports.
type NodeCapabilities struct {
	Backup   map[string]interface{} `json:"backup,omitempty"`
	Migrate  map[string]interface{} `json:"migrate,omitempty"`
	Snapshot map[string]interface{} `json:"snapshot,omitempty"`
}

// NodeSyslogEntry is a single line from the node syslog.
type NodeSyslogEntry struct {
	N    int    `json:"n"`
	Text string `json:"t"`
}

// RRDDataPoint is a single time-series data point from the /rrddata endpoints.
type RRDDataPoint struct {
	Time       int64   `json:"time"`
	CPU        float64 `json:"cpu"`
	MaxCPU     float64 `json:"maxcpu"`
	Mem        float64 `json:"mem"`
	MaxMem     float64 `json:"maxmem"`
	NetIn      float64 `json:"netin"`
	NetOut     float64 `json:"netout"`
	DiskRead   float64 `json:"diskread"`
	DiskWrite  float64 `json:"diskwrite"`
	LoadAvg    float64 `json:"loadavg"`
	SwapTotal  float64 `json:"swaptotal"`
	SwapUsed   float64 `json:"swapused"`
}

// NodeNetstatEntry is per-interface network traffic statistics for a node.
type NodeNetstatEntry struct {
	Iface   string `json:"dev"`
	RxPkts  int64  `json:"in"`
	TxPkts  int64  `json:"out"`
	RxBytes int64  `json:"rxbytes"`
	TxBytes int64  `json:"txbytes"`
	RxErr   int64  `json:"rxerr"`
	TxErr   int64  `json:"txerr"`
	RxDrop  int64  `json:"rxdrop"`
	TxDrop  int64  `json:"txdrop"`
}

// NodeAPTUpdate is a single pending package update from the APT update list.
type NodeAPTUpdate struct {
	Package      string `json:"Package"`
	Title        string `json:"Title"`
	Version      string `json:"Version"`
	OldVersion   string `json:"OldVersion"`
	Priority     string `json:"Priority"`
	Section      string `json:"Section"`
	Architecture string `json:"Architecture"`
	Description  string `json:"Description"`
}
