package node_journal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &NodeJournalDataSource{}
var _ datasource.DataSourceWithConfigure = &NodeJournalDataSource{}

type NodeJournalDataSource struct {
	client *client.Client
}

type NodeJournalDataSourceModel struct {
	ID       types.String   `tfsdk:"id"`
	NodeName types.String   `tfsdk:"node_name"`
	Entries  []types.String `tfsdk:"entries"`
}

func NewDataSource() datasource.DataSource {
	return &NodeJournalDataSource{}
}

func (d *NodeJournalDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_journal"
}

func (d *NodeJournalDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the systemd journal entries from a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox VE node.",
				Required:    true,
			},
			"entries": schema.ListAttribute{
				Description: "The journal log entries.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *NodeJournalDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NodeJournalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config NodeJournalDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading Proxmox VE node journal", map[string]any{"node": node})

	lines, err := d.client.GetNodeJournal(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading node journal", err.Error())
		return
	}

	entries := make([]types.String, len(lines))
	for i, line := range lines {
		entries[i] = types.StringValue(line)
	}

	state := NodeJournalDataSourceModel{
		ID:       types.StringValue(fmt.Sprintf("%s/journal", node)),
		NodeName: config.NodeName,
		Entries:  entries,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
