package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetVMConsoleProxy returns VNC websocket connection info for a VM's console.
func (c *Client) GetVMConsoleProxy(ctx context.Context, node string, vmid int, vncticket, port string) (*models.ConsoleProxyInfo, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/vncwebsocket", url.PathEscape(node), vmid)
	params := map[string]interface{}{
		"vncticket": vncticket,
		"port":      port,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	var result models.APIResponse[models.ConsoleProxyInfo]
	if err := c.PostJSON(ctx, path, bytes.NewReader(body), &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetContainerConsoleProxy returns VNC websocket connection info for a container's console.
func (c *Client) GetContainerConsoleProxy(ctx context.Context, node string, vmid int, vncticket, port string) (*models.ConsoleProxyInfo, error) {
	path := fmt.Sprintf("/nodes/%s/lxc/%d/vncwebsocket", url.PathEscape(node), vmid)
	params := map[string]interface{}{
		"vncticket": vncticket,
		"port":      port,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	var result models.APIResponse[models.ConsoleProxyInfo]
	if err := c.PostJSON(ctx, path, bytes.NewReader(body), &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetQEMUMonitor fetches info about the QEMU monitor interface on a VM.
func (c *Client) GetQEMUMonitor(ctx context.Context, node string, vmid int) (*models.QEMUMonitorInfo, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/monitor", url.PathEscape(node), vmid)
	var result models.APIResponse[models.QEMUMonitorInfo]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// ExecuteQEMUMonitor runs a raw QEMU monitor command on a VM.
func (c *Client) ExecuteQEMUMonitor(ctx context.Context, node string, vmid int, command string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/monitor", url.PathEscape(node), vmid)
	params := map[string]interface{}{
		"command": command,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	return c.PostTask(ctx, path, bytes.NewReader(body))
}
