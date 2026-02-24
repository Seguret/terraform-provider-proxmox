package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// --- PCI Hardware Mappings ---

// GetPCIHardwareMappings lists all PCI hardware mappings defined at the cluster level.
func (c *Client) GetPCIHardwareMappings(ctx context.Context) ([]models.PCIHardwareMapping, error) {
	var result models.APIResponse[[]models.PCIHardwareMapping]
	if err := c.Get(ctx, "/cluster/mapping/pci", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetPCIHardwareMapping fetches a single PCI hardware mapping by ID.
func (c *Client) GetPCIHardwareMapping(ctx context.Context, id string) (*models.PCIHardwareMapping, error) {
	path := fmt.Sprintf("/cluster/mapping/pci/%s", url.PathEscape(id))
	var result models.APIResponse[models.PCIHardwareMapping]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreatePCIHardwareMapping creates a new cluster-level PCI mapping.
func (c *Client) CreatePCIHardwareMapping(ctx context.Context, req *models.PCIHardwareMappingCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/mapping/pci", bytes.NewReader(body))
}

// UpdatePCIHardwareMapping updates an existing PCI hardware mapping.
func (c *Client) UpdatePCIHardwareMapping(ctx context.Context, id string, req *models.PCIHardwareMappingUpdateRequest) error {
	path := fmt.Sprintf("/cluster/mapping/pci/%s", url.PathEscape(id))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeletePCIHardwareMapping removes a PCI hardware mapping.
func (c *Client) DeletePCIHardwareMapping(ctx context.Context, id string) error {
	path := fmt.Sprintf("/cluster/mapping/pci/%s", url.PathEscape(id))
	return c.Delete(ctx, path)
}

// --- USB Hardware Mappings ---

// GetUSBHardwareMappings lists all USB hardware mappings at the cluster level.
func (c *Client) GetUSBHardwareMappings(ctx context.Context) ([]models.USBHardwareMapping, error) {
	var result models.APIResponse[[]models.USBHardwareMapping]
	if err := c.Get(ctx, "/cluster/mapping/usb", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetUSBHardwareMapping fetches a single USB hardware mapping by ID.
func (c *Client) GetUSBHardwareMapping(ctx context.Context, id string) (*models.USBHardwareMapping, error) {
	path := fmt.Sprintf("/cluster/mapping/usb/%s", url.PathEscape(id))
	var result models.APIResponse[models.USBHardwareMapping]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateUSBHardwareMapping adds a new USB device mapping at the cluster level.
func (c *Client) CreateUSBHardwareMapping(ctx context.Context, req *models.USBHardwareMappingCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/mapping/usb", bytes.NewReader(body))
}

// UpdateUSBHardwareMapping modifies an existing USB hardware mapping.
func (c *Client) UpdateUSBHardwareMapping(ctx context.Context, id string, req *models.USBHardwareMappingUpdateRequest) error {
	path := fmt.Sprintf("/cluster/mapping/usb/%s", url.PathEscape(id))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteUSBHardwareMapping removes a USB hardware mapping from the cluster.
func (c *Client) DeleteUSBHardwareMapping(ctx context.Context, id string) error {
	path := fmt.Sprintf("/cluster/mapping/usb/%s", url.PathEscape(id))
	return c.Delete(ctx, path)
}
