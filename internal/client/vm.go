package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetVMs lists all QEMU VMs running on a specific node.
func (c *Client) GetVMs(ctx context.Context, node string) ([]models.VMListEntry, error) {
	path := fmt.Sprintf("/nodes/%s/qemu", url.PathEscape(node))
	var result models.APIResponse[[]models.VMListEntry]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetVMConfig fetches the current configuration of a VM.
func (c *Client) GetVMConfig(ctx context.Context, node string, vmid int) (*models.VMConfig, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/config", url.PathEscape(node), vmid)
	var result models.APIResponse[models.VMConfig]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetVMStatus returns the current running status of a VM (running, stopped, etc).
func (c *Client) GetVMStatus(ctx context.Context, node string, vmid int) (*models.VMStatus, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/status/current", url.PathEscape(node), vmid)
	var result models.APIResponse[models.VMStatus]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateVM provisions a new virtual machine on the given node. Returns the task UPID.
func (c *Client) CreateVM(ctx context.Context, node string, req *models.VMCreateRequest) (string, error) {
	path := fmt.Sprintf("/nodes/%s/qemu", url.PathEscape(node))
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

// UpdateVMConfig applies config changes to an existing VM. Accepts a raw map for flexiblity.
func (c *Client) UpdateVMConfig(ctx context.Context, node string, vmid int, configMap map[string]interface{}) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/config", url.PathEscape(node), vmid)
	body, err := json.Marshal(configMap)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteVM destroys a VM and its associated resources. Async - returns the UPID.
func (c *Client) DeleteVM(ctx context.Context, node string, vmid int) (string, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d", url.PathEscape(node), vmid)

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

// StartVM powers on a VM. Returns the UPID of the start task.
func (c *Client) StartVM(ctx context.Context, node string, vmid int) (string, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/status/start", url.PathEscape(node), vmid)
	var result models.APIResponse[string]
	if err := c.Post(ctx, path, nil, &result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// StopVM forcefully powers off a VM (hard stop). Returns the UPID.
func (c *Client) StopVM(ctx context.Context, node string, vmid int) (string, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/status/stop", url.PathEscape(node), vmid)
	var result models.APIResponse[string]
	if err := c.Post(ctx, path, nil, &result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// ShutdownVM sends an ACPI shutdown signal to the VM for a clean poweroff. Returns the UPID.
func (c *Client) ShutdownVM(ctx context.Context, node string, vmid int) (string, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/status/shutdown", url.PathEscape(node), vmid)
	var result models.APIResponse[string]
	if err := c.Post(ctx, path, nil, &result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// CloneVM creates a copy of an existing VM. Returns the UPID since cloning is async.
func (c *Client) CloneVM(ctx context.Context, node string, vmid int, req *models.VMCloneRequest) (string, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/clone", url.PathEscape(node), vmid)
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

// ResizeVMDisk grows (or shrinks) a disk attached to a VM.
func (c *Client) ResizeVMDisk(ctx context.Context, node string, vmid int, req *models.VMResizeRequest) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/resize", url.PathEscape(node), vmid)
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// GetNextVMID asks the cluster for the next free VMID to use.
func (c *Client) GetNextVMID(ctx context.Context) (int, error) {
	var result models.APIResponse[json.Number]
	if err := c.Get(ctx, "/cluster/nextid", &result); err != nil {
		return 0, err
	}
	id, err := result.Data.Int64()
	if err != nil {
		return 0, fmt.Errorf("invalid VMID: %w", err)
	}
	return int(id), nil
}

// GetVMVNCProxy sets up a VNC proxy session for a VM and returns connection info.
func (c *Client) GetVMVNCProxy(ctx context.Context, node string, vmid int) (*models.VMProxyInfo, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/vncproxy", url.PathEscape(node), vmid)
	params := map[string]interface{}{}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	var result models.APIResponse[models.VMProxyInfo]
	if err := c.PostJSON(ctx, path, bytes.NewReader(body), &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetVMSPICEProxy sets up a SPICE proxy session for a VM. Pass a proxy hostname if needed.
func (c *Client) GetVMSPICEProxy(ctx context.Context, node string, vmid int, proxyParam string) (*models.VMProxyInfo, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/spiceproxy", url.PathEscape(node), vmid)
	params := map[string]interface{}{}
	if proxyParam != "" {
		params["proxy"] = proxyParam
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	var result models.APIResponse[models.VMProxyInfo]
	if err := c.PostJSON(ctx, path, bytes.NewReader(body), &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetVMTermProxy opens a terminal proxy session for a VM. Optionally target a specific serial port.
func (c *Client) GetVMTermProxy(ctx context.Context, node string, vmid int, serial string) (*models.VMProxyInfo, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/termproxy", url.PathEscape(node), vmid)
	params := map[string]interface{}{}
	if serial != "" {
		params["serial"] = serial
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	var result models.APIResponse[models.VMProxyInfo]
	if err := c.PostJSON(ctx, path, bytes.NewReader(body), &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// MigrateVM moves a VM to another cluster node. Set online=true for live migration.
func (c *Client) MigrateVM(ctx context.Context, node string, vmid int, target string, online bool) (string, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/migrate", url.PathEscape(node), vmid)
	params := map[string]interface{}{
		"target": target,
	}
	if online {
		params["online"] = 1
	}
	body, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	return c.PostTask(ctx, path, bytes.NewReader(body))
}

// SendKeyVM injects a key event into a running VM (useful for ctrl-alt-del etc).
func (c *Client) SendKeyVM(ctx context.Context, node string, vmid int, key string) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/sendkey", url.PathEscape(node), vmid)
	params := map[string]string{"key": key}
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// GetVMFeature checks whether a specific feature (e.g. snapshot) is supported for a VM.
func (c *Client) GetVMFeature(ctx context.Context, node string, vmid int, feature string) (*models.VMFeatureInfo, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/feature?feature=%s", url.PathEscape(node), vmid, url.QueryEscape(feature))
	var result models.APIResponse[models.VMFeatureInfo]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetVMRRDData pulls historical performance metrics for a VM over the given timeframe.
func (c *Client) GetVMRRDData(ctx context.Context, node string, vmid int, timeframe string) ([]models.RRDDataPoint, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/rrddata?timeframe=%s", url.PathEscape(node), vmid, url.QueryEscape(timeframe))
	var result models.APIResponse[[]models.RRDDataPoint]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}
