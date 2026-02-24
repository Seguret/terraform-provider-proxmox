package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetDirHardwareMappings lists all directory-type hardware mappings in the cluster.
func (c *Client) GetDirHardwareMappings(ctx context.Context) ([]models.DirHardwareMapping, error) {
	var result models.APIResponse[[]models.DirHardwareMapping]
	if err := c.Get(ctx, "/cluster/mapping/dir", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetDirHardwareMapping fetches a single directory mapping by ID.
func (c *Client) GetDirHardwareMapping(ctx context.Context, id string) (*models.DirHardwareMapping, error) {
	path := fmt.Sprintf("/cluster/mapping/dir/%s", url.PathEscape(id))
	var result models.APIResponse[models.DirHardwareMapping]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateDirHardwareMapping registers a new directory hardware mapping.
func (c *Client) CreateDirHardwareMapping(ctx context.Context, req *models.DirHardwareMappingCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/mapping/dir", bytes.NewReader(body))
}

// UpdateDirHardwareMapping updates an existing directory mapping.
func (c *Client) UpdateDirHardwareMapping(ctx context.Context, id string, req *models.DirHardwareMappingUpdateRequest) error {
	path := fmt.Sprintf("/cluster/mapping/dir/%s", url.PathEscape(id))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteDirHardwareMapping removes a directory mapping from the cluster.
func (c *Client) DeleteDirHardwareMapping(ctx context.Context, id string) error {
	path := fmt.Sprintf("/cluster/mapping/dir/%s", url.PathEscape(id))
	return c.Delete(ctx, path)
}
