package models

// MetricsServer is an external metrics target (InfluxDB, Graphite, etc).
type MetricsServer struct {
	ID           string  `json:"id"`
	Type         string  `json:"type"`
	Server       string  `json:"server"`
	Port         int     `json:"port"`
	Disable      *int    `json:"disable,omitempty"`
	MTU          int     `json:"mtu,omitempty"`
	Path         string  `json:"path,omitempty"`
	Proto        string  `json:"proto,omitempty"`
	Timeout      int     `json:"timeout,omitempty"`
	Bucket       string  `json:"bucket,omitempty"`
	InfluxDBProto string  `json:"influxdbproto,omitempty"`
	Organization string  `json:"organization,omitempty"`
	Token        string  `json:"token,omitempty"`
	MaxBodySize  int     `json:"max-body-size,omitempty"`
}

// MetricsServerCreateRequest is sent when adding a new external metrics server.
type MetricsServerCreateRequest struct {
	Type         string `json:"type"`
	Server       string `json:"server"`
	Port         int    `json:"port"`
	Disable      *int   `json:"disable,omitempty"`
	MTU          int    `json:"mtu,omitempty"`
	Path         string `json:"path,omitempty"`
	Proto        string `json:"proto,omitempty"`
	Timeout      int    `json:"timeout,omitempty"`
	Bucket       string `json:"bucket,omitempty"`
	InfluxDBProto string `json:"influxdbproto,omitempty"`
	Organization string `json:"organization,omitempty"`
	Token        string `json:"token,omitempty"`
	MaxBodySize  int    `json:"max-body-size,omitempty"`
}

// MetricsServerUpdateRequest is sent to update an existing metrics server config.
type MetricsServerUpdateRequest struct {
	Server       string `json:"server,omitempty"`
	Port         int    `json:"port,omitempty"`
	Disable      *int   `json:"disable,omitempty"`
	MTU          int    `json:"mtu,omitempty"`
	Path         string `json:"path,omitempty"`
	Proto        string `json:"proto,omitempty"`
	Timeout      int    `json:"timeout,omitempty"`
	Bucket       string `json:"bucket,omitempty"`
	InfluxDBProto string `json:"influxdbproto,omitempty"`
	Organization string `json:"organization,omitempty"`
	Token        string `json:"token,omitempty"`
	MaxBodySize  int    `json:"max-body-size,omitempty"`
}
