package sdn_zone_vlan

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

var _ resource.Resource = &SDNZoneVLANResource{}
var _ resource.ResourceWithConfigure = &SDNZoneVLANResource{}
var _ resource.ResourceWithImportState = &SDNZoneVLANResource{}

type SDNZoneVLANResource struct {
	client *client.Client
}

type SDNZoneVLANModel struct {
	ID         types.String `tfsdk:"id"`
	Zone       types.String `tfsdk:"zone"`
	Bridge     types.String `tfsdk:"bridge"`
	Nodes      types.String `tfsdk:"nodes"`
	IPAM       types.String `tfsdk:"ipam"`
	DNS        types.String `tfsdk:"dns"`
	DNSZone    types.String `tfsdk:"dns_zone"`
	ReverseDNS types.String `tfsdk:"reverse_dns"`
}

func NewResource() resource.Resource {
	return &SDNZoneVLANResource{}
}

func (r *SDNZoneVLANResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_sdn_zone_vlan"
}

func (r *SDNZoneVLANResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE SDN zone of type 'vlan'.",
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
			"bridge": schema.StringAttribute{
				Description: "The Linux bridge interface to use for VLAN tagging.",
				Required:    true,
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

func (r *SDNZoneVLANResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SDNZoneVLANResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SDNZoneVLANModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating SDN VLAN zone", map[string]any{"zone": plan.Zone.ValueString()})

	if err := r.client.CreateSDNZone(ctx, &models.SDNZoneCreateRequest{
		Zone:       plan.Zone.ValueString(),
		Type:       "vlan",
		Bridge:     plan.Bridge.ValueString(),
		Nodes:      plan.Nodes.ValueString(),
		IPAM:       plan.IPAM.ValueString(),
		DNS:        plan.DNS.ValueString(),
		DNSZone:    plan.DNSZone.ValueString(),
		ReverseDNS: plan.ReverseDNS.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error creating SDN VLAN zone", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	plan.ID = plan.Zone
	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading SDN VLAN zone after create", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNZoneVLANResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SDNZoneVLANModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading SDN VLAN zone", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNZoneVLANResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SDNZoneVLANModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateSDNZone(ctx, plan.Zone.ValueString(), &models.SDNZoneUpdateRequest{
		Bridge:     plan.Bridge.ValueString(),
		Nodes:      plan.Nodes.ValueString(),
		IPAM:       plan.IPAM.ValueString(),
		DNS:        plan.DNS.ValueString(),
		DNSZone:    plan.DNSZone.ValueString(),
		ReverseDNS: plan.ReverseDNS.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error updating SDN VLAN zone", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading SDN VLAN zone after update", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNZoneVLANResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SDNZoneVLANModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSDNZone(ctx, state.Zone.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting SDN VLAN zone", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)
}

func (r *SDNZoneVLANResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := SDNZoneVLANModel{
		ID:   types.StringValue(req.ID),
		Zone: types.StringValue(req.ID),
	}
	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing SDN VLAN zone", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNZoneVLANResource) readIntoModel(ctx context.Context, model *SDNZoneVLANModel) error {
	zone, err := r.client.GetSDNZone(ctx, model.Zone.ValueString())
	if err != nil {
		return err
	}
	if zone.Type != "vlan" {
		return &client.APIError{StatusCode: 404, Status: "404 Not Found"}
	}

	model.Bridge = types.StringValue(zone.Bridge)
	model.Nodes = types.StringValue(zone.Nodes)
	model.IPAM = types.StringValue(zone.IPAM)
	model.DNS = types.StringValue(zone.DNS)
	model.DNSZone = types.StringValue(zone.DNSZone)
	model.ReverseDNS = types.StringValue(zone.ReverseDNS)
	return nil
}
