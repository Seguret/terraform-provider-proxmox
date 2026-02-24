package user_tfa

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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

var _ resource.Resource = &UserTFAResource{}
var _ resource.ResourceWithConfigure = &UserTFAResource{}
var _ resource.ResourceWithImportState = &UserTFAResource{}

// UserTFAResource handles TFA entries for Proxmox users.
type UserTFAResource struct {
	client *client.Client
}

type UserTFAResourceModel struct {
	// id is the TFA entry ID proxmox assigns on creation.
	ID types.String `tfsdk:"id"`
	// user_id is the proxmox user (e.g. "root@pam").
	UserID types.String `tfsdk:"user_id"`
	// type is one of: totp, recovery, webauthn, yubico.
	Type types.String `tfsdk:"type"`
	// description is an optional label for the entry.
	Description types.String `tfsdk:"description"`
	// totp is the provisioning URI used during setup — write-only, never returned on Read.
	TOTP types.String `tfsdk:"totp"`
	// value is the proof value (e.g. current TOTP code) used to validate setup. Sensitive.
	Value types.String `tfsdk:"value"`
	// enabled controls whether this TFA entry is active.
	Enabled types.Bool `tfsdk:"enabled"`
	// qrcode is returned on creation for totp/yubico types. Sensitive.
	QRCode types.String `tfsdk:"qrcode"`
	// url is the provisioning URL returned for totp entries on creation.
	URL types.String `tfsdk:"url"`
}

func NewResource() resource.Resource {
	return &UserTFAResource{}
}

func (r *UserTFAResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_user_tfa"
}

func (r *UserTFAResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a two-factor authentication (TFA) entry for a Proxmox VE user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The TFA entry identifier assigned by Proxmox.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "The Proxmox user ID (e.g. 'root@pam').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The TFA type. Valid values: 'totp', 'recovery', 'webauthn', 'yubico'.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A human-readable description for the TFA entry.",
				Optional:    true,
				Computed:    true,
			},
			"totp": schema.StringAttribute{
				Description: "The TOTP provisioning URI or secret used during initial setup. Not returned after creation.",
				Optional:    true,
				Sensitive:   true,
			},
			"value": schema.StringAttribute{
				Description: "The proof value (e.g. current TOTP code) used to validate the TFA entry during creation.",
				Optional:    true,
				Sensitive:   true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the TFA entry is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"qrcode": schema.StringAttribute{
				Description: "The QR code data URI returned on creation (TOTP/yubico types). Sensitive.",
				Computed:    true,
				Sensitive:   true,
			},
			"url": schema.StringAttribute{
				Description: "The provisioning URL returned on creation for TOTP types.",
				Computed:    true,
			},
		},
	}
}

func (r *UserTFAResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserTFAResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserTFAResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &models.TFACreateRequest{
		Type:        plan.Type.ValueString(),
		Description: plan.Description.ValueString(),
		TOTP:        plan.TOTP.ValueString(),
		Value:       plan.Value.ValueString(),
	}

	tflog.Debug(ctx, "Creating user TFA entry", map[string]any{
		"user_id": plan.UserID.ValueString(),
		"type":    plan.Type.ValueString(),
	})

	result, err := r.client.CreateTFAEntry(ctx, plan.UserID.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating TFA entry", err.Error())
		return
	}

	plan.ID = types.StringValue(result.ID)
	// save creation-only fields before they disappear
	plan.QRCode = types.StringValue(result.QRCode)
	plan.URL = types.StringValue(result.URL)

	// read back to get description and enabled state
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserTFAResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserTFAResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// qrcode and url are write-only — keep whatever is already in state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserTFAResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserTFAResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state UserTFAResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// only description and enabled are mutable after creation
	enabled := plan.Enabled.ValueBool()
	updateReq := &models.TFAUpdateRequest{
		Description: plan.Description.ValueString(),
		Enable:      &enabled,
	}

	tflog.Debug(ctx, "Updating user TFA entry", map[string]any{
		"user_id": plan.UserID.ValueString(),
		"id":      plan.ID.ValueString(),
	})

	if err := r.client.UpdateTFAEntry(ctx, plan.UserID.ValueString(), plan.ID.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating TFA entry", err.Error())
		return
	}

	// carry over creation-time fields from existing state
	plan.QRCode = state.QRCode
	plan.URL = state.URL

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserTFAResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserTFAResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting user TFA entry", map[string]any{
		"user_id": state.UserID.ValueString(),
		"id":      state.ID.ValueString(),
	})

	if err := r.client.DeleteTFAEntry(ctx, state.UserID.ValueString(), state.ID.ValueString()); err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return
		}
		resp.Diagnostics.AddError("Error deleting TFA entry", err.Error())
	}
}

func (r *UserTFAResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// expected import format: "user_id:tfa_id" e.g. "root@pam:TOTP-abc123"
	sep := strings.LastIndex(req.ID, ":")
	if sep < 0 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Format must be '<user_id>:<tfa_id>' (e.g. 'root@pam:TOTP-abc123')",
		)
		return
	}

	userID := req.ID[:sep]
	tfaID := req.ID[sep+1:]

	state := UserTFAResourceModel{
		ID:     types.StringValue(tfaID),
		UserID: types.StringValue(userID),
		// creation-time fields like qrcode/url are gone at this point — zero them out
		QRCode: types.StringValue(""),
		URL:    types.StringValue(""),
		TOTP:   types.StringValue(""),
		Value:  types.StringValue(""),
	}

	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// readIntoModel fetches the current TFA entry state and updates mutable fields.
// Doesnt touch write-only fields like totp, value, qrcode, or url.
func (r *UserTFAResource) readIntoModel(ctx context.Context, model *UserTFAResourceModel, diagnostics *diag.Diagnostics) {
	entry, err := r.client.GetTFAEntry(ctx, model.UserID.ValueString(), model.ID.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diagnostics.AddWarning("TFA entry not found", "The TFA entry no longer exists in Proxmox.")
			return
		}
		diagnostics.AddError("Error reading TFA entry", err.Error())
		return
	}

	model.Type = types.StringValue(entry.Type)
	model.Description = types.StringValue(entry.Description)
	model.Enabled = types.BoolValue(entry.Enable)
	// dont touch totp, value, qrcode, url — they are write-only or creation-time only
}
