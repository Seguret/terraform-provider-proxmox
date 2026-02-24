package version

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &VersionDataSource{}
var _ datasource.DataSourceWithConfigure = &VersionDataSource{}

type VersionDataSource struct {
	client *client.Client
}

type VersionDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Release types.String `tfsdk:"release"`
	RepoID  types.String `tfsdk:"repo_id"`
	Version types.String `tfsdk:"version"`
}

func NewDataSource() datasource.DataSource {
	return &VersionDataSource{}
}

func (d *VersionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_version"
}

func (d *VersionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the Proxmox VE version information.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"release": schema.StringAttribute{
				Description: "The Proxmox VE release string (e.g., '8.1').",
				Computed:    true,
			},
			"repo_id": schema.StringAttribute{
				Description: "The Proxmox VE repository ID.",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "The Proxmox VE version string (e.g., '8.1.3').",
				Computed:    true,
			},
		},
	}
}

func (d *VersionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = c
}

func (d *VersionDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Proxmox VE version")

	ver, err := d.client.GetVersion(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Proxmox VE Version",
			"An error occurred while reading the Proxmox VE version: "+err.Error(),
		)
		return
	}

	state := VersionDataSourceModel{
		ID:      types.StringValue("version"),
		Release: types.StringValue(ver.Release),
		RepoID:  types.StringValue(ver.RepoID),
		Version: types.StringValue(ver.Version),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
