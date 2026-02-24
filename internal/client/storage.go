package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetStorageConfigs lists all global storage definitions configured on the cluster.
func (c *Client) GetStorageConfigs(ctx context.Context) ([]models.StorageConfig, error) {
	var result models.APIResponse[[]models.StorageConfig]
	if err := c.Get(ctx, "/storage", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetStorageConfig fetches the config for a single storage backend by ID.
func (c *Client) GetStorageConfig(ctx context.Context, storageID string) (*models.StorageConfig, error) {
	path := fmt.Sprintf("/storage/%s", url.PathEscape(storageID))
	var result models.APIResponse[models.StorageConfig]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateStorage adds a new storage backend definition to the cluster.
func (c *Client) CreateStorage(ctx context.Context, req *models.StorageCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/storage", bytes.NewReader(body))
}

// UpdateStorage modifies an existing storage backend's settings.
func (c *Client) UpdateStorage(ctx context.Context, storageID string, req *models.StorageUpdateRequest) error {
	path := fmt.Sprintf("/storage/%s", url.PathEscape(storageID))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteStorage removes a storage backend definition. Doesnt destroy the actual data.
func (c *Client) DeleteStorage(ctx context.Context, storageID string) error {
	path := fmt.Sprintf("/storage/%s", url.PathEscape(storageID))
	return c.Delete(ctx, path)
}
