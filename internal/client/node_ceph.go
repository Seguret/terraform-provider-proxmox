package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// GetNodeCephStatus returns the Ceph cluster status as seen from a specific node.
func (c *Client) GetNodeCephStatus(ctx context.Context, node string) (*models.NodeCephStatus, error) {
	path := fmt.Sprintf("/nodes/%s/ceph/status", url.PathEscape(node))
	var result models.APIResponse[models.NodeCephStatus]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetNodeCephOSD returns the list of Ceph OSDs on a specific node.
func (c *Client) GetNodeCephOSD(ctx context.Context, node string) ([]models.CephOSD, error) {
	path := fmt.Sprintf("/nodes/%s/ceph/osd", url.PathEscape(node))
	var result models.APIResponse[[]models.CephOSD]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetNodeCephMON lists Ceph monitor daemons on a node.
func (c *Client) GetNodeCephMON(ctx context.Context, node string) ([]models.CephMON, error) {
	path := fmt.Sprintf("/nodes/%s/ceph/mon", url.PathEscape(node))
	var result models.APIResponse[[]models.CephMON]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetNodeCephConfig fetches the ceph.conf contents from a node.
func (c *Client) GetNodeCephConfig(ctx context.Context, node string) (*models.CephConfig, error) {
	path := fmt.Sprintf("/nodes/%s/ceph/config", url.PathEscape(node))
	var result models.APIResponse[models.CephConfig]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// --- Ceph Pool ---

// GetCephPool fetches details for a single Ceph pool by name.
func (c *Client) GetCephPool(ctx context.Context, node, name string) (*models.CephPool, error) {
	path := fmt.Sprintf("/nodes/%s/ceph/pools/%s", url.PathEscape(node), url.PathEscape(name))
	var result models.APIResponse[models.CephPool]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetCephPools lists all Ceph pools on a node.
func (c *Client) GetCephPools(ctx context.Context, node string) ([]models.CephPool, error) {
	path := fmt.Sprintf("/nodes/%s/ceph/pools", url.PathEscape(node))
	var result models.APIResponse[[]models.CephPool]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// CreateCephPool creates a new Ceph pool on a node. This one is synchronous.
func (c *Client) CreateCephPool(ctx context.Context, node string, req *models.CephPoolCreateRequest) error {
	path := fmt.Sprintf("/nodes/%s/ceph/pools", url.PathEscape(node))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.PostNoResponse(ctx, path, bytes.NewReader(body))
}

// UpdateCephPool updates settings on an existing Ceph pool.
func (c *Client) UpdateCephPool(ctx context.Context, node, name string, req *models.CephPoolUpdateRequest) error {
	path := fmt.Sprintf("/nodes/%s/ceph/pools/%s", url.PathEscape(node), url.PathEscape(name))
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return c.Put(ctx, path, bytes.NewReader(body))
}

// DeleteCephPool removes a Ceph pool (async operation, waits for task completion).
func (c *Client) DeleteCephPool(ctx context.Context, node, name string) error {
	path := fmt.Sprintf("/nodes/%s/ceph/pools/%s", url.PathEscape(node), url.PathEscape(name))
	_, err := c.DeleteTask(ctx, path)
	return err
}

// --- Ceph OSD ---

// GetCephOSDs returns the full OSD tree for a node (includes nested disk info).
func (c *Client) GetCephOSDs(ctx context.Context, node string) (*models.CephOSDListResponse, error) {
	path := fmt.Sprintf("/nodes/%s/ceph/osd", url.PathEscape(node))
	var result models.APIResponse[models.CephOSDListResponse]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateCephOSD provisions a new Ceph OSD on the given disk. Async, returns UPID.
func (c *Client) CreateCephOSD(ctx context.Context, node string, req *models.CephOSDCreateRequest) (string, error) {
	path := fmt.Sprintf("/nodes/%s/ceph/osd", url.PathEscape(node))
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	return c.PostTask(ctx, path, bytes.NewReader(body))
}

// DeleteCephOSD removes a Ceph OSD. Passes cleanup=1 to wipe the disk too.
func (c *Client) DeleteCephOSD(ctx context.Context, node string, osdID int) error {
	path := fmt.Sprintf("/nodes/%s/ceph/osd/%d?cleanup=1", url.PathEscape(node), osdID)
	_, err := c.DeleteTask(ctx, path)
	return err
}

// --- Ceph MON ---

// CreateCephMON sets up a new Ceph monitor daemon on a node (async).
func (c *Client) CreateCephMON(ctx context.Context, node, monID string) error {
	path := fmt.Sprintf("/nodes/%s/ceph/mon/%s", url.PathEscape(node), url.PathEscape(monID))
	_, err := c.PostTask(ctx, path, nil)
	return err
}

// DeleteCephMON tears down a Ceph monitor daemon (async).
func (c *Client) DeleteCephMON(ctx context.Context, node, monID string) error {
	path := fmt.Sprintf("/nodes/%s/ceph/mon/%s", url.PathEscape(node), url.PathEscape(monID))
	_, err := c.DeleteTask(ctx, path)
	return err
}

// --- Ceph MDS ---

// GetCephMDSList lists MDS daemons on a node.
// The API returns a map keyed by daemon name — we convert that to a plain slice.
func (c *Client) GetCephMDSList(ctx context.Context, node string) ([]models.CephMDS, error) {
	path := fmt.Sprintf("/nodes/%s/ceph/mds", url.PathEscape(node))
	var result models.APIResponse[map[string]map[string]interface{}]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	mdsList := make([]models.CephMDS, 0, len(result.Data))
	for name, props := range result.Data {
		mds := models.CephMDS{Name: name}
		if v, ok := props["state"].(string); ok {
			mds.State = v
		}
		if v, ok := props["addr"].(string); ok {
			mds.Addr = v
		}
		mdsList = append(mdsList, mds)
	}
	return mdsList, nil
}

// CreateCephMDS starts a new Ceph MDS daemon on a node (async).
func (c *Client) CreateCephMDS(ctx context.Context, node, name string) error {
	path := fmt.Sprintf("/nodes/%s/ceph/mds/%s", url.PathEscape(node), url.PathEscape(name))
	_, err := c.PostTask(ctx, path, nil)
	return err
}

// DeleteCephMDS stops and removes an MDS daemon (async).
func (c *Client) DeleteCephMDS(ctx context.Context, node, name string) error {
	path := fmt.Sprintf("/nodes/%s/ceph/mds/%s", url.PathEscape(node), url.PathEscape(name))
	_, err := c.DeleteTask(ctx, path)
	return err
}

// --- Ceph MGR ---

// GetCephMGRList returns the Ceph manager daemon list for a node.
func (c *Client) GetCephMGRList(ctx context.Context, node string) (*models.CephMGRListResponse, error) {
	path := fmt.Sprintf("/nodes/%s/ceph/mgr", url.PathEscape(node))
	var result models.APIResponse[models.CephMGRListResponse]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateCephMGR starts a new Ceph MGR daemon on the given node (async).
func (c *Client) CreateCephMGR(ctx context.Context, node, id string) error {
	path := fmt.Sprintf("/nodes/%s/ceph/mgr/%s", url.PathEscape(node), url.PathEscape(id))
	_, err := c.PostTask(ctx, path, nil)
	return err
}

// DeleteCephMGR removes a Ceph MGR daemon (async).
func (c *Client) DeleteCephMGR(ctx context.Context, node, id string) error {
	path := fmt.Sprintf("/nodes/%s/ceph/mgr/%s", url.PathEscape(node), url.PathEscape(id))
	_, err := c.DeleteTask(ctx, path)
	return err
}

// --- Ceph FS ---

// GetCephFSList returns all CephFS filesystems available on a node.
func (c *Client) GetCephFSList(ctx context.Context, node string) ([]models.CephFS, error) {
	path := fmt.Sprintf("/nodes/%s/ceph/fs", url.PathEscape(node))
	var result models.APIResponse[[]models.CephFS]
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// CreateCephFS creates a new CephFS filesystem. Optionally set pg_num via pgNum.
func (c *Client) CreateCephFS(ctx context.Context, node, name string, pgNum *int) error {
	path := fmt.Sprintf("/nodes/%s/ceph/fs/%s", url.PathEscape(node), url.PathEscape(name))
	var body io.Reader
	if pgNum != nil {
		data, err := json.Marshal(map[string]int{"pg_num": *pgNum})
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	}
	return c.PostNoResponse(ctx, path, body)
}

// DeleteCephFS removes a CephFS filesystem from a node.
func (c *Client) DeleteCephFS(ctx context.Context, node, name string) error {
	path := fmt.Sprintf("/nodes/%s/ceph/fs/%s", url.PathEscape(node), url.PathEscape(name))
	return c.Delete(ctx, path)
}
