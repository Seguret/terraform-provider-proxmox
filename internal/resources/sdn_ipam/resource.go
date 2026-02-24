package sdn_ipam

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

var _ resource.Resource = &SDNIpamResource{}
var _ resource.ResourceWithConfigure = &SDNIpamResource{}
var _ resource.ResourceWithImportState = &SDNIpamResource{}

type SDNIpamResource struct {
	client *client.Client
}

type SDNIpamResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Ipam    types.String `tfsdk:"ipam"`
	Type    types.String `tfsdk:"type"`
	URL     types.String `tfsdk:"url"`
	Token   types.String `tfsdk:"token"`
	Section types.Int64  `tfsdk:"section"`
}

func NewResource() resource.Resource {
	return &SDNIpamResource{}
}

func (r *SDNIpamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_sdn_ipam"
}

func (r *SDNIpamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE SDN IPAM provider.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ipam": schema.StringAttribute{
				Description: "The SDN IPAM provider name (identifier).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The IPAM plugin type (netbox, phpipam, or pve).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Description: "The IPAM server API URL.",
				Optional:    true,
				Computed:    true,
			},
			"token": schema.StringAttribute{
				Description: "The IPAM server API token. Sensitive.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
			},
			"section": schema.Int64Attribute{
				Description: "The phpIPAM section ID.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *SDNIpamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SDNIpamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SDNIpamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating SDN IPAM provider", map[string]any{
		"ipam": plan.Ipam.ValueString(),
		"type": plan.Type.ValueString(),
	})

	if err := r.client.CreateSDNIpam(ctx, &models.SDNIpamCreateRequest{
		Ipam:    plan.Ipam.ValueString(),
		Type:    plan.Type.ValueString(),
		URL:     plan.URL.ValueString(),
		Token:   plan.Token.ValueString(),
		Section: int(plan.Section.ValueInt64()),
	}); err != nil {
		resp.Diagnostics.AddError("Error creating SDN IPAM provider", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	plan.ID = plan.Ipam
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNIpamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SDNIpamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if isNotFound := r.readIntoModelWithNotFound(ctx, &state, &resp.Diagnostics); isNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNIpamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SDNIpamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateSDNIpam(ctx, plan.Ipam.ValueString(), &models.SDNIpamUpdateRequest{
		URL:     plan.URL.ValueString(),
		Token:   plan.Token.ValueString(),
		Section: int(plan.Section.ValueInt64()),
	}); err != nil {
		resp.Diagnostics.AddError("Error updating SDN IPAM provider", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNIpamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SDNIpamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSDNIpam(ctx, state.Ipam.ValueString()); err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return
		}
		resp.Diagnostics.AddError("Error deleting SDN IPAM provider", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)
}

func (r *SDNIpamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := SDNIpamResourceModel{
		ID:   types.StringValue(req.ID),
		Ipam: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// readIntoModelWithNotFound fetches IPAM state; returns true if it doesnt exist.
func (r *SDNIpamResource) readIntoModelWithNotFound(ctx context.Context, model *SDNIpamResourceModel, diags interface{ AddError(string, string) }) bool {
	ipam, err := r.client.GetSDNIpam(ctx, model.Ipam.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return true
		}
		diags.AddError("Error reading SDN IPAM provider", err.Error())
		return false
	}
	populateModel(model, ipam)
	return false
}

// readIntoModel fetches IPAM state and adds diagnostic errors on failure.
func (r *SDNIpamResource) readIntoModel(ctx context.Context, model *SDNIpamResourceModel, diags interface{ AddError(string, string) }) {
	ipam, err := r.client.GetSDNIpam(ctx, model.Ipam.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diags.AddError("SDN IPAM provider not found", "The SDN IPAM provider no longer exists.")
			return
		}
		diags.AddError("Error reading SDN IPAM provider", err.Error())
		return
	}
	populateModel(model, ipam)
}

func populateModel(model *SDNIpamResourceModel, ipam *models.SDNIpam) {
	model.Type = types.StringValue(ipam.Type)
	model.URL = types.StringValue(ipam.URL)
	// token is sensitive — only overwrite if the API actually returned something
	if ipam.Token != "" {
		model.Token = types.StringValue(ipam.Token)
	}
	model.Section = types.Int64Value(int64(ipam.Section))
}
