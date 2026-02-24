package acme_account

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &ACMEAccountResource{}
var _ resource.ResourceWithConfigure = &ACMEAccountResource{}
var _ resource.ResourceWithImportState = &ACMEAccountResource{}

type ACMEAccountResource struct {
	client *client.Client
}

type ACMEAccountResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Contact    types.String `tfsdk:"contact"`
	Directory  types.String `tfsdk:"directory"`
	TosURL     types.String `tfsdk:"tos_url"`
	AccountURL types.String `tfsdk:"account_url"`
}

func NewResource() resource.Resource {
	return &ACMEAccountResource{}
}

func (r *ACMEAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_acme_account"
}

func (r *ACMEAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE ACME account registration (e.g. Let's Encrypt).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Account name (default: 'default').",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"contact": schema.StringAttribute{
				Description: "Contact email address for the ACME account.",
				Required:    true,
			},
			"directory": schema.StringAttribute{
				Description: "ACME directory URL. Defaults to Let's Encrypt production.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tos_url": schema.StringAttribute{
				Description: "URL of the ACME Terms of Service (must be accepted on account creation).",
				Optional:    true,
				Computed:    true,
			},
			"account_url": schema.StringAttribute{
				Description: "The registered ACME account URL (assigned by the CA).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ACMEAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ACMEAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ACMEAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	if name == "" {
		name = "default"
	}

	tflog.Debug(ctx, "Creating ACME account", map[string]any{"name": name})

	upid, err := r.client.CreateACMEAccount(ctx, &models.ACMEAccountCreateRequest{
		Name:      name,
		Contact:   plan.Contact.ValueString(),
		Directory: plan.Directory.ValueString(),
		TosURL:    plan.TosURL.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating ACME account", err.Error())
		return
	}

	if upid != "" {
		if err := r.client.WaitForUPID(ctx, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for ACME account registration", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(name)
	plan.Name = types.StringValue(name)
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ACMEAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ACMEAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ACMEAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ACMEAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upid, err := r.client.UpdateACMEAccount(ctx, plan.Name.ValueString(), &models.ACMEAccountUpdateRequest{
		Contact: plan.Contact.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating ACME account", err.Error())
		return
	}

	if upid != "" {
		if err := r.client.WaitForUPID(ctx, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for ACME account update", err.Error())
			return
		}
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ACMEAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ACMEAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upid, err := r.client.DeleteACMEAccount(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting ACME account", err.Error())
		return
	}

	if upid != "" {
		_ = r.client.WaitForUPID(ctx, upid)
	}
}

func (r *ACMEAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := ACMEAccountResourceModel{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ACMEAccountResource) readIntoModel(ctx context.Context, model *ACMEAccountResourceModel, diagnostics interface{ AddError(string, string) }) {
	account, err := r.client.GetACMEAccount(ctx, model.Name.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diagnostics.AddError("ACME account not found", "The ACME account no longer exists.")
			return
		}
		diagnostics.AddError("Error reading ACME account", err.Error())
		return
	}

	model.Contact = types.StringValue(account.Contact)
	model.Directory = types.StringValue(account.Directory)
	model.TosURL = types.StringValue(account.TosURL)
	model.AccountURL = types.StringValue(account.AccountURL)
}
