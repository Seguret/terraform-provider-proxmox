package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// --- VM Snapshots ---

// GetVMSnapshots returns all snapshots for a given VM.
func (c *Client) GetVMSnapshots(ctx context.Context, node string, vmid int) ([]models.VMSnapshot, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/snapshot", url.PathEscape(node), vmid)
	var result models.APIResponse[[]models.VMSnapshot]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetVMSnapshot fetches the config for a specific VM snapshot by name.
func (c *Client) GetVMSnapshot(ctx context.Context, node string, vmid int, snapname string) (*models.VMSnapshot, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/snapshot/%s/config", url.PathEscape(node), vmid, url.PathEscape(snapname))
	var result models.APIResponse[models.VMSnapshot]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	result.Data.Name = snapname
	return &result.Data, nil
}

// CreateVMSnapshot takes a snapshot of a VM — async, returns the UPID.
func (c *Client) CreateVMSnapshot(ctx context.Context, node string, vmid int, req *models.VMSnapshotCreateRequest) (string, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/snapshot", url.PathEscape(node), vmid)
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

// UpdateVMSnapshot updates the description on an existing VM snapshot.
func (c *Client) UpdateVMSnapshot(ctx context.Context, node string, vmid int, snapname string, req *models.VMSnapshotUpdateRequest) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/snapshot/%s/config", url.PathEscape(node), vmid, url.PathEscape(snapname))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteVMSnapshot removes a VM snapshot. Async operation - returns the UPID to track progress.
func (c *Client) DeleteVMSnapshot(ctx context.Context, node string, vmid int, snapname string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/snapshot/%s", url.PathEscape(node), vmid, url.PathEscape(snapname))
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
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Data, nil
}

// RollbackVMSnapshot reverts a VM back to a named snapshot. Returns UPID since this runs async.
func (c *Client) RollbackVMSnapshot(ctx context.Context, node string, vmid int, snapname string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/snapshot/%s/rollback", url.PathEscape(node), vmid, url.PathEscape(snapname))
	var result models.APIResponse[string]
	if err := c.Post(ctx, path, nil, &result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// --- Container Snapshots ---

// GetContainerSnapshots lists all snapshots for a given LXC container.
func (c *Client) GetContainerSnapshots(ctx context.Context, node string, vmid int) ([]models.ContainerSnapshot, error) {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/snapshot", url.PathEscape(node), vmid)
	var result models.APIResponse[[]models.ContainerSnapshot]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetContainerSnapshot retrieves a specific container snapshot config.
func (c *Client) GetContainerSnapshot(ctx context.Context, node string, vmid int, snapname string) (*models.ContainerSnapshot, error) {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/snapshot/%s/config", url.PathEscape(node), vmid, url.PathEscape(snapname))
	var result models.APIResponse[models.ContainerSnapshot]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	result.Data.Name = snapname
	return &result.Data, nil
}

// CreateContainerSnapshot takes a snapshot of a container. Async - returns UPID.
func (c *Client) CreateContainerSnapshot(ctx context.Context, node string, vmid int, req *models.ContainerSnapshotCreateRequest) (string, error) {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/snapshot", url.PathEscape(node), vmid)
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

// UpdateContainerSnapshot updates a container snapshot description.
func (c *Client) UpdateContainerSnapshot(ctx context.Context, node string, vmid int, snapname string, req *models.ContainerSnapshotUpdateRequest) error {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/snapshot/%s/config", url.PathEscape(node), vmid, url.PathEscape(snapname))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteContainerSnapshot deletes a container snapshot. Returns UPID.
func (c *Client) DeleteContainerSnapshot(ctx context.Context, node string, vmid int, snapname string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/snapshot/%s", url.PathEscape(node), vmid, url.PathEscape(snapname))
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
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Data, nil
}

// RollbackContainerSnapshot rolls back a container to a snapshot. Returns UPID.
func (c *Client) RollbackContainerSnapshot(ctx context.Context, node string, vmid int, snapname string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/snapshot/%s/rollback", url.PathEscape(node), vmid, url.PathEscape(snapname))
	var result models.APIResponse[string]
	if err := c.Post(ctx, path, nil, &result); err != nil {
		return "", err
	}
	return result.Data, nil
}
