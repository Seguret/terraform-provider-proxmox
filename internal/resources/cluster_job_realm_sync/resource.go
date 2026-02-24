package cluster_job_realm_sync

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &ClusterJobRealmSyncResource{}
var _ resource.ResourceWithConfigure = &ClusterJobRealmSyncResource{}
var _ resource.ResourceWithImportState = &ClusterJobRealmSyncResource{}

type ClusterJobRealmSyncResource struct {
	client *client.Client
}

type ClusterJobRealmSyncResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Realm          types.String `tfsdk:"realm"`
	Schedule       types.String `tfsdk:"schedule"`
	Scope          types.String `tfsdk:"scope"`
	RemoveVanished types.String `tfsdk:"remove_vanished"`
	EnableNew      types.Bool   `tfsdk:"enable_new"`
	Comment        types.String `tfsdk:"comment"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

func NewResource() resource.Resource {
	return &ClusterJobRealmSyncResource{}
}

func (r *ClusterJobRealmSyncResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_cluster_job_realm_sync"
}

func (r *ClusterJobRealmSyncResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE cluster realm-sync scheduled job.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The realm-sync job identifier.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"realm": schema.StringAttribute{
				Description: "The authentication realm to synchronize.",
				Required:    true,
			},
			"schedule": schema.StringAttribute{
				Description: "The job schedule in systemd calendar format (e.g. 'daily', '0 3 * * *').",
				Required:    true,
			},
			"scope": schema.StringAttribute{
				Description: "What to sync: 'users', 'groups', or 'both'.",
				Optional:    true,
				Computed:    true,
			},
			"remove_vanished": schema.StringAttribute{
				Description: "Comma-separated list of objects to remove when they vanish from the LDAP directory (acl, entry, properties).",
				Optional:    true,
				Computed:    true,
			},
			"enable_new": schema.BoolAttribute{
				Description: "Whether to enable newly synced users.",
				Optional:    true,
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "A description for this sync job.",
				Optional:    true,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the sync job is enabled.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *ClusterJobRealmSyncResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ClusterJobRealmSyncResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ClusterJobRealmSyncResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// use the realm name as job ID — proxmox needs a unique ID per job
	jobID := plan.Realm.ValueString()

	tflog.Debug(ctx, "Creating cluster realm-sync job", map[string]any{
		"id":    jobID,
		"realm": plan.Realm.ValueString(),
	})

	if err := r.client.CreateClusterJobRealmSync(ctx, jobID, &models.ClusterJobRealmSyncCreateRequest{
		Realm:          plan.Realm.ValueString(),
		Schedule:       plan.Schedule.ValueString(),
		Scope:          plan.Scope.ValueString(),
		RemoveVanished: plan.RemoveVanished.ValueString(),
		EnableNew:      plan.EnableNew.ValueBool(),
		Comment:        plan.Comment.ValueString(),
		Enabled:        plan.Enabled.ValueBool(),
	}); err != nil {
		resp.Diagnostics.AddError("Error creating cluster realm-sync job", err.Error())
		return
	}

	plan.ID = types.StringValue(jobID)
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ClusterJobRealmSyncResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ClusterJobRealmSyncResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if isNotFound := r.readIntoModelWithNotFound(ctx, &state, &resp.Diagnostics); isNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ClusterJobRealmSyncResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ClusterJobRealmSyncResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enableNew := plan.EnableNew.ValueBool()
	enabled := plan.Enabled.ValueBool()

	if err := r.client.UpdateClusterJobRealmSync(ctx, plan.ID.ValueString(), &models.ClusterJobRealmSyncUpdateRequest{
		Schedule:       plan.Schedule.ValueString(),
		Scope:          plan.Scope.ValueString(),
		RemoveVanished: plan.RemoveVanished.ValueString(),
		EnableNew:      &enableNew,
		Comment:        plan.Comment.ValueString(),
		Enabled:        &enabled,
	}); err != nil {
		resp.Diagnostics.AddError("Error updating cluster realm-sync job", err.Error())
		return
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ClusterJobRealmSyncResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ClusterJobRealmSyncResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteClusterJobRealmSync(ctx, state.ID.ValueString()); err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return
		}
		resp.Diagnostics.AddError("Error deleting cluster realm-sync job", err.Error())
	}
}

func (r *ClusterJobRealmSyncResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := ClusterJobRealmSyncResourceModel{
		ID: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// readIntoModelWithNotFound fetches realm-sync job state and returns true if the job doesnt exist.
func (r *ClusterJobRealmSyncResource) readIntoModelWithNotFound(ctx context.Context, model *ClusterJobRealmSyncResourceModel, diags interface{ AddError(string, string) }) bool {
	job, err := r.client.GetClusterJobRealmSync(ctx, model.ID.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return true
		}
		diags.AddError("Error reading cluster realm-sync job", err.Error())
		return false
	}
	populateModel(model, job)
	return false
}

// readIntoModel fetches the realm-sync job and populates the model.
func (r *ClusterJobRealmSyncResource) readIntoModel(ctx context.Context, model *ClusterJobRealmSyncResourceModel, diags interface{ AddError(string, string) }) {
	job, err := r.client.GetClusterJobRealmSync(ctx, model.ID.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diags.AddError("Cluster realm-sync job not found", "The cluster realm-sync job no longer exists.")
			return
		}
		diags.AddError("Error reading cluster realm-sync job", err.Error())
		return
	}
	populateModel(model, job)
}

func populateModel(model *ClusterJobRealmSyncResourceModel, job *models.ClusterJobRealmSync) {
	model.Realm = types.StringValue(job.Realm)
	model.Schedule = types.StringValue(job.Schedule)
	model.Scope = types.StringValue(job.Scope)
	model.RemoveVanished = types.StringValue(job.RemoveVanished)
	model.EnableNew = types.BoolValue(job.EnableNew)
	model.Comment = types.StringValue(job.Comment)
	model.Enabled = types.BoolValue(job.Enabled)
}
