package sdn_vnet

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

var _ resource.Resource = &SDNVnetResource{}
var _ resource.ResourceWithConfigure = &SDNVnetResource{}
var _ resource.ResourceWithImportState = &SDNVnetResource{}

type SDNVnetResource struct {
	client *client.Client
}

type SDNVnetResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Vnet      types.String `tfsdk:"vnet"`
	Zone      types.String `tfsdk:"zone"`
	Alias     types.String `tfsdk:"alias"`
	Tag       types.Int64  `tfsdk:"tag"`
	VlanAware types.Bool   `tfsdk:"vlan_aware"`
	Comment   types.String `tfsdk:"comment"`
}

func NewResource() resource.Resource {
	return &SDNVnetResource{}
}

func (r *SDNVnetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_sdn_vnet"
}

func (r *SDNVnetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE SDN VNet.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vnet": schema.StringAttribute{
				Description: "The VNet name (identifier, max 8 characters).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zone": schema.StringAttribute{
				Description: "The SDN zone this VNet belongs to.",
				Required:    true,
			},
			"alias": schema.StringAttribute{
				Description: "An optional alias/description for the VNet.",
				Optional:    true,
				Computed:    true,
			},
			"tag": schema.Int64Attribute{
				Description: "VLAN tag (for VLAN-aware zones).",
				Optional:    true,
				Computed:    true,
			},
			"vlan_aware": schema.BoolAttribute{
				Description: "Whether the VNet is VLAN-aware.",
				Optional:    true,
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "A description for this VNet.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *SDNVnetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SDNVnetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SDNVnetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vlanAwareInt := boolToIntPtr(plan.VlanAware.ValueBool())

	tflog.Debug(ctx, "Creating SDN VNet", map[string]any{"vnet": plan.Vnet.ValueString(), "zone": plan.Zone.ValueString()})

	if err := r.client.CreateSDNVnet(ctx, &models.SDNVnetCreateRequest{
		Vnet:      plan.Vnet.ValueString(),
		Zone:      plan.Zone.ValueString(),
		Alias:     plan.Alias.ValueString(),
		Tag:       int(plan.Tag.ValueInt64()),
		VlanAware: vlanAwareInt,
		Comment:   plan.Comment.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error creating SDN VNet", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	plan.ID = plan.Vnet
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNVnetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SDNVnetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNVnetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SDNVnetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vlanAwareInt := boolToIntPtr(plan.VlanAware.ValueBool())

	if err := r.client.UpdateSDNVnet(ctx, plan.Vnet.ValueString(), &models.SDNVnetUpdateRequest{
		Zone:      plan.Zone.ValueString(),
		Alias:     plan.Alias.ValueString(),
		Tag:       int(plan.Tag.ValueInt64()),
		VlanAware: vlanAwareInt,
		Comment:   plan.Comment.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error updating SDN VNet", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNVnetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SDNVnetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSDNVnet(ctx, state.Vnet.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting SDN VNet", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)
}

func (r *SDNVnetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := SDNVnetResourceModel{
		ID:   types.StringValue(req.ID),
		Vnet: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNVnetResource) readIntoModel(ctx context.Context, model *SDNVnetResourceModel, diagnostics interface{ AddError(string, string) }) {
	vnet, err := r.client.GetSDNVnet(ctx, model.Vnet.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diagnostics.AddError("SDN VNet not found", "The SDN VNet no longer exists.")
			return
		}
		diagnostics.AddError("Error reading SDN VNet", err.Error())
		return
	}

	model.Zone = types.StringValue(vnet.Zone)
	model.Alias = types.StringValue(vnet.Alias)
	model.Tag = types.Int64Value(int64(vnet.Tag))
	model.Comment = types.StringValue(vnet.Comment)

	if vnet.VlanAware != nil {
		model.VlanAware = types.BoolValue(*vnet.VlanAware == 1)
	} else {
		model.VlanAware = types.BoolValue(false)
	}
}

func boolToIntPtr(b bool) *int {
	v := 0
	if b {
		v = 1
	}
	return &v
}
