package replication

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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

var _ resource.Resource = &ReplicationResource{}
var _ resource.ResourceWithConfigure = &ReplicationResource{}
var _ resource.ResourceWithImportState = &ReplicationResource{}

type ReplicationResource struct {
	client *client.Client
}

type ReplicationResourceModel struct {
	ID         types.String  `tfsdk:"id"`
	JobID      types.Int64   `tfsdk:"job_id"`
	SourceVMID types.Int64   `tfsdk:"source_vm_id"`
	TargetNode types.String  `tfsdk:"target_node"`
	Schedule   types.String  `tfsdk:"schedule"`
	Rate       types.Float64 `tfsdk:"rate"`
	Comment    types.String  `tfsdk:"comment"`
	Enabled    types.Bool    `tfsdk:"enabled"`
}

func NewResource() resource.Resource {
	return &ReplicationResource{}
}

func (r *ReplicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_replication"
}

func (r *ReplicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE replication job. The job ID is formed as '{vmid}-{job_id}'.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The replication job ID (format: vmid-jobid, e.g. '100-0').",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"job_id": schema.Int64Attribute{
				Description: "The numeric job ID within the VM (e.g. 0, 1).",
				Required:    true,
				// Int64 planmodifier.RequiresReplace equivalent
			},
			"source_vm_id": schema.Int64Attribute{
				Description: "The VMID to replicate.",
				Required:    true,
			},
			"target_node": schema.StringAttribute{
				Description: "The destination node name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schedule": schema.StringAttribute{
				Description: "Cron schedule for the replication job (e.g. '*/15').",
				Optional:    true,
				Computed:    true,
			},
			"rate": schema.Float64Attribute{
				Description: "Bandwidth limit in MB/s.",
				Optional:    true,
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "A comment for the replication job.",
				Optional:    true,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the replication job is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *ReplicationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ReplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ReplicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jobFullID := fmt.Sprintf("%d-%d", plan.SourceVMID.ValueInt64(), plan.JobID.ValueInt64())
	disableInt := boolToInt(!plan.Enabled.ValueBool())

	createReq := &models.ReplicationJobCreateRequest{
		ID:       jobFullID,
		Type:     "local",
		Target:   plan.TargetNode.ValueString(),
		Guest:    int(plan.SourceVMID.ValueInt64()),
		Schedule: plan.Schedule.ValueString(),
		Rate:     plan.Rate.ValueFloat64(),
		Comment:  plan.Comment.ValueString(),
	}
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() && !plan.Enabled.ValueBool() {
		createReq.Disable = &disableInt
	}

	tflog.Debug(ctx, "Creating replication job", map[string]any{"id": jobFullID, "target": createReq.Target})

	if err := r.client.CreateReplicationJob(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating replication job", err.Error())
		return
	}

	plan.ID = types.StringValue(jobFullID)
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ReplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ReplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ReplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ReplicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jobFullID := plan.ID.ValueString()
	disableInt := boolToInt(!plan.Enabled.ValueBool())

	updateReq := &models.ReplicationJobUpdateRequest{
		Schedule: plan.Schedule.ValueString(),
		Rate:     plan.Rate.ValueFloat64(),
		Comment:  plan.Comment.ValueString(),
	}
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() && !plan.Enabled.ValueBool() {
		updateReq.Disable = &disableInt
	}

	if err := r.client.UpdateReplicationJob(ctx, jobFullID, updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating replication job", err.Error())
		return
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ReplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ReplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteReplicationJob(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting replication job", err.Error())
	}
}

func (r *ReplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// expected format: "{vmid}-{jobid}" e.g. "100-0"
	parts := strings.SplitN(req.ID, "-", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Format: <vmid>-<job_id> (e.g. '100-0')")
		return
	}

	vmid, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid VMID", err.Error())
		return
	}
	jobID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid job ID", err.Error())
		return
	}

	state := ReplicationResourceModel{
		ID:         types.StringValue(req.ID),
		JobID:      types.Int64Value(jobID),
		SourceVMID: types.Int64Value(vmid),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ReplicationResource) readIntoModel(ctx context.Context, model *ReplicationResourceModel, diagnostics *diag.Diagnostics) {
	job, err := r.client.GetReplicationJob(ctx, model.ID.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok {
			if apiErr.IsNotFound() {
				diagnostics.AddWarning("Replication job not found", "The replication job no longer exists.")
				return
			}
		}
		diagnostics.AddError("Error reading replication job", err.Error())
		return
	}

	model.TargetNode = types.StringValue(job.Target)
	model.SourceVMID = types.Int64Value(int64(job.Guest))
	model.Schedule = types.StringValue(job.Schedule)
	model.Rate = types.Float64Value(job.Rate)
	model.Comment = types.StringValue(job.Comment)
	if job.Disable != nil {
		model.Enabled = types.BoolValue(*job.Disable == 0)
	} else {
		model.Enabled = types.BoolValue(true)
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
