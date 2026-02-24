package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetMetricsServers lists all external metrics servers the cluster is pushing to.
func (c *Client) GetMetricsServers(ctx context.Context) ([]models.MetricsServer, error) {
	var result models.APIResponse[[]models.MetricsServer]
	if err := c.Get(ctx, "/cluster/metrics/server", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetMetricsServer fetches a single metrics server config by ID.
func (c *Client) GetMetricsServer(ctx context.Context, id string) (*models.MetricsServer, error) {
	path := fmt.Sprintf("/cluster/metrics/server/%s", url.PathEscape(id))
	var result models.APIResponse[models.MetricsServer]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateMetricsServer adds a new external metrics server target (influxdb, graphite, etc).
func (c *Client) CreateMetricsServer(ctx context.Context, id string, req *models.MetricsServerCreateRequest) error {
	path := fmt.Sprintf("/cluster/metrics/server/%s", url.PathEscape(id))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, path, bytes.NewReader(body))
}

// UpdateMetricsServer updates an existing metrics server configuration.
func (c *Client) UpdateMetricsServer(ctx context.Context, id string, req *models.MetricsServerUpdateRequest) error {
	path := fmt.Sprintf("/cluster/metrics/server/%s", url.PathEscape(id))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteMetricsServer removes a metrics push target from the cluster.
func (c *Client) DeleteMetricsServer(ctx context.Context, id string) error {
	path := fmt.Sprintf("/cluster/metrics/server/%s", url.PathEscape(id))
	return c.Delete(ctx, path)
}
