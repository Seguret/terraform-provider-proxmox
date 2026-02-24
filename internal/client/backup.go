package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetBackupJobs returns all configured backup schedules from the cluster.
func (c *Client) GetBackupJobs(ctx context.Context) ([]models.BackupJob, error) {
	var result models.APIResponse[[]models.BackupJob]
	if err := c.Get(ctx, "/cluster/backup", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetBackupJob fetches a single backup job by ID.
func (c *Client) GetBackupJob(ctx context.Context, id string) (*models.BackupJob, error) {
	path := fmt.Sprintf("/cluster/backup/%s", url.PathEscape(id))
	var result models.APIResponse[models.BackupJob]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateBackupJob creates a new backup schedule and returns the ID that Proxmox assigned to it.
func (c *Client) CreateBackupJob(ctx context.Context, req *models.BackupJobCreateRequest) (string, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	var result models.APIResponse[struct {
		ID string `json:"id"`
	}]
	if err := c.Post(ctx, "/cluster/backup", bytes.NewReader(body), &result); err != nil {
		return "", err
	}
	return result.Data.ID, nil
}

// UpdateBackupJob updates an existing backup job schedule.
func (c *Client) UpdateBackupJob(ctx context.Context, id string, req *models.BackupJobUpdateRequest) error {
	path := fmt.Sprintf("/cluster/backup/%s", url.PathEscape(id))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteBackupJob removes a backup schedule by ID.
func (c *Client) DeleteBackupJob(ctx context.Context, id string) error {
	path := fmt.Sprintf("/cluster/backup/%s", url.PathEscape(id))
	return c.Delete(ctx, path)
}
