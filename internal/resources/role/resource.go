package role

import (
	"context"
	"fmt"
	"sort"
	"strings"

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

var _ resource.Resource = &RoleResource{}
var _ resource.ResourceWithConfigure = &RoleResource{}
var _ resource.ResourceWithImportState = &RoleResource{}

type RoleResource struct {
	client *client.Client
}

type RoleResourceModel struct {
	ID         types.String `tfsdk:"id"`
	RoleID     types.String `tfsdk:"role_id"`
	Privileges types.String `tfsdk:"privileges"`
}

func NewResource() resource.Resource {
	return &RoleResource{}
}

func (r *RoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_role"
}

func (r *RoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE role.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"role_id": schema.StringAttribute{
				Description: "The role identifier.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"privileges": schema.StringAttribute{
				Description: "Comma-separated list of privileges (e.g., 'VM.Audit,VM.Console').",
				Required:    true,
			},
		},
	}
}

func (r *RoleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RoleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// sort privileges so the order matches what the API returns
	privs := normalizePrivileges(plan.Privileges.ValueString())

	createReq := &models.RoleCreateRequest{
		RoleID: plan.RoleID.ValueString(),
		Privs:  privs,
	}

	tflog.Debug(ctx, "Creating role", map[string]any{"role_id": createReq.RoleID})

	if err := r.client.CreateRole(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating role", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.RoleID.ValueString())
	plan.Privileges = types.StringValue(privs)
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RoleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RoleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// sort privileges to match API response order
	privs := normalizePrivileges(plan.Privileges.ValueString())
	updateReq := &models.RoleUpdateRequest{Privs: privs}

	if err := r.client.UpdateRole(ctx, plan.RoleID.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating role", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.RoleID.ValueString())
	plan.Privileges = types.StringValue(privs)
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RoleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteRole(ctx, state.RoleID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting role", err.Error())
	}
}

func (r *RoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := RoleResourceModel{
		ID:     types.StringValue(req.ID),
		RoleID: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RoleResource) readIntoModel(ctx context.Context, model *RoleResourceModel, diagnostics *diag.Diagnostics) {
	// the single-role endpoint doesnt return privs, so grab the full list and filter
	roles, err := r.client.GetRoles(ctx)
	if err != nil {
		diagnostics.AddError("Error reading roles", err.Error())
		return
	}
	
	roleID := model.RoleID.ValueString()
	for _, role := range roles {
		if role.RoleID == roleID {
			tflog.Debug(ctx, "Found role", map[string]any{
				"role_id": role.RoleID,
				"privs":   role.Privs,
			})
			model.Privileges = types.StringValue(role.Privs)
			return
		}
	}
	
	diagnostics.AddError("Role not found", fmt.Sprintf("Role %s not found", roleID))
}

// normalizePrivileges sorts the comma-separated privilege list alphabetically.
// proxmox returns them sorted so we need to match that for clean diffs.
func normalizePrivileges(privs string) string {
	if privs == "" {
		return ""
	}
	parts := strings.Split(privs, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}
