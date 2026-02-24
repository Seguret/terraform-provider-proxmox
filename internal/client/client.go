package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

const (
	defaultTimeout     = 30 * time.Second
	taskPollInterval   = 2 * time.Second
	taskDefaultTimeout = 30 * time.Minute
)

// Config holds all the settings needed to talk to the Proxmox API.
type Config struct {
	Endpoint string
	APIToken string
	Username string
	Password string
	Insecure bool
	Timeout  time.Duration
}

// Client is the main HTTP client for the Proxmox VE API.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiToken   string

	// fields for ticket-based auth (username/password login)
	ticket    string
	csrfToken string
}

// New builds a new Proxmox client from the given config.
func New(cfg Config) (*Client, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("proxmox endpoint is required")
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: cfg.Insecure, //nolint:gosec // user-configured
		},
	}

	c := &Client{
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		baseURL: strings.TrimRight(cfg.Endpoint, "/") + "/api2/json",
	}

	if cfg.APIToken != "" {
		c.apiToken = cfg.APIToken
	} else if cfg.Username != "" && cfg.Password != "" {
		if err := c.authenticate(context.Background(), cfg.Username, cfg.Password); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	} else {
		return nil, fmt.Errorf("either api_token or username/password must be provided")
	}

	return c, nil
}

// authenticate does the username/password login and stores the ticket + CSRF token.
func (c *Client) authenticate(ctx context.Context, username, password string) error {
	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/access/ticket", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.parseError(resp)
	}

	var result models.APIResponse[struct {
		Ticket              string `json:"ticket"`
		CSRFPreventionToken string `json:"CSRFPreventionToken"`
	}]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode auth response: %w", err)
	}

	c.ticket = result.Data.Ticket
	c.csrfToken = result.Data.CSRFPreventionToken
	return nil
}

// DoRequest fires a raw HTTP request to the Proxmox API with the right auth headers.
func (c *Client) DoRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	reqURL := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return nil, err
	}

	if c.apiToken != "" {
		req.Header.Set("Authorization", "PVEAPIToken="+c.apiToken)
	} else if c.ticket != "" {
		req.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: c.ticket})
		if method != http.MethodGet {
			req.Header.Set("CSRFPreventionToken", c.csrfToken)
		}
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

// Get does a GET and decodes the JSON response into target.
func (c *Client) Get(ctx context.Context, path string, target any) error {
	resp, err := c.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.parseError(resp)
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	return nil
}

// Post sends a POST with a JSON body and decodes whatever comes back.
func (c *Client) Post(ctx context.Context, path string, body io.Reader, target any) error {
	resp, err := c.DoRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.parseError(resp)
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	return nil
}

// PostJSON is just an alias for Post, kept for backwards compat.
func (c *Client) PostJSON(ctx context.Context, path string, body io.Reader, target any) error {
	return c.Post(ctx, path, body, target)
}

// PostTask posts to an endpoint that kicks off an async task, waits for it to finish, and returns the UPID.
func (c *Client) PostTask(ctx context.Context, path string, body io.Reader) (string, error) {
	var result models.APIResponse[string]
	if err := c.Post(ctx, path, body, &result); err != nil {
		return "", err
	}
	upid := result.Data
	if upid != "" {
		if err := c.WaitForUPID(ctx, upid); err != nil {
			return "", err
		}
	}
	return upid, nil
}


// PostNoResponse fires a POST but doesnt bother decoding the response body.
func (c *Client) PostNoResponse(ctx context.Context, path string, body io.Reader) error {
	resp, err := c.DoRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.parseError(resp)
	}
	return nil
}

// Put sends a PUT with a JSON body, used for updates.
func (c *Client) Put(ctx context.Context, path string, body io.Reader) error {
	resp, err := c.DoRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.parseError(resp)
	}
	return nil
}

// Delete sends a DELETE request to the given path.
func (c *Client) Delete(ctx context.Context, path string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.parseError(resp)
	}
	return nil
}

// DeleteTask deletes a resource that triggers an async task, waits for completion, returns the UPID.
func (c *Client) DeleteTask(ctx context.Context, path string) (string, error) {
	resp, err := c.DoRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", c.parseError(resp)
	}

	var result models.APIResponse[string]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	upid := result.Data
	if upid != "" {
		if err := c.WaitForUPID(ctx, upid); err != nil {
			return "", err
		}
	}
	return upid, nil
}


// WaitForUPID parses the node out of a UPID string and waits for the task to finish.
// UPID format is: UPID:<node>:<pid>:<pstart>:<starttime>:<type>:<id>:<user>:
func (c *Client) WaitForUPID(ctx context.Context, upid string) error {
	parts := strings.SplitN(upid, ":", 4)
	if len(parts) < 3 || parts[0] != "UPID" {
		return fmt.Errorf("invalid UPID format: %s", upid)
	}
	return c.WaitForTask(ctx, parts[1], upid)
}

// WaitForTask polls a task until it finishes or times out (uses the default timeout).
func (c *Client) WaitForTask(ctx context.Context, node, upid string) error {
	return c.WaitForTaskWithTimeout(ctx, node, upid, taskDefaultTimeout)
}

// WaitForTaskWithTimeout polls a task status until done or the custom timeout is hit.
func (c *Client) WaitForTaskWithTimeout(ctx context.Context, node, upid string, timeout time.Duration) error {
	ticker := time.NewTicker(taskPollInterval)
	defer ticker.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	path := fmt.Sprintf("/nodes/%s/tasks/%s/status", url.PathEscape(node), url.PathEscape(upid))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return fmt.Errorf("timeout waiting for task %s after %s", upid, timeout)
		case <-ticker.C:
			var result models.APIResponse[models.TaskStatus]
			if err := c.Get(ctx, path, &result); err != nil {
				return fmt.Errorf("failed to poll task status: %w", err)
			}

			if result.Data.Status == "stopped" {
				if result.Data.ExitStatus != "OK" {
					return &TaskError{
						UPID:       upid,
						Node:       node,
						ExitStatus: result.Data.ExitStatus,
					}
				}
				return nil
			}
		}
	}
}

// GetVersion fetches the PVE API version info.
func (c *Client) GetVersion(ctx context.Context) (*models.Version, error) {
	var result models.APIResponse[models.Version]
	if err := c.Get(ctx, "/version", &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetNodes returns the list of nodes in the cluster.
func (c *Client) GetNodes(ctx context.Context) ([]models.NodeListEntry, error) {
	var result models.APIResponse[[]models.NodeListEntry]
	if err := c.Get(ctx, "/nodes", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetNodeStatus grabs current status for a specific node.
func (c *Client) GetNodeStatus(ctx context.Context, node string) (*models.NodeStatus, error) {
	path := fmt.Sprintf("/nodes/%s/status", url.PathEscape(node))
	var result models.APIResponse[models.NodeStatus]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetNodeStorage lists storage pools available on a given node.
func (c *Client) GetNodeStorage(ctx context.Context, node string) ([]models.StorageListEntry, error) {
	path := fmt.Sprintf("/nodes/%s/storage", url.PathEscape(node))
	var result models.APIResponse[[]models.StorageListEntry]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// parseError reads the response body and wraps it in an APIError.
func (c *Client) parseError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
	}

	// try to parse a structured error body if the API sent one
	var errResp struct {
		Errors map[string]string `json:"errors"`
		Data   *string           `json:"data"`
	}
	if json.Unmarshal(body, &errResp) == nil {
		apiErr.Errors = errResp.Errors
		if errResp.Data != nil {
			apiErr.Message = *errResp.Data
		}
	}

	if apiErr.Message == "" && len(body) > 0 {
		apiErr.Message = string(body)
	}

	return apiErr
}
