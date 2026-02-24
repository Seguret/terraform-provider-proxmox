package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetClusterJobsRealmSync lists all realm-sync jobs configured in the cluster.
func (c *Client) GetClusterJobsRealmSync(ctx context.Context) ([]models.ClusterJobRealmSync, error) {
	var result models.APIResponse[[]models.ClusterJobRealmSync]
	if err := c.Get(ctx, "/cluster/jobs/realm-sync", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetClusterJobRealmSync fetches a single realm-sync job by ID.
func (c *Client) GetClusterJobRealmSync(ctx context.Context, id string) (*models.ClusterJobRealmSync, error) {
	path := fmt.Sprintf("/cluster/jobs/realm-sync/%s", url.PathEscape(id))
	var result models.APIResponse[models.ClusterJobRealmSync]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateClusterJobRealmSync sets up a new realm-sync job to periodically sync users from a realm.
func (c *Client) CreateClusterJobRealmSync(ctx context.Context, id string, req *models.ClusterJobRealmSyncCreateRequest) error {
	path := fmt.Sprintf("/cluster/jobs/realm-sync/%s", url.PathEscape(id))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, path, bytes.NewReader(body))
}

// UpdateClusterJobRealmSync updates settings on an existing realm-sync job.
func (c *Client) UpdateClusterJobRealmSync(ctx context.Context, id string, req *models.ClusterJobRealmSyncUpdateRequest) error {
	path := fmt.Sprintf("/cluster/jobs/realm-sync/%s", url.PathEscape(id))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteClusterJobRealmSync removes a realm-sync job from the cluster.
func (c *Client) DeleteClusterJobRealmSync(ctx context.Context, id string) error {
	path := fmt.Sprintf("/cluster/jobs/realm-sync/%s", url.PathEscape(id))
	return c.Delete(ctx, path)
}
