package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// ListNodeDisks returns all physical disks visible on a node.
func (c *Client) ListNodeDisks(ctx context.Context, node string) ([]models.NodeDisk, error) {
	path := fmt.Sprintf("/nodes/%s/disks/list", url.PathEscape(node))
	var result models.APIResponse[[]models.NodeDisk]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetNodeDiskSmart pulls SMART health data for a specific disk on a node.
func (c *Client) GetNodeDiskSmart(ctx context.Context, node, disk string) (*models.NodeDiskSmart, error) {
	path := fmt.Sprintf("/nodes/%s/disks/smart?disk=%s", url.PathEscape(node), url.QueryEscape(disk))
	var result models.APIResponse[models.NodeDiskSmart]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateNodeDiskDirectory formats a disk and creates a directory-based storage on it.
// Returns a UPID for the async task.
func (c *Client) CreateNodeDiskDirectory(ctx context.Context, node string, req *models.NodeDiskDirectoryCreateRequest) (string, error) {
	path := fmt.Sprintf("/nodes/%s/disks/directory", url.PathEscape(node))
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

// DeleteNodeDiskDirectory removes a directory storage from a node, returns the UPID.
func (c *Client) DeleteNodeDiskDirectory(ctx context.Context, node, name string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/disks/directory/%s", url.PathEscape(node), url.PathEscape(name))
	resp, err := c.DoRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result models.APIResponse[string]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// CreateNodeDiskLVM creates an LVM volume group on a disk. Returns UPID for the async task.
func (c *Client) CreateNodeDiskLVM(ctx context.Context, node string, req *models.NodeDiskLVMCreateRequest) (string, error) {
	path := fmt.Sprintf("/nodes/%s/disks/lvm", url.PathEscape(node))
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

// DeleteNodeDiskLVM removes an LVM VG from a node. Returns the async task UPID.
func (c *Client) DeleteNodeDiskLVM(ctx context.Context, node, name string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/disks/lvm/%s", url.PathEscape(node), url.PathEscape(name))
	resp, err := c.DoRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result models.APIResponse[string]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// CreateNodeDiskLVMThin creates an LVM thin pool on a disk. Async, returns UPID.
func (c *Client) CreateNodeDiskLVMThin(ctx context.Context, node string, req *models.NodeDiskLVMThinCreateRequest) (string, error) {
	path := fmt.Sprintf("/nodes/%s/disks/lvmthin", url.PathEscape(node))
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

// DeleteNodeDiskLVMThin removes an LVM thin pool. Note: you also need to pass the volume group name.
func (c *Client) DeleteNodeDiskLVMThin(ctx context.Context, node, name, volumegroup string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/disks/lvmthin/%s?volume-group=%s", url.PathEscape(node), url.PathEscape(name), url.QueryEscape(volumegroup))
	resp, err := c.DoRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result models.APIResponse[string]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// CreateNodeDiskZFS creates a ZFS pool using one or more disks. Returns UPID for the task.
func (c *Client) CreateNodeDiskZFS(ctx context.Context, node string, req *models.NodeDiskZFSCreateRequest) (string, error) {
	path := fmt.Sprintf("/nodes/%s/disks/zfs", url.PathEscape(node))
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

// DeleteNodeDiskZFS destroys a ZFS pool on a node. Returns UPID for the async task.
func (c *Client) DeleteNodeDiskZFS(ctx context.Context, node, name string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/disks/zfs/%s", url.PathEscape(node), url.PathEscape(name))
	resp, err := c.DoRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result models.APIResponse[string]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Data, nil
}
