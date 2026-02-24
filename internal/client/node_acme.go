package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// OrderNodeACMECertificate requests a new ACME cert for a node.
// Set force=true to re-order even if a valid cert already exists. Returns the task UPID.
func (c *Client) OrderNodeACMECertificate(ctx context.Context, node string, force bool) (string, error) {
	path := fmt.Sprintf("/nodes/%s/certificates/acme/certificate", url.PathEscape(node))
	req := models.ACMECertOrderRequest{Force: force}
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

// RenewNodeACMECertificate triggers an ACME cert renewal on a node.
// Returns the UPID of the async task.
func (c *Client) RenewNodeACMECertificate(ctx context.Context, node string, force bool) (string, error) {
	path := fmt.Sprintf("/nodes/%s/certificates/acme/certificate", url.PathEscape(node))
	req := models.ACMERenewRequest{Force: force}
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	resp, err := c.DoRequest(ctx, "PUT", path, bytes.NewReader(body))
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

// RevokeNodeACMECertificate revokes and removes the ACME cert from a node.
// Returns the UPID, which may be empty if the operation finishes synchronously.
func (c *Client) RevokeNodeACMECertificate(ctx context.Context, node string) (string, error) {
	path := fmt.Sprintf("/nodes/%s/certificates/acme/certificate", url.PathEscape(node))
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
		// body might be empty on delete, just swallow the decode error
		return "", nil
	}
	return result.Data, nil
}
