package vm_snapshots

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &VMSnapshotsDataSource{}
var _ datasource.DataSourceWithConfigure = &VMSnapshotsDataSource{}

type VMSnapshotsDataSource struct {
	client *client.Client
}

type VMSnapshotsDataSourceModel struct {
	ID           types.String   `tfsdk:"id"`
	NodeName     types.String   `tfsdk:"node_name"`
	VMID         types.Int64    `tfsdk:"vmid"`
	SnapNames    []types.String `tfsdk:"snap_names"`
	Descriptions []types.String `tfsdk:"descriptions"`
	Snaptimes    []types.Int64  `tfsdk:"snaptimes"`
}

func NewDataSource() datasource.DataSource {
	return &VMSnapshotsDataSource{}
}

func (d *VMSnapshotsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_vm_snapshots"
}

func (d *VMSnapshotsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of snapshots for a Proxmox VE virtual machine.",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true},
			"node_name":    schema.StringAttribute{Description: "The node name.", Required: true},
			"vmid":         schema.Int64Attribute{Description: "The VM ID.", Required: true},
			"snap_names":   schema.ListAttribute{Description: "Snapshot names.", Computed: true, ElementType: types.StringType},
			"descriptions": schema.ListAttribute{Description: "Snapshot descriptions.", Computed: true, ElementType: types.StringType},
			"snaptimes":    schema.ListAttribute{Description: "Snapshot creation times as UNIX timestamps.", Computed: true, ElementType: types.Int64Type},
		},
	}
}

func (d *VMSnapshotsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VMSnapshotsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config VMSnapshotsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	vmid := int(config.VMID.ValueInt64())

	tflog.Debug(ctx, "Reading VM snapshots", map[string]any{"node": node, "vmid": vmid})

	snapshots, err := d.client.GetVMSnapshots(ctx, node, vmid)
	if err != nil {
		resp.Diagnostics.AddError("Error reading VM snapshots", err.Error())
		return
	}

	state := VMSnapshotsDataSourceModel{
		ID:           types.StringValue(fmt.Sprintf("%s/%d/snapshots", node, vmid)),
		NodeName:     config.NodeName,
		VMID:         config.VMID,
		SnapNames:    make([]types.String, len(snapshots)),
		Descriptions: make([]types.String, len(snapshots)),
		Snaptimes:    make([]types.Int64, len(snapshots)),
	}

	for i, s := range snapshots {
		state.SnapNames[i] = types.StringValue(s.Name)
		state.Descriptions[i] = types.StringValue(s.Description)
		state.Snaptimes[i] = types.Int64Value(s.Snaptime)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
