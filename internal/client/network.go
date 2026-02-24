package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetNetworkInterfaces returns all network interfaces configured on a node.
func (c *Client) GetNetworkInterfaces(ctx context.Context, node string) ([]models.NetworkInterface, error) {
	path := fmt.Sprintf("/nodes/%s/network", url.PathEscape(node))
	var result models.APIResponse[[]models.NetworkInterface]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetNetworkInterface fetches a single network interface by name.
func (c *Client) GetNetworkInterface(ctx context.Context, node, iface string) (*models.NetworkInterface, error) {
	path := fmt.Sprintf("/nodes/%s/network/%s", url.PathEscape(node), url.PathEscape(iface))
	var result models.APIResponse[models.NetworkInterface]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateNetworkInterface adds a new network interface to a node (bridge, bond, vlan, etc).
func (c *Client) CreateNetworkInterface(ctx context.Context, node string, req *models.NetworkInterfaceCreateRequest) error {
	path := fmt.Sprintf("/nodes/%s/network", url.PathEscape(node))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, path, bytes.NewReader(body))
}

// UpdateNetworkInterface updates settings on an existing network interface.
func (c *Client) UpdateNetworkInterface(ctx context.Context, node, iface string, req *models.NetworkInterfaceCreateRequest) error {
	path := fmt.Sprintf("/nodes/%s/network/%s", url.PathEscape(node), url.PathEscape(iface))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteNetworkInterface removes a network interface from the node config.
func (c *Client) DeleteNetworkInterface(ctx context.Context, node, iface string) error {
	path := fmt.Sprintf("/nodes/%s/network/%s", url.PathEscape(node), url.PathEscape(iface))
	return c.Delete(ctx, path)
}

// ApplyNetworkConfig commits any pending network config changes on a node.
func (c *Client) ApplyNetworkConfig(ctx context.Context, node string) error {
	path := fmt.Sprintf("/nodes/%s/network", url.PathEscape(node))
	return c.Put(ctx, path, nil)
}

// RevertNetworkConfig discards any unsaved network changes on a node.
func (c *Client) RevertNetworkConfig(ctx context.Context, node string) error {
	path := fmt.Sprintf("/nodes/%s/network", url.PathEscape(node))
	return c.Delete(ctx, path)
}
