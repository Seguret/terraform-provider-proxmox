package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetContainers returns all LXC containers running on a node.
func (c *Client) GetContainers(ctx context.Context, node string) ([]models.ContainerListEntry, error) {
	path := fmt.Sprintf("/nodes/%s/lxc", url.PathEscape(node))
	var result models.APIResponse[[]models.ContainerListEntry]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetContainerConfig fetches the full configuration for a container.
func (c *Client) GetContainerConfig(ctx context.Context, node string, vmid int) (*models.ContainerConfig, error) {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/config", url.PathEscape(node), vmid)
	var result models.APIResponse[models.ContainerConfig]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetContainerStatus returns the current running status of a container (running, stopped, etc).
func (c *Client) GetContainerStatus(ctx context.Context, node string, vmid int) (*models.ContainerStatus, error) {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/status/current", url.PathEscape(node), vmid)
	var result models.APIResponse[models.ContainerStatus]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateContainer kicks off LXC container creation and returns the task UPID.
func (c *Client) CreateContainer(ctx context.Context, node string, req *models.ContainerCreateRequest) (string, error) {
	path := fmt.Sprintf("/nodes/%s/lxc", url.PathEscape(node))
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	var result models.APIResponse[string]
	if err := c.Post(ctx, path, bytes.NewReader(body), &result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// UpdateContainerConfig applies config changes to a container, takes a generic map.
func (c *Client) UpdateContainerConfig(ctx context.Context, node string, vmid int, configMap map[string]interface{}) error {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/config", url.PathEscape(node), vmid)
	body, err := json.Marshal(configMap)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteContainer destroys a container and returns the UPID for the async task.
func (c *Client) DeleteContainer(ctx context.Context, node string, vmid int) (string, error) {
	path := fmt.Sprintf("/nodes/%s/lxc/%d", url.PathEscape(node), vmid)
	resp, err := c.DoRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", c.parseError(resp)
	}
	var result models.APIResponse[string]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode delete response: %w", err)
	}
	return result.Data, nil
}

// StartContainer boots up a container, returns the task UPID.
func (c *Client) StartContainer(ctx context.Context, node string, vmid int) (string, error) {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/status/start", url.PathEscape(node), vmid)
	var result models.APIResponse[string]
	if err := c.Post(ctx, path, nil, &result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// StopContainer forcefully halts a container, returns the task UPID.
func (c *Client) StopContainer(ctx context.Context, node string, vmid int) (string, error) {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/status/stop", url.PathEscape(node), vmid)
	var result models.APIResponse[string]
	if err := c.Post(ctx, path, nil, &result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// ShutdownContainer sends a graceful shutdown signal to a container.
func (c *Client) ShutdownContainer(ctx context.Context, node string, vmid int) (string, error) {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/status/shutdown", url.PathEscape(node), vmid)
	var result models.APIResponse[string]
	if err := c.Post(ctx, path, nil, &result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// MigrateContainer moves a container to another node. Pass restart=true to restart on arrival.
func (c *Client) MigrateContainer(ctx context.Context, node string, vmid int, target string, restart bool) (string, error) {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/migrate", url.PathEscape(node), vmid)
	params := map[string]interface{}{
		"target": target,
	}
	if restart {
		params["restart"] = 1
	}
	body, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	return c.PostTask(ctx, path, bytes.NewReader(body))
}

// GetContainerRRDData pulls performance metrics (RRD) for a container over a given timeframe.
func (c *Client) GetContainerRRDData(ctx context.Context, node string, vmid int, timeframe string) ([]models.RRDDataPoint, error) {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/rrddata?timeframe=%s", url.PathEscape(node), vmid, url.QueryEscape(timeframe))
	var result models.APIResponse[[]models.RRDDataPoint]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}
