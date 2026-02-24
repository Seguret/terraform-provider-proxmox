package network_linux_vlan

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

var _ resource.Resource = &NetworkLinuxVLANResource{}
var _ resource.ResourceWithConfigure = &NetworkLinuxVLANResource{}
var _ resource.ResourceWithImportState = &NetworkLinuxVLANResource{}

type NetworkLinuxVLANResource struct {
	client *client.Client
}

type NetworkLinuxVLANResourceModel struct {
	ID            types.String `tfsdk:"id"`
	NodeName      types.String `tfsdk:"node_name"`
	Name          types.String `tfsdk:"name"`
	VLANID        types.Int64  `tfsdk:"vlan_id"`
	VLANRawDevice types.String `tfsdk:"vlan_raw_device"`
	Address       types.String `tfsdk:"address"`
	Autostart     types.Bool   `tfsdk:"autostart"`
	Comments      types.String `tfsdk:"comments"`
	Gateway       types.String `tfsdk:"gateway"`
	MTU           types.Int64  `tfsdk:"mtu"`
}

func NewResource() resource.Resource {
	return &NetworkLinuxVLANResource{}
}

func (r *NetworkLinuxVLANResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_network_linux_vlan"
}

func (r *NetworkLinuxVLANResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Linux VLAN network interface on a Proxmox VE node.",
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
				Description:   "The VLAN interface name (e.g., 'ens18.100').",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"vlan_id": schema.Int64Attribute{
				Description: "The VLAN tag (1-4094).",
				Optional:    true,
				Computed:    true,
			},
			"vlan_raw_device": schema.StringAttribute{
				Description: "The parent interface for this VLAN.",
				Optional:    true,
				Computed:    true,
			},
			"address": schema.StringAttribute{
				Description: "IPv4 address in CIDR notation (e.g., '192.168.100.1/24').",
				Optional:    true,
				Computed:    true,
			},
			"autostart": schema.BoolAttribute{
				Description: "Whether to bring up the interface at boot.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"comments": schema.StringAttribute{
				Description: "Comments for the VLAN interface.",
				Optional:    true,
				Computed:    true,
			},
			"gateway": schema.StringAttribute{
				Description: "IPv4 gateway address.",
				Optional:    true,
				Computed:    true,
			},
			"mtu": schema.Int64Attribute{
				Description: "The interface MTU.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *NetworkLinuxVLANResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworkLinuxVLANResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworkLinuxVLANResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	name := plan.Name.ValueString()

	tflog.Debug(ctx, "Creating Linux VLAN interface", map[string]any{"node": node, "name": name})

	createReq := r.modelToCreateRequest(&plan)
	if err := r.client.CreateNetworkInterface(ctx, node, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating Linux VLAN interface", err.Error())
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

func (r *NetworkLinuxVLANResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkLinuxVLANResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkLinuxVLANResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkLinuxVLANResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	name := plan.Name.ValueString()

	tflog.Debug(ctx, "Updating Linux VLAN interface", map[string]any{"node": node, "name": name})

	updateReq := r.modelToCreateRequest(&plan)
	if err := r.client.UpdateNetworkInterface(ctx, node, name, updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating Linux VLAN interface", err.Error())
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

func (r *NetworkLinuxVLANResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NetworkLinuxVLANResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := state.NodeName.ValueString()
	name := state.Name.ValueString()

	tflog.Debug(ctx, "Deleting Linux VLAN interface", map[string]any{"node": node, "name": name})

	if err := r.client.DeleteNetworkInterface(ctx, node, name); err != nil {
		resp.Diagnostics.AddError("Error deleting Linux VLAN interface", err.Error())
		return
	}

	if err := r.client.ApplyNetworkConfig(ctx, node); err != nil {
		resp.Diagnostics.AddError("Error applying network configuration", err.Error())
	}
}

func (r *NetworkLinuxVLANResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Import ID must be in format 'node_name/vlan_name'")
		return
	}
	state := NetworkLinuxVLANResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(parts[0]),
		Name:     types.StringValue(parts[1]),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkLinuxVLANResource) readIntoModel(ctx context.Context, model *NetworkLinuxVLANResourceModel, diagnostics *diag.Diagnostics) {
	ni, err := r.client.GetNetworkInterface(ctx, model.NodeName.ValueString(), model.Name.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diagnostics.AddWarning("Linux VLAN interface not found",
				fmt.Sprintf("Interface '%s' on node '%s' no longer exists.", model.Name.ValueString(), model.NodeName.ValueString()))
			return
		}
		diagnostics.AddError("Error reading Linux VLAN interface", err.Error())
		return
	}

	if ni.Type != "vlan" {
		diagnostics.AddError("Unexpected interface type",
			fmt.Sprintf("Expected type 'vlan', got '%s' for interface '%s'.", ni.Type, model.Name.ValueString()))
		return
	}

	model.VLANID = types.Int64Value(int64(ni.VLANID))
	model.VLANRawDevice = types.StringValue(ni.VLANRawDev)
	model.Address = types.StringValue(ni.CIDR)
	model.Autostart = types.BoolValue(ni.Autostart == 1)
	model.Comments = types.StringValue(ni.Comments)
	model.Gateway = types.StringValue(ni.Gateway)
	model.MTU = types.Int64Value(int64(ni.MTU))
}

func (r *NetworkLinuxVLANResource) modelToCreateRequest(model *NetworkLinuxVLANResourceModel) *models.NetworkInterfaceCreateRequest {
	req := &models.NetworkInterfaceCreateRequest{
		Iface:      model.Name.ValueString(),
		Type:       "vlan",
		CIDR:       model.Address.ValueString(),
		VLANRawDev: model.VLANRawDevice.ValueString(),
		Comments:   model.Comments.ValueString(),
		Gateway:    model.Gateway.ValueString(),
	}

	autostart := boolToInt(model.Autostart.ValueBool())
	req.Autostart = &autostart

	if !model.VLANID.IsNull() && !model.VLANID.IsUnknown() && model.VLANID.ValueInt64() > 0 {
		v := int(model.VLANID.ValueInt64())
		req.VLANID = &v
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
