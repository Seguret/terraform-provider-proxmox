package network_interface

import (
	"context"
	"fmt"
	"strings"

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

var _ resource.Resource = &NetworkInterfaceResource{}
var _ resource.ResourceWithConfigure = &NetworkInterfaceResource{}
var _ resource.ResourceWithImportState = &NetworkInterfaceResource{}

type NetworkInterfaceResource struct {
	client *client.Client
}

type NetworkInterfaceResourceModel struct {
	ID       types.String `tfsdk:"id"`
	NodeName types.String `tfsdk:"node_name"`
	Iface    types.String `tfsdk:"iface"`
	Type     types.String `tfsdk:"type"`

	// IP
	Method   types.String `tfsdk:"method"`
	Method6  types.String `tfsdk:"method6"`
	Address  types.String `tfsdk:"address"`
	Netmask  types.String `tfsdk:"netmask"`
	Gateway  types.String `tfsdk:"gateway"`
	Address6 types.String `tfsdk:"address6"`
	Netmask6 types.Int64  `tfsdk:"netmask6"`
	Gateway6 types.String `tfsdk:"gateway6"`
	CIDR     types.String `tfsdk:"cidr"`
	CIDR6    types.String `tfsdk:"cidr6"`

	// Bridge
	BridgePorts     types.String `tfsdk:"bridge_ports"`
	BridgeSTP       types.String `tfsdk:"bridge_stp"`
	BridgeFD        types.Int64  `tfsdk:"bridge_fd"`
	BridgeVLANAware types.Bool   `tfsdk:"bridge_vlan_aware"`

	// Bond
	BondPrimary        types.String `tfsdk:"bond_primary"`
	BondMode           types.String `tfsdk:"bond_mode"`
	BondXmitHashPolicy types.String `tfsdk:"bond_xmit_hash_policy"`
	Slaves             types.String `tfsdk:"slaves"`

	// VLAN
	VLANRawDev types.String `tfsdk:"vlan_raw_device"`
	VLANID     types.Int64  `tfsdk:"vlan_id"`

	// Misc
	Autostart types.Bool   `tfsdk:"autostart"`
	MTU       types.Int64  `tfsdk:"mtu"`
	Comments  types.String `tfsdk:"comments"`
	Comments6 types.String `tfsdk:"comments6"`

	// Apply after changes
	ApplyConfig types.Bool `tfsdk:"apply_config"`
}

func NewResource() resource.Resource {
	return &NetworkInterfaceResource{}
}

func (r *NetworkInterfaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_network_interface"
}

func (r *NetworkInterfaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a network interface on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"node_name": schema.StringAttribute{
				Description: "The node name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"iface": schema.StringAttribute{
				Description: "The interface name (e.g., 'vmbr0', 'bond0', 'eth0.100').",
				Required:    true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"type": schema.StringAttribute{
				Description: "The interface type (bridge, bond, eth, vlan, OVSBridge, OVSBond, OVSPort, OVSIntPort).",
				Required:    true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			// IP
			"method":   schema.StringAttribute{Description: "IPv4 method (static, dhcp, manual).", Optional: true, Computed: true},
			"method6":  schema.StringAttribute{Description: "IPv6 method (static, dhcp, manual).", Optional: true, Computed: true},
			"address":  schema.StringAttribute{Description: "IPv4 address.", Optional: true, Computed: true},
			"netmask":  schema.StringAttribute{Description: "IPv4 netmask.", Optional: true, Computed: true},
			"gateway":  schema.StringAttribute{Description: "IPv4 gateway.", Optional: true, Computed: true},
			"address6": schema.StringAttribute{Description: "IPv6 address.", Optional: true, Computed: true},
			"netmask6": schema.Int64Attribute{Description: "IPv6 prefix length.", Optional: true, Computed: true},
			"gateway6": schema.StringAttribute{Description: "IPv6 gateway.", Optional: true, Computed: true},
			"cidr":     schema.StringAttribute{Description: "IPv4 CIDR (e.g., '192.168.1.1/24').", Optional: true, Computed: true},
			"cidr6":    schema.StringAttribute{Description: "IPv6 CIDR.", Optional: true, Computed: true},
			// Bridge
			"bridge_ports": schema.StringAttribute{Description: "Bridge ports (space-separated interface names).", Optional: true, Computed: true},
			"bridge_stp":   schema.StringAttribute{Description: "STP mode (on or off).", Optional: true, Computed: true},
			"bridge_fd":    schema.Int64Attribute{Description: "Bridge forward delay.", Optional: true, Computed: true, Default: int64default.StaticInt64(0)},
			"bridge_vlan_aware": schema.BoolAttribute{Description: "Whether to enable VLAN-aware bridge.", Optional: true, Computed: true, Default: booldefault.StaticBool(false)},
			// Bond
			"bond_primary":          schema.StringAttribute{Description: "Primary bond interface.", Optional: true, Computed: true},
			"bond_mode":             schema.StringAttribute{Description: "Bond mode (balance-rr, active-backup, balance-xor, etc.).", Optional: true, Computed: true},
			"bond_xmit_hash_policy": schema.StringAttribute{Description: "Bond transmit hash policy.", Optional: true, Computed: true},
			"slaves":                schema.StringAttribute{Description: "Slave interfaces for bond.", Optional: true, Computed: true},
			// VLAN
			"vlan_raw_device": schema.StringAttribute{Description: "The underlying device for the VLAN interface.", Optional: true, Computed: true},
			"vlan_id":         schema.Int64Attribute{Description: "VLAN ID.", Optional: true, Computed: true},
			// Misc
			"autostart": schema.BoolAttribute{Description: "Whether to bring up the interface at boot.", Optional: true, Computed: true, Default: booldefault.StaticBool(true)},
			"mtu":       schema.Int64Attribute{Description: "Interface MTU.", Optional: true, Computed: true},
			"comments":  schema.StringAttribute{Description: "Comments for the interface.", Optional: true, Computed: true},
			"comments6": schema.StringAttribute{Description: "IPv6 comments.", Optional: true, Computed: true},
			"apply_config": schema.BoolAttribute{
				Description: "Whether to apply the network configuration immediately after changes.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *NetworkInterfaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworkInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworkInterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	iface := plan.Iface.ValueString()

	autostart := boolToInt(plan.Autostart.ValueBool())
	createReq := r.modelToCreateRequest(&plan)
	createReq.Autostart = &autostart

	tflog.Debug(ctx, "Creating network interface", map[string]any{"node": node, "iface": iface, "type": plan.Type.ValueString()})

	if err := r.client.CreateNetworkInterface(ctx, node, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating network interface", err.Error())
		return
	}

	if plan.ApplyConfig.ValueBool() {
		if err := r.client.ApplyNetworkConfig(ctx, node); err != nil {
			resp.Diagnostics.AddError("Error applying network config", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", node, iface))
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkInterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkInterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	iface := plan.Iface.ValueString()
	updateReq := r.modelToCreateRequest(&plan)

	if err := r.client.UpdateNetworkInterface(ctx, node, iface, updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating network interface", err.Error())
		return
	}

	if plan.ApplyConfig.ValueBool() {
		if err := r.client.ApplyNetworkConfig(ctx, node); err != nil {
			resp.Diagnostics.AddError("Error applying network config", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", node, iface))
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NetworkInterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := state.NodeName.ValueString()
	iface := state.Iface.ValueString()

	if err := r.client.DeleteNetworkInterface(ctx, node, iface); err != nil {
		resp.Diagnostics.AddError("Error deleting network interface", err.Error())
		return
	}

	if state.ApplyConfig.ValueBool() {
		if err := r.client.ApplyNetworkConfig(ctx, node); err != nil {
			resp.Diagnostics.AddError("Error applying network config", err.Error())
		}
	}
}

func (r *NetworkInterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Import ID must be in format 'node_name/iface_name'")
		return
	}
	state := NetworkInterfaceResourceModel{
		ID:          types.StringValue(req.ID),
		NodeName:    types.StringValue(parts[0]),
		Iface:       types.StringValue(parts[1]),
		ApplyConfig: types.BoolValue(false),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkInterfaceResource) readIntoModel(ctx context.Context, model *NetworkInterfaceResourceModel, diagnostics *diag.Diagnostics) {
	ni, err := r.client.GetNetworkInterface(ctx, model.NodeName.ValueString(), model.Iface.ValueString())
	if err != nil {
		diagnostics.AddError("Error reading network interface", err.Error())
		return
	}

	model.Type = types.StringValue(ni.Type)
	model.Autostart = types.BoolValue(ni.Autostart == 1)
	model.Method = types.StringValue(ni.Method)
	model.Method6 = types.StringValue(ni.Method6)
	model.Address = types.StringValue(ni.Address)
	model.Netmask = types.StringValue(ni.Netmask)
	model.Gateway = types.StringValue(ni.Gateway)
	model.Address6 = types.StringValue(ni.Address6)
	model.Netmask6 = types.Int64Value(int64(ni.Netmask6))
	model.Gateway6 = types.StringValue(ni.Gateway6)
	model.CIDR = types.StringValue(ni.CIDR)
	model.CIDR6 = types.StringValue(ni.CIDR6)
	model.BridgePorts = types.StringValue(ni.BridgePorts)
	model.BridgeSTP = types.StringValue(ni.BridgeSTP)
	model.BridgeFD = types.Int64Value(int64(ni.BridgeFD))
	model.BridgeVLANAware = types.BoolValue(ni.BridgeVLANAware == 1)
	model.BondPrimary = types.StringValue(ni.BondPrimary)
	model.BondMode = types.StringValue(ni.BondMode)
	model.BondXmitHashPolicy = types.StringValue(ni.BondXmitHashPolicy)
	model.Slaves = types.StringValue(ni.Slaves)
	model.VLANRawDev = types.StringValue(ni.VLANRawDev)
	model.VLANID = types.Int64Value(int64(ni.VLANID))
	model.MTU = types.Int64Value(int64(ni.MTU))
	model.Comments = types.StringValue(ni.Comments)
	model.Comments6 = types.StringValue(ni.Comments6)
}

func (r *NetworkInterfaceResource) modelToCreateRequest(model *NetworkInterfaceResourceModel) *models.NetworkInterfaceCreateRequest {
	req := &models.NetworkInterfaceCreateRequest{
		Iface:              model.Iface.ValueString(),
		Type:               model.Type.ValueString(),
		Method:             model.Method.ValueString(),
		Method6:            model.Method6.ValueString(),
		Address:            model.Address.ValueString(),
		Netmask:            model.Netmask.ValueString(),
		Gateway:            model.Gateway.ValueString(),
		Address6:           model.Address6.ValueString(),
		Gateway6:           model.Gateway6.ValueString(),
		CIDR:               model.CIDR.ValueString(),
		CIDR6:              model.CIDR6.ValueString(),
		BridgePorts:        model.BridgePorts.ValueString(),
		BridgeSTP:          model.BridgeSTP.ValueString(),
		BondPrimary:        model.BondPrimary.ValueString(),
		BondMode:           model.BondMode.ValueString(),
		BondXmitHashPolicy: model.BondXmitHashPolicy.ValueString(),
		Slaves:             model.Slaves.ValueString(),
		VLANRawDev:         model.VLANRawDev.ValueString(),
		Comments:           model.Comments.ValueString(),
		Comments6:          model.Comments6.ValueString(),
	}

	if !model.BridgeFD.IsNull() {
		v := int(model.BridgeFD.ValueInt64())
		req.BridgeFD = &v
	}
	if !model.BridgeVLANAware.IsNull() {
		v := boolToInt(model.BridgeVLANAware.ValueBool())
		req.BridgeVLANAware = &v
	}
	if !model.Netmask6.IsNull() {
		v := int(model.Netmask6.ValueInt64())
		req.Netmask6 = &v
	}
	if !model.VLANID.IsNull() && model.VLANID.ValueInt64() > 0 {
		v := int(model.VLANID.ValueInt64())
		req.VLANID = &v
	}
	if !model.MTU.IsNull() && model.MTU.ValueInt64() > 0 {
		v := int(model.MTU.ValueInt64())
		req.MTU = &v
	}
	return req
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
