package ha_resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &HAResourceResource{}
var _ resource.ResourceWithConfigure = &HAResourceResource{}
var _ resource.ResourceWithImportState = &HAResourceResource{}

type HAResourceResource struct {
	client *client.Client
}

type HAResourceResourceModel struct {
	ID          types.String `tfsdk:"id"`
	SID         types.String `tfsdk:"sid"`
	Type        types.String `tfsdk:"type"`
	State       types.String `tfsdk:"state"`
	Group       types.String `tfsdk:"group"`
	MaxRestart  types.Int64  `tfsdk:"max_restart"`
	MaxRelocate types.Int64  `tfsdk:"max_relocate"`
	Comment     types.String `tfsdk:"comment"`
}

func NewResource() resource.Resource {
	return &HAResourceResource{}
}

func (r *HAResourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_ha_resource"
}

func (r *HAResourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE High Availability resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sid": schema.StringAttribute{
				Description: "The HA resource SID (e.g., 'vm:100' or 'ct:200').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The resource type (vm or ct).",
				Optional:    true,
				Computed:    true,
			},
			"state": schema.StringAttribute{
				Description: "The desired state (started, stopped, enabled, disabled, ignored).",
				Optional:    true,
				Computed:    true,
			},
			"group": schema.StringAttribute{
				Description: "The HA group name.",
				Optional:    true,
				Computed:    true,
			},
			"max_restart": schema.Int64Attribute{
				Description: "Maximum number of restart attempts.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
			"max_relocate": schema.Int64Attribute{
				Description: "Maximum number of relocation attempts.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
			"comment": schema.StringAttribute{
				Description: "HA resource description.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *HAResourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *HAResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan HAResourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &models.HAResourceCreateRequest{
		SID:         plan.SID.ValueString(),
		Type:        plan.Type.ValueString(),
		State:       plan.State.ValueString(),
		Group:       plan.Group.ValueString(),
		MaxRestart:  int(plan.MaxRestart.ValueInt64()),
		MaxRelocate: int(plan.MaxRelocate.ValueInt64()),
		Comment:     plan.Comment.ValueString(),
	}

	tflog.Debug(ctx, "Creating HA resource", map[string]any{"sid": plan.SID.ValueString()})

	if err := r.client.CreateHAResource(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating HA resource", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.SID.ValueString())

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading HA resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *HAResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state HAResourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading HA resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *HAResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan HAResourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	maxRestart := int(plan.MaxRestart.ValueInt64())
	maxRelocate := int(plan.MaxRelocate.ValueInt64())

	updateReq := &models.HAResourceUpdateRequest{
		State:       plan.State.ValueString(),
		Group:       plan.Group.ValueString(),
		MaxRestart:  &maxRestart,
		MaxRelocate: &maxRelocate,
		Comment:     plan.Comment.ValueString(),
	}

	if err := r.client.UpdateHAResource(ctx, plan.SID.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating HA resource", err.Error())
		return
	}

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading HA resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *HAResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state HAResourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteHAResource(ctx, state.SID.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting HA resource", err.Error())
	}
}

func (r *HAResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := HAResourceResourceModel{
		ID:  types.StringValue(req.ID),
		SID: types.StringValue(req.ID),
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing HA resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *HAResourceResource) readIntoModel(ctx context.Context, model *HAResourceResourceModel) error {
	ha, err := r.client.GetHAResource(ctx, model.SID.ValueString())
	if err != nil {
		return err
	}
	model.SID = types.StringValue(ha.SID)
	model.Type = types.StringValue(ha.Type)
	model.State = types.StringValue(ha.State)
	model.Group = types.StringValue(ha.Group)
	model.MaxRestart = types.Int64Value(int64(ha.MaxRestart))
	model.MaxRelocate = types.Int64Value(int64(ha.MaxRelocate))
	model.Comment = types.StringValue(ha.Comment)
	return nil
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
