package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// --- Node DNS ---

// GetNodeDNS returns the DNS resolver config for a node.
func (c *Client) GetNodeDNS(ctx context.Context, node string) (*models.NodeDNS, error) {
	path := fmt.Sprintf("/nodes/%s/dns", url.PathEscape(node))
	var result models.APIResponse[models.NodeDNS]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// UpdateNodeDNS changes the DNS server settings on a node.
func (c *Client) UpdateNodeDNS(ctx context.Context, node string, req *models.NodeDNSUpdateRequest) error {
	path := fmt.Sprintf("/nodes/%s/dns", url.PathEscape(node))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// --- Node Hosts ---

// GetNodeHosts reads the /etc/hosts content from a node.
func (c *Client) GetNodeHosts(ctx context.Context, node string) (*models.NodeHosts, error) {
	path := fmt.Sprintf("/nodes/%s/hosts", url.PathEscape(node))
	var result models.APIResponse[models.NodeHosts]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// UpdateNodeHosts overwrites the /etc/hosts file on a node with new content.
func (c *Client) UpdateNodeHosts(ctx context.Context, node string, req *models.NodeHostsUpdateRequest) error {
	path := fmt.Sprintf("/nodes/%s/hosts", url.PathEscape(node))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, path, bytes.NewReader(body))
}

// --- Node Certificates ---

// GetNodeCertificates returns all TLS certs installed on a node.
func (c *Client) GetNodeCertificates(ctx context.Context, node string) ([]models.NodeCertificate, error) {
	path := fmt.Sprintf("/nodes/%s/certificates/info", url.PathEscape(node))
	var result models.APIResponse[[]models.NodeCertificate]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// UploadNodeCertificate installs a custom TLS certificate on a node.
func (c *Client) UploadNodeCertificate(ctx context.Context, node string, req *models.NodeCertificateUploadRequest) (*models.NodeCertificate, error) {
	path := fmt.Sprintf("/nodes/%s/certificates/custom", url.PathEscape(node))
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var result models.APIResponse[models.NodeCertificate]
	if err := c.Post(ctx, path, bytes.NewReader(body), &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// DeleteNodeCertificate removes the custom TLS cert from a node, falling back to self-signed.
func (c *Client) DeleteNodeCertificate(ctx context.Context, node string) error {
	path := fmt.Sprintf("/nodes/%s/certificates/custom", url.PathEscape(node))
	return c.Delete(ctx, path)
}

// --- Node Time ---

// GetNodeTime returns the current time and configured timezone on a node.
func (c *Client) GetNodeTime(ctx context.Context, node string) (map[string]any, error) {
	path := fmt.Sprintf("/nodes/%s/time", url.PathEscape(node))
	var result models.APIResponse[map[string]any]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// SetNodeTimezone changes the timezone on a node (e.g. "Europe/Berlin").
func (c *Client) SetNodeTimezone(ctx context.Context, node, timezone string) error {
	path := fmt.Sprintf("/nodes/%s/time", url.PathEscape(node))
	body, err := json.Marshal(map[string]string{"timezone": timezone})
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}
