package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// ApplySDN pushes any pending SDN config changes live — call this after any SDN modification.
func (c *Client) ApplySDN(ctx context.Context) error {
	return c.Put(ctx, "/cluster/sdn", nil)
}

// --- SDN Zones ---

// GetSDNZones lists all SDN zones configured in the cluster.
func (c *Client) GetSDNZones(ctx context.Context) ([]models.SDNZone, error) {
	var result models.APIResponse[[]models.SDNZone]
	if err := c.Get(ctx, "/cluster/sdn/zones", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetSDNZone fetches a single SDN zone by name.
func (c *Client) GetSDNZone(ctx context.Context, zone string) (*models.SDNZone, error) {
	path := fmt.Sprintf("/cluster/sdn/zones/%s", url.PathEscape(zone))
	var result models.APIResponse[models.SDNZone]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateSDNZone adds a new SDN zone (vlan, vxlan, evpn, etc).
func (c *Client) CreateSDNZone(ctx context.Context, req *models.SDNZoneCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/sdn/zones", bytes.NewReader(body))
}

// UpdateSDNZone updates an existing SDN zone's configuration.
func (c *Client) UpdateSDNZone(ctx context.Context, zone string, req *models.SDNZoneUpdateRequest) error {
	path := fmt.Sprintf("/cluster/sdn/zones/%s", url.PathEscape(zone))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteSDNZone removes an SDN zone. Remember to call ApplySDN after.
func (c *Client) DeleteSDNZone(ctx context.Context, zone string) error {
	path := fmt.Sprintf("/cluster/sdn/zones/%s", url.PathEscape(zone))
	return c.Delete(ctx, path)
}

// --- SDN VNets ---

// GetSDNVnets lists all virtual networks defined across SDN zones.
func (c *Client) GetSDNVnets(ctx context.Context) ([]models.SDNVnet, error) {
	var result models.APIResponse[[]models.SDNVnet]
	if err := c.Get(ctx, "/cluster/sdn/vnets", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetSDNVnet fetches a single virtual network by name.
func (c *Client) GetSDNVnet(ctx context.Context, vnet string) (*models.SDNVnet, error) {
	path := fmt.Sprintf("/cluster/sdn/vnets/%s", url.PathEscape(vnet))
	var result models.APIResponse[models.SDNVnet]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateSDNVnet creates a new virtual network in the SDN.
func (c *Client) CreateSDNVnet(ctx context.Context, req *models.SDNVnetCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/sdn/vnets", bytes.NewReader(body))
}

// UpdateSDNVnet updates a virtual network's configuration.
func (c *Client) UpdateSDNVnet(ctx context.Context, vnet string, req *models.SDNVnetUpdateRequest) error {
	path := fmt.Sprintf("/cluster/sdn/vnets/%s", url.PathEscape(vnet))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteSDNVnet removes a virtual network from the SDN config.
func (c *Client) DeleteSDNVnet(ctx context.Context, vnet string) error {
	path := fmt.Sprintf("/cluster/sdn/vnets/%s", url.PathEscape(vnet))
	return c.Delete(ctx, path)
}

// --- SDN Subnets ---

// GetSDNSubnets returns all subnets attached to a specific VNet.
func (c *Client) GetSDNSubnets(ctx context.Context, vnet string) ([]models.SDNSubnet, error) {
	path := fmt.Sprintf("/cluster/sdn/vnets/%s/subnets", url.PathEscape(vnet))
	var result models.APIResponse[[]models.SDNSubnet]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetSDNSubnet fetches a specific subnet from a VNet.
// The subnet param is a CIDR in URL-safe form - use dashes not slashes (e.g. "10.0.0.0-24").
func (c *Client) GetSDNSubnet(ctx context.Context, vnet, subnet string) (*models.SDNSubnet, error) {
	path := fmt.Sprintf("/cluster/sdn/vnets/%s/subnets/%s", url.PathEscape(vnet), url.PathEscape(subnet))
	var result models.APIResponse[models.SDNSubnet]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateSDNSubnet adds a new subnet to a VNet.
func (c *Client) CreateSDNSubnet(ctx context.Context, vnet string, req *models.SDNSubnetCreateRequest) error {
	path := fmt.Sprintf("/cluster/sdn/vnets/%s/subnets", url.PathEscape(vnet))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, path, bytes.NewReader(body))
}

// UpdateSDNSubnet updates an SDN subnet (gateway, snat settings, etc).
func (c *Client) UpdateSDNSubnet(ctx context.Context, vnet, subnet string, req *models.SDNSubnetUpdateRequest) error {
	path := fmt.Sprintf("/cluster/sdn/vnets/%s/subnets/%s", url.PathEscape(vnet), url.PathEscape(subnet))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteSDNSubnet removes a subnet from a VNet.
func (c *Client) DeleteSDNSubnet(ctx context.Context, vnet, subnet string) error {
	path := fmt.Sprintf("/cluster/sdn/vnets/%s/subnets/%s", url.PathEscape(vnet), url.PathEscape(subnet))
	return c.Delete(ctx, path)
}
