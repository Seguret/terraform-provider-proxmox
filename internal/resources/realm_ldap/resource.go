package realm_ldap

import (
	"context"
	"fmt"

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

var _ resource.Resource = &RealmLDAPResource{}
var _ resource.ResourceWithConfigure = &RealmLDAPResource{}
var _ resource.ResourceWithImportState = &RealmLDAPResource{}

type RealmLDAPResource struct {
	client *client.Client
}

type RealmLDAPResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Realm   types.String `tfsdk:"realm"`
	Server1 types.String `tfsdk:"server1"`
	Port    types.Int64  `tfsdk:"port"`
	BaseDN  types.String `tfsdk:"base_dn"`
	UserAttr types.String `tfsdk:"user_attr"`
	Domain  types.String `tfsdk:"domain"`
	Secure  types.Bool   `tfsdk:"secure"`
	Comment types.String `tfsdk:"comment"`
	Default types.Bool   `tfsdk:"default"`
}

func NewResource() resource.Resource {
	return &RealmLDAPResource{}
}

func (r *RealmLDAPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_realm_ldap"
}

func (r *RealmLDAPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE LDAP authentication realm.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"realm": schema.StringAttribute{
				Description:   "The realm identifier (e.g., 'my-ldap').",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"server1": schema.StringAttribute{
				Description: "The primary LDAP server hostname or IP address.",
				Required:    true,
			},
			"port": schema.Int64Attribute{
				Description: "The LDAP server port.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(389),
			},
			"base_dn": schema.StringAttribute{
				Description: "The LDAP base distinguished name for user searches.",
				Required:    true,
			},
			"user_attr": schema.StringAttribute{
				Description: "The LDAP attribute used to identify users (e.g., 'uid' or 'sAMAccountName').",
				Optional:    true,
				Computed:    true,
			},
			"domain": schema.StringAttribute{
				Description: "The optional AD domain name.",
				Optional:    true,
				Computed:    true,
			},
			"secure": schema.BoolAttribute{
				Description: "Whether to use LDAPS (TLS/SSL).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"comment": schema.StringAttribute{
				Description: "A comment for the realm.",
				Optional:    true,
				Computed:    true,
			},
			"default": schema.BoolAttribute{
				Description: "Whether this realm is the default login realm.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *RealmLDAPResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RealmLDAPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RealmLDAPResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating LDAP realm", map[string]any{"realm": plan.Realm.ValueString()})

	createReq := r.buildCreateRequest(&plan)
	if err := r.client.CreateRealm(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating LDAP realm", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Realm.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RealmLDAPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RealmLDAPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RealmLDAPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RealmLDAPResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating LDAP realm", map[string]any{"realm": plan.Realm.ValueString()})

	updateReq := r.buildUpdateRequest(&plan)
	if err := r.client.UpdateRealm(ctx, plan.Realm.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating LDAP realm", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Realm.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RealmLDAPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RealmLDAPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting LDAP realm", map[string]any{"realm": state.Realm.ValueString()})

	if err := r.client.DeleteRealm(ctx, state.Realm.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting LDAP realm", err.Error())
	}
}

func (r *RealmLDAPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := RealmLDAPResourceModel{
		ID:    types.StringValue(req.ID),
		Realm: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RealmLDAPResource) readIntoModel(ctx context.Context, model *RealmLDAPResourceModel, diagnostics *diag.Diagnostics) {
	realm, err := r.client.GetRealm(ctx, model.Realm.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diagnostics.AddWarning("LDAP realm not found", "The realm no longer exists.")
			return
		}
		diagnostics.AddError("Error reading LDAP realm", err.Error())
		return
	}

	if realm.Type != "ldap" {
		diagnostics.AddError("Unexpected realm type",
			fmt.Sprintf("Expected type 'ldap', got '%s' for realm '%s'.", realm.Type, model.Realm.ValueString()))
		return
	}

	model.Server1 = types.StringValue(realm.Server1)
	model.BaseDN = types.StringValue(realm.BaseDN)
	model.UserAttr = types.StringValue(realm.UserAttr)
	model.Comment = types.StringValue(realm.Comment)

	if realm.Port != nil {
		model.Port = types.Int64Value(int64(*realm.Port))
	} else {
		model.Port = types.Int64Value(389)
	}

	if realm.Secure != nil {
		model.Secure = types.BoolValue(*realm.Secure == 1)
	} else {
		model.Secure = types.BoolValue(false)
	}

	if realm.Default != nil {
		model.Default = types.BoolValue(*realm.Default == 1)
	} else {
		model.Default = types.BoolValue(false)
	}

	// domain isnt in the proxmox ldap model, so just leave the existing value alone
}

func (r *RealmLDAPResource) buildCreateRequest(plan *RealmLDAPResourceModel) *models.AuthRealmCreateRequest {
	req := &models.AuthRealmCreateRequest{
		Realm:    plan.Realm.ValueString(),
		Type:     "ldap",
		Server1:  plan.Server1.ValueString(),
		BaseDN:   plan.BaseDN.ValueString(),
		UserAttr: plan.UserAttr.ValueString(),
		Comment:  plan.Comment.ValueString(),
	}

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() && plan.Port.ValueInt64() > 0 {
		v := int(plan.Port.ValueInt64())
		req.Port = &v
	}

	if !plan.Secure.IsNull() && !plan.Secure.IsUnknown() {
		v := boolToInt(plan.Secure.ValueBool())
		req.Secure = &v
	}

	if !plan.Default.IsNull() && !plan.Default.IsUnknown() {
		v := boolToInt(plan.Default.ValueBool())
		req.Default = &v
	}

	return req
}

func (r *RealmLDAPResource) buildUpdateRequest(plan *RealmLDAPResourceModel) *models.AuthRealmUpdateRequest {
	req := &models.AuthRealmUpdateRequest{
		Server1:  plan.Server1.ValueString(),
		BaseDN:   plan.BaseDN.ValueString(),
		UserAttr: plan.UserAttr.ValueString(),
		Comment:  plan.Comment.ValueString(),
	}

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() && plan.Port.ValueInt64() > 0 {
		v := int(plan.Port.ValueInt64())
		req.Port = &v
	}

	if !plan.Secure.IsNull() && !plan.Secure.IsUnknown() {
		v := boolToInt(plan.Secure.ValueBool())
		req.Secure = &v
	}

	if !plan.Default.IsNull() && !plan.Default.IsUnknown() {
		v := boolToInt(plan.Default.ValueBool())
		req.Default = &v
	}

	return req
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
