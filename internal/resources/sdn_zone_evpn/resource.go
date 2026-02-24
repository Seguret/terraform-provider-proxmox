package sdn_zone_evpn

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

var _ resource.Resource = &SDNZoneEVPNResource{}
var _ resource.ResourceWithConfigure = &SDNZoneEVPNResource{}
var _ resource.ResourceWithImportState = &SDNZoneEVPNResource{}

type SDNZoneEVPNResource struct {
	client *client.Client
}

type SDNZoneEVPNModel struct {
	ID                    types.String `tfsdk:"id"`
	Zone                  types.String `tfsdk:"zone"`
	Controller            types.String `tfsdk:"controller"`
	VRFVxlan              types.Int64  `tfsdk:"vrf_vxlan"`
	ExitNodes             types.String `tfsdk:"exit_nodes"`
	ExitNodesLocalRouting types.Bool   `tfsdk:"exit_nodes_local_routing"`
	AdvertiseSubnets      types.Bool   `tfsdk:"advertise_subnets"`
	Nodes                 types.String `tfsdk:"nodes"`
	IPAM                  types.String `tfsdk:"ipam"`
	DNS                   types.String `tfsdk:"dns"`
	DNSZone               types.String `tfsdk:"dns_zone"`
	ReverseDNS            types.String `tfsdk:"reverse_dns"`
}

func NewResource() resource.Resource {
	return &SDNZoneEVPNResource{}
}

func (r *SDNZoneEVPNResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_sdn_zone_evpn"
}

func (r *SDNZoneEVPNResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE SDN zone of type 'evpn'.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone": schema.StringAttribute{
				Description: "The SDN zone identifier.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"controller": schema.StringAttribute{
				Description: "EVPN controller name.",
				Required:    true,
			},
			"vrf_vxlan": schema.Int64Attribute{
				Description: "VRF VxLAN tag number.",
				Required:    true,
			},
			"exit_nodes": schema.StringAttribute{
				Description: "Comma-separated list of nodes acting as exit nodes for the zone.",
				Optional:    true,
			},
			"exit_nodes_local_routing": schema.BoolAttribute{
				Description: "Enable local routing on exit nodes.",
				Optional:    true,
			},
			"advertise_subnets": schema.BoolAttribute{
				Description: "Advertise subnets via EVPN.",
				Optional:    true,
			},
			"nodes": schema.StringAttribute{
				Description: "Comma-separated list of nodes where the zone is deployed.",
				Optional:    true,
			},
			"ipam": schema.StringAttribute{
				Description: "IPAM plugin name.",
				Optional:    true,
			},
			"dns": schema.StringAttribute{
				Description: "DNS plugin name.",
				Optional:    true,
			},
			"dns_zone": schema.StringAttribute{
				Description: "DNS domain.",
				Optional:    true,
			},
			"reverse_dns": schema.StringAttribute{
				Description: "Reverse DNS plugin name.",
				Optional:    true,
			},
		},
	}
}

func (r *SDNZoneEVPNResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (r *SDNZoneEVPNResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SDNZoneEVPNModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating SDN EVPN zone", map[string]any{"zone": plan.Zone.ValueString()})

	if err := r.client.CreateSDNZone(ctx, &models.SDNZoneCreateRequest{
		Zone:                  plan.Zone.ValueString(),
		Type:                  "evpn",
		Controller:            plan.Controller.ValueString(),
		VRFVxlan:              int(plan.VRFVxlan.ValueInt64()),
		ExitNodes:             plan.ExitNodes.ValueString(),
		ExitNodesLocalRouting: boolToInt(plan.ExitNodesLocalRouting.ValueBool()),
		AdvertiseSubnets:      boolToInt(plan.AdvertiseSubnets.ValueBool()),
		Nodes:                 plan.Nodes.ValueString(),
		IPAM:                  plan.IPAM.ValueString(),
		DNS:                   plan.DNS.ValueString(),
		DNSZone:               plan.DNSZone.ValueString(),
		ReverseDNS:            plan.ReverseDNS.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error creating SDN EVPN zone", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	plan.ID = plan.Zone
	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading SDN EVPN zone after create", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNZoneEVPNResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SDNZoneEVPNModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading SDN EVPN zone", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNZoneEVPNResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SDNZoneEVPNModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateSDNZone(ctx, plan.Zone.ValueString(), &models.SDNZoneUpdateRequest{
		Controller:            plan.Controller.ValueString(),
		VRFVxlan:              int(plan.VRFVxlan.ValueInt64()),
		ExitNodes:             plan.ExitNodes.ValueString(),
		ExitNodesLocalRouting: boolToInt(plan.ExitNodesLocalRouting.ValueBool()),
		AdvertiseSubnets:      boolToInt(plan.AdvertiseSubnets.ValueBool()),
		Nodes:                 plan.Nodes.ValueString(),
		IPAM:                  plan.IPAM.ValueString(),
		DNS:                   plan.DNS.ValueString(),
		DNSZone:               plan.DNSZone.ValueString(),
		ReverseDNS:            plan.ReverseDNS.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error updating SDN EVPN zone", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading SDN EVPN zone after update", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNZoneEVPNResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SDNZoneEVPNModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSDNZone(ctx, state.Zone.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting SDN EVPN zone", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)
}

func (r *SDNZoneEVPNResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := SDNZoneEVPNModel{
		ID:   types.StringValue(req.ID),
		Zone: types.StringValue(req.ID),
	}
	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing SDN EVPN zone", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNZoneEVPNResource) readIntoModel(ctx context.Context, model *SDNZoneEVPNModel) error {
	zone, err := r.client.GetSDNZone(ctx, model.Zone.ValueString())
	if err != nil {
		return err
	}
	if zone.Type != "evpn" {
		return &client.APIError{StatusCode: 404, Status: "404 Not Found"}
	}

	model.Controller = types.StringValue(zone.Controller)
	model.VRFVxlan = types.Int64Value(int64(zone.VRFVxlan))
	model.ExitNodes = types.StringValue(zone.ExitNodes)
	model.ExitNodesLocalRouting = types.BoolValue(zone.ExitNodesLocalRouting != 0)
	model.AdvertiseSubnets = types.BoolValue(zone.AdvertiseSubnets != 0)
	model.Nodes = types.StringValue(zone.Nodes)
	model.IPAM = types.StringValue(zone.IPAM)
	model.DNS = types.StringValue(zone.DNS)
	model.DNSZone = types.StringValue(zone.DNSZone)
	model.ReverseDNS = types.StringValue(zone.ReverseDNS)
	return nil
}
