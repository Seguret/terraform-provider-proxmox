package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetClusterStatus returns the cluster membership status - nodes, quorum, etc.
func (c *Client) GetClusterStatus(ctx context.Context) ([]models.ClusterStatusEntry, error) {
	var result models.APIResponse[[]models.ClusterStatusEntry]
	if err := c.Get(ctx, "/cluster/status", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetClusterResources lists all cluster resources, with optional type filter (vm, storage, node, etc).
func (c *Client) GetClusterResources(ctx context.Context, resourceType string) ([]models.ClusterResource, error) {
	path := "/cluster/resources"
	if resourceType != "" {
		path += "?type=" + url.QueryEscape(resourceType)
	}
	var result models.APIResponse[[]models.ClusterResource]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetClusterTasks grabs the recent task list across all cluster nodes.
func (c *Client) GetClusterTasks(ctx context.Context) ([]models.ClusterTask, error) {
	var result models.APIResponse[[]models.ClusterTask]
	if err := c.Get(ctx, "/cluster/tasks", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetHAStatus returns current HA manager status for the cluster.
func (c *Client) GetHAStatus(ctx context.Context) ([]models.HAStatusEntry, error) {
	var result models.APIResponse[[]models.HAStatusEntry]
	if err := c.Get(ctx, "/cluster/ha/status/current", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// ListNodeHardwarePCI enumerates PCI devices present on a node.
func (c *Client) ListNodeHardwarePCI(ctx context.Context, node string) ([]models.NodeHardwarePCI, error) {
	path := fmt.Sprintf("/nodes/%s/hardware/pci", url.PathEscape(node))
	var result models.APIResponse[[]models.NodeHardwarePCI]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// ListNodeHardwareUSB enumerates USB devices attached to a node.
func (c *Client) ListNodeHardwareUSB(ctx context.Context, node string) ([]models.NodeHardwareUSB, error) {
	path := fmt.Sprintf("/nodes/%s/hardware/usb", url.PathEscape(node))
	var result models.APIResponse[[]models.NodeHardwareUSB]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// ListNodeTasks returns the task history for a given node.
func (c *Client) ListNodeTasks(ctx context.Context, node string) ([]models.NodeTask, error) {
	path := fmt.Sprintf("/nodes/%s/tasks", url.PathEscape(node))
	var result models.APIResponse[[]models.NodeTask]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetNodeJournal pulls systemd journal lines from a node.
func (c *Client) GetNodeJournal(ctx context.Context, node string) ([]string, error) {
	path := fmt.Sprintf("/nodes/%s/journal", url.PathEscape(node))
	var result models.APIResponse[[]string]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetAPTChangelog fetches the changelog text for a specific APT package on a node.
func (c *Client) GetAPTChangelog(ctx context.Context, node, pkg string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/apt/changelog?name=%s", url.PathEscape(node), url.QueryEscape(pkg))
	var result models.APIResponse[string]
	if err := c.Get(ctx, path, &result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// GetAPTVersions returns all installed package versions on a node.
func (c *Client) GetAPTVersions(ctx context.Context, node string) ([]models.APTPackageVersion, error) {
	path := fmt.Sprintf("/nodes/%s/apt/versions", url.PathEscape(node))
	var result models.APIResponse[[]models.APTPackageVersion]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetACMEDirectories returns the known ACME CA directory URLs.
func (c *Client) GetACMEDirectories(ctx context.Context) ([]models.ACMEDirectory, error) {
	var result models.APIResponse[[]models.ACMEDirectory]
	if err := c.Get(ctx, "/cluster/acme/directories", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetClusterLog fetches recent cluster log events, up to max entries.
func (c *Client) GetClusterLog(ctx context.Context, max int) ([]models.ClusterLogEntry, error) {
	path := fmt.Sprintf("/cluster/log?max=%d", max)
	var result models.APIResponse[[]models.ClusterLogEntry]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetClusterBackupInfo returns VMs that arent covered by any backup schedule.
func (c *Client) GetClusterBackupInfo(ctx context.Context) ([]models.ClusterBackupInfoEntry, error) {
	var result models.APIResponse[[]models.ClusterBackupInfoEntry]
	if err := c.Get(ctx, "/cluster/backup-info/not-backed-up", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}
