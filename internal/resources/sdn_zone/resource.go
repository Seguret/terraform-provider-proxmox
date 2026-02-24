package sdn_zone

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

var _ resource.Resource = &SDNZoneResource{}
var _ resource.ResourceWithConfigure = &SDNZoneResource{}
var _ resource.ResourceWithImportState = &SDNZoneResource{}

type SDNZoneResource struct {
	client *client.Client
}

type SDNZoneResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Zone       types.String `tfsdk:"zone"`
	Type       types.String `tfsdk:"type"`
	Comment    types.String `tfsdk:"comment"`
	Bridge     types.String `tfsdk:"bridge"`
	Tag        types.Int64  `tfsdk:"tag"`
	Peers      types.String `tfsdk:"peers"`
	VRFVxlan   types.Int64  `tfsdk:"vrf_vxlan"`
	Controller types.String `tfsdk:"controller"`
	MTU        types.Int64  `tfsdk:"mtu"`
	DNS        types.String `tfsdk:"dns"`
	DNSZone    types.String `tfsdk:"dns_zone"`
	ReverseDNS types.String `tfsdk:"reverse_dns"`
	IPAM       types.String `tfsdk:"ipam"`
}

func NewResource() resource.Resource {
	return &SDNZoneResource{}
}

func (r *SDNZoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_sdn_zone"
}

func (r *SDNZoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE SDN zone.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone": schema.StringAttribute{
				Description: "The SDN zone name (identifier).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "Zone type ('simple', 'vlan', 'qinq', 'vxlan', 'evpn').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Description: "A description for this zone.",
				Optional:    true,
				Computed:    true,
			},
			"bridge": schema.StringAttribute{
				Description: "The bridge interface (for vlan/qinq zones).",
				Optional:    true,
				Computed:    true,
			},
			"tag": schema.Int64Attribute{
				Description: "VLAN tag (for vlan/qinq zones).",
				Optional:    true,
				Computed:    true,
			},
			"peers": schema.StringAttribute{
				Description: "Comma-separated list of VXLAN peer addresses (for vxlan zones).",
				Optional:    true,
				Computed:    true,
			},
			"vrf_vxlan": schema.Int64Attribute{
				Description: "VRF VXLAN tag (for evpn zones).",
				Optional:    true,
				Computed:    true,
			},
			"controller": schema.StringAttribute{
				Description: "EVPN controller name (for evpn zones).",
				Optional:    true,
				Computed:    true,
			},
			"mtu": schema.Int64Attribute{
				Description: "MTU value for the zone.",
				Optional:    true,
				Computed:    true,
			},
			"dns": schema.StringAttribute{
				Description: "DNS plugin name.",
				Optional:    true,
				Computed:    true,
			},
			"dns_zone": schema.StringAttribute{
				Description: "DNS domain.",
				Optional:    true,
				Computed:    true,
			},
			"reverse_dns": schema.StringAttribute{
				Description: "Reverse DNS plugin name.",
				Optional:    true,
				Computed:    true,
			},
			"ipam": schema.StringAttribute{
				Description: "IPAM plugin name.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *SDNZoneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SDNZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SDNZoneResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating SDN zone", map[string]any{"zone": plan.Zone.ValueString(), "type": plan.Type.ValueString()})

	if err := r.client.CreateSDNZone(ctx, &models.SDNZoneCreateRequest{
		Zone:       plan.Zone.ValueString(),
		Type:       plan.Type.ValueString(),
		Comment:    plan.Comment.ValueString(),
		Bridge:     plan.Bridge.ValueString(),
		Tag:        int(plan.Tag.ValueInt64()),
		Peers:      plan.Peers.ValueString(),
		VRFVxlan:   int(plan.VRFVxlan.ValueInt64()),
		Controller: plan.Controller.ValueString(),
		MTU:        int(plan.MTU.ValueInt64()),
		DNS:        plan.DNS.ValueString(),
		DNSZone:    plan.DNSZone.ValueString(),
		ReverseDNS: plan.ReverseDNS.ValueString(),
		IPAM:       plan.IPAM.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error creating SDN zone", err.Error())
		return
	}

	// push pending SDN changes live
	_ = r.client.ApplySDN(ctx)

	plan.ID = plan.Zone
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SDNZoneResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SDNZoneResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateSDNZone(ctx, plan.Zone.ValueString(), &models.SDNZoneUpdateRequest{
		Comment:    plan.Comment.ValueString(),
		Bridge:     plan.Bridge.ValueString(),
		Tag:        int(plan.Tag.ValueInt64()),
		Peers:      plan.Peers.ValueString(),
		VRFVxlan:   int(plan.VRFVxlan.ValueInt64()),
		Controller: plan.Controller.ValueString(),
		MTU:        int(plan.MTU.ValueInt64()),
		DNS:        plan.DNS.ValueString(),
		DNSZone:    plan.DNSZone.ValueString(),
		ReverseDNS: plan.ReverseDNS.ValueString(),
		IPAM:       plan.IPAM.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error updating SDN zone", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SDNZoneResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSDNZone(ctx, state.Zone.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting SDN zone", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)
}

func (r *SDNZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := SDNZoneResourceModel{
		ID:   types.StringValue(req.ID),
		Zone: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNZoneResource) readIntoModel(ctx context.Context, model *SDNZoneResourceModel, diagnostics interface{ AddError(string, string) }) {
	zone, err := r.client.GetSDNZone(ctx, model.Zone.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diagnostics.AddError("SDN zone not found", "The SDN zone no longer exists.")
			return
		}
		diagnostics.AddError("Error reading SDN zone", err.Error())
		return
	}

	model.Type = types.StringValue(zone.Type)
	model.Comment = types.StringValue(zone.Comment)
	model.Bridge = types.StringValue(zone.Bridge)
	model.Tag = types.Int64Value(int64(zone.Tag))
	model.Peers = types.StringValue(zone.Peers)
	model.VRFVxlan = types.Int64Value(int64(zone.VRFVxlan))
	model.Controller = types.StringValue(zone.Controller)
	model.MTU = types.Int64Value(int64(zone.MTU))
	model.DNS = types.StringValue(zone.DNS)
	model.DNSZone = types.StringValue(zone.DNSZone)
	model.ReverseDNS = types.StringValue(zone.ReverseDNS)
	model.IPAM = types.StringValue(zone.IPAM)
}
