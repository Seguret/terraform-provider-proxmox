package sdn_zone_qinq

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &SDNZoneQinQResource{}
var _ resource.ResourceWithConfigure = &SDNZoneQinQResource{}
var _ resource.ResourceWithImportState = &SDNZoneQinQResource{}

type SDNZoneQinQResource struct {
	client *client.Client
}

type SDNZoneQinQModel struct {
	ID           types.String `tfsdk:"id"`
	Zone         types.String `tfsdk:"zone"`
	Bridge       types.String `tfsdk:"bridge"`
	Tag          types.Int64  `tfsdk:"tag"`
	VlanProtocol types.String `tfsdk:"vlan_protocol"`
	Nodes        types.String `tfsdk:"nodes"`
	IPAM         types.String `tfsdk:"ipam"`
	DNS          types.String `tfsdk:"dns"`
	DNSZone      types.String `tfsdk:"dns_zone"`
	ReverseDNS   types.String `tfsdk:"reverse_dns"`
}

func NewResource() resource.Resource {
	return &SDNZoneQinQResource{}
}

func (r *SDNZoneQinQResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_sdn_zone_qinq"
}

func (r *SDNZoneQinQResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE SDN zone of type 'qinq'.",
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
				Description: "The Linux bridge interface to use for QinQ double-tagging.",
				Required:    true,
			},
			"tag": schema.Int64Attribute{
				Description: "Outer VLAN tag (1–4094) for the QinQ zone.",
				Required:    true,
			},
			"vlan_protocol": schema.StringAttribute{
				Description: "VLAN protocol to use. Valid values: '802.1q' (default), '802.1ad'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("802.1q"),
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

func (r *SDNZoneQinQResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SDNZoneQinQResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SDNZoneQinQModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating SDN QinQ zone", map[string]any{"zone": plan.Zone.ValueString()})

	if err := r.client.CreateSDNZone(ctx, &models.SDNZoneCreateRequest{
		Zone:         plan.Zone.ValueString(),
		Type:         "qinq",
		Bridge:       plan.Bridge.ValueString(),
		Tag:          int(plan.Tag.ValueInt64()),
		VlanProtocol: plan.VlanProtocol.ValueString(),
		Nodes:        plan.Nodes.ValueString(),
		IPAM:         plan.IPAM.ValueString(),
		DNS:          plan.DNS.ValueString(),
		DNSZone:      plan.DNSZone.ValueString(),
		ReverseDNS:   plan.ReverseDNS.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error creating SDN QinQ zone", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	plan.ID = plan.Zone
	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading SDN QinQ zone after create", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNZoneQinQResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SDNZoneQinQModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading SDN QinQ zone", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNZoneQinQResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SDNZoneQinQModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateSDNZone(ctx, plan.Zone.ValueString(), &models.SDNZoneUpdateRequest{
		Bridge:       plan.Bridge.ValueString(),
		Tag:          int(plan.Tag.ValueInt64()),
		VlanProtocol: plan.VlanProtocol.ValueString(),
		Nodes:        plan.Nodes.ValueString(),
		IPAM:         plan.IPAM.ValueString(),
		DNS:          plan.DNS.ValueString(),
		DNSZone:      plan.DNSZone.ValueString(),
		ReverseDNS:   plan.ReverseDNS.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error updating SDN QinQ zone", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading SDN QinQ zone after update", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNZoneQinQResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SDNZoneQinQModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSDNZone(ctx, state.Zone.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting SDN QinQ zone", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)
}

func (r *SDNZoneQinQResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := SDNZoneQinQModel{
		ID:   types.StringValue(req.ID),
		Zone: types.StringValue(req.ID),
	}
	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing SDN QinQ zone", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNZoneQinQResource) readIntoModel(ctx context.Context, model *SDNZoneQinQModel) error {
	zone, err := r.client.GetSDNZone(ctx, model.Zone.ValueString())
	if err != nil {
		return err
	}
	if zone.Type != "qinq" {
		return &client.APIError{StatusCode: 404, Status: "404 Not Found"}
	}

	model.Bridge = types.StringValue(zone.Bridge)
	model.Tag = types.Int64Value(int64(zone.Tag))
	model.VlanProtocol = types.StringValue(zone.VlanProtocol)
	model.Nodes = types.StringValue(zone.Nodes)
	model.IPAM = types.StringValue(zone.IPAM)
	model.DNS = types.StringValue(zone.DNS)
	model.DNSZone = types.StringValue(zone.DNSZone)
	model.ReverseDNS = types.StringValue(zone.ReverseDNS)
	return nil
}
