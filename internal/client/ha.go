package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// --- HA Resources ---

// GetHAResources returns all resources managed by the HA subsystem.
func (c *Client) GetHAResources(ctx context.Context) ([]models.HAResource, error) {
	var result models.APIResponse[[]models.HAResource]
	if err := c.Get(ctx, "/cluster/ha/resources", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetHAResource fetches a single HA-managed resource by its service ID.
func (c *Client) GetHAResource(ctx context.Context, sid string) (*models.HAResource, error) {
	path := fmt.Sprintf("/cluster/ha/resources/%s", url.PathEscape(sid))
	var result models.APIResponse[models.HAResource]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateHAResource registers a VM or container with the HA manager.
func (c *Client) CreateHAResource(ctx context.Context, req *models.HAResourceCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/ha/resources", bytes.NewReader(body))
}

// UpdateHAResource changes settings on an existing HA resource.
func (c *Client) UpdateHAResource(ctx context.Context, sid string, req *models.HAResourceUpdateRequest) error {
	path := fmt.Sprintf("/cluster/ha/resources/%s", url.PathEscape(sid))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteHAResource unregisters a resource from the HA manager.
func (c *Client) DeleteHAResource(ctx context.Context, sid string) error {
	path := fmt.Sprintf("/cluster/ha/resources/%s", url.PathEscape(sid))
	return c.Delete(ctx, path)
}

// --- HA Groups ---

// GetHAGroups lists all HA groups configured on the cluster.
func (c *Client) GetHAGroups(ctx context.Context) ([]models.HAGroup, error) {
	var result models.APIResponse[[]models.HAGroup]
	if err := c.Get(ctx, "/cluster/ha/groups", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetHAGroup fetches a single HA group by name.
func (c *Client) GetHAGroup(ctx context.Context, group string) (*models.HAGroup, error) {
	path := fmt.Sprintf("/cluster/ha/groups/%s", url.PathEscape(group))
	var result models.APIResponse[models.HAGroup]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateHAGroup creates a new HA group with a set of member nodes.
func (c *Client) CreateHAGroup(ctx context.Context, req *models.HAGroupCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/ha/groups", bytes.NewReader(body))
}

// UpdateHAGroup modifies an existing HA group's configuration.
func (c *Client) UpdateHAGroup(ctx context.Context, group string, req *models.HAGroupUpdateRequest) error {
	path := fmt.Sprintf("/cluster/ha/groups/%s", url.PathEscape(group))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteHAGroup removes an HA group from the cluster.
func (c *Client) DeleteHAGroup(ctx context.Context, group string) error {
	path := fmt.Sprintf("/cluster/ha/groups/%s", url.PathEscape(group))
	return c.Delete(ctx, path)
}
