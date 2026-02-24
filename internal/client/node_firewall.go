package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetNodeFirewallRules returns all firewall rules configured on a specific node.
func (c *Client) GetNodeFirewallRules(ctx context.Context, node string) ([]models.FirewallRule, error) {
	path := fmt.Sprintf("/nodes/%s/firewall/rules", url.PathEscape(node))
	var result models.APIResponse[[]models.FirewallRule]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetVMFirewallRules returns all firewall rules for a specific VM.
func (c *Client) GetVMFirewallRules(ctx context.Context, node string, vmid int) ([]models.FirewallRule, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/firewall/rules", url.PathEscape(node), vmid)
	var result models.APIResponse[[]models.FirewallRule]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetContainerFirewallRules returns all firewall rules for a given container.
func (c *Client) GetContainerFirewallRules(ctx context.Context, node string, vmid int) ([]models.FirewallRule, error) {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/firewall/rules", url.PathEscape(node), vmid)
	var result models.APIResponse[[]models.FirewallRule]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

