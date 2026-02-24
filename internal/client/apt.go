package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetAptRepositories fetches all configured APT repos for a node.
func (c *Client) GetAptRepositories(ctx context.Context, node string) (*models.AptRepositoriesResponse, error) {
	path := fmt.Sprintf("/nodes/%s/apt/repositories", url.PathEscape(node))
	var result models.APIResponse[models.AptRepositoriesResponse]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// AddAptRepository adds a standard APT repo to a node.
func (c *Client) AddAptRepository(ctx context.Context, node string, req *models.AptRepositoryAddRequest) error {
	path := fmt.Sprintf("/nodes/%s/apt/repositories", url.PathEscape(node))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, path, bytes.NewReader(body))
}

// ChangeAptRepository updates an existing APT repository entry (enable/disable etc).
func (c *Client) ChangeAptRepository(ctx context.Context, node string, req *models.AptRepositoryChangeRequest) error {
	path := fmt.Sprintf("/nodes/%s/apt/repositories", url.PathEscape(node))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}
