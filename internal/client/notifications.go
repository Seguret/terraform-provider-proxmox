package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// --- Sendmail Endpoints ---

// GetNotificationEndpointSendmails returns all sendmail notification endpoints.
func (c *Client) GetNotificationEndpointSendmails(ctx context.Context) ([]models.NotificationEndpointSendmail, error) {
	var result models.APIResponse[[]models.NotificationEndpointSendmail]
	if err := c.Get(ctx, "/cluster/notifications/endpoints/sendmail", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetNotificationEndpointSendmail fetches a sendmail endpoint config by name.
func (c *Client) GetNotificationEndpointSendmail(ctx context.Context, name string) (*models.NotificationEndpointSendmail, error) {
	path := fmt.Sprintf("/cluster/notifications/endpoints/sendmail/%s", url.PathEscape(name))
	var result models.APIResponse[models.NotificationEndpointSendmail]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateNotificationEndpointSendmail adds a new sendmail notification target.
func (c *Client) CreateNotificationEndpointSendmail(ctx context.Context, req *models.NotificationEndpointSendmailCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/notifications/endpoints/sendmail", bytes.NewReader(body))
}

// UpdateNotificationEndpointSendmail updates an existing sendmail endpoint.
func (c *Client) UpdateNotificationEndpointSendmail(ctx context.Context, name string, req *models.NotificationEndpointSendmailUpdateRequest) error {
	path := fmt.Sprintf("/cluster/notifications/endpoints/sendmail/%s", url.PathEscape(name))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteNotificationEndpointSendmail removes a sendmail endpoint.
func (c *Client) DeleteNotificationEndpointSendmail(ctx context.Context, name string) error {
	path := fmt.Sprintf("/cluster/notifications/endpoints/sendmail/%s", url.PathEscape(name))
	return c.Delete(ctx, path)
}

// --- Gotify Endpoints ---

// GetNotificationEndpointGotifys lists all configured Gotify notification endpoints.
func (c *Client) GetNotificationEndpointGotifys(ctx context.Context) ([]models.NotificationEndpointGotify, error) {
	var result models.APIResponse[[]models.NotificationEndpointGotify]
	if err := c.Get(ctx, "/cluster/notifications/endpoints/gotify", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetNotificationEndpointGotify fetches a Gotify endpoint config by name.
func (c *Client) GetNotificationEndpointGotify(ctx context.Context, name string) (*models.NotificationEndpointGotify, error) {
	path := fmt.Sprintf("/cluster/notifications/endpoints/gotify/%s", url.PathEscape(name))
	var result models.APIResponse[models.NotificationEndpointGotify]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateNotificationEndpointGotify registers a new Gotify push notification target.
func (c *Client) CreateNotificationEndpointGotify(ctx context.Context, req *models.NotificationEndpointGotifyCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/notifications/endpoints/gotify", bytes.NewReader(body))
}

// UpdateNotificationEndpointGotify updates a Gotify endpoint config.
func (c *Client) UpdateNotificationEndpointGotify(ctx context.Context, name string, req *models.NotificationEndpointGotifyUpdateRequest) error {
	path := fmt.Sprintf("/cluster/notifications/endpoints/gotify/%s", url.PathEscape(name))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteNotificationEndpointGotify removes a Gotify endpoint.
func (c *Client) DeleteNotificationEndpointGotify(ctx context.Context, name string) error {
	path := fmt.Sprintf("/cluster/notifications/endpoints/gotify/%s", url.PathEscape(name))
	return c.Delete(ctx, path)
}

// --- SMTP Endpoints ---

// GetNotificationEndpointSmtps lists all SMTP notification endpoints.
func (c *Client) GetNotificationEndpointSmtps(ctx context.Context) ([]models.NotificationEndpointSmtp, error) {
	var result models.APIResponse[[]models.NotificationEndpointSmtp]
	if err := c.Get(ctx, "/cluster/notifications/endpoints/smtp", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetNotificationEndpointSmtp fetches an SMTP endpoint config by name.
func (c *Client) GetNotificationEndpointSmtp(ctx context.Context, name string) (*models.NotificationEndpointSmtp, error) {
	path := fmt.Sprintf("/cluster/notifications/endpoints/smtp/%s", url.PathEscape(name))
	var result models.APIResponse[models.NotificationEndpointSmtp]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateNotificationEndpointSmtp adds a new SMTP mail delivery endpoint.
func (c *Client) CreateNotificationEndpointSmtp(ctx context.Context, req *models.NotificationEndpointSmtpCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/notifications/endpoints/smtp", bytes.NewReader(body))
}

// UpdateNotificationEndpointSmtp updates an existing SMTP endpoint config.
func (c *Client) UpdateNotificationEndpointSmtp(ctx context.Context, name string, req *models.NotificationEndpointSmtpUpdateRequest) error {
	path := fmt.Sprintf("/cluster/notifications/endpoints/smtp/%s", url.PathEscape(name))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteNotificationEndpointSmtp removes an SMTP notification endpoint.
func (c *Client) DeleteNotificationEndpointSmtp(ctx context.Context, name string) error {
	path := fmt.Sprintf("/cluster/notifications/endpoints/smtp/%s", url.PathEscape(name))
	return c.Delete(ctx, path)
}

// --- Webhook Endpoints ---

// GetNotificationEndpointWebhooks lists all webhook notification endpoints.
func (c *Client) GetNotificationEndpointWebhooks(ctx context.Context) ([]models.NotificationEndpointWebhook, error) {
	var result models.APIResponse[[]models.NotificationEndpointWebhook]
	if err := c.Get(ctx, "/cluster/notifications/endpoints/webhook", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetNotificationEndpointWebhook fetches a webhook endpoint config by name.
func (c *Client) GetNotificationEndpointWebhook(ctx context.Context, name string) (*models.NotificationEndpointWebhook, error) {
	path := fmt.Sprintf("/cluster/notifications/endpoints/webhook/%s", url.PathEscape(name))
	var result models.APIResponse[models.NotificationEndpointWebhook]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateNotificationEndpointWebhook adds a new HTTP webhook notification target.
func (c *Client) CreateNotificationEndpointWebhook(ctx context.Context, req *models.NotificationEndpointWebhookCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/notifications/endpoints/webhook", bytes.NewReader(body))
}

// UpdateNotificationEndpointWebhook updates an existing webhook endpoint.
func (c *Client) UpdateNotificationEndpointWebhook(ctx context.Context, name string, req *models.NotificationEndpointWebhookUpdateRequest) error {
	path := fmt.Sprintf("/cluster/notifications/endpoints/webhook/%s", url.PathEscape(name))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteNotificationEndpointWebhook removes a webhook endpoint.
func (c *Client) DeleteNotificationEndpointWebhook(ctx context.Context, name string) error {
	path := fmt.Sprintf("/cluster/notifications/endpoints/webhook/%s", url.PathEscape(name))
	return c.Delete(ctx, path)
}

// --- Notification Filters ---

// GetNotificationFilters lists all notification filter rules.
func (c *Client) GetNotificationFilters(ctx context.Context) ([]models.NotificationFilter, error) {
	var result models.APIResponse[[]models.NotificationFilter]
	if err := c.Get(ctx, "/cluster/notifications/filters", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetNotificationFilter fetches a single notification filter by name.
func (c *Client) GetNotificationFilter(ctx context.Context, name string) (*models.NotificationFilter, error) {
	path := fmt.Sprintf("/cluster/notifications/filters/%s", url.PathEscape(name))
	var result models.APIResponse[models.NotificationFilter]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateNotificationFilter adds a new notification filter rule.
func (c *Client) CreateNotificationFilter(ctx context.Context, req *models.NotificationFilterCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/notifications/filters", bytes.NewReader(body))
}

// UpdateNotificationFilter updates an existing filter rule.
func (c *Client) UpdateNotificationFilter(ctx context.Context, name string, req *models.NotificationFilterUpdateRequest) error {
	path := fmt.Sprintf("/cluster/notifications/filters/%s", url.PathEscape(name))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteNotificationFilter removes a notification filter rule.
func (c *Client) DeleteNotificationFilter(ctx context.Context, name string) error {
	path := fmt.Sprintf("/cluster/notifications/filters/%s", url.PathEscape(name))
	return c.Delete(ctx, path)
}

// --- Notification Matchers ---

// GetNotificationMatchers lists all notification matchers (routing rules).
func (c *Client) GetNotificationMatchers(ctx context.Context) ([]models.NotificationMatcher, error) {
	var result models.APIResponse[[]models.NotificationMatcher]
	if err := c.Get(ctx, "/cluster/notifications/matchers", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetNotificationMatcher fetches a single matcher rule by name.
func (c *Client) GetNotificationMatcher(ctx context.Context, name string) (*models.NotificationMatcher, error) {
	path := fmt.Sprintf("/cluster/notifications/matchers/%s", url.PathEscape(name))
	var result models.APIResponse[models.NotificationMatcher]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateNotificationMatcher adds a new matcher that routes notifications to endpoints.
func (c *Client) CreateNotificationMatcher(ctx context.Context, req *models.NotificationMatcherCreateRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, "/cluster/notifications/matchers", bytes.NewReader(body))
}

// UpdateNotificationMatcher updates an existing notification matcher.
func (c *Client) UpdateNotificationMatcher(ctx context.Context, name string, req *models.NotificationMatcherUpdateRequest) error {
	path := fmt.Sprintf("/cluster/notifications/matchers/%s", url.PathEscape(name))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteNotificationMatcher removes a notification matcher by name.
func (c *Client) DeleteNotificationMatcher(ctx context.Context, name string) error {
	path := fmt.Sprintf("/cluster/notifications/matchers/%s", url.PathEscape(name))
	return c.Delete(ctx, path)
}
