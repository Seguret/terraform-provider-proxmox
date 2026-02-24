package apt_versions

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &APTVersionsDataSource{}
var _ datasource.DataSourceWithConfigure = &APTVersionsDataSource{}

type APTVersionsDataSource struct {
	client *client.Client
}

type APTPackageModel struct {
	Package    types.String `tfsdk:"package"`
	Version    types.String `tfsdk:"version"`
	OldVersion types.String `tfsdk:"old_version"`
	Priority   types.String `tfsdk:"priority"`
	Section    types.String `tfsdk:"section"`
	Title      types.String `tfsdk:"title"`
}

type APTVersionsDataSourceModel struct {
	ID       types.String      `tfsdk:"id"`
	NodeName types.String      `tfsdk:"node_name"`
	Packages []APTPackageModel `tfsdk:"packages"`
}

func NewDataSource() datasource.DataSource {
	return &APTVersionsDataSource{}
}

func (d *APTVersionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_apt_versions"
}

func (d *APTVersionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of installed APT package versions on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox VE node.",
				Required:    true,
			},
			"packages": schema.ListNestedAttribute{
				Description: "The list of installed packages and their versions.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"package": schema.StringAttribute{
							Description: "The package name.",
							Computed:    true,
						},
						"version": schema.StringAttribute{
							Description: "The currently installed version.",
							Computed:    true,
						},
						"old_version": schema.StringAttribute{
							Description: "The previous installed version (if an upgrade is pending).",
							Computed:    true,
						},
						"priority": schema.StringAttribute{
							Description: "The package priority.",
							Computed:    true,
						},
						"section": schema.StringAttribute{
							Description: "The package section.",
							Computed:    true,
						},
						"title": schema.StringAttribute{
							Description: "The package short description.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *APTVersionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *APTVersionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config APTVersionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading Proxmox VE APT package versions", map[string]any{"node": node})

	pkgs, err := d.client.GetAPTVersions(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading APT package versions", err.Error())
		return
	}

	state := APTVersionsDataSourceModel{
		ID:       types.StringValue(fmt.Sprintf("%s/apt/versions", node)),
		NodeName: config.NodeName,
		Packages: make([]APTPackageModel, len(pkgs)),
	}

	for i, p := range pkgs {
		state.Packages[i] = APTPackageModel{
			Package:    types.StringValue(p.Package),
			Version:    types.StringValue(p.Version),
			OldVersion: types.StringValue(p.OldVersion),
			Priority:   types.StringValue(p.Priority),
			Section:    types.StringValue(p.Section),
			Title:      types.StringValue(p.Title),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
