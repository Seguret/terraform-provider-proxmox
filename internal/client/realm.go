package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetRealms lists all authentication realms (pve, pam, ad, ldap, etc).
func (c *Client) GetRealms(ctx context.Context) ([]models.AuthRealm, error) {
	var result models.APIResponse[[]models.AuthRealm]
	if err := c.Get(ctx, "/access/domains", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetRealm fetches a single auth realm by name.
func (c *Client) GetRealm(ctx context.Context, realm string) (*models.AuthRealm, error) {
	path := fmt.Sprintf("/access/domains/%s", url.PathEscape(realm))
	var result models.APIResponse[models.AuthRealm]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	result.Data.Realm = realm
	return &result.Data, nil
}

// CreateRealm registers a new authentication realm.
func (c *Client) CreateRealm(ctx context.Context, req *models.AuthRealmCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/access/domains", bytes.NewReader(body))
}

// UpdateRealm modifies an existing auth realm's configuration.
func (c *Client) UpdateRealm(ctx context.Context, realm string, req *models.AuthRealmUpdateRequest) error {
	path := fmt.Sprintf("/access/domains/%s", url.PathEscape(realm))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteRealm removes an authentication realm.
func (c *Client) DeleteRealm(ctx context.Context, realm string) error {
	path := fmt.Sprintf("/access/domains/%s", url.PathEscape(realm))
	return c.Delete(ctx, path)
}
