package datastores

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &DatastoresDataSource{}
var _ datasource.DataSourceWithConfigure = &DatastoresDataSource{}

type DatastoresDataSource struct {
	client *client.Client
}

type DatastoresDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	NodeName types.String `tfsdk:"node_name"`

	Names        []types.String  `tfsdk:"names"`
	Types        []types.String  `tfsdk:"types"`
	ContentTypes []types.String  `tfsdk:"content_types"`
	Active       []types.Bool    `tfsdk:"active"`
	Enabled      []types.Bool    `tfsdk:"enabled"`
	Shared       []types.Bool    `tfsdk:"shared"`
	Total        []types.Int64   `tfsdk:"total"`
	Used         []types.Int64   `tfsdk:"used"`
	Available    []types.Int64   `tfsdk:"available"`
	UsedFraction []types.Float64 `tfsdk:"used_fraction"`
}

func NewDataSource() datasource.DataSource {
	return &DatastoresDataSource{}
}

func (d *DatastoresDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_datastores"
}

func (d *DatastoresDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of datastores (storage) available on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox VE node.",
				Required:    true,
			},
			"names": schema.ListAttribute{
				Description: "The storage names.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"types": schema.ListAttribute{
				Description: "The storage types (e.g., 'dir', 'lvm', 'zfspool', 'nfs', 'cifs').",
				Computed:    true,
				ElementType: types.StringType,
			},
			"content_types": schema.ListAttribute{
				Description: "The content types supported by each storage (comma-separated, e.g., 'images,rootdir,vztmpl,iso,backup').",
				Computed:    true,
				ElementType: types.StringType,
			},
			"active": schema.ListAttribute{
				Description: "Whether each storage is active.",
				Computed:    true,
				ElementType: types.BoolType,
			},
			"enabled": schema.ListAttribute{
				Description: "Whether each storage is enabled.",
				Computed:    true,
				ElementType: types.BoolType,
			},
			"shared": schema.ListAttribute{
				Description: "Whether each storage is shared across nodes.",
				Computed:    true,
				ElementType: types.BoolType,
			},
			"total": schema.ListAttribute{
				Description: "Total size in bytes for each storage.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"used": schema.ListAttribute{
				Description: "Used size in bytes for each storage.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"available": schema.ListAttribute{
				Description: "Available size in bytes for each storage.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"used_fraction": schema.ListAttribute{
				Description: "The used fraction (0.0-1.0) for each storage.",
				Computed:    true,
				ElementType: types.Float64Type,
			},
		},
	}
}

func (d *DatastoresDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DatastoresDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DatastoresDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodeName := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading Proxmox VE datastores", map[string]any{"node": nodeName})

	storageList, err := d.client.GetNodeStorage(ctx, nodeName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Proxmox VE Datastores",
			fmt.Sprintf("An error occurred while reading storage for node '%s': %s", nodeName, err.Error()),
		)
		return
	}

	state := DatastoresDataSourceModel{
		ID:           types.StringValue(fmt.Sprintf("datastores/%s", nodeName)),
		NodeName:     types.StringValue(nodeName),
		Names:        make([]types.String, len(storageList)),
		Types:        make([]types.String, len(storageList)),
		ContentTypes: make([]types.String, len(storageList)),
		Active:       make([]types.Bool, len(storageList)),
		Enabled:      make([]types.Bool, len(storageList)),
		Shared:       make([]types.Bool, len(storageList)),
		Total:        make([]types.Int64, len(storageList)),
		Used:         make([]types.Int64, len(storageList)),
		Available:    make([]types.Int64, len(storageList)),
		UsedFraction: make([]types.Float64, len(storageList)),
	}

	for i, s := range storageList {
		state.Names[i] = types.StringValue(s.Storage)
		state.Types[i] = types.StringValue(s.Type)
		state.ContentTypes[i] = types.StringValue(s.Content)
		state.Active[i] = types.BoolValue(s.Active == 1)
		state.Enabled[i] = types.BoolValue(s.Enabled == 1)
		state.Shared[i] = types.BoolValue(s.Shared == 1)
		state.Total[i] = types.Int64Value(s.Total)
		state.Used[i] = types.Int64Value(s.Used)
		state.Available[i] = types.Int64Value(s.Avail)
		state.UsedFraction[i] = types.Float64Value(s.UsedFrac)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
