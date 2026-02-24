package node_tasks

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &NodeTasksDataSource{}
var _ datasource.DataSourceWithConfigure = &NodeTasksDataSource{}

type NodeTasksDataSource struct {
	client *client.Client
}

type NodeTaskModel struct {
	UPID      types.String `tfsdk:"upid"`
	Type      types.String `tfsdk:"type"`
	TaskID    types.String `tfsdk:"task_id"`
	User      types.String `tfsdk:"user"`
	Status    types.String `tfsdk:"status"`
	StartTime types.Int64  `tfsdk:"start_time"`
	EndTime   types.Int64  `tfsdk:"end_time"`
}

type NodeTasksDataSourceModel struct {
	ID       types.String    `tfsdk:"id"`
	NodeName types.String    `tfsdk:"node_name"`
	Tasks    []NodeTaskModel `tfsdk:"tasks"`
}

func NewDataSource() datasource.DataSource {
	return &NodeTasksDataSource{}
}

func (d *NodeTasksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_tasks"
}

func (d *NodeTasksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of recent tasks on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox VE node.",
				Required:    true,
			},
			"tasks": schema.ListNestedAttribute{
				Description: "The list of tasks on the node.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"upid": schema.StringAttribute{
							Description: "The unique task process ID.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The task type.",
							Computed:    true,
						},
						"task_id": schema.StringAttribute{
							Description: "The task-specific identifier (e.g., VM ID).",
							Computed:    true,
						},
						"user": schema.StringAttribute{
							Description: "The user who initiated the task.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The task status (e.g., 'OK', error message).",
							Computed:    true,
						},
						"start_time": schema.Int64Attribute{
							Description: "The task start time as a Unix timestamp.",
							Computed:    true,
						},
						"end_time": schema.Int64Attribute{
							Description: "The task end time as a Unix timestamp.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *NodeTasksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NodeTasksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config NodeTasksDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading Proxmox VE node tasks", map[string]any{"node": node})

	tasks, err := d.client.ListNodeTasks(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading node tasks", err.Error())
		return
	}

	state := NodeTasksDataSourceModel{
		ID:       types.StringValue(fmt.Sprintf("%s/tasks", node)),
		NodeName: config.NodeName,
		Tasks:    make([]NodeTaskModel, len(tasks)),
	}

	for i, t := range tasks {
		state.Tasks[i] = NodeTaskModel{
			UPID:      types.StringValue(t.UPID),
			Type:      types.StringValue(t.Type),
			TaskID:    types.StringValue(t.ID),
			User:      types.StringValue(t.User),
			Status:    types.StringValue(t.Status),
			StartTime: types.Int64Value(t.StartTime),
			EndTime:   types.Int64Value(t.EndTime),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
