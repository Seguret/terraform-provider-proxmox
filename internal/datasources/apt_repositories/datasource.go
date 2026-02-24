package apt_repositories

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &AptRepositoriesDataSource{}
var _ datasource.DataSourceWithConfigure = &AptRepositoriesDataSource{}

type AptRepositoriesDataSource struct {
	client *client.Client
}

type AptRepositoriesDataSourceModel struct {
	ID        types.String   `tfsdk:"id"`
	NodeName  types.String   `tfsdk:"node_name"`
	Files     []types.String `tfsdk:"files"`
	URIs      []types.String `tfsdk:"uris"`
	Suites    []types.String `tfsdk:"suites"`
	Enabled   []types.Bool   `tfsdk:"enabled"`
}

func NewDataSource() datasource.DataSource {
	return &AptRepositoriesDataSource{}
}

func (d *AptRepositoriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_apt_repositories"
}

func (d *AptRepositoriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves all APT repositories configured on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Computed: true},
			"node_name": schema.StringAttribute{Description: "The node name.", Required: true},
			"files":     schema.ListAttribute{Description: "Source file path for each repository.", Computed: true, ElementType: types.StringType},
			"uris":      schema.ListAttribute{Description: "URI of each repository.", Computed: true, ElementType: types.StringType},
			"suites":    schema.ListAttribute{Description: "Suite (e.g. bookworm) of each repository.", Computed: true, ElementType: types.StringType},
			"enabled":   schema.ListAttribute{Description: "Whether each repository is enabled.", Computed: true, ElementType: types.BoolType},
		},
	}
}

func (d *AptRepositoriesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AptRepositoriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config AptRepositoriesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading APT repositories", map[string]any{"node": node})

	repos, err := d.client.GetAptRepositories(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading APT repositories", err.Error())
		return
	}

	state := AptRepositoriesDataSourceModel{
		ID:       types.StringValue(node + "/apt_repositories"),
		NodeName: config.NodeName,
	}

	for _, f := range repos.Files {
		for _, r := range f.Repositories {
			state.Files = append(state.Files, types.StringValue(f.Filename))
			uri := ""
			if len(r.URIs) > 0 {
				uri = r.URIs[0]
			}
			state.URIs = append(state.URIs, types.StringValue(uri))
			suite := ""
			if len(r.Suites) > 0 {
				suite = r.Suites[0]
			}
			state.Suites = append(state.Suites, types.StringValue(suite))
			state.Enabled = append(state.Enabled, types.BoolValue(r.Enabled))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
