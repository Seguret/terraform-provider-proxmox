package user_password

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ resource.Resource = &UserPasswordResource{}

type UserPasswordResource struct {
	client *client.Client
}

type UserPasswordResourceModel struct {
	ID       types.String `tfsdk:"id"`
	UserID   types.String `tfsdk:"user_id"`
	Password types.String `tfsdk:"password"`
}

func NewResource() resource.Resource {
	return &UserPasswordResource{}
}

func (r *UserPasswordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_password"
}

func (r *UserPasswordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE user password. This resource allows changing user passwords via the API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "The user ID (e.g., 'root@pam' or 'user@pve').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Description: "The new password for the user.",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func (r *UserPasswordResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cl, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = cl
}

func (r *UserPasswordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserPasswordResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userID := plan.UserID.ValueString()
	password := plan.Password.ValueString()

	if err := r.client.ChangeUserPassword(ctx, userID, password); err != nil {
		resp.Diagnostics.AddError("Error changing user password", err.Error())
		return
	}

	plan.ID = types.StringValue(userID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserPasswordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserPasswordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// passwords cant be read back from the API, just keep what we have in state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserPasswordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserPasswordResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userID := plan.UserID.ValueString()
	password := plan.Password.ValueString()

	if err := r.client.ChangeUserPassword(ctx, userID, password); err != nil {
		resp.Diagnostics.AddError("Error changing user password", err.Error())
		return
	}

	plan.ID = types.StringValue(userID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserPasswordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// cant really delete a password — just drop it from state
}
