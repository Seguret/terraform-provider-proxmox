package apt_standard_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// handleURIs maps well-known Proxmox standard repository handles to their canonical URIs.
var handleURIs = map[string]string{
	"pve-enterprise":              "https://enterprise.proxmox.com/debian/pve",
	"pve-no-subscription":         "http://download.proxmox.com/debian/pve",
	"pvetest":                     "http://download.proxmox.com/debian/pvetest",
	"ceph-quincy-enterprise":      "https://enterprise.proxmox.com/debian/ceph-quincy",
	"ceph-quincy-no-subscription": "http://download.proxmox.com/debian/ceph-quincy",
	"ceph-reef-enterprise":        "https://enterprise.proxmox.com/debian/ceph-reef",
	"ceph-reef-no-subscription":   "http://download.proxmox.com/debian/ceph-reef",
}

var _ datasource.DataSource = &AptStandardRepositoryDataSource{}
var _ datasource.DataSourceWithConfigure = &AptStandardRepositoryDataSource{}

type AptStandardRepositoryDataSource struct {
	client *client.Client
}

type AptStandardRepositoryDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	NodeName    types.String `tfsdk:"node_name"`
	Handle      types.String `tfsdk:"handle"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func NewDataSource() datasource.DataSource {
	return &AptStandardRepositoryDataSource{}
}

func (d *AptStandardRepositoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_apt_standard_repository"
}

func (d *AptStandardRepositoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads the state of a standard (built-in) Proxmox VE APT repository on a node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for this data source (node/handle).",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox VE node.",
				Required:    true,
			},
			"handle": schema.StringAttribute{
				Description: "The standard repository handle (e.g. 'pve-no-subscription', 'pve-enterprise').",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the repository is currently enabled.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The repository name / handle alias.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of the repository.",
				Computed:    true,
			},
		},
	}
}

func (d *AptStandardRepositoryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	d.client = c
}

func (d *AptStandardRepositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config AptStandardRepositoryDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	handle := config.Handle.ValueString()

	tflog.Debug(ctx, "Reading APT standard repository", map[string]any{"node": node, "handle": handle})

	repos, err := d.client.GetAptRepositories(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading APT repositories", err.Error())
		return
	}

	filePath, index, found := findRepoByHandle(handle, repos)
	if !found {
		resp.Diagnostics.AddError("APT standard repository not found",
			fmt.Sprintf("Could not find handle '%s' on node '%s'.", handle, node))
		return
	}

	// find the matching entry and read its enabled state
	var enabled bool
	for _, f := range repos.Files {
		if f.Filename != filePath {
			continue
		}
		if index < len(f.Repositories) {
			enabled = f.Repositories[index].Enabled
		}
		break
	}

	state := AptStandardRepositoryDataSourceModel{
		ID:          types.StringValue(fmt.Sprintf("%s/%s", node, handle)),
		NodeName:    types.StringValue(node),
		Handle:      types.StringValue(handle),
		Enabled:     types.BoolValue(enabled),
		Name:        types.StringValue(handle),
		Description: types.StringValue(fmt.Sprintf("Proxmox standard APT repository: %s", handle)),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// findRepoByHandle searches repos for an entry matching the handle's known URI.
// Returns (filePath, index, found).
func findRepoByHandle(handle string, repos *models.AptRepositoriesResponse) (string, int, bool) {
	uri, known := handleURIs[handle]
	for _, f := range repos.Files {
		for i, repo := range f.Repositories {
			for _, u := range repo.URIs {
				if known && strings.TrimRight(u, "/") == strings.TrimRight(uri, "/") {
					return f.Filename, i, true
				}
				if !known && strings.Contains(u, handle) {
					return f.Filename, i, true
				}
			}
		}
	}
	return "", 0, false
}
