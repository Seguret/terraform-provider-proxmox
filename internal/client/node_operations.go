package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetNodeCapabilities returns what features/capabilities are supported on a node.
func (c *Client) GetNodeCapabilities(ctx context.Context, node string) (*models.NodeCapabilities, error) {
	path := fmt.Sprintf("/nodes/%s/capabilities", url.PathEscape(node))
	var result models.APIResponse[models.NodeCapabilities]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// ExecuteNodeCommand runs one or more shell commands on a node.
func (c *Client) ExecuteNodeCommand(ctx context.Context, node string, commands []string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/execute", url.PathEscape(node))
	params := map[string]interface{}{
		"commands": commands,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	return c.PostTask(ctx, path, bytes.NewReader(body))
}

// GetNodeSyslog fetches syslog entries from a node with pagination via start/limit.
func (c *Client) GetNodeSyslog(ctx context.Context, node string, start, limit int) ([]models.NodeSyslogEntry, error) {
	path := fmt.Sprintf("/nodes/%s/syslog?start=%d&limit=%d", url.PathEscape(node), start, limit)
	var result models.APIResponse[[]models.NodeSyslogEntry]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetNodeRRDData fetches RRD performance metrics for a node over a given timeframe.
func (c *Client) GetNodeRRDData(ctx context.Context, node, timeframe string) ([]models.RRDDataPoint, error) {
	path := fmt.Sprintf("/nodes/%s/rrddata?timeframe=%s", url.PathEscape(node), url.QueryEscape(timeframe))
	var result models.APIResponse[[]models.RRDDataPoint]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetNodeNetstat returns per-interface network stats for a node.
func (c *Client) GetNodeNetstat(ctx context.Context, node string) ([]models.NodeNetstatEntry, error) {
	path := fmt.Sprintf("/nodes/%s/netstat", url.PathEscape(node))
	var result models.APIResponse[[]models.NodeNetstatEntry]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// ScanNode probes a node for available storage sources of a given type (nfs, zfs, etc).
func (c *Client) ScanNode(ctx context.Context, node, scanType string, params map[string]string) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/nodes/%s/scan/%s", url.PathEscape(node), url.PathEscape(scanType))
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			q.Set(k, v)
		}
		path += "?" + q.Encode()
	}
	var result models.APIResponse[[]map[string]interface{}]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// UploadToStorage uploads a file to node storage.
func (c *Client) UploadToStorage(ctx context.Context, node, storage, content, filename string, data []byte) (string, error) {
	path := fmt.Sprintf("/nodes/%s/storage/%s/upload", url.PathEscape(node), url.PathEscape(storage))
	// TODO: real upload needs multipart form, this is a placeholder
	params := map[string]interface{}{
		"content":  content,
		"filename": filename,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	return c.PostTask(ctx, path, bytes.NewReader(body))
}

// PurgeStorage prunes old backups from a storage pool based on the prune policy.
func (c *Client) PurgeStorage(ctx context.Context, node, storage, pruneBackups string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/storage/%s/prune-backups?prune-backups=%s", 
		url.PathEscape(node), url.PathEscape(storage), url.QueryEscape(pruneBackups))
	return c.DeleteTask(ctx, path)
}

// GetNodeAPTUpdate lists packages that have available updates on a node.
func (c *Client) GetNodeAPTUpdate(ctx context.Context, node string) ([]models.NodeAPTUpdate, error) {
	path := fmt.Sprintf("/nodes/%s/apt/update", url.PathEscape(node))
	var result models.APIResponse[[]models.NodeAPTUpdate]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// RefreshNodeAPT triggers an apt-get update on a node. Returns UPID since it runs async.
func (c *Client) RefreshNodeAPT(ctx context.Context, node string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/apt/update", url.PathEscape(node))
	return c.PostTask(ctx, path, nil)
}
