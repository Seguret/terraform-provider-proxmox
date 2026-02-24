package models

// SDNController is an SDN routing controller (BGP/EVPN/ISIS).
type SDNController struct {
	Controller string `json:"controller"`
	Type       string `json:"type"`
	ASN        int    `json:"asn,omitempty"`
	EBGP       bool   `json:"ebgp,omitempty"`
	ISISIface  string `json:"isis-iface,omitempty"`
	ISISDomain string `json:"isis-domain,omitempty"`
	ISISNet    string `json:"isis-net,omitempty"`
	Loopback   bool   `json:"loopback,omitempty"`
	Node       string `json:"node,omitempty"`
	Peers      string `json:"peers,omitempty"`
}

// SDNControllerCreateRequest is sent when adding a new SDN routing controller.
type SDNControllerCreateRequest struct {
	Controller string `json:"controller"`
	Type       string `json:"type"`
	ASN        int    `json:"asn,omitempty"`
	EBGP       bool   `json:"ebgp,omitempty"`
	Node       string `json:"node,omitempty"`
	Peers      string `json:"peers,omitempty"`
}

// SDNControllerUpdateRequest is sent to update an existing SDN controller.
type SDNControllerUpdateRequest struct {
	ASN   int    `json:"asn,omitempty"`
	EBGP  bool   `json:"ebgp,omitempty"`
	Node  string `json:"node,omitempty"`
	Peers string `json:"peers,omitempty"`
}

// SDNDns is an SDN DNS integration provider (e.g. PowerDNS).
type SDNDns struct {
	Dns  string `json:"dns"`
	Type string `json:"type"`
	URL  string `json:"url,omitempty"`
	Key  string `json:"key,omitempty"`
	TTL  int    `json:"ttl,omitempty"`
}

// SDNDnsCreateRequest is sent when registering a new DNS provider.
type SDNDnsCreateRequest struct {
	Dns  string `json:"dns"`
	Type string `json:"type"`
	URL  string `json:"url,omitempty"`
	Key  string `json:"key,omitempty"`
	TTL  int    `json:"ttl,omitempty"`
}

// SDNDnsUpdateRequest is sent to update DNS provider settings.
type SDNDnsUpdateRequest struct {
	URL string `json:"url,omitempty"`
	Key string `json:"key,omitempty"`
	TTL int    `json:"ttl,omitempty"`
}

// SDNIpam is an SDN IP address management provider.
type SDNIpam struct {
	Ipam    string `json:"ipam"`
	Type    string `json:"type"`
	URL     string `json:"url,omitempty"`
	Token   string `json:"token,omitempty"`
	Section int    `json:"section,omitempty"`
}

// SDNIpamCreateRequest is sent when adding a new IPAM provider.
type SDNIpamCreateRequest struct {
	Ipam    string `json:"ipam"`
	Type    string `json:"type"`
	URL     string `json:"url,omitempty"`
	Token   string `json:"token,omitempty"`
	Section int    `json:"section,omitempty"`
}

// SDNIpamUpdateRequest is sent to update an existing IPAM provider config.
type SDNIpamUpdateRequest struct {
	URL     string `json:"url,omitempty"`
	Token   string `json:"token,omitempty"`
	Section int    `json:"section,omitempty"`
}

// ClusterJobRealmSync is a scheduled job that syncs users from an auth realm.
type ClusterJobRealmSync struct {
	ID             string `json:"id"`
	Realm          string `json:"realm"`
	Schedule       string `json:"schedule"`
	Scope          string `json:"scope,omitempty"`
	RemoveVanished string `json:"remove-vanished,omitempty"`
	EnableNew      bool   `json:"enable-new,omitempty"`
	Comment        string `json:"comment,omitempty"`
	Enabled        bool   `json:"enabled,omitempty"`
}

// ClusterJobRealmSyncCreateRequest is sent when scheduling a new realm-sync job.
type ClusterJobRealmSyncCreateRequest struct {
	Realm          string `json:"realm"`
	Schedule       string `json:"schedule"`
	Scope          string `json:"scope,omitempty"`
	RemoveVanished string `json:"remove-vanished,omitempty"`
	EnableNew      bool   `json:"enable-new,omitempty"`
	Comment        string `json:"comment,omitempty"`
	Enabled        bool   `json:"enabled,omitempty"`
}

// ClusterJobRealmSyncUpdateRequest is sent to modify an existing realm-sync schedule.
type ClusterJobRealmSyncUpdateRequest struct {
	Schedule       string `json:"schedule,omitempty"`
	Scope          string `json:"scope,omitempty"`
	RemoveVanished string `json:"remove-vanished,omitempty"`
	EnableNew      *bool  `json:"enable-new,omitempty"`
	Comment        string `json:"comment,omitempty"`
	Enabled        *bool  `json:"enabled,omitempty"`
}
