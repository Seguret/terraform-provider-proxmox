package acme_directories

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &ACMEDirectoriesDataSource{}
var _ datasource.DataSourceWithConfigure = &ACMEDirectoriesDataSource{}

type ACMEDirectoriesDataSource struct {
	client *client.Client
}

type ACMEDirectoryModel struct {
	Name types.String `tfsdk:"name"`
	URL  types.String `tfsdk:"url"`
}

type ACMEDirectoriesDataSourceModel struct {
	ID          types.String         `tfsdk:"id"`
	Directories []ACMEDirectoryModel `tfsdk:"directories"`
}

func NewDataSource() datasource.DataSource {
	return &ACMEDirectoriesDataSource{}
}

func (d *ACMEDirectoriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_acme_directories"
}

func (d *ACMEDirectoriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of ACME directory endpoints available in Proxmox VE.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"directories": schema.ListNestedAttribute{
				Description: "The list of ACME directory entries.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The human-readable name of the ACME directory.",
							Computed:    true,
						},
						"url": schema.StringAttribute{
							Description: "The ACME directory URL.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *ACMEDirectoriesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ACMEDirectoriesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Proxmox VE ACME directories")

	dirs, err := d.client.GetACMEDirectories(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ACME directories", err.Error())
		return
	}

	state := ACMEDirectoriesDataSourceModel{
		ID:          types.StringValue("acme_directories"),
		Directories: make([]ACMEDirectoryModel, len(dirs)),
	}

	for i, dir := range dirs {
		state.Directories[i] = ACMEDirectoryModel{
			Name: types.StringValue(dir.Name),
			URL:  types.StringValue(dir.URL),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
