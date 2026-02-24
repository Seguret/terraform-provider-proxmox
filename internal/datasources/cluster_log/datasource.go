package cluster_log

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &ClusterLogDataSource{}
var _ datasource.DataSourceWithConfigure = &ClusterLogDataSource{}

type ClusterLogDataSource struct {
	client *client.Client
}

type ClusterLogEntryModel struct {
	PID      types.Int64  `tfsdk:"pid"`
	UID      types.Int64  `tfsdk:"uid"`
	GID      types.Int64  `tfsdk:"gid"`
	Node     types.String `tfsdk:"node"`
	UserID   types.String `tfsdk:"user_id"`
	Tag      types.String `tfsdk:"tag"`
	Severity types.String `tfsdk:"severity"`
	Msg      types.String `tfsdk:"msg"`
	Time     types.Int64  `tfsdk:"time"`
}

type ClusterLogDataSourceModel struct {
	ID      types.String           `tfsdk:"id"`
	Max     types.Int64            `tfsdk:"max"`
	Entries []ClusterLogEntryModel `tfsdk:"entries"`
}

func NewDataSource() datasource.DataSource {
	return &ClusterLogDataSource{}
}

func (d *ClusterLogDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_cluster_log"
}

func (d *ClusterLogDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the Proxmox VE cluster event log.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"max": schema.Int64Attribute{
				Description: "Maximum number of log entries to return (default 50).",
				Optional:    true,
				Computed:    true,
			},
			"entries": schema.ListNestedAttribute{
				Description: "Cluster log entries.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"pid":      schema.Int64Attribute{Computed: true, Description: "Process ID."},
						"uid":      schema.Int64Attribute{Computed: true, Description: "User ID."},
						"gid":      schema.Int64Attribute{Computed: true, Description: "Group ID."},
						"node":     schema.StringAttribute{Computed: true, Description: "Node name."},
						"user_id":  schema.StringAttribute{Computed: true, Description: "User who triggered the event."},
						"tag":      schema.StringAttribute{Computed: true, Description: "Log tag/service."},
						"severity": schema.StringAttribute{Computed: true, Description: "Severity level."},
						"msg":      schema.StringAttribute{Computed: true, Description: "Log message."},
						"time":     schema.Int64Attribute{Computed: true, Description: "Unix timestamp."},
					},
				},
			},
		},
	}
}

func (d *ClusterLogDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ClusterLogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ClusterLogDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	max := 50
	if !config.Max.IsNull() && !config.Max.IsUnknown() {
		max = int(config.Max.ValueInt64())
	}

	tflog.Debug(ctx, "Reading cluster log", map[string]any{"max": max})

	entries, err := d.client.GetClusterLog(ctx, max)
	if err != nil {
		resp.Diagnostics.AddError("Error reading cluster log", err.Error())
		return
	}

	state := ClusterLogDataSourceModel{
		ID:      types.StringValue("cluster-log"),
		Max:     types.Int64Value(int64(max)),
		Entries: make([]ClusterLogEntryModel, len(entries)),
	}

	for i, e := range entries {
		state.Entries[i] = ClusterLogEntryModel{
			PID:      types.Int64Value(e.PID),
			UID:      types.Int64Value(e.UID),
			GID:      types.Int64Value(e.GID),
			Node:     types.StringValue(e.Node),
			UserID:   types.StringValue(e.UserID),
			Tag:      types.StringValue(e.Tag),
			Severity: types.StringValue(e.Severity),
			Msg:      types.StringValue(e.Msg),
			Time:     types.Int64Value(e.Time),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
