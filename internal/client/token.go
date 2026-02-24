package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetUserTokens lists all API tokens belonging to a user.
func (c *Client) GetUserTokens(ctx context.Context, userID string) ([]models.UserToken, error) {
	path := fmt.Sprintf("/access/users/%s/token", url.PathEscape(userID))
	var result models.APIResponse[[]models.UserToken]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetUserToken fetches a specific API token by its ID. Note: the secret value is not returned here.
func (c *Client) GetUserToken(ctx context.Context, userID, tokenID string) (*models.UserToken, error) {
	path := fmt.Sprintf("/access/users/%s/token/%s", url.PathEscape(userID), url.PathEscape(tokenID))
	var result models.APIResponse[models.UserToken]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	result.Data.TokenID = tokenID
	return &result.Data, nil
}

// CreateUserToken creates a new API token for a user. The secret is only returned once on creation.
func (c *Client) CreateUserToken(ctx context.Context, userID, tokenID string, req *models.UserTokenCreateRequest) (*models.UserTokenCreateResponse, error) {
	path := fmt.Sprintf("/access/users/%s/token/%s", url.PathEscape(userID), url.PathEscape(tokenID))
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var result models.APIResponse[models.UserTokenCreateResponse]
	if err := c.Post(ctx, path, bytes.NewReader(body), &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// UpdateUserToken updates metadata on an existing API token (comment, expiry, privileges).
func (c *Client) UpdateUserToken(ctx context.Context, userID, tokenID string, req *models.UserTokenUpdateRequest) error {
	path := fmt.Sprintf("/access/users/%s/token/%s", url.PathEscape(userID), url.PathEscape(tokenID))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteUserToken revokes and removes an API token.
func (c *Client) DeleteUserToken(ctx context.Context, userID, tokenID string) error {
	path := fmt.Sprintf("/access/users/%s/token/%s", url.PathEscape(userID), url.PathEscape(tokenID))
	return c.Delete(ctx, path)
}
