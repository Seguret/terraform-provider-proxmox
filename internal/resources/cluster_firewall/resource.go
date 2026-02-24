package cluster_firewall

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &ClusterFirewallResource{}
var _ resource.ResourceWithConfigure = &ClusterFirewallResource{}

type ClusterFirewallResource struct {
	client *client.Client
}

type ClusterFirewallResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Enable               types.Bool   `tfsdk:"enable"`
	PolicyIn             types.String `tfsdk:"policy_in"`
	PolicyOut            types.String `tfsdk:"policy_out"`
	LogRatelimitEnable   types.Bool   `tfsdk:"log_ratelimit_enable"`
	LogRatelimitBurst    types.Int64  `tfsdk:"log_ratelimit_burst"`
	LogRatelimitRate     types.String `tfsdk:"log_ratelimit_rate"`
}

func NewResource() resource.Resource {
	return &ClusterFirewallResource{}
}

func (r *ClusterFirewallResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_cluster_firewall"
}

func (r *ClusterFirewallResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages cluster-level firewall options for Proxmox VE. This is a singleton resource — one per cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enable": schema.BoolAttribute{
				Description: "Whether the cluster firewall is enabled.",
				Optional:    true,
				Computed:    true,
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
			"log_ratelimit_enable": schema.BoolAttribute{
				Description: "Whether log rate limiting is enabled.",
				Optional:    true,
				Computed:    true,
			},
			"log_ratelimit_burst": schema.Int64Attribute{
				Description: "Initial burst value for log rate limiting.",
				Optional:    true,
				Computed:    true,
			},
			"log_ratelimit_rate": schema.StringAttribute{
				Description: "Rate for log rate limiting (e.g. '1/second').",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *ClusterFirewallResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ClusterFirewallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ClusterFirewallResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateFirewallOptions(ctx, client.ClusterFirewallPath(), r.modelToOpts(&plan)); err != nil {
		resp.Diagnostics.AddError("Error setting cluster firewall options", err.Error())
		return
	}

	plan.ID = types.StringValue("cluster")

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading cluster firewall options", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ClusterFirewallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ClusterFirewallResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error reading cluster firewall options", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ClusterFirewallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ClusterFirewallResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateFirewallOptions(ctx, client.ClusterFirewallPath(), r.modelToOpts(&plan)); err != nil {
		resp.Diagnostics.AddError("Error updating cluster firewall options", err.Error())
		return
	}

	plan.ID = types.StringValue("cluster")

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading cluster firewall options", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ClusterFirewallResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// cant delete cluster firewall config — just remove from state
}

func (r *ClusterFirewallResource) readIntoModel(ctx context.Context, model *ClusterFirewallResourceModel) error {
	opts, err := r.client.GetFirewallOptions(ctx, client.ClusterFirewallPath())
	if err != nil {
		return err
	}

	if opts.Enable != nil {
		model.Enable = types.BoolValue(*opts.Enable == 1)
	} else {
		model.Enable = types.BoolValue(false)
	}

	model.PolicyIn = types.StringValue(opts.PolicyIn)
	model.PolicyOut = types.StringValue(opts.PolicyOut)

	rlEnable, rlBurst, rlRate := parseLogRatelimit(opts.LogRatelimit)
	model.LogRatelimitEnable = rlEnable
	model.LogRatelimitBurst = rlBurst
	model.LogRatelimitRate = rlRate

	return nil
}

func (r *ClusterFirewallResource) modelToOpts(model *ClusterFirewallResourceModel) *models.FirewallOptions {
	enableInt := boolToInt(model.Enable.ValueBool())
	opts := &models.FirewallOptions{
		Enable:       &enableInt,
		PolicyIn:     model.PolicyIn.ValueString(),
		PolicyOut:    model.PolicyOut.ValueString(),
		LogRatelimit: formatLogRatelimit(model.LogRatelimitEnable, model.LogRatelimitBurst, model.LogRatelimitRate),
	}
	return opts
}

// parseLogRatelimit splits a proxmox log_ratelimit string like "enable=1,rate=1/second,burst=5" into its pieces.
func parseLogRatelimit(s string) (enable types.Bool, burst types.Int64, rate types.String) {
	enable = types.BoolValue(false)
	burst = types.Int64Value(0)
	rate = types.StringValue("")

	if s == "" {
		return
	}

	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		switch key {
		case "enable":
			enable = types.BoolValue(val == "1")
		case "burst":
			if n, err := strconv.ParseInt(val, 10, 64); err == nil {
				burst = types.Int64Value(n)
			}
		case "rate":
			rate = types.StringValue(val)
		}
	}
	return
}

// formatLogRatelimit packs the three log ratelimit fields back into the proxmox string format.
func formatLogRatelimit(enable types.Bool, burst types.Int64, rate types.String) string {
	if enable.IsNull() && burst.IsNull() && rate.IsNull() {
		return ""
	}
	if enable.IsUnknown() || burst.IsUnknown() || rate.IsUnknown() {
		return ""
	}

	var parts []string

	enableVal := "0"
	if !enable.IsNull() && enable.ValueBool() {
		enableVal = "1"
	}
	parts = append(parts, "enable="+enableVal)

	if !rate.IsNull() && rate.ValueString() != "" {
		parts = append(parts, "rate="+rate.ValueString())
	}
	if !burst.IsNull() && burst.ValueInt64() != 0 {
		parts = append(parts, "burst="+strconv.FormatInt(burst.ValueInt64(), 10))
	}

	if len(parts) == 1 && parts[0] == "enable=0" {
		return ""
	}

	return strings.Join(parts, ",")
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
