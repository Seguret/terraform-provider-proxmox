package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// --- ACME Accounts ---

// GetACMEAccounts returns all registered ACME accounts.
func (c *Client) GetACMEAccounts(ctx context.Context) ([]models.ACMEAccountListEntry, error) {
	var result models.APIResponse[[]models.ACMEAccountListEntry]
	if err := c.Get(ctx, "/cluster/acme/account", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetACMEAccount fetches a single ACME account by name.
func (c *Client) GetACMEAccount(ctx context.Context, name string) (*models.ACMEAccount, error) {
	path := fmt.Sprintf("/cluster/acme/account/%s", url.PathEscape(name))
	var result models.APIResponse[models.ACMEAccount]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateACMEAccount registers a new ACME account — async, returns the UPID.
func (c *Client) CreateACMEAccount(ctx context.Context, req *models.ACMEAccountCreateRequest) (string, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	var result models.APIResponse[string]
	if err := c.Post(ctx, "/cluster/acme/account", bytes.NewReader(body), &result); err != nil {
		return "", err
	}
	return result.Data, nil
}

// UpdateACMEAccount updates the contact info on an ACME account. Also async, returns UPID.
func (c *Client) UpdateACMEAccount(ctx context.Context, name string, req *models.ACMEAccountUpdateRequest) (string, error) {
	path := fmt.Sprintf("/cluster/acme/account/%s", url.PathEscape(name))
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

// DeleteACMEAccount deregisters an ACME account and returns the UPID if available.
func (c *Client) DeleteACMEAccount(ctx context.Context, name string) (string, error) {
	path := fmt.Sprintf("/cluster/acme/account/%s", url.PathEscape(name))
	var result models.APIResponse[string]
	if err := c.Get(ctx, path+"?delete=1", &result); err != nil {
		// fall back to a plain DELETE if the get trick didnt work
		if delErr := c.Delete(ctx, path); delErr != nil {
			return "", delErr
		}
		return "", nil
	}
	return result.Data, nil
}

// --- ACME Plugins ---

// GetACMEPlugins returns all configured ACME DNS plugins.
func (c *Client) GetACMEPlugins(ctx context.Context) ([]models.ACMEPlugin, error) {
	var result models.APIResponse[[]models.ACMEPlugin]
	if err := c.Get(ctx, "/cluster/acme/plugins", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetACMEPlugin fetches a single ACME plugin by ID.
func (c *Client) GetACMEPlugin(ctx context.Context, id string) (*models.ACMEPlugin, error) {
	path := fmt.Sprintf("/cluster/acme/plugins/%s", url.PathEscape(id))
	var result models.APIResponse[models.ACMEPlugin]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateACMEPlugin adds a new ACME DNS validation plugin.
func (c *Client) CreateACMEPlugin(ctx context.Context, req *models.ACMEPluginCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/acme/plugins", bytes.NewReader(body))
}

// UpdateACMEPlugin updates an existing ACME plugin configuration.
func (c *Client) UpdateACMEPlugin(ctx context.Context, id string, req *models.ACMEPluginUpdateRequest) error {
	path := fmt.Sprintf("/cluster/acme/plugins/%s", url.PathEscape(id))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteACMEPlugin removes an ACME DNS plugin by ID.
func (c *Client) DeleteACMEPlugin(ctx context.Context, id string) error {
	path := fmt.Sprintf("/cluster/acme/plugins/%s", url.PathEscape(id))
	return c.Delete(ctx, path)
}
