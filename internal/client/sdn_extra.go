package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// --- SDN Controllers ---

// GetSDNControllers lists all SDN routing controllers (e.g. EVPN BGP controllers).
func (c *Client) GetSDNControllers(ctx context.Context) ([]models.SDNController, error) {
	var result models.APIResponse[[]models.SDNController]
	if err := c.Get(ctx, "/cluster/sdn/controllers", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetSDNController fetches a single SDN controller config by name.
func (c *Client) GetSDNController(ctx context.Context, controller string) (*models.SDNController, error) {
	path := fmt.Sprintf("/cluster/sdn/controllers/%s", url.PathEscape(controller))
	var result models.APIResponse[models.SDNController]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateSDNController adds a new SDN routing controller.
func (c *Client) CreateSDNController(ctx context.Context, req *models.SDNControllerCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/sdn/controllers", bytes.NewReader(body))
}

// UpdateSDNController updates an existing SDN controller.
func (c *Client) UpdateSDNController(ctx context.Context, controller string, req *models.SDNControllerUpdateRequest) error {
	path := fmt.Sprintf("/cluster/sdn/controllers/%s", url.PathEscape(controller))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteSDNController removes an SDN controller.
func (c *Client) DeleteSDNController(ctx context.Context, controller string) error {
	path := fmt.Sprintf("/cluster/sdn/controllers/%s", url.PathEscape(controller))
	return c.Delete(ctx, path)
}

// --- SDN DNS ---

// GetSDNDnsList returns all SDN DNS integration providers.
func (c *Client) GetSDNDnsList(ctx context.Context) ([]models.SDNDns, error) {
	var result models.APIResponse[[]models.SDNDns]
	if err := c.Get(ctx, "/cluster/sdn/dns", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetSDNDns fetches a single SDN DNS provider config.
func (c *Client) GetSDNDns(ctx context.Context, dns string) (*models.SDNDns, error) {
	path := fmt.Sprintf("/cluster/sdn/dns/%s", url.PathEscape(dns))
	var result models.APIResponse[models.SDNDns]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateSDNDns adds a new SDN DNS provider to the cluster.
func (c *Client) CreateSDNDns(ctx context.Context, req *models.SDNDnsCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/sdn/dns", bytes.NewReader(body))
}

// UpdateSDNDns updates the config for an existing SDN DNS provider.
func (c *Client) UpdateSDNDns(ctx context.Context, dns string, req *models.SDNDnsUpdateRequest) error {
	path := fmt.Sprintf("/cluster/sdn/dns/%s", url.PathEscape(dns))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteSDNDns removes an SDN DNS provider.
func (c *Client) DeleteSDNDns(ctx context.Context, dns string) error {
	path := fmt.Sprintf("/cluster/sdn/dns/%s", url.PathEscape(dns))
	return c.Delete(ctx, path)
}

// --- SDN IPAM ---

// GetSDNIpams lists all SDN IPAM providers configured in the cluster.
func (c *Client) GetSDNIpams(ctx context.Context) ([]models.SDNIpam, error) {
	var result models.APIResponse[[]models.SDNIpam]
	if err := c.Get(ctx, "/cluster/sdn/ipams", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetSDNIpam fetches a single IPAM provider by name.
func (c *Client) GetSDNIpam(ctx context.Context, ipam string) (*models.SDNIpam, error) {
	path := fmt.Sprintf("/cluster/sdn/ipams/%s", url.PathEscape(ipam))
	var result models.APIResponse[models.SDNIpam]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateSDNIpam adds a new IPAM provider to the SDN config.
func (c *Client) CreateSDNIpam(ctx context.Context, req *models.SDNIpamCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/sdn/ipams", bytes.NewReader(body))
}

// UpdateSDNIpam updates settings on an exsiting IPAM provider.
func (c *Client) UpdateSDNIpam(ctx context.Context, ipam string, req *models.SDNIpamUpdateRequest) error {
	path := fmt.Sprintf("/cluster/sdn/ipams/%s", url.PathEscape(ipam))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteSDNIpam removes an IPAM provider from the SDN config.
func (c *Client) DeleteSDNIpam(ctx context.Context, ipam string) error {
	path := fmt.Sprintf("/cluster/sdn/ipams/%s", url.PathEscape(ipam))
	return c.Delete(ctx, path)
}
