package user_token

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

var _ resource.Resource = &UserTokenResource{}
var _ resource.ResourceWithConfigure = &UserTokenResource{}
var _ resource.ResourceWithImportState = &UserTokenResource{}

type UserTokenResource struct {
	client *client.Client
}

type UserTokenResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	UserID              types.String `tfsdk:"user_id"`
	TokenID             types.String `tfsdk:"token_id"`
	Comment             types.String `tfsdk:"comment"`
	Expire              types.Int64  `tfsdk:"expire"`
	PrivilegesSeparated types.Bool   `tfsdk:"privileges_separation"`
	Value               types.String `tfsdk:"value"`
}

func NewResource() resource.Resource {
	return &UserTokenResource{}
}

func (r *UserTokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_user_token"
}

func (r *UserTokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE user API token.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The full token identifier (userid/tokenid).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "The user ID (e.g., 'root@pam').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"token_id": schema.StringAttribute{
				Description: "The token name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Description: "A comment for the token.",
				Optional:    true,
				Computed:    true,
			},
			"expire": schema.Int64Attribute{
				Description: "Token expiration date as a UNIX timestamp. 0 means no expiration.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"privileges_separation": schema.BoolAttribute{
				Description: "Whether the token has privilege separation (tokens cannot exceed user privileges when true).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"value": schema.StringAttribute{
				Description: "The token secret. Only available after creation; stored in state.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (r *UserTokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserTokenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	privsep := boolToInt(plan.PrivilegesSeparated.ValueBool())
	createReq := &models.UserTokenCreateRequest{
		Comment: plan.Comment.ValueString(),
		Expire:  plan.Expire.ValueInt64(),
		Privsep: &privsep,
	}

	tflog.Debug(ctx, "Creating user token", map[string]any{
		"user_id":  plan.UserID.ValueString(),
		"token_id": plan.TokenID.ValueString(),
	})

	result, err := r.client.CreateUserToken(ctx, plan.UserID.ValueString(), plan.TokenID.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user token", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.UserID.ValueString(), plan.TokenID.ValueString()))
	// token value is only returned once on creation — grab it now
	plan.Value = types.StringValue(result.Value)

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserTokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// value stays from state — the API wont give it back
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserTokenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state UserTokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	privsep := boolToInt(plan.PrivilegesSeparated.ValueBool())
	expire := plan.Expire.ValueInt64()
	updateReq := &models.UserTokenUpdateRequest{
		Comment: plan.Comment.ValueString(),
		Expire:  &expire,
		Privsep: &privsep,
	}

	if err := r.client.UpdateUserToken(ctx, plan.UserID.ValueString(), plan.TokenID.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating user token", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.UserID.ValueString(), plan.TokenID.ValueString()))
	// keep the secret value from state since the API doesnt return it
	plan.Value = state.Value

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserTokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteUserToken(ctx, state.UserID.ValueString(), state.TokenID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting user token", err.Error())
	}
}

func (r *UserTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// expected format: "userid/tokenid" e.g. "root@pam/my-token"
	// userID can have "@" in it but tokenid wont have "/"
	// so split on the last "/" to isolate the tokenid
	lastSlash := strings.LastIndex(req.ID, "/")
	if lastSlash < 0 {
		resp.Diagnostics.AddError("Invalid import ID", "Format: <user_id>/<token_id> (e.g. 'root@pam/my-token')")
		return
	}

	userID := req.ID[:lastSlash]
	tokenID := req.ID[lastSlash+1:]

	state := UserTokenResourceModel{
		ID:      types.StringValue(req.ID),
		UserID:  types.StringValue(userID),
		TokenID: types.StringValue(tokenID),
		// value cant be recovered on import — zero it out
		Value: types.StringValue(""),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserTokenResource) readIntoModel(ctx context.Context, model *UserTokenResourceModel, diagnostics *diag.Diagnostics) {
	token, err := r.client.GetUserToken(ctx, model.UserID.ValueString(), model.TokenID.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok {
			if apiErr.IsNotFound() {
				diagnostics.AddWarning("User token not found", "The token no longer exists.")
				return
			}
		}
		diagnostics.AddError("Error reading user token", err.Error())
		return
	}

	model.Comment = types.StringValue(token.Comment)
	model.Expire = types.Int64Value(token.Expire)
	if token.Privsep != nil {
		model.PrivilegesSeparated = types.BoolValue(*token.Privsep == 1)
	} else {
		model.PrivilegesSeparated = types.BoolValue(true)
	}
	// dont touch the stored secret — only update the other fields
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
