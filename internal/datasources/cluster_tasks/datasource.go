package cluster_tasks

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &ClusterTasksDataSource{}
var _ datasource.DataSourceWithConfigure = &ClusterTasksDataSource{}

type ClusterTasksDataSource struct {
	client *client.Client
}

type ClusterTaskModel struct {
	UPID      types.String `tfsdk:"upid"`
	Node      types.String `tfsdk:"node"`
	PID       types.Int64  `tfsdk:"pid"`
	StartTime types.Int64  `tfsdk:"start_time"`
	EndTime   types.Int64  `tfsdk:"end_time"`
	Type      types.String `tfsdk:"type"`
	TaskID    types.String `tfsdk:"task_id"`
	User      types.String `tfsdk:"user"`
	Status    types.String `tfsdk:"status"`
}

type ClusterTasksDataSourceModel struct {
	ID    types.String       `tfsdk:"id"`
	Tasks []ClusterTaskModel `tfsdk:"tasks"`
}

func NewDataSource() datasource.DataSource {
	return &ClusterTasksDataSource{}
}

func (d *ClusterTasksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_cluster_tasks"
}

func (d *ClusterTasksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of recent tasks across the Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"tasks": schema.ListNestedAttribute{
				Description: "The list of cluster tasks.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"upid": schema.StringAttribute{
							Description: "The unique task process ID.",
							Computed:    true,
						},
						"node": schema.StringAttribute{
							Description: "The node on which the task ran.",
							Computed:    true,
						},
						"pid": schema.Int64Attribute{
							Description: "The process ID of the task.",
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
					},
				},
			},
		},
	}
}

func (d *ClusterTasksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ClusterTasksDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Proxmox VE cluster tasks")

	tasks, err := d.client.GetClusterTasks(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading cluster tasks", err.Error())
		return
	}

	state := ClusterTasksDataSourceModel{
		ID:    types.StringValue("cluster_tasks"),
		Tasks: make([]ClusterTaskModel, len(tasks)),
	}

	for i, t := range tasks {
		state.Tasks[i] = ClusterTaskModel{
			UPID:      types.StringValue(t.UPID),
			Node:      types.StringValue(t.Node),
			PID:       types.Int64Value(int64(t.PID)),
			StartTime: types.Int64Value(t.StartTime),
			EndTime:   types.Int64Value(t.EndTime),
			Type:      types.StringValue(t.Type),
			TaskID:    types.StringValue(t.ID),
			User:      types.StringValue(t.User),
			Status:    types.StringValue(t.Status),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
