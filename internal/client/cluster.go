package client

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// --- Cluster Options ---

// GetClusterOptions fetches the global cluster configuration options.
func (c *Client) GetClusterOptions(ctx context.Context) (*models.ClusterOptions, error) {
	var result models.APIResponse[models.ClusterOptions]
	if err := c.Get(ctx, "/cluster/options", &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// UpdateClusterOptions applies changes to cluster-wide settings.
func (c *Client) UpdateClusterOptions(ctx context.Context, req *models.ClusterOptionsUpdateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, "/cluster/options", bytes.NewReader(body))
}

// --- Cluster Config ---

// GetClusterConfig returns the cluster config (name, nodes, etc).
func (c *Client) GetClusterConfig(ctx context.Context) (*models.ClusterConfig, error) {
	var result models.APIResponse[models.ClusterConfig]
	if err := c.Get(ctx, "/cluster/config", &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// --- Hardware Mapping ---

// GetHardwareMappingPCI returns all PCI hardware mappings defined in the cluster.
func (c *Client) GetHardwareMappingPCI(ctx context.Context) ([]models.HardwareMappingPCI, error) {
	var result models.APIResponse[[]models.HardwareMappingPCI]
	if err := c.Get(ctx, "/cluster/mapping/pci", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetHardwareMappingUSB returns all USB hardware mappings defined in the cluster.
func (c *Client) GetHardwareMappingUSB(ctx context.Context) ([]models.HardwareMappingUSB, error) {
	var result models.APIResponse[[]models.HardwareMappingUSB]
	if err := c.Get(ctx, "/cluster/mapping/usb", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// --- Ceph Cluster (Read-only) ---

// GetCephStatus grabs the overall Ceph cluster status from the cluster endpoint.
func (c *Client) GetCephStatus(ctx context.Context) (*models.CephStatus, error) {
	var result models.APIResponse[models.CephStatus]
	if err := c.Get(ctx, "/cluster/ceph/status", &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}
