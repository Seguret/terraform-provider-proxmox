package network_linux_bridge

import (
	"context"
	"fmt"
	"strings"

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

var _ resource.Resource = &NetworkLinuxBridgeResource{}
var _ resource.ResourceWithConfigure = &NetworkLinuxBridgeResource{}
var _ resource.ResourceWithImportState = &NetworkLinuxBridgeResource{}

type NetworkLinuxBridgeResource struct {
	client *client.Client
}

type NetworkLinuxBridgeResourceModel struct {
	ID             types.String `tfsdk:"id"`
	NodeName       types.String `tfsdk:"node_name"`
	Name           types.String `tfsdk:"name"`
	Address        types.String `tfsdk:"address"`
	Address6       types.String `tfsdk:"address6"`
	Autostart      types.Bool   `tfsdk:"autostart"`
	BridgePorts    types.String `tfsdk:"bridge_ports"`
	BridgeVLANAware types.Bool  `tfsdk:"bridge_vlan_aware"`
	Comments       types.String `tfsdk:"comments"`
	MTU            types.Int64  `tfsdk:"mtu"`
	Gateway        types.String `tfsdk:"gateway"`
	Gateway6       types.String `tfsdk:"gateway6"`
}

func NewResource() resource.Resource {
	return &NetworkLinuxBridgeResource{}
}

func (r *NetworkLinuxBridgeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_network_linux_bridge"
}

func (r *NetworkLinuxBridgeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Linux bridge network interface on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"node_name": schema.StringAttribute{
				Description:   "The name of the Proxmox VE node.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "The bridge interface name (e.g., 'vmbr0').",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"address": schema.StringAttribute{
				Description: "IPv4 address in CIDR notation (e.g., '192.168.1.1/24').",
				Optional:    true,
				Computed:    true,
			},
			"address6": schema.StringAttribute{
				Description: "IPv6 address in CIDR notation.",
				Optional:    true,
				Computed:    true,
			},
			"autostart": schema.BoolAttribute{
				Description: "Whether to bring up the interface at boot.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"bridge_ports": schema.StringAttribute{
				Description: "Space-separated list of ports to add to the bridge.",
				Optional:    true,
				Computed:    true,
			},
			"bridge_vlan_aware": schema.BoolAttribute{
				Description: "Whether the bridge is VLAN aware.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"comments": schema.StringAttribute{
				Description: "Comments for the bridge interface.",
				Optional:    true,
				Computed:    true,
			},
			"mtu": schema.Int64Attribute{
				Description: "The interface MTU.",
				Optional:    true,
				Computed:    true,
			},
			"gateway": schema.StringAttribute{
				Description: "IPv4 gateway address.",
				Optional:    true,
				Computed:    true,
			},
			"gateway6": schema.StringAttribute{
				Description: "IPv6 gateway address.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *NetworkLinuxBridgeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworkLinuxBridgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworkLinuxBridgeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	name := plan.Name.ValueString()

	tflog.Debug(ctx, "Creating Linux bridge interface", map[string]any{"node": node, "name": name})

	createReq := r.modelToCreateRequest(&plan)
	if err := r.client.CreateNetworkInterface(ctx, node, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating Linux bridge interface", err.Error())
		return
	}

	if err := r.client.ApplyNetworkConfig(ctx, node); err != nil {
		resp.Diagnostics.AddError("Error applying network configuration", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", node, name))
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkLinuxBridgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkLinuxBridgeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkLinuxBridgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkLinuxBridgeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	name := plan.Name.ValueString()

	tflog.Debug(ctx, "Updating Linux bridge interface", map[string]any{"node": node, "name": name})

	updateReq := r.modelToCreateRequest(&plan)
	if err := r.client.UpdateNetworkInterface(ctx, node, name, updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating Linux bridge interface", err.Error())
		return
	}

	if err := r.client.ApplyNetworkConfig(ctx, node); err != nil {
		resp.Diagnostics.AddError("Error applying network configuration", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", node, name))
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkLinuxBridgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NetworkLinuxBridgeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := state.NodeName.ValueString()
	name := state.Name.ValueString()

	tflog.Debug(ctx, "Deleting Linux bridge interface", map[string]any{"node": node, "name": name})

	if err := r.client.DeleteNetworkInterface(ctx, node, name); err != nil {
		resp.Diagnostics.AddError("Error deleting Linux bridge interface", err.Error())
		return
	}

	if err := r.client.ApplyNetworkConfig(ctx, node); err != nil {
		resp.Diagnostics.AddError("Error applying network configuration", err.Error())
	}
}

func (r *NetworkLinuxBridgeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Import ID must be in format 'node_name/bridge_name'")
		return
	}
	state := NetworkLinuxBridgeResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(parts[0]),
		Name:     types.StringValue(parts[1]),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkLinuxBridgeResource) readIntoModel(ctx context.Context, model *NetworkLinuxBridgeResourceModel, diagnostics *diag.Diagnostics) {
	ni, err := r.client.GetNetworkInterface(ctx, model.NodeName.ValueString(), model.Name.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diagnostics.AddWarning("Linux bridge interface not found",
				fmt.Sprintf("Interface '%s' on node '%s' no longer exists.", model.Name.ValueString(), model.NodeName.ValueString()))
			return
		}
		diagnostics.AddError("Error reading Linux bridge interface", err.Error())
		return
	}

	if ni.Type != "bridge" {
		diagnostics.AddError("Unexpected interface type",
			fmt.Sprintf("Expected type 'bridge', got '%s' for interface '%s'.", ni.Type, model.Name.ValueString()))
		return
	}

	model.Address = types.StringValue(ni.CIDR)
	model.Address6 = types.StringValue(ni.CIDR6)
	model.Autostart = types.BoolValue(ni.Autostart == 1)
	model.BridgePorts = types.StringValue(ni.BridgePorts)
	model.BridgeVLANAware = types.BoolValue(ni.BridgeVLANAware == 1)
	model.Comments = types.StringValue(ni.Comments)
	model.MTU = types.Int64Value(int64(ni.MTU))
	model.Gateway = types.StringValue(ni.Gateway)
	model.Gateway6 = types.StringValue(ni.Gateway6)
}

func (r *NetworkLinuxBridgeResource) modelToCreateRequest(model *NetworkLinuxBridgeResourceModel) *models.NetworkInterfaceCreateRequest {
	req := &models.NetworkInterfaceCreateRequest{
		Iface:       model.Name.ValueString(),
		Type:        "bridge",
		CIDR:        model.Address.ValueString(),
		CIDR6:       model.Address6.ValueString(),
		BridgePorts: model.BridgePorts.ValueString(),
		Comments:    model.Comments.ValueString(),
		Gateway:     model.Gateway.ValueString(),
		Gateway6:    model.Gateway6.ValueString(),
	}

	autostart := boolToInt(model.Autostart.ValueBool())
	req.Autostart = &autostart

	if !model.BridgeVLANAware.IsNull() && !model.BridgeVLANAware.IsUnknown() {
		v := boolToInt(model.BridgeVLANAware.ValueBool())
		req.BridgeVLANAware = &v
	}

	if !model.MTU.IsNull() && !model.MTU.IsUnknown() && model.MTU.ValueInt64() > 0 {
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
