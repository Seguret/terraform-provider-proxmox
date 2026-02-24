package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithConfigure = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}

type UserResource struct {
	client *client.Client
}

type UserResourceModel struct {
	ID        types.String `tfsdk:"id"`
	UserID    types.String `tfsdk:"user_id"`
	Password  types.String `tfsdk:"password"`
	Email     types.String `tfsdk:"email"`
	Enabled   types.Bool   `tfsdk:"enabled"`
	Expire    types.Int64  `tfsdk:"expire"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	Comment   types.String `tfsdk:"comment"`
	Groups    types.String `tfsdk:"groups"`
	Keys      types.String `tfsdk:"keys"`
}

func NewResource() resource.Resource {
	return &UserResource{}
}

func (r *UserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_user"
}

func (r *UserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "The user identifier (e.g., 'user@pam' or 'user@pve').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Description: "The user password. Only used during creation for PVE realm users.",
				Optional:    true,
				Sensitive:   true,
			},
			"email": schema.StringAttribute{
				Description: "The user's email address.",
				Optional:    true,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the user account is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"expire": schema.Int64Attribute{
				Description: "Account expiration date (UNIX epoch). 0 means no expiration.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"first_name": schema.StringAttribute{
				Description: "The user's first name.",
				Optional:    true,
				Computed:    true,
			},
			"last_name": schema.StringAttribute{
				Description: "The user's last name.",
				Optional:    true,
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "A comment for the user.",
				Optional:    true,
				Computed:    true,
			},
			"groups": schema.StringAttribute{
				Description: "Comma-separated list of groups.",
				Optional:    true,
				Computed:    true,
			},
			"keys": schema.StringAttribute{
				Description: "Keys for two-factor authentication.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *UserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enableInt := boolToInt(plan.Enabled.ValueBool())

	createReq := &models.UserCreateRequest{
		UserID:    plan.UserID.ValueString(),
		Password:  plan.Password.ValueString(),
		Email:     plan.Email.ValueString(),
		Enable:    &enableInt,
		Expire:    plan.Expire.ValueInt64(),
		FirstName: plan.FirstName.ValueString(),
		LastName:  plan.LastName.ValueString(),
		Comment:   plan.Comment.ValueString(),
		Groups:    plan.Groups.ValueString(),
		Keys:      plan.Keys.ValueString(),
	}

	tflog.Debug(ctx, "Creating user", map[string]any{"user_id": createReq.UserID})

	if err := r.client.CreateUser(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.UserID.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enableInt := boolToInt(plan.Enabled.ValueBool())
	email := plan.Email.ValueString()
	expire := plan.Expire.ValueInt64()
	firstName := plan.FirstName.ValueString()
	lastName := plan.LastName.ValueString()
	comment := plan.Comment.ValueString()
	groups := plan.Groups.ValueString()
	keys := plan.Keys.ValueString()

	updateReq := &models.UserUpdateRequest{
		Email:     &email,
		Enable:    &enableInt,
		Expire:    &expire,
		FirstName: &firstName,
		LastName:  &lastName,
		Comment:   &comment,
		Groups:    &groups,
		Keys:      &keys,
	}

	if err := r.client.UpdateUser(ctx, plan.UserID.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.UserID.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteUser(ctx, state.UserID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	userID := req.ID
	if !strings.Contains(userID, "@") {
		resp.Diagnostics.AddError("Invalid user ID", "User ID must be in format 'user@realm'")
		return
	}

	state := UserResourceModel{
		ID:     types.StringValue(userID),
		UserID: types.StringValue(userID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserResource) readIntoModel(ctx context.Context, model *UserResourceModel, diagnostics *diag.Diagnostics) {
	user, err := r.client.GetUser(ctx, model.UserID.ValueString())
	if err != nil {
		diagnostics.AddError("Error reading user", err.Error())
		return
	}

	model.Email = types.StringValue(user.Email)
	if user.Enable != nil {
		model.Enabled = types.BoolValue(*user.Enable == 1)
	} else {
		model.Enabled = types.BoolValue(true)
	}
	model.Expire = types.Int64Value(user.Expire)
	model.FirstName = types.StringValue(user.FirstName)
	model.LastName = types.StringValue(user.LastName)
	model.Comment = types.StringValue(user.Comment)
	model.Groups = types.StringValue(user.Groups)
	model.Keys = types.StringValue(user.Keys)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
