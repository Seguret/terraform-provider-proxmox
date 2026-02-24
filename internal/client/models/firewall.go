package models

// FirewallRule is a single rule in a firewall ruleset (cluster, node, VM, or container).
type FirewallRule struct {
	Pos     int    `json:"pos"`
	Type    string `json:"type"`    // in or out
	Action  string `json:"action"`  // ACCEPT, DROP, REJECT
	Enable  *int   `json:"enable,omitempty"`
	Macro   string `json:"macro,omitempty"`
	Proto   string `json:"proto,omitempty"`
	Source  string `json:"source,omitempty"`
	Dest    string `json:"dest,omitempty"`
	DPort   string `json:"dport,omitempty"`
	Sport   string `json:"sport,omitempty"`
	IFace   string `json:"iface,omitempty"`
	Log     string `json:"log,omitempty"`
	Comment string `json:"comment,omitempty"`
	ICMPType string `json:"icmp-type,omitempty"`
}

// FirewallRuleCreateRequest is sent when adding a new firewall rule.
type FirewallRuleCreateRequest struct {
	Type    string `json:"type"`
	Action  string `json:"action"`
	Enable  *int   `json:"enable,omitempty"`
	Macro   string `json:"macro,omitempty"`
	Proto   string `json:"proto,omitempty"`
	Source  string `json:"source,omitempty"`
	Dest    string `json:"dest,omitempty"`
	DPort   string `json:"dport,omitempty"`
	Sport   string `json:"sport,omitempty"`
	IFace   string `json:"iface,omitempty"`
	Log     string `json:"log,omitempty"`
	Comment string `json:"comment,omitempty"`
	Pos     *int   `json:"pos,omitempty"`
}

// FirewallRuleUpdateRequest is sent to modify an existing firewall rule at a given position.
type FirewallRuleUpdateRequest struct {
	Type    *string `json:"type,omitempty"`
	Action  *string `json:"action,omitempty"`
	Enable  *int    `json:"enable,omitempty"`
	Macro   *string `json:"macro,omitempty"`
	Proto   *string `json:"proto,omitempty"`
	Source  *string `json:"source,omitempty"`
	Dest    *string `json:"dest,omitempty"`
	DPort   *string `json:"dport,omitempty"`
	Sport   *string `json:"sport,omitempty"`
	IFace   *string `json:"iface,omitempty"`
	Log     *string `json:"log,omitempty"`
	Comment *string `json:"comment,omitempty"`
	MoveTo  *int    `json:"moveto,omitempty"`
}

// FirewallIPSet is a named collection of IP addresses or CIDRs used in firewall rules.
type FirewallIPSet struct {
	Name    string `json:"name"`
	Comment string `json:"comment,omitempty"`
	Digest  string `json:"digest,omitempty"`
}

// FirewallIPSetEntry is a single CIDR entry inside an IP set.
type FirewallIPSetEntry struct {
	CIDR    string `json:"cidr"`
	Comment string `json:"comment,omitempty"`
	NoMatch *int   `json:"nomatch,omitempty"`
}

// FirewallOptions holds the enable/policy settings for a firewall level.
type FirewallOptions struct {
	Enable         *int    `json:"enable,omitempty"`
	PolicyIn       string  `json:"policy_in,omitempty"`
	PolicyOut      string  `json:"policy_out,omitempty"`
	LogRatelimit   string  `json:"log_ratelimit,omitempty"`
	Ebtables       *int    `json:"ebtables,omitempty"`
	Macfilter      *int    `json:"macfilter,omitempty"`
	IPFilter       *int    `json:"ipfilter,omitempty"`
	NDPProxy       *int    `json:"ndp,omitempty"`
	DHCPFilter     *int    `json:"dhcp,omitempty"`
	RadVert        *int    `json:"radv,omitempty"`
}

// FirewallAlias is a named shorthand for an IP address or CIDR range in firewall rules.
type FirewallAlias struct {
	Name    string `json:"name"`
	CIDR    string `json:"cidr"`
	Comment string `json:"comment,omitempty"`
	Digest  string `json:"digest,omitempty"`
}

// FirewallAliasCreateRequest is sent when adding a new firewall alias.
type FirewallAliasCreateRequest struct {
	Name    string `json:"name"`
	CIDR    string `json:"cidr"`
	Comment string `json:"comment,omitempty"`
}

// FirewallAliasUpdateRequest is sent to rename or change the CIDR of a firewall alias.
type FirewallAliasUpdateRequest struct {
	Rename  string `json:"rename,omitempty"`
	CIDR    string `json:"cidr"`
	Comment string `json:"comment,omitempty"`
}

// FirewallSecurityGroup is a reusable named set of firewall rules at the cluster level.
type FirewallSecurityGroup struct {
	Group   string `json:"group"`
	Comment string `json:"comment,omitempty"`
	Digest  string `json:"digest,omitempty"`
}

// FirewallSecurityGroupCreateRequest is sent when creating a new security group.
type FirewallSecurityGroupCreateRequest struct {
	Group   string `json:"group"`
	Comment string `json:"comment,omitempty"`
}

// FirewallSecurityGroupUpdateRequest is sent to rename or update a security group.
type FirewallSecurityGroupUpdateRequest struct {
	Rename  string `json:"rename,omitempty"`
	Comment string `json:"comment,omitempty"`
}
