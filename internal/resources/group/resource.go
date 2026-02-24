package group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &GroupResource{}
var _ resource.ResourceWithConfigure = &GroupResource{}
var _ resource.ResourceWithImportState = &GroupResource{}

type GroupResource struct {
	client *client.Client
}

type GroupResourceModel struct {
	ID      types.String `tfsdk:"id"`
	GroupID types.String `tfsdk:"group_id"`
	Comment types.String `tfsdk:"comment"`
	Members types.String `tfsdk:"members"`
}

func NewResource() resource.Resource {
	return &GroupResource{}
}

func (r *GroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_group"
}

func (r *GroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Description: "The group identifier.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Description: "A comment for the group.",
				Optional:    true,
				Computed:    true,
			},
			"members": schema.StringAttribute{
				Description: "The group members (comma-separated user IDs). Read-only.",
				Computed:    true,
			},
		},
	}
}

func (r *GroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan GroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &models.GroupCreateRequest{
		GroupID: plan.GroupID.ValueString(),
		Comment: plan.Comment.ValueString(),
	}

	tflog.Debug(ctx, "Creating group", map[string]any{"group_id": createReq.GroupID})

	if err := r.client.CreateGroup(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating group", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.GroupID.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	comment := plan.Comment.ValueString()
	updateReq := &models.GroupUpdateRequest{Comment: &comment}

	if err := r.client.UpdateGroup(ctx, plan.GroupID.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating group", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.GroupID.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteGroup(ctx, state.GroupID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting group", err.Error())
	}
}

func (r *GroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := GroupResourceModel{
		ID:      types.StringValue(req.ID),
		GroupID: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *GroupResource) readIntoModel(ctx context.Context, model *GroupResourceModel, diagnostics *diag.Diagnostics) {
	grp, err := r.client.GetGroup(ctx, model.GroupID.ValueString())
	if err != nil {
		diagnostics.AddError("Error reading group", err.Error())
		return
	}
	model.Comment = types.StringValue(grp.Comment)
	
	// members can come from either field depending on the API version
	members := grp.Members
	if len(members) == 0 {
		members = grp.Users
	}

	var membersStr string
	if len(members) > 0 {
		// build comma-separated string from the array
		for i, member := range members {
			if i > 0 {
				membersStr += ","
			}
			membersStr += member
		}
	}
	model.Members = types.StringValue(membersStr)
}
