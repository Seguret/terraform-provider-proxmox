package ha_group

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

var _ resource.Resource = &HAGroupResource{}
var _ resource.ResourceWithConfigure = &HAGroupResource{}
var _ resource.ResourceWithImportState = &HAGroupResource{}

type HAGroupResource struct {
	client *client.Client
}

type HAGroupResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Group      types.String `tfsdk:"group"`
	Nodes      types.String `tfsdk:"nodes"`
	Restricted types.Bool   `tfsdk:"restricted"`
	NoFailback types.Bool   `tfsdk:"no_failback"`
	Comment    types.String `tfsdk:"comment"`
}

func NewResource() resource.Resource {
	return &HAGroupResource{}
}

func (r *HAGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_ha_group"
}

func (r *HAGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE High Availability group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group": schema.StringAttribute{
				Description: "The HA group name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"nodes": schema.StringAttribute{
				Description: "Comma-separated list of nodes with optional priority (e.g., 'node1:2,node2:1').",
				Required:    true,
			},
			"restricted": schema.BoolAttribute{
				Description: "Whether the HA group is restricted to its members.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"no_failback": schema.BoolAttribute{
				Description: "Whether to disable automatic failback.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"comment": schema.StringAttribute{
				Description: "HA group description.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *HAGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func boolToIntPtr(b bool) *int {
	v := 0
	if b {
		v = 1
	}
	return &v
}

func (r *HAGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan HAGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &models.HAGroupCreateRequest{
		Group:      plan.Group.ValueString(),
		Nodes:      plan.Nodes.ValueString(),
		Restricted: boolToIntPtr(plan.Restricted.ValueBool()),
		NoFailback: boolToIntPtr(plan.NoFailback.ValueBool()),
		Comment:    plan.Comment.ValueString(),
	}

	tflog.Debug(ctx, "Creating HA group", map[string]any{"group": plan.Group.ValueString()})

	if err := r.client.CreateHAGroup(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating HA group", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Group.ValueString())

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading HA group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *HAGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state HAGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading HA group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *HAGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan HAGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &models.HAGroupUpdateRequest{
		Nodes:      plan.Nodes.ValueString(),
		Restricted: boolToIntPtr(plan.Restricted.ValueBool()),
		NoFailback: boolToIntPtr(plan.NoFailback.ValueBool()),
		Comment:    plan.Comment.ValueString(),
	}

	if err := r.client.UpdateHAGroup(ctx, plan.Group.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating HA group", err.Error())
		return
	}

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading HA group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *HAGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state HAGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteHAGroup(ctx, state.Group.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting HA group", err.Error())
	}
}

func (r *HAGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := HAGroupResourceModel{
		ID:    types.StringValue(req.ID),
		Group: types.StringValue(req.ID),
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing HA group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *HAGroupResource) readIntoModel(ctx context.Context, model *HAGroupResourceModel) error {
	g, err := r.client.GetHAGroup(ctx, model.Group.ValueString())
	if err != nil {
		return err
	}
	model.Group = types.StringValue(g.Group)
	model.Nodes = types.StringValue(g.Nodes)
	if g.Restricted != nil {
		model.Restricted = types.BoolValue(*g.Restricted == 1)
	} else {
		model.Restricted = types.BoolValue(false)
	}
	if g.NoFailback != nil {
		model.NoFailback = types.BoolValue(*g.NoFailback == 1)
	} else {
		model.NoFailback = types.BoolValue(false)
	}
	model.Comment = types.StringValue(g.Comment)
	return nil
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
