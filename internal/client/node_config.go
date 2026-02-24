package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetNodeConfig fetches the current node-level config (wake-on-lan, description, etc).
func (c *Client) GetNodeConfig(ctx context.Context, node string) (*models.NodeConfig, error) {
	path := fmt.Sprintf("/nodes/%s/config", url.PathEscape(node))
	var result models.APIResponse[models.NodeConfig]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// UpdateNodeConfig applies changes to node-level configuration.
func (c *Client) UpdateNodeConfig(ctx context.Context, node string, req *models.NodeConfigUpdateRequest) error {
	path := fmt.Sprintf("/nodes/%s/config", url.PathEscape(node))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// GetNodeSubscription returns subscription status and expiry info for a node.
func (c *Client) GetNodeSubscription(ctx context.Context, node string) (*models.NodeSubscription, error) {
	path := fmt.Sprintf("/nodes/%s/subscription", url.PathEscape(node))
	var result models.APIResponse[models.NodeSubscription]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// SetNodeSubscription activates a subscription key on a node.
func (c *Client) SetNodeSubscription(ctx context.Context, node, key string) error {
	path := fmt.Sprintf("/nodes/%s/subscription", url.PathEscape(node))
	body, err := json.Marshal(map[string]string{"key": key})
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, path, bytes.NewReader(body))
}

// DeleteNodeSubscription removes the subscription key from a node.
func (c *Client) DeleteNodeSubscription(ctx context.Context, node string) error {
	path := fmt.Sprintf("/nodes/%s/subscription", url.PathEscape(node))
	return c.Delete(ctx, path)
}
