package acl

import (
	"context"
	"fmt"
	"strings"

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

var _ resource.Resource = &ACLResource{}
var _ resource.ResourceWithConfigure = &ACLResource{}

type ACLResource struct {
	client *client.Client
}

type ACLResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Path      types.String `tfsdk:"path"`
	RoleID    types.String `tfsdk:"role_id"`
	UserID    types.String `tfsdk:"user_id"`
	GroupID   types.String `tfsdk:"group_id"`
	Propagate types.Bool   `tfsdk:"propagate"`
}

func NewResource() resource.Resource {
	return &ACLResource{}
}

func (r *ACLResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_acl"
}

func (r *ACLResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE access control list entry.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"path": schema.StringAttribute{
				Description: "The access control path (e.g., '/', '/vms/100', '/storage/local').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_id": schema.StringAttribute{
				Description: "The role to assign.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "The user to assign the role to. Mutually exclusive with group_id.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_id": schema.StringAttribute{
				Description: "The group to assign the role to. Mutually exclusive with user_id.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"propagate": schema.BoolAttribute{
				Description: "Whether to propagate the ACL to child objects.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *ACLResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ACLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ACLResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.UserID.IsNull() && plan.GroupID.IsNull() {
		resp.Diagnostics.AddError("Invalid ACL", "Either user_id or group_id must be specified.")
		return
	}

	propagate := 1
	if !plan.Propagate.IsNull() && !plan.Propagate.ValueBool() {
		propagate = 0
	}

	aclReq := &models.ACLUpdateRequest{
		Path:      plan.Path.ValueString(),
		Roles:     plan.RoleID.ValueString(),
		Users:     plan.UserID.ValueString(),
		Groups:    plan.GroupID.ValueString(),
		Propagate: &propagate,
	}

	tflog.Debug(ctx, "Creating ACL", map[string]any{"path": aclReq.Path, "role": aclReq.Roles})

	if err := r.client.UpdateACL(ctx, aclReq); err != nil {
		resp.Diagnostics.AddError("Error creating ACL", err.Error())
		return
	}

	plan.ID = types.StringValue(r.buildID(&plan))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ACLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ACLResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	aclList, err := r.client.GetACL(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ACL", err.Error())
		return
	}

	found := false
	for _, entry := range aclList {
		if entry.Path == state.Path.ValueString() && entry.RoleID == state.RoleID.ValueString() {
			if (entry.Type == "user" && entry.UGid == state.UserID.ValueString()) ||
				(entry.Type == "group" && entry.UGid == state.GroupID.ValueString()) {
				if entry.Propagate != nil {
					state.Propagate = types.BoolValue(*entry.Propagate == 1)
				}
				found = true
				break
			}
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ACLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// all key fields are ForceNew, so only propagate can actually change here
	var plan ACLResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	propagate := 1
	if !plan.Propagate.IsNull() && !plan.Propagate.ValueBool() {
		propagate = 0
	}

	aclReq := &models.ACLUpdateRequest{
		Path:      plan.Path.ValueString(),
		Roles:     plan.RoleID.ValueString(),
		Users:     plan.UserID.ValueString(),
		Groups:    plan.GroupID.ValueString(),
		Propagate: &propagate,
	}

	if err := r.client.UpdateACL(ctx, aclReq); err != nil {
		resp.Diagnostics.AddError("Error updating ACL", err.Error())
		return
	}

	plan.ID = types.StringValue(r.buildID(&plan))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ACLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ACLResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteFlag := 1
	aclReq := &models.ACLUpdateRequest{
		Path:   state.Path.ValueString(),
		Roles:  state.RoleID.ValueString(),
		Users:  state.UserID.ValueString(),
		Groups: state.GroupID.ValueString(),
		Delete: &deleteFlag,
	}

	if err := r.client.UpdateACL(ctx, aclReq); err != nil {
		resp.Diagnostics.AddError("Error deleting ACL", err.Error())
	}
}

func (r *ACLResource) buildID(model *ACLResourceModel) string {
	ugid := model.UserID.ValueString()
	if ugid == "" {
		ugid = model.GroupID.ValueString()
	}
	return strings.Join([]string{model.Path.ValueString(), model.RoleID.ValueString(), ugid}, "::")
}
