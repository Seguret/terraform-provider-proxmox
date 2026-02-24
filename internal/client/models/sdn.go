package models

// SDNZone is an SDN network zone (vlan, vxlan, evpn, etc).
type SDNZone struct {
	Zone                  string `json:"zone"`
	Type                  string `json:"type"`
	Comment               string `json:"comment,omitempty"`
	Nodes                 string `json:"nodes,omitempty"`
	Bridge                string `json:"bridge,omitempty"`
	Tag                   int    `json:"tag,omitempty"`
	VlanProtocol          string `json:"vlan-protocol,omitempty"`
	Peers                 string `json:"peers,omitempty"`
	VRFVxlan              int    `json:"vrf-vxlan,omitempty"`
	Controller            string `json:"controller,omitempty"`
	ExitNodes             string `json:"exitnodes,omitempty"`
	ExitNodesLocalRouting int    `json:"exitnodes-local-routing,omitempty"`
	AdvertiseSubnets      int    `json:"advertise-subnets,omitempty"`
	MTU                   int    `json:"mtu,omitempty"`
	DNS                   string `json:"dns,omitempty"`
	DNSZone               string `json:"dnszone,omitempty"`
	ReverseDNS            string `json:"reversedns,omitempty"`
	IPAM                  string `json:"ipam,omitempty"`
	Pending               *int   `json:"pending,omitempty"`
}

// SDNZoneCreateRequest is sent when adding a new SDN zone.
type SDNZoneCreateRequest struct {
	Zone                  string `json:"zone"`
	Type                  string `json:"type"`
	Comment               string `json:"comment,omitempty"`
	Nodes                 string `json:"nodes,omitempty"`
	Bridge                string `json:"bridge,omitempty"`
	Tag                   int    `json:"tag,omitempty"`
	VlanProtocol          string `json:"vlan-protocol,omitempty"`
	Peers                 string `json:"peers,omitempty"`
	VRFVxlan              int    `json:"vrf-vxlan,omitempty"`
	Controller            string `json:"controller,omitempty"`
	ExitNodes             string `json:"exitnodes,omitempty"`
	ExitNodesLocalRouting int    `json:"exitnodes-local-routing,omitempty"`
	AdvertiseSubnets      int    `json:"advertise-subnets,omitempty"`
	MTU                   int    `json:"mtu,omitempty"`
	DNS                   string `json:"dns,omitempty"`
	DNSZone               string `json:"dnszone,omitempty"`
	ReverseDNS            string `json:"reversedns,omitempty"`
	IPAM                  string `json:"ipam,omitempty"`
}

// SDNZoneUpdateRequest is sent to modify an existing SDN zone.
type SDNZoneUpdateRequest struct {
	Comment               string `json:"comment,omitempty"`
	Nodes                 string `json:"nodes,omitempty"`
	Bridge                string `json:"bridge,omitempty"`
	Tag                   int    `json:"tag,omitempty"`
	VlanProtocol          string `json:"vlan-protocol,omitempty"`
	Peers                 string `json:"peers,omitempty"`
	VRFVxlan              int    `json:"vrf-vxlan,omitempty"`
	Controller            string `json:"controller,omitempty"`
	ExitNodes             string `json:"exitnodes,omitempty"`
	ExitNodesLocalRouting int    `json:"exitnodes-local-routing,omitempty"`
	AdvertiseSubnets      int    `json:"advertise-subnets,omitempty"`
	MTU                   int    `json:"mtu,omitempty"`
	DNS                   string `json:"dns,omitempty"`
	DNSZone               string `json:"dnszone,omitempty"`
	ReverseDNS            string `json:"reversedns,omitempty"`
	IPAM                  string `json:"ipam,omitempty"`
}

// SDNVnet is a virtual network within an SDN zone.
type SDNVnet struct {
	Vnet      string `json:"vnet"`
	Zone      string `json:"zone"`
	Alias     string `json:"alias,omitempty"`
	Tag       int    `json:"tag,omitempty"`
	VlanAware *int   `json:"vlanaware,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

// SDNVnetCreateRequest is sent when creating a new virtual network.
type SDNVnetCreateRequest struct {
	Vnet      string `json:"vnet"`
	Zone      string `json:"zone"`
	Alias     string `json:"alias,omitempty"`
	Tag       int    `json:"tag,omitempty"`
	VlanAware *int   `json:"vlanaware,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

// SDNVnetUpdateRequest is sent to modify an existing VNet.
type SDNVnetUpdateRequest struct {
	Zone      string `json:"zone,omitempty"`
	Alias     string `json:"alias,omitempty"`
	Tag       int    `json:"tag,omitempty"`
	VlanAware *int   `json:"vlanaware,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

// SDNSubnet is an IP subnet attached to an SDN VNet.
type SDNSubnet struct {
	Subnet          string `json:"subnet"`
	Type            string `json:"type,omitempty"`
	Vnet            string `json:"vnet,omitempty"`
	Gateway         string `json:"gateway,omitempty"`
	Snat            *int   `json:"snat,omitempty"`
	DHCPDNSServer   string `json:"dhcp-dns-server,omitempty"`
	DNSZonePrefix   string `json:"dnszoneprefix,omitempty"`
}

// SDNSubnetCreateRequest is sent when adding a subnet to an SDN VNet.
type SDNSubnetCreateRequest struct {
	Subnet          string `json:"subnet"`
	Type            string `json:"type"`
	Gateway         string `json:"gateway,omitempty"`
	Snat            *int   `json:"snat,omitempty"`
	DHCPDNSServer   string `json:"dhcp-dns-server,omitempty"`
	DNSZonePrefix   string `json:"dnszoneprefix,omitempty"`
}

// SDNSubnetUpdateRequest is sent to modify gateway, SNAT or DNS settings on a subnet.
type SDNSubnetUpdateRequest struct {
	Gateway         string `json:"gateway,omitempty"`
	Snat            *int   `json:"snat,omitempty"`
	DHCPDNSServer   string `json:"dhcp-dns-server,omitempty"`
	DNSZonePrefix   string `json:"dnszoneprefix,omitempty"`
}
