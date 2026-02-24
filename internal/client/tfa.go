package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// ListUserTFA returns all two-factor auth entries registered for a user.
func (c *Client) ListUserTFA(ctx context.Context, userid string) ([]models.TFAEntry, error) {
	path := fmt.Sprintf("/access/users/%s/tfa", url.PathEscape(userid))
	var result models.APIResponse[[]models.TFAEntry]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetTFAEntry looks up a single TFA entry by its ID.
// Proxmox doesnt have a single-entry GET endpoint, so we list all and filter locally.
func (c *Client) GetTFAEntry(ctx context.Context, userid, id string) (*models.TFAEntry, error) {
	entries, err := c.ListUserTFA(ctx, userid)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, &APIError{StatusCode: 404, Status: "404 Not Found", Message: "TFA entry not found"}
}

// CreateTFAEntry registers a new 2FA device or secret for the given user.
func (c *Client) CreateTFAEntry(ctx context.Context, userid string, req *models.TFACreateRequest) (*models.TFACreateResponse, error) {
	path := fmt.Sprintf("/access/users/%s/tfa", url.PathEscape(userid))
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var result models.APIResponse[models.TFACreateResponse]
	if err := c.Post(ctx, path, bytes.NewReader(body), &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// UpdateTFAEntry changes the description or enabled/disabled state of a TFA entry.
func (c *Client) UpdateTFAEntry(ctx context.Context, userid, id string, req *models.TFAUpdateRequest) error {
	path := fmt.Sprintf("/access/users/%s/tfa/%s", url.PathEscape(userid), url.PathEscape(id))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteTFAEntry removes a 2FA entry from a user's account.
func (c *Client) DeleteTFAEntry(ctx context.Context, userid, id string) error {
	path := fmt.Sprintf("/access/users/%s/tfa/%s", url.PathEscape(userid), url.PathEscape(id))
	return c.Delete(ctx, path)
}
