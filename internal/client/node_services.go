package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// ListNodeServices returns all systemd services managed through Proxmox on a node.
func (c *Client) ListNodeServices(ctx context.Context, node string) ([]models.NodeService, error) {
	path := fmt.Sprintf("/nodes/%s/services", url.PathEscape(node))
	var result models.APIResponse[[]models.NodeService]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetNodeService checks the state (running/stopped) of a specific service on a node.
func (c *Client) GetNodeService(ctx context.Context, node, service string) (*models.NodeService, error) {
	path := fmt.Sprintf("/nodes/%s/services/%s/state", url.PathEscape(node), url.PathEscape(service))
	var result models.APIResponse[models.NodeService]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// StartNodeService starts a service on a node, returns task UPID.
func (c *Client) StartNodeService(ctx context.Context, node, service string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/services/%s/start", url.PathEscape(node), url.PathEscape(service))
	var result models.APIResponse[string]
	if err := c.Post(ctx, path, nil, &result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// StopNodeService stops a running service on a node, returns task UPID.
func (c *Client) StopNodeService(ctx context.Context, node, service string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/services/%s/stop", url.PathEscape(node), url.PathEscape(service))
	var result models.APIResponse[string]
	if err := c.Post(ctx, path, nil, &result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// RestartNodeService restarts a service on a node, returns task UPID.
func (c *Client) RestartNodeService(ctx context.Context, node, service string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/services/%s/restart", url.PathEscape(node), url.PathEscape(service))
	var result models.APIResponse[string]
	if err := c.Post(ctx, path, nil, &result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// ReloadNodeService sends a reload signal to a service without full restart, returns task UPID.
func (c *Client) ReloadNodeService(ctx context.Context, node, service string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/services/%s/reload", url.PathEscape(node), url.PathEscape(service))
	var result models.APIResponse[string]
	if err := c.Post(ctx, path, nil, &result); err != nil {
		return "", err
	}
	return result.Data, nil
}
