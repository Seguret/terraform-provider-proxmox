package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// --- Firewall Rules ---

// GetFirewallRules returns all firewall rules under a given path prefix (cluster/node/vm).
func (c *Client) GetFirewallRules(ctx context.Context, pathPrefix string) ([]models.FirewallRule, error) {
	path := pathPrefix + "/rules"
	var result models.APIResponse[[]models.FirewallRule]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetFirewallRule fetches a single firewall rule by its position index.
func (c *Client) GetFirewallRule(ctx context.Context, pathPrefix string, pos int) (*models.FirewallRule, error) {
	path := fmt.Sprintf("%s/rules/%d", pathPrefix, pos)
	var result models.APIResponse[models.FirewallRule]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateFirewallRule adds a new firewall rule at the given path prefix.
func (c *Client) CreateFirewallRule(ctx context.Context, pathPrefix string, req *models.FirewallRuleCreateRequest) error {
	path := pathPrefix + "/rules"
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, path, bytes.NewReader(body))
}

// UpdateFirewallRule updates the rule at the specified position.
func (c *Client) UpdateFirewallRule(ctx context.Context, pathPrefix string, pos int, req *models.FirewallRuleUpdateRequest) error {
	path := fmt.Sprintf("%s/rules/%d", pathPrefix, pos)
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteFirewallRule removes the rule at the given position.
func (c *Client) DeleteFirewallRule(ctx context.Context, pathPrefix string, pos int) error {
	path := fmt.Sprintf("%s/rules/%d", pathPrefix, pos)
	return c.Delete(ctx, path)
}

// --- IP Sets ---

// GetFirewallIPSets lists all IP sets at the given path prefix.
func (c *Client) GetFirewallIPSets(ctx context.Context, pathPrefix string) ([]models.FirewallIPSet, error) {
	path := pathPrefix + "/ipset"
	var result models.APIResponse[[]models.FirewallIPSet]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// CreateFirewallIPSet creates a named IP set with an optional comment.
func (c *Client) CreateFirewallIPSet(ctx context.Context, pathPrefix, name, comment string) error {
	path := pathPrefix + "/ipset"
	body, err := json.Marshal(map[string]string{"name": name, "comment": comment})
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, path, bytes.NewReader(body))
}

// DeleteFirewallIPSet removes an IP set by name.
func (c *Client) DeleteFirewallIPSet(ctx context.Context, pathPrefix, name string) error {
	path := fmt.Sprintf("%s/ipset/%s", pathPrefix, url.PathEscape(name))
	return c.Delete(ctx, path)
}

// GetFirewallIPSetEntries returns all CIDR entries belonging to an IP set.
func (c *Client) GetFirewallIPSetEntries(ctx context.Context, pathPrefix, name string) ([]models.FirewallIPSetEntry, error) {
	path := fmt.Sprintf("%s/ipset/%s", pathPrefix, url.PathEscape(name))
	var result models.APIResponse[[]models.FirewallIPSetEntry]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// CreateFirewallIPSetEntry adds a CIDR to an existing IP set.
func (c *Client) CreateFirewallIPSetEntry(ctx context.Context, pathPrefix, name string, entry *models.FirewallIPSetEntry) error {
	path := fmt.Sprintf("%s/ipset/%s", pathPrefix, url.PathEscape(name))
	body, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, path, bytes.NewReader(body))
}

// DeleteFirewallIPSetEntry removes a specific CIDR from an IP set.
func (c *Client) DeleteFirewallIPSetEntry(ctx context.Context, pathPrefix, name, cidr string) error {
	path := fmt.Sprintf("%s/ipset/%s/%s", pathPrefix, url.PathEscape(name), url.PathEscape(cidr))
	return c.Delete(ctx, path)
}

// --- Firewall Options ---

// GetFirewallOptions fetches the firewall options (enable/disable, logging level, etc).
func (c *Client) GetFirewallOptions(ctx context.Context, pathPrefix string) (*models.FirewallOptions, error) {
	path := pathPrefix + "/options"
	var result models.APIResponse[models.FirewallOptions]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// UpdateFirewallOptions saves updated firewall settings.
func (c *Client) UpdateFirewallOptions(ctx context.Context, pathPrefix string, opts *models.FirewallOptions) error {
	path := pathPrefix + "/options"
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// --- Helper path builders ---

// ClusterFirewallPath returns the base path for cluster-level firewall operations.
func ClusterFirewallPath() string { return "/cluster/firewall" }

// NodeFirewallPath builds the API path for a node's firewall.
func NodeFirewallPath(node string) string {
	return fmt.Sprintf("/nodes/%s/firewall", url.PathEscape(node))
}

// VMFirewallPath builds the API path for a VM's firewall rules.
func VMFirewallPath(node string, vmid int) string {
	return fmt.Sprintf("/nodes/%s/qemu/%d/firewall", url.PathEscape(node), vmid)
}

// ContainerFirewallPath builds the API path for an LXC container's firewall.
func ContainerFirewallPath(node string, vmid int) string {
	return fmt.Sprintf("/nodes/%s/lxc/%d/firewall", url.PathEscape(node), vmid)
}

// --- Firewall Aliases ---

// GetFirewallAliases lists all named firewall aliases under a path prefix.
func (c *Client) GetFirewallAliases(ctx context.Context, pathPrefix string) ([]models.FirewallAlias, error) {
	path := pathPrefix + "/aliases"
	var result models.APIResponse[[]models.FirewallAlias]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetFirewallAlias looks up a single firewall alias by name.
func (c *Client) GetFirewallAlias(ctx context.Context, pathPrefix, name string) (*models.FirewallAlias, error) {
	path := fmt.Sprintf("%s/aliases/%s", pathPrefix, url.PathEscape(name))
	var result models.APIResponse[models.FirewallAlias]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateFirewallAlias adds a new named IP alias to the firewall config.
func (c *Client) CreateFirewallAlias(ctx context.Context, pathPrefix string, req *models.FirewallAliasCreateRequest) error {
	path := pathPrefix + "/aliases"
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, path, bytes.NewReader(body))
}

// UpdateFirewallAlias changes the CIDR or comment on an existing alias.
func (c *Client) UpdateFirewallAlias(ctx context.Context, pathPrefix, name string, req *models.FirewallAliasUpdateRequest) error {
	path := fmt.Sprintf("%s/aliases/%s", pathPrefix, url.PathEscape(name))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteFirewallAlias removes a named firewall alias.
func (c *Client) DeleteFirewallAlias(ctx context.Context, pathPrefix, name string) error {
	path := fmt.Sprintf("%s/aliases/%s", pathPrefix, url.PathEscape(name))
	return c.Delete(ctx, path)
}

// --- Firewall Security Groups ---

// GetFirewallSecurityGroups returns all security groups defined at the cluster level.
func (c *Client) GetFirewallSecurityGroups(ctx context.Context) ([]models.FirewallSecurityGroup, error) {
	var result models.APIResponse[[]models.FirewallSecurityGroup]
	if err := c.Get(ctx, "/cluster/firewall/groups", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetFirewallSecurityGroup fetches a single security group by name.
func (c *Client) GetFirewallSecurityGroup(ctx context.Context, group string) (*models.FirewallSecurityGroup, error) {
	path := fmt.Sprintf("/cluster/firewall/groups/%s", url.PathEscape(group))
	var result models.APIResponse[models.FirewallSecurityGroup]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateFirewallSecurityGroup adds a new security group to the cluster firewall.
func (c *Client) CreateFirewallSecurityGroup(ctx context.Context, req *models.FirewallSecurityGroupCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/firewall/groups", bytes.NewReader(body))
}

// UpdateFirewallSecurityGroup updates the comment on a security group.
func (c *Client) UpdateFirewallSecurityGroup(ctx context.Context, group string, req *models.FirewallSecurityGroupUpdateRequest) error {
	path := fmt.Sprintf("/cluster/firewall/groups/%s", url.PathEscape(group))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteFirewallSecurityGroup removes a security group and all its rules.
func (c *Client) DeleteFirewallSecurityGroup(ctx context.Context, group string) error {
	path := fmt.Sprintf("/cluster/firewall/groups/%s", url.PathEscape(group))
	return c.Delete(ctx, path)
}
