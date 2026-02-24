package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// --- Users ---

func (c *Client) GetUsers(ctx context.Context) ([]models.User, error) {
	var result models.APIResponse[[]models.User]
	if err := c.Get(ctx, "/access/users", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) GetUser(ctx context.Context, userID string) (*models.User, error) {
	path := fmt.Sprintf("/access/users/%s", url.PathEscape(userID))
	var result models.APIResponse[models.User]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) CreateUser(ctx context.Context, req *models.UserCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/access/users", bytes.NewReader(body))
}

func (c *Client) UpdateUser(ctx context.Context, userID string, req *models.UserUpdateRequest) error {
	path := fmt.Sprintf("/access/users/%s", url.PathEscape(userID))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	path := fmt.Sprintf("/access/users/%s", url.PathEscape(userID))
	return c.Delete(ctx, path)
}

// --- Groups ---

func (c *Client) GetGroups(ctx context.Context) ([]models.Group, error) {
	var result models.APIResponse[[]models.Group]
	if err := c.Get(ctx, "/access/groups", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) GetGroup(ctx context.Context, groupID string) (*models.Group, error) {
	path := fmt.Sprintf("/access/groups/%s", url.PathEscape(groupID))
	var result models.APIResponse[models.Group]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) CreateGroup(ctx context.Context, req *models.GroupCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/access/groups", bytes.NewReader(body))
}

func (c *Client) UpdateGroup(ctx context.Context, groupID string, req *models.GroupUpdateRequest) error {
	path := fmt.Sprintf("/access/groups/%s", url.PathEscape(groupID))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

func (c *Client) DeleteGroup(ctx context.Context, groupID string) error {
	path := fmt.Sprintf("/access/groups/%s", url.PathEscape(groupID))
	return c.Delete(ctx, path)
}

// --- Roles ---

func (c *Client) GetRoles(ctx context.Context) ([]models.Role, error) {
	var result models.APIResponse[[]models.Role]
	if err := c.Get(ctx, "/access/roles", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) GetRole(ctx context.Context, roleID string) (*models.Role, error) {
	path := fmt.Sprintf("/access/roles/%s", url.PathEscape(roleID))
	var result models.APIResponse[models.Role]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) CreateRole(ctx context.Context, req *models.RoleCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/access/roles", bytes.NewReader(body))
}

func (c *Client) UpdateRole(ctx context.Context, roleID string, req *models.RoleUpdateRequest) error {
	path := fmt.Sprintf("/access/roles/%s", url.PathEscape(roleID))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

func (c *Client) DeleteRole(ctx context.Context, roleID string) error {
	path := fmt.Sprintf("/access/roles/%s", url.PathEscape(roleID))
	return c.Delete(ctx, path)
}

// --- ACL ---

func (c *Client) GetACL(ctx context.Context) ([]models.ACLEntry, error) {
	var result models.APIResponse[[]models.ACLEntry]
	if err := c.Get(ctx, "/access/acl", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) UpdateACL(ctx context.Context, req *models.ACLUpdateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, "/access/acl", bytes.NewReader(body))
}

// --- Pools ---

func (c *Client) GetPools(ctx context.Context) ([]models.Pool, error) {
	var result models.APIResponse[[]models.Pool]
	if err := c.Get(ctx, "/pools", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) GetPool(ctx context.Context, poolID string) (*models.Pool, error) {
	path := fmt.Sprintf("/pools/%s", url.PathEscape(poolID))
	var result models.APIResponse[models.Pool]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) CreatePool(ctx context.Context, req *models.PoolCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/pools", bytes.NewReader(body))
}

func (c *Client) UpdatePool(ctx context.Context, poolID string, req *models.PoolUpdateRequest) error {
	path := fmt.Sprintf("/pools/%s", url.PathEscape(poolID))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

func (c *Client) DeletePool(ctx context.Context, poolID string) error {
	path := fmt.Sprintf("/pools/%s", url.PathEscape(poolID))
	return c.Delete(ctx, path)
}

// --- Password ---

// ChangeUserPassword updates the password for a given user account.
func (c *Client) ChangeUserPassword(ctx context.Context, userID, password string) error {
	path := fmt.Sprintf("/access/password")
	
	body, err := json.Marshal(map[string]string{
		"userid":   userID,
		"password": password,
	})
	if err != nil {
		return err
	}
	
	return c.Put(ctx, path, bytes.NewReader(body))
}

// --- Permissions ---

// GetUserPermissions fetches effective permissions for a user, optionally scoped to a path.
func (c *Client) GetUserPermissions(ctx context.Context, userID, path string) (*models.UserPermissions, error) {
	apiPath := fmt.Sprintf("/access/permissions?userid=%s", url.QueryEscape(userID))
	if path != "" {
		apiPath += fmt.Sprintf("&path=%s", url.QueryEscape(path))
	}
	
	var result models.APIResponse[models.UserPermissions]
	if err := c.Get(ctx, apiPath, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// --- OpenID ---

// GetOpenIDConfig fetches the OpenID Connect configuration for a realm.
func (c *Client) GetOpenIDConfig(ctx context.Context, realm string) (*models.OpenIDConfig, error) {
	path := fmt.Sprintf("/access/openid/%s", url.PathEscape(realm))
	var result models.APIResponse[models.OpenIDConfig]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// AuthorizeOpenID kicks off the OIDC authorization flow and returns the auth URL.
func (c *Client) AuthorizeOpenID(ctx context.Context, realm, redirectURL string) (*models.OpenIDAuthResponse, error) {
	path := fmt.Sprintf("/access/openid/auth-url")
	
	body, err := json.Marshal(map[string]string{
		"realm":        realm,
		"redirect-url": redirectURL,
	})
	if err != nil {
		return nil, err
	}
	
	var result models.APIResponse[models.OpenIDAuthResponse]
	if err := c.Post(ctx, path, bytes.NewReader(body), &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// LoginOpenID completes the OIDC flow using the authorization code from the callback.
func (c *Client) LoginOpenID(ctx context.Context, realm, code, redirectURL string) (*models.OpenIDLoginResponse, error) {
	path := fmt.Sprintf("/access/openid/login")
	
	body, err := json.Marshal(map[string]string{
		"realm":        realm,
		"code":         code,
		"redirect-url": redirectURL,
	})
	if err != nil {
		return nil, err
	}
	
	var result models.APIResponse[models.OpenIDLoginResponse]
	if err := c.Post(ctx, path, bytes.NewReader(body), &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}
