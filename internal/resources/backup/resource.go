package backup

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

var _ resource.Resource = &BackupResource{}
var _ resource.ResourceWithConfigure = &BackupResource{}
var _ resource.ResourceWithImportState = &BackupResource{}

type BackupResource struct {
	client *client.Client
}

type BackupResourceModel struct {
	ID               types.String  `tfsdk:"id"`
	Storage          types.String  `tfsdk:"storage"`
	Schedule         types.String  `tfsdk:"schedule"`
	VMIDs            types.String  `tfsdk:"vmids"`
	Nodes            types.String  `tfsdk:"nodes"`
	All              types.Bool    `tfsdk:"all"`
	Compress         types.String  `tfsdk:"compress"`
	Mode             types.String  `tfsdk:"mode"`
	Comment          types.String  `tfsdk:"comment"`
	Mailto           types.String  `tfsdk:"mailto"`
	MailNotification types.String  `tfsdk:"mail_notification"`
	MaxFiles         types.Int64   `tfsdk:"max_files"`
	Enabled          types.Bool    `tfsdk:"enabled"`
	BWLimit          types.Float64 `tfsdk:"bw_limit"`
	NotesTemplate    types.String  `tfsdk:"notes_template"`
}

func NewResource() resource.Resource {
	return &BackupResource{}
}

func (r *BackupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_backup"
}

func (r *BackupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE vzdump backup schedule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The backup job ID (assigned by Proxmox).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"storage": schema.StringAttribute{
				Description: "The target storage for backups.",
				Required:    true,
			},
			"schedule": schema.StringAttribute{
				Description: "The backup schedule in systemd calendar format (e.g. 'daily', 'Mon,Tue 02:00').",
				Required:    true,
			},
			"vmids": schema.StringAttribute{
				Description: "Comma-separated list of VMID(s) to back up. Leave empty when 'all' is true.",
				Optional:    true,
				Computed:    true,
			},
			"nodes": schema.StringAttribute{
				Description: "Comma-separated list of nodes to run the job on. Empty means all nodes.",
				Optional:    true,
				Computed:    true,
			},
			"all": schema.BoolAttribute{
				Description: "Back up all VMs and containers.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"compress": schema.StringAttribute{
				Description: "Compression algorithm ('0', '1', 'gzip', 'lzo', 'zstd').",
				Optional:    true,
				Computed:    true,
			},
			"mode": schema.StringAttribute{
				Description: "Backup mode ('snapshot', 'suspend', 'stop').",
				Optional:    true,
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "A description for this backup job.",
				Optional:    true,
				Computed:    true,
			},
			"mailto": schema.StringAttribute{
				Description: "Comma-separated list of email addresses to notify.",
				Optional:    true,
				Computed:    true,
			},
			"mail_notification": schema.StringAttribute{
				Description: "Email notification mode ('always', 'failure').",
				Optional:    true,
				Computed:    true,
			},
			"max_files": schema.Int64Attribute{
				Description: "Maximum number of backups to keep (deprecated; use prune settings).",
				Optional:    true,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the backup job is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"bw_limit": schema.Float64Attribute{
				Description: "Bandwidth limit in KiB/s.",
				Optional:    true,
				Computed:    true,
			},
			"notes_template": schema.StringAttribute{
				Description: "Template for backup notes.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *BackupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BackupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	allInt := boolToIntPtr(plan.All.ValueBool())
	enabledInt := boolToIntPtr(plan.Enabled.ValueBool())

	createReq := &models.BackupJobCreateRequest{
		Storage:          plan.Storage.ValueString(),
		Schedule:         plan.Schedule.ValueString(),
		VMIDs:            plan.VMIDs.ValueString(),
		Nodes:            plan.Nodes.ValueString(),
		All:              allInt,
		Compress:         plan.Compress.ValueString(),
		Mode:             plan.Mode.ValueString(),
		Comment:          plan.Comment.ValueString(),
		Mailto:           plan.Mailto.ValueString(),
		MailNotification: plan.MailNotification.ValueString(),
		MaxFiles:         int(plan.MaxFiles.ValueInt64()),
		Remove:           enabledInt,
		BWLimit:          plan.BWLimit.ValueFloat64(),
		NotesTemplate:    plan.NotesTemplate.ValueString(),
	}

	tflog.Debug(ctx, "Creating backup job", map[string]any{"storage": createReq.Storage, "schedule": createReq.Schedule})

	id, err := r.client.CreateBackupJob(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating backup job", err.Error())
		return
	}

	plan.ID = types.StringValue(id)
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state BackupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BackupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	allInt := boolToIntPtr(plan.All.ValueBool())
	enabledInt := boolToIntPtr(plan.Enabled.ValueBool())

	updateReq := &models.BackupJobUpdateRequest{
		Storage:          plan.Storage.ValueString(),
		Schedule:         plan.Schedule.ValueString(),
		VMIDs:            plan.VMIDs.ValueString(),
		Nodes:            plan.Nodes.ValueString(),
		All:              allInt,
		Compress:         plan.Compress.ValueString(),
		Mode:             plan.Mode.ValueString(),
		Comment:          plan.Comment.ValueString(),
		Mailto:           plan.Mailto.ValueString(),
		MailNotification: plan.MailNotification.ValueString(),
		MaxFiles:         int(plan.MaxFiles.ValueInt64()),
		Enabled:          enabledInt,
		BWLimit:          plan.BWLimit.ValueFloat64(),
		NotesTemplate:    plan.NotesTemplate.ValueString(),
	}

	if err := r.client.UpdateBackupJob(ctx, plan.ID.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating backup job", err.Error())
		return
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state BackupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteBackupJob(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting backup job", err.Error())
	}
}

func (r *BackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := BackupResourceModel{ID: types.StringValue(req.ID)}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BackupResource) readIntoModel(ctx context.Context, model *BackupResourceModel, diagnostics interface{ AddError(string, string) }) {
	job, err := r.client.GetBackupJob(ctx, model.ID.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diagnostics.AddError("Backup job not found", "The backup job no longer exists.")
			return
		}
		diagnostics.AddError("Error reading backup job", err.Error())
		return
	}

	model.Storage = types.StringValue(job.Storage)
	model.Schedule = types.StringValue(job.Schedule)
	model.VMIDs = types.StringValue(job.VMIDs)
	model.Nodes = types.StringValue(job.Nodes)
	model.Compress = types.StringValue(job.Compress)
	model.Mode = types.StringValue(job.Mode)
	model.Comment = types.StringValue(job.Comment)
	model.Mailto = types.StringValue(job.Mailto)
	model.MailNotification = types.StringValue(job.MailNotification)
	model.MaxFiles = types.Int64Value(int64(job.MaxFiles))
	model.BWLimit = types.Float64Value(job.BWLimit)
	model.NotesTemplate = types.StringValue(job.NotesTemplate)

	if job.All != nil {
		model.All = types.BoolValue(*job.All == 1)
	} else {
		model.All = types.BoolValue(false)
	}

	if job.Enabled != nil {
		model.Enabled = types.BoolValue(*job.Enabled == 1)
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
