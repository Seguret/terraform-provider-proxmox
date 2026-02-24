package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetReplicationJobs lists all replication jobs configured in the cluster.
func (c *Client) GetReplicationJobs(ctx context.Context) ([]models.ReplicationJob, error) {
	var result models.APIResponse[[]models.ReplicationJob]
	if err := c.Get(ctx, "/cluster/replication", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetReplicationJob fetches a single replication job by its compound ID (vmid-jobid).
func (c *Client) GetReplicationJob(ctx context.Context, id string) (*models.ReplicationJob, error) {
	path := fmt.Sprintf("/cluster/replication/%s", url.PathEscape(id))
	var result models.APIResponse[models.ReplicationJob]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateReplicationJob sets up a new ZFS replication job to a target node.
func (c *Client) CreateReplicationJob(ctx context.Context, req *models.ReplicationJobCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/replication", bytes.NewReader(body))
}

// UpdateReplicationJob updates settings on an existing replication job.
func (c *Client) UpdateReplicationJob(ctx context.Context, id string, req *models.ReplicationJobUpdateRequest) error {
	path := fmt.Sprintf("/cluster/replication/%s", url.PathEscape(id))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteReplicationJob removes a replication job.
func (c *Client) DeleteReplicationJob(ctx context.Context, id string) error {
	path := fmt.Sprintf("/cluster/replication/%s", url.PathEscape(id))
	return c.Delete(ctx, path)
}
