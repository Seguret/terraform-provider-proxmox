package sdn_subnet

import (
	"context"
	"fmt"
	"strings"

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

var _ resource.Resource = &SDNSubnetResource{}
var _ resource.ResourceWithConfigure = &SDNSubnetResource{}
var _ resource.ResourceWithImportState = &SDNSubnetResource{}

type SDNSubnetResource struct {
	client *client.Client
}

type SDNSubnetResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Vnet          types.String `tfsdk:"vnet"`
	Subnet        types.String `tfsdk:"subnet"`
	Gateway       types.String `tfsdk:"gateway"`
	Snat          types.Bool   `tfsdk:"snat"`
	DHCPDNSServer types.String `tfsdk:"dhcp_dns_server"`
	DNSZonePrefix types.String `tfsdk:"dns_zone_prefix"`
}

func NewResource() resource.Resource {
	return &SDNSubnetResource{}
}

func (r *SDNSubnetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_sdn_subnet"
}

func (r *SDNSubnetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE SDN subnet within a VNet.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vnet": schema.StringAttribute{
				Description: "The parent VNet name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet": schema.StringAttribute{
				Description: "The subnet CIDR (e.g. '10.0.0.0/24').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"gateway": schema.StringAttribute{
				Description: "The subnet gateway IP address.",
				Optional:    true,
				Computed:    true,
			},
			"snat": schema.BoolAttribute{
				Description: "Enable SNAT for outbound traffic.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"dhcp_dns_server": schema.StringAttribute{
				Description: "DNS server pushed to DHCP clients.",
				Optional:    true,
				Computed:    true,
			},
			"dns_zone_prefix": schema.StringAttribute{
				Description: "DNS zone prefix for PTR records.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *SDNSubnetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SDNSubnetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SDNSubnetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	snatInt := boolToIntPtr(plan.Snat.ValueBool())
	vnet := plan.Vnet.ValueString()
	subnet := plan.Subnet.ValueString()

	tflog.Debug(ctx, "Creating SDN subnet", map[string]any{"vnet": vnet, "subnet": subnet})

	if err := r.client.CreateSDNSubnet(ctx, vnet, &models.SDNSubnetCreateRequest{
		Subnet:        subnet,
		Type:          "subnet",
		Gateway:       plan.Gateway.ValueString(),
		Snat:          snatInt,
		DHCPDNSServer: plan.DHCPDNSServer.ValueString(),
		DNSZonePrefix: plan.DNSZonePrefix.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error creating SDN subnet", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", vnet, subnet))
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNSubnetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SDNSubnetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNSubnetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SDNSubnetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	snatInt := boolToIntPtr(plan.Snat.ValueBool())

	if err := r.client.UpdateSDNSubnet(ctx, plan.Vnet.ValueString(), cidrToAPIID(plan.Subnet.ValueString()), &models.SDNSubnetUpdateRequest{
		Gateway:       plan.Gateway.ValueString(),
		Snat:          snatInt,
		DHCPDNSServer: plan.DHCPDNSServer.ValueString(),
		DNSZonePrefix: plan.DNSZonePrefix.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error updating SDN subnet", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNSubnetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SDNSubnetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSDNSubnet(ctx, state.Vnet.ValueString(), cidrToAPIID(state.Subnet.ValueString())); err != nil {
		resp.Diagnostics.AddError("Error deleting SDN subnet", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)
}

func (r *SDNSubnetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// format: "<vnet>/<subnet-cidr>" e.g. "myvnet/10.0.0.0/24"
	// subnet CIDRs have their own '/', so only split on the first one
	idx := strings.Index(req.ID, "/")
	if idx < 0 {
		resp.Diagnostics.AddError("Invalid import ID", "Format: <vnet>/<subnet-cidr> (e.g. 'myvnet/10.0.0.0/24')")
		return
	}

	vnet := req.ID[:idx]
	subnet := req.ID[idx+1:]

	state := SDNSubnetResourceModel{
		ID:     types.StringValue(req.ID),
		Vnet:   types.StringValue(vnet),
		Subnet: types.StringValue(subnet),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNSubnetResource) readIntoModel(ctx context.Context, model *SDNSubnetResourceModel, diagnostics interface{ AddError(string, string) }) {
	s, err := r.client.GetSDNSubnet(ctx, model.Vnet.ValueString(), cidrToAPIID(model.Subnet.ValueString()))
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diagnostics.AddError("SDN subnet not found", "The SDN subnet no longer exists.")
			return
		}
		diagnostics.AddError("Error reading SDN subnet", err.Error())
		return
	}

	model.Gateway = types.StringValue(s.Gateway)
	model.DHCPDNSServer = types.StringValue(s.DHCPDNSServer)
	model.DNSZonePrefix = types.StringValue(s.DNSZonePrefix)

	if s.Snat != nil {
		model.Snat = types.BoolValue(*s.Snat == 1)
	} else {
		model.Snat = types.BoolValue(false)
	}
}

// cidrToAPIID converts "10.0.0.0/24" → "10.0.0.0-24" for the Proxmox API.
func cidrToAPIID(cidr string) string {
	return strings.ReplaceAll(cidr, "/", "-")
}

func boolToIntPtr(b bool) *int {
	v := 0
	if b {
		v = 1
	}
	return &v
}
