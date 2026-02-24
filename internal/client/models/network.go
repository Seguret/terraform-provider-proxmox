package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// IntOrString handles fields that can come back as either a JSON number or a string.
type IntOrString int

func (i *IntOrString) UnmarshalJSON(b []byte) error {
	if string(b) == "null" || len(b) == 0 {
		*i = 0
		return nil
	}

	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if s == "" {
			*i = 0
			return nil
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*i = IntOrString(v)
		return nil
	}

	var num json.Number
	if err := json.Unmarshal(b, &num); err == nil {
		if v, err := num.Int64(); err == nil {
			*i = IntOrString(v)
			return nil
		}
	}

	var v int
	if err := json.Unmarshal(b, &v); err == nil {
		*i = IntOrString(v)
		return nil
	}

	return fmt.Errorf("invalid int value: %s", string(b))
}

// NetworkInterface is a node network interface (bridge, bond, VLAN, etc).
type NetworkInterface struct {
	Iface              string      `json:"iface"`
	Type               string      `json:"type"` // bridge, bond, eth, vlan, OVSBridge, OVSBond, OVSPort, OVSIntPort
	Active             int         `json:"active,omitempty"`
	Autostart          int         `json:"autostart,omitempty"`
	Method             string      `json:"method,omitempty"`  // static, dhcp, manual
	Method6            string      `json:"method6,omitempty"` // static, dhcp, manual
	Address            string      `json:"address,omitempty"`
	Netmask            string      `json:"netmask,omitempty"`
	Gateway            string      `json:"gateway,omitempty"`
	Address6           string      `json:"address6,omitempty"`
	Netmask6           int         `json:"netmask6,omitempty"`
	Gateway6           string      `json:"gateway6,omitempty"`
	CIDR               string      `json:"cidr,omitempty"`
	CIDR6              string      `json:"cidr6,omitempty"`
	BridgePorts        string      `json:"bridge_ports,omitempty"`
	BridgeSTP          string      `json:"bridge_stp,omitempty"`
	BridgeFD           IntOrString `json:"bridge_fd,omitempty"`
	BridgeVLANAware    int         `json:"bridge_vlan_aware,omitempty"`
	BondPrimary        string      `json:"bond_primary,omitempty"`
	BondMode           string      `json:"bond_mode,omitempty"`
	BondXmitHashPolicy string      `json:"bond_xmit_hash_policy,omitempty"`
	Slaves             string      `json:"slaves,omitempty"`
	VLANRawDev         string      `json:"vlan-raw-device,omitempty"`
	VLANID             int         `json:"vlan-id,omitempty"`
	MTU                int         `json:"mtu,omitempty"`
	Comments           string      `json:"comments,omitempty"`
	Comments6          string      `json:"comments6,omitempty"`
	Families           []string    `json:"families,omitempty"`
}

// NetworkInterfaceCreateRequest is sent when adding a new network interface to a node.
type NetworkInterfaceCreateRequest struct {
	Iface              string `json:"iface"`
	Type               string `json:"type"`
	Autostart          *int   `json:"autostart,omitempty"`
	Method             string `json:"method,omitempty"`
	Method6            string `json:"method6,omitempty"`
	Address            string `json:"address,omitempty"`
	Netmask            string `json:"netmask,omitempty"`
	Gateway            string `json:"gateway,omitempty"`
	Address6           string `json:"address6,omitempty"`
	Netmask6           *int   `json:"netmask6,omitempty"`
	Gateway6           string `json:"gateway6,omitempty"`
	CIDR               string `json:"cidr,omitempty"`
	CIDR6              string `json:"cidr6,omitempty"`
	BridgePorts        string `json:"bridge_ports,omitempty"`
	BridgeSTP          string `json:"bridge_stp,omitempty"`
	BridgeFD           *int   `json:"bridge_fd,omitempty"`
	BridgeVLANAware    *int   `json:"bridge_vlan_aware,omitempty"`
	BondPrimary        string `json:"bond_primary,omitempty"`
	BondMode           string `json:"bond_mode,omitempty"`
	BondXmitHashPolicy string `json:"bond_xmit_hash_policy,omitempty"`
	Slaves             string `json:"slaves,omitempty"`
	VLANRawDev         string `json:"vlan-raw-device,omitempty"`
	VLANID             *int   `json:"vlan-id,omitempty"`
	MTU                *int   `json:"mtu,omitempty"`
	Comments           string `json:"comments,omitempty"`
	Comments6          string `json:"comments6,omitempty"`
}
