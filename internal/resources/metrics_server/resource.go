package metrics_server

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &MetricsServerResource{}
var _ resource.ResourceWithConfigure = &MetricsServerResource{}
var _ resource.ResourceWithImportState = &MetricsServerResource{}

type MetricsServerResource struct {
	client *client.Client
}

type MetricsServerResourceModel struct {
	ID            types.String  `tfsdk:"id"`
	Name          types.String  `tfsdk:"name"`
	Type          types.String  `tfsdk:"type"`
	Server        types.String  `tfsdk:"server"`
	Port          types.Int64   `tfsdk:"port"`
	Enabled       types.Bool    `tfsdk:"enabled"`
	MTU           types.Int64   `tfsdk:"mtu"`
	Path          types.String  `tfsdk:"path"`
	Proto         types.String  `tfsdk:"proto"`
	Timeout       types.Int64   `tfsdk:"timeout"`
	Bucket        types.String  `tfsdk:"bucket"`
	InfluxDBProto types.String  `tfsdk:"influxdb_proto"`
	Organization  types.String  `tfsdk:"organization"`
	Token         types.String  `tfsdk:"token"`
	MaxBodySize   types.Int64   `tfsdk:"max_body_size"`
}

func NewResource() resource.Resource {
	return &MetricsServerResource{}
}

func (r *MetricsServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_metrics_server"
}

func (r *MetricsServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE external metrics server (Graphite or InfluxDB).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The metrics server identifier.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "Server type ('graphite' or 'influxdb').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server": schema.StringAttribute{
				Description: "The hostname or IP address of the metrics server.",
				Required:    true,
			},
			"port": schema.Int64Attribute{
				Description: "The port of the metrics server.",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the metrics server is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"mtu": schema.Int64Attribute{
				Description: "MTU (for InfluxDB UDP).",
				Optional:    true,
				Computed:    true,
			},
			"path": schema.StringAttribute{
				Description: "Root graphite path (for Graphite type).",
				Optional:    true,
				Computed:    true,
			},
			"proto": schema.StringAttribute{
				Description: "Protocol ('udp' or 'tcp', for Graphite type).",
				Optional:    true,
				Computed:    true,
			},
			"timeout": schema.Int64Attribute{
				Description: "TCP socket connection timeout in seconds.",
				Optional:    true,
				Computed:    true,
			},
			"bucket": schema.StringAttribute{
				Description: "InfluxDB bucket/database name.",
				Optional:    true,
				Computed:    true,
			},
			"influxdb_proto": schema.StringAttribute{
				Description: "InfluxDB protocol ('udp', 'http', 'https').",
				Optional:    true,
				Computed:    true,
			},
			"organization": schema.StringAttribute{
				Description: "InfluxDB organization name.",
				Optional:    true,
				Computed:    true,
			},
			"token": schema.StringAttribute{
				Description: "InfluxDB access token. Sensitive.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
			},
			"max_body_size": schema.Int64Attribute{
				Description: "Maximum body size in bytes for InfluxDB HTTP(S).",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *MetricsServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *MetricsServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MetricsServerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	disableInt := boolToIntPtr(!plan.Enabled.ValueBool())

	tflog.Debug(ctx, "Creating metrics server", map[string]any{"name": plan.Name.ValueString(), "type": plan.Type.ValueString()})

	if err := r.client.CreateMetricsServer(ctx, plan.Name.ValueString(), &models.MetricsServerCreateRequest{
		Type:          plan.Type.ValueString(),
		Server:        plan.Server.ValueString(),
		Port:          int(plan.Port.ValueInt64()),
		Disable:       disableInt,
		MTU:           int(plan.MTU.ValueInt64()),
		Path:          plan.Path.ValueString(),
		Proto:         plan.Proto.ValueString(),
		Timeout:       int(plan.Timeout.ValueInt64()),
		Bucket:        plan.Bucket.ValueString(),
		InfluxDBProto: plan.InfluxDBProto.ValueString(),
		Organization:  plan.Organization.ValueString(),
		Token:         plan.Token.ValueString(),
		MaxBodySize:   int(plan.MaxBodySize.ValueInt64()),
	}); err != nil {
		resp.Diagnostics.AddError("Error creating metrics server", err.Error())
		return
	}

	plan.ID = plan.Name
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MetricsServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MetricsServerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MetricsServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MetricsServerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	disableInt := boolToIntPtr(!plan.Enabled.ValueBool())

	if err := r.client.UpdateMetricsServer(ctx, plan.Name.ValueString(), &models.MetricsServerUpdateRequest{
		Server:        plan.Server.ValueString(),
		Port:          int(plan.Port.ValueInt64()),
		Disable:       disableInt,
		MTU:           int(plan.MTU.ValueInt64()),
		Path:          plan.Path.ValueString(),
		Proto:         plan.Proto.ValueString(),
		Timeout:       int(plan.Timeout.ValueInt64()),
		Bucket:        plan.Bucket.ValueString(),
		InfluxDBProto: plan.InfluxDBProto.ValueString(),
		Organization:  plan.Organization.ValueString(),
		Token:         plan.Token.ValueString(),
		MaxBodySize:   int(plan.MaxBodySize.ValueInt64()),
	}); err != nil {
		resp.Diagnostics.AddError("Error updating metrics server", err.Error())
		return
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MetricsServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MetricsServerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteMetricsServer(ctx, state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting metrics server", err.Error())
	}
}

func (r *MetricsServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := MetricsServerResourceModel{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MetricsServerResource) readIntoModel(ctx context.Context, model *MetricsServerResourceModel, diagnostics interface{ AddError(string, string) }) {
	srv, err := r.client.GetMetricsServer(ctx, model.Name.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diagnostics.AddError("Metrics server not found", "The metrics server no longer exists.")
			return
		}
		diagnostics.AddError("Error reading metrics server", err.Error())
		return
	}

	model.Type = types.StringValue(srv.Type)
	model.Server = types.StringValue(srv.Server)
	model.Port = types.Int64Value(int64(srv.Port))
	model.MTU = types.Int64Value(int64(srv.MTU))
	model.Path = types.StringValue(srv.Path)
	model.Proto = types.StringValue(srv.Proto)
	model.Timeout = types.Int64Value(int64(srv.Timeout))
	model.Bucket = types.StringValue(srv.Bucket)
	model.InfluxDBProto = types.StringValue(srv.InfluxDBProto)
	model.Organization = types.StringValue(srv.Organization)
	model.MaxBodySize = types.Int64Value(int64(srv.MaxBodySize))
	// token is sensitive — only overwrite if the API returned something
	if srv.Token != "" {
		model.Token = types.StringValue(srv.Token)
	}

	if srv.Disable != nil {
		model.Enabled = types.BoolValue(*srv.Disable == 0)
	} else {
		model.Enabled = types.BoolValue(true)
	}
}

func boolToIntPtr(b bool) *int {
	v := 0
	if b {
		v = 1
	}
	return &v
}
