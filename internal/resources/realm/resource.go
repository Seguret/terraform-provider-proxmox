package realm

import (
	"context"
	"fmt"

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

var _ resource.Resource = &RealmResource{}
var _ resource.ResourceWithConfigure = &RealmResource{}
var _ resource.ResourceWithImportState = &RealmResource{}

type RealmResource struct {
	client *client.Client
}

type RealmResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Realm         types.String `tfsdk:"realm"`
	Type          types.String `tfsdk:"type"`
	Comment       types.String `tfsdk:"comment"`
	Default       types.Bool   `tfsdk:"default"`
	Server1       types.String `tfsdk:"server1"`
	Server2       types.String `tfsdk:"server2"`
	Port          types.Int64  `tfsdk:"port"`
	BaseDN        types.String `tfsdk:"base_dn"`
	BindDN        types.String `tfsdk:"bind_dn"`
	Password      types.String `tfsdk:"password"`
	UserAttr      types.String `tfsdk:"user_attr"`
	Secure        types.Bool   `tfsdk:"secure"`
	Verify        types.Bool   `tfsdk:"verify"`
	IssuerURL     types.String `tfsdk:"issuer_url"`
	ClientID      types.String `tfsdk:"client_id"`
	ClientKey     types.String `tfsdk:"client_key"`
	UsernameClaim types.String `tfsdk:"username_claim"`
	AutoCreate    types.Bool   `tfsdk:"auto_create"`
	TFA           types.String `tfsdk:"tfa"`
}

func NewResource() resource.Resource {
	return &RealmResource{}
}

func (r *RealmResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_realm"
}

func (r *RealmResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE authentication realm.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"realm": schema.StringAttribute{
				Description: "The realm identifier (e.g., 'my-ldap').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The realm type (pam, pve, ad, ldap, openid).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Description: "A comment for the realm.",
				Optional:    true,
				Computed:    true,
			},
			"default": schema.BoolAttribute{
				Description: "Whether this is the default realm.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"server1": schema.StringAttribute{
				Description: "Primary server address (for ldap/ad).",
				Optional:    true,
				Computed:    true,
			},
			"server2": schema.StringAttribute{
				Description: "Secondary server address (for ldap/ad).",
				Optional:    true,
				Computed:    true,
			},
			"port": schema.Int64Attribute{
				Description: "Server port (for ldap/ad).",
				Optional:    true,
				Computed:    true,
			},
			"base_dn": schema.StringAttribute{
				Description: "LDAP base distinguished name.",
				Optional:    true,
				Computed:    true,
			},
			"bind_dn": schema.StringAttribute{
				Description: "LDAP bind distinguished name.",
				Optional:    true,
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "LDAP bind password.",
				Optional:    true,
				Sensitive:   true,
			},
			"user_attr": schema.StringAttribute{
				Description: "LDAP user attribute name.",
				Optional:    true,
				Computed:    true,
			},
			"secure": schema.BoolAttribute{
				Description: "Use LDAPS/TLS.",
				Optional:    true,
				Computed:    true,
			},
			"verify": schema.BoolAttribute{
				Description: "Verify server TLS certificate.",
				Optional:    true,
				Computed:    true,
			},
			"issuer_url": schema.StringAttribute{
				Description: "OpenID Connect issuer URL.",
				Optional:    true,
				Computed:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "OpenID Connect client ID.",
				Optional:    true,
				Computed:    true,
			},
			"client_key": schema.StringAttribute{
				Description: "OpenID Connect client secret.",
				Optional:    true,
				Sensitive:   true,
			},
			"username_claim": schema.StringAttribute{
				Description: "OpenID Connect claim used as the username.",
				Optional:    true,
				Computed:    true,
			},
			"auto_create": schema.BoolAttribute{
				Description: "Automatically create users on first login.",
				Optional:    true,
				Computed:    true,
			},
			"tfa": schema.StringAttribute{
				Description: "Two-factor authentication provider.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *RealmResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RealmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RealmResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := r.buildCreateRequest(&plan)

	tflog.Debug(ctx, "Creating realm", map[string]any{"realm": createReq.Realm, "type": createReq.Type})

	if err := r.client.CreateRealm(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating realm", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Realm.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RealmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RealmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RealmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RealmResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := r.buildUpdateRequest(&plan)

	if err := r.client.UpdateRealm(ctx, plan.Realm.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating realm", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Realm.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RealmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RealmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteRealm(ctx, state.Realm.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting realm", err.Error())
	}
}

func (r *RealmResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := RealmResourceModel{
		ID:    types.StringValue(req.ID),
		Realm: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RealmResource) readIntoModel(ctx context.Context, model *RealmResourceModel, diagnostics *diag.Diagnostics) {
	realm, err := r.client.GetRealm(ctx, model.Realm.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok {
			if apiErr.IsNotFound() {
				diagnostics.AddWarning("Realm not found", "The realm no longer exists.")
				return
			}
		}
		diagnostics.AddError("Error reading realm", err.Error())
		return
	}

	model.Type = types.StringValue(realm.Type)
	model.Comment = types.StringValue(realm.Comment)
	if realm.Default != nil {
		model.Default = types.BoolValue(*realm.Default == 1)
	} else {
		model.Default = types.BoolValue(false)
	}
	model.Server1 = types.StringValue(realm.Server1)
	model.Server2 = types.StringValue(realm.Server2)
	if realm.Port != nil {
		model.Port = types.Int64Value(int64(*realm.Port))
	} else {
		model.Port = types.Int64Value(0)
	}
	model.BaseDN = types.StringValue(realm.BaseDN)
	model.BindDN = types.StringValue(realm.BindDN)
	model.UserAttr = types.StringValue(realm.UserAttr)
	if realm.Secure != nil {
		model.Secure = types.BoolValue(*realm.Secure == 1)
	} else {
		model.Secure = types.BoolValue(false)
	}
	if realm.Verify != nil {
		model.Verify = types.BoolValue(*realm.Verify == 1)
	} else {
		model.Verify = types.BoolValue(false)
	}
	model.IssuerURL = types.StringValue(realm.IssuerURL)
	model.ClientID = types.StringValue(realm.ClientID)
	model.UsernameClaim = types.StringValue(realm.Username)
	if realm.AutoCreate != nil {
		model.AutoCreate = types.BoolValue(*realm.AutoCreate == 1)
	} else {
		model.AutoCreate = types.BoolValue(false)
	}
	model.TFA = types.StringValue(realm.TFAType)
	// password and client_key are write-only — dont overwrite them on read
}

func (r *RealmResource) buildCreateRequest(plan *RealmResourceModel) *models.AuthRealmCreateRequest {
	req := &models.AuthRealmCreateRequest{
		Realm:         plan.Realm.ValueString(),
		Type:          plan.Type.ValueString(),
		Comment:       plan.Comment.ValueString(),
		Server1:       plan.Server1.ValueString(),
		Server2:       plan.Server2.ValueString(),
		BaseDN:        plan.BaseDN.ValueString(),
		BindDN:        plan.BindDN.ValueString(),
		Password:      plan.Password.ValueString(),
		UserAttr:      plan.UserAttr.ValueString(),
		IssuerURL:     plan.IssuerURL.ValueString(),
		ClientID:      plan.ClientID.ValueString(),
		ClientKey:     plan.ClientKey.ValueString(),
		UsernameClaim: plan.UsernameClaim.ValueString(),
		TFAType:       plan.TFA.ValueString(),
	}
	if !plan.Default.IsNull() && !plan.Default.IsUnknown() {
		v := boolToInt(plan.Default.ValueBool())
		req.Default = &v
	}
	if !plan.Port.IsNull() && !plan.Port.IsUnknown() && plan.Port.ValueInt64() != 0 {
		v := int(plan.Port.ValueInt64())
		req.Port = &v
	}
	if !plan.Secure.IsNull() && !plan.Secure.IsUnknown() {
		v := boolToInt(plan.Secure.ValueBool())
		req.Secure = &v
	}
	if !plan.Verify.IsNull() && !plan.Verify.IsUnknown() {
		v := boolToInt(plan.Verify.ValueBool())
		req.Verify = &v
	}
	if !plan.AutoCreate.IsNull() && !plan.AutoCreate.IsUnknown() {
		v := boolToInt(plan.AutoCreate.ValueBool())
		req.AutoCreate = &v
	}
	return req
}

func (r *RealmResource) buildUpdateRequest(plan *RealmResourceModel) *models.AuthRealmUpdateRequest {
	req := &models.AuthRealmUpdateRequest{
		Comment:       plan.Comment.ValueString(),
		Server1:       plan.Server1.ValueString(),
		Server2:       plan.Server2.ValueString(),
		BaseDN:        plan.BaseDN.ValueString(),
		BindDN:        plan.BindDN.ValueString(),
		Password:      plan.Password.ValueString(),
		UserAttr:      plan.UserAttr.ValueString(),
		IssuerURL:     plan.IssuerURL.ValueString(),
		ClientID:      plan.ClientID.ValueString(),
		ClientKey:     plan.ClientKey.ValueString(),
		UsernameClaim: plan.UsernameClaim.ValueString(),
		TFAType:       plan.TFA.ValueString(),
	}
	if !plan.Default.IsNull() && !plan.Default.IsUnknown() {
		v := boolToInt(plan.Default.ValueBool())
		req.Default = &v
	}
	if !plan.Port.IsNull() && !plan.Port.IsUnknown() && plan.Port.ValueInt64() != 0 {
		v := int(plan.Port.ValueInt64())
		req.Port = &v
	}
	if !plan.Secure.IsNull() && !plan.Secure.IsUnknown() {
		v := boolToInt(plan.Secure.ValueBool())
		req.Secure = &v
	}
	if !plan.Verify.IsNull() && !plan.Verify.IsUnknown() {
		v := boolToInt(plan.Verify.ValueBool())
		req.Verify = &v
	}
	if !plan.AutoCreate.IsNull() && !plan.AutoCreate.IsUnknown() {
		v := boolToInt(plan.AutoCreate.ValueBool())
		req.AutoCreate = &v
	}
	return req
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
