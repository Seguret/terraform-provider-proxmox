package firewall_options

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &FirewallOptionsResource{}
var _ resource.ResourceWithConfigure = &FirewallOptionsResource{}

type FirewallOptionsResource struct {
	client *client.Client
}

type FirewallOptionsResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Scope     types.String `tfsdk:"scope"`
	Enabled   types.Bool   `tfsdk:"enabled"`
	PolicyIn  types.String `tfsdk:"policy_in"`
	PolicyOut types.String `tfsdk:"policy_out"`
	Ebtables  types.Bool   `tfsdk:"ebtables"`
	IPFilter  types.Bool   `tfsdk:"ip_filter"`
	Macfilter types.Bool   `tfsdk:"mac_filter"`
	NDPProxy  types.Bool   `tfsdk:"ndp"`
	DHCP      types.Bool   `tfsdk:"dhcp"`
	LogRatelimit types.String `tfsdk:"log_ratelimit"`
}

func NewResource() resource.Resource {
	return &FirewallOptionsResource{}
}

func (r *FirewallOptionsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_firewall_options"
}

func (r *FirewallOptionsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages firewall options for cluster, node, VM, or container.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"scope": schema.StringAttribute{
				Description: "The firewall scope: 'cluster', 'node/<node>', 'vm/<node>/<vmid>', 'ct/<node>/<vmid>'.",
				Required:    true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the firewall is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"policy_in": schema.StringAttribute{
				Description: "Default inbound policy (ACCEPT, DROP, REJECT).",
				Optional:    true,
				Computed:    true,
			},
			"policy_out": schema.StringAttribute{
				Description: "Default outbound policy (ACCEPT, DROP, REJECT).",
				Optional:    true,
				Computed:    true,
			},
			"ebtables": schema.BoolAttribute{
				Description: "Whether to enable Ethernet bridge filtering.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"ip_filter": schema.BoolAttribute{
				Description: "Whether to enable IP filtering.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"mac_filter": schema.BoolAttribute{
				Description: "Whether to enable MAC address filtering.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"ndp": schema.BoolAttribute{
				Description: "Whether to enable NDP (Neighbor Discovery Protocol).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"dhcp": schema.BoolAttribute{
				Description: "Whether to allow DHCP traffic.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"log_ratelimit": schema.StringAttribute{
				Description: "Log rate limit (e.g., 'enable=1,rate=1/second,burst=5').",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *FirewallOptionsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FirewallOptionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallOptionsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := r.scopeToPath(plan.Scope.ValueString())
	if err := r.client.UpdateFirewallOptions(ctx, pathPrefix, r.modelToOpts(&plan)); err != nil {
		resp.Diagnostics.AddError("Error setting firewall options", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Scope.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallOptionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FirewallOptionsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FirewallOptionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FirewallOptionsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := r.scopeToPath(plan.Scope.ValueString())
	if err := r.client.UpdateFirewallOptions(ctx, pathPrefix, r.modelToOpts(&plan)); err != nil {
		resp.Diagnostics.AddError("Error updating firewall options", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Scope.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallOptionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// cant really delete firewall options — just disable it instead
	var state FirewallOptionsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	disabled := 0
	pathPrefix := r.scopeToPath(state.Scope.ValueString())
	_ = r.client.UpdateFirewallOptions(ctx, pathPrefix, &models.FirewallOptions{Enable: &disabled})
}

func (r *FirewallOptionsResource) readIntoModel(ctx context.Context, model *FirewallOptionsResourceModel, diagnostics *diag.Diagnostics) {
	pathPrefix := r.scopeToPath(model.Scope.ValueString())
	opts, err := r.client.GetFirewallOptions(ctx, pathPrefix)
	if err != nil {
		diagnostics.AddError("Error reading firewall options", err.Error())
		return
	}

	if opts.Enable != nil {
		model.Enabled = types.BoolValue(*opts.Enable == 1)
	} else {
		model.Enabled = types.BoolValue(false)
	}
	model.PolicyIn = types.StringValue(opts.PolicyIn)
	model.PolicyOut = types.StringValue(opts.PolicyOut)
	model.LogRatelimit = types.StringValue(opts.LogRatelimit)

	if opts.Ebtables != nil {
		model.Ebtables = types.BoolValue(*opts.Ebtables == 1)
	}
	if opts.Macfilter != nil {
		model.Macfilter = types.BoolValue(*opts.Macfilter == 1)
	}
	if opts.IPFilter != nil {
		model.IPFilter = types.BoolValue(*opts.IPFilter == 1)
	}
	if opts.NDPProxy != nil {
		model.NDPProxy = types.BoolValue(*opts.NDPProxy == 1)
	}
	if opts.DHCPFilter != nil {
		model.DHCP = types.BoolValue(*opts.DHCPFilter == 1)
	}
}

func (r *FirewallOptionsResource) modelToOpts(model *FirewallOptionsResourceModel) *models.FirewallOptions {
	enableInt := boolToInt(model.Enabled.ValueBool())
	ebtablesInt := boolToInt(model.Ebtables.ValueBool())
	macfilterInt := boolToInt(model.Macfilter.ValueBool())
	ipfilterInt := boolToInt(model.IPFilter.ValueBool())
	ndpInt := boolToInt(model.NDPProxy.ValueBool())
	dhcpInt := boolToInt(model.DHCP.ValueBool())

	opts := &models.FirewallOptions{
		Enable:       &enableInt,
		PolicyIn:     model.PolicyIn.ValueString(),
		PolicyOut:    model.PolicyOut.ValueString(),
		Ebtables:     &ebtablesInt,
		Macfilter:    &macfilterInt,
		IPFilter:     &ipfilterInt,
		NDPProxy:     &ndpInt,
		DHCPFilter:   &dhcpInt,
		LogRatelimit: model.LogRatelimit.ValueString(),
	}
	return opts
}

func (r *FirewallOptionsResource) scopeToPath(scope string) string {
	parts := strings.Split(scope, "/")
	switch parts[0] {
	case "cluster":
		return client.ClusterFirewallPath()
	case "node":
		if len(parts) >= 2 {
			return client.NodeFirewallPath(parts[1])
		}
	case "vm":
		if len(parts) >= 3 {
			vmid, _ := strconv.Atoi(parts[2])
			return client.VMFirewallPath(parts[1], vmid)
		}
	case "ct":
		if len(parts) >= 3 {
			vmid, _ := strconv.Atoi(parts[2])
			return client.ContainerFirewallPath(parts[1], vmid)
		}
	}
	return client.ClusterFirewallPath()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
