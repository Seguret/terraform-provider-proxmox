package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// DownloadFile tells a node to download a file from a URL directly into storage.
// The actual download happens async - returns the UPID to track it.
func (c *Client) DownloadFile(ctx context.Context, node, storage string, req *models.DownloadURLRequest) (string, error) {
	path := fmt.Sprintf("/nodes/%s/storage/%s/download-url", url.PathEscape(node), url.PathEscape(storage))
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	var result models.APIResponse[string]
	if err := c.Post(ctx, path, bytes.NewReader(body), &result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// GetStorageContent lists files/volumes in a given node storage.
func (c *Client) GetStorageContent(ctx context.Context, node, storage string) ([]models.StorageContent, error) {
	path := fmt.Sprintf("/nodes/%s/storage/%s/content", url.PathEscape(node), url.PathEscape(storage))
	var result models.APIResponse[[]models.StorageContent]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// DeleteStorageContent removes a volume from node storage. Returns a UPID since this runs async.
func (c *Client) DeleteStorageContent(ctx context.Context, node, storage, volid string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/storage/%s/content/%s", url.PathEscape(node), url.PathEscape(storage), url.PathEscape(volid))
	resp, err := c.DoRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", c.parseError(resp)
	}
	var result models.APIResponse[string]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Data, nil
}
