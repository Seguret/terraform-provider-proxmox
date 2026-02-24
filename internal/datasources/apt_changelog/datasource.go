package apt_changelog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &APTChangelogDataSource{}
var _ datasource.DataSourceWithConfigure = &APTChangelogDataSource{}

type APTChangelogDataSource struct {
	client *client.Client
}

type APTChangelogDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	NodeName  types.String `tfsdk:"node_name"`
	Package   types.String `tfsdk:"package"`
	Changelog types.String `tfsdk:"changelog"`
}

func NewDataSource() datasource.DataSource {
	return &APTChangelogDataSource{}
}

func (d *APTChangelogDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_apt_changelog"
}

func (d *APTChangelogDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the APT changelog for a package on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox VE node.",
				Required:    true,
			},
			"package": schema.StringAttribute{
				Description: "The name of the APT package to retrieve the changelog for.",
				Required:    true,
			},
			"changelog": schema.StringAttribute{
				Description: "The changelog content for the package.",
				Computed:    true,
			},
		},
	}
}

func (d *APTChangelogDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *APTChangelogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config APTChangelogDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	pkg := config.Package.ValueString()
	tflog.Debug(ctx, "Reading Proxmox VE APT changelog", map[string]any{"node": node, "package": pkg})

	changelog, err := d.client.GetAPTChangelog(ctx, node, pkg)
	if err != nil {
		resp.Diagnostics.AddError("Error reading APT changelog", err.Error())
		return
	}

	state := APTChangelogDataSourceModel{
		ID:        types.StringValue(fmt.Sprintf("%s/apt/changelog/%s", node, pkg)),
		NodeName:  config.NodeName,
		Package:   config.Package,
		Changelog: types.StringValue(changelog),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
