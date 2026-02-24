package firewall_rule

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
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &FirewallRuleResource{}
var _ resource.ResourceWithConfigure = &FirewallRuleResource{}
var _ resource.ResourceWithImportState = &FirewallRuleResource{}

type FirewallRuleResource struct {
	client *client.Client
}

// FirewallRuleResourceModel - scope encodes where the rule lives:
//   - "cluster"          -> cluster firewall
//   - "node/pve1"        -> node firewall
//   - "vm/pve1/100"      -> VM firewall
//   - "ct/pve1/100"      -> container firewall
type FirewallRuleResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Scope   types.String `tfsdk:"scope"`
	Pos     types.Int64  `tfsdk:"pos"`
	Type    types.String `tfsdk:"type"`
	Action  types.String `tfsdk:"action"`
	Enabled types.Bool   `tfsdk:"enabled"`
	Macro   types.String `tfsdk:"macro"`
	Proto   types.String `tfsdk:"proto"`
	Source  types.String `tfsdk:"source"`
	Dest    types.String `tfsdk:"dest"`
	DPort   types.String `tfsdk:"dport"`
	Sport   types.String `tfsdk:"sport"`
	IFace   types.String `tfsdk:"iface"`
	Log     types.String `tfsdk:"log"`
	Comment types.String `tfsdk:"comment"`
}

func NewResource() resource.Resource {
	return &FirewallRuleResource{}
}

func (r *FirewallRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_firewall_rule"
}

func (r *FirewallRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE firewall rule. Scope determines whether the rule applies to cluster, node, VM, or container.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"scope": schema.StringAttribute{
				Description: "The firewall scope: 'cluster', 'node/<node>', 'vm/<node>/<vmid>', 'ct/<node>/<vmid>'.",
				Required:    true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"pos": schema.Int64Attribute{
				Description: "The rule position. Computed after creation.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The rule direction (in or out).",
				Required:    true,
			},
			"action": schema.StringAttribute{
				Description: "The rule action (ACCEPT, DROP, REJECT).",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the rule is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"macro":   schema.StringAttribute{Description: "Macro name (e.g., 'SSH', 'HTTP').", Optional: true, Computed: true},
			"proto":   schema.StringAttribute{Description: "Protocol (tcp, udp, icmp, etc.).", Optional: true, Computed: true},
			"source":  schema.StringAttribute{Description: "Source address/CIDR/IPset.", Optional: true, Computed: true},
			"dest":    schema.StringAttribute{Description: "Destination address/CIDR/IPset.", Optional: true, Computed: true},
			"dport":   schema.StringAttribute{Description: "Destination port(s).", Optional: true, Computed: true},
			"sport":   schema.StringAttribute{Description: "Source port(s).", Optional: true, Computed: true},
			"iface":   schema.StringAttribute{Description: "Network interface.", Optional: true, Computed: true},
			"log":     schema.StringAttribute{Description: "Log level (emerg, alert, crit, err, warning, notice, info, debug, nolog).", Optional: true, Computed: true},
			"comment": schema.StringAttribute{Description: "Rule comment.", Optional: true, Computed: true},
		},
	}
}

func (r *FirewallRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FirewallRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := scopeToPath(plan.Scope.ValueString())
	enableInt := boolToInt(plan.Enabled.ValueBool())

	createReq := &models.FirewallRuleCreateRequest{
		Type:    plan.Type.ValueString(),
		Action:  plan.Action.ValueString(),
		Enable:  &enableInt,
		Macro:   plan.Macro.ValueString(),
		Proto:   plan.Proto.ValueString(),
		Source:  plan.Source.ValueString(),
		Dest:    plan.Dest.ValueString(),
		DPort:   plan.DPort.ValueString(),
		Sport:   plan.Sport.ValueString(),
		IFace:   plan.IFace.ValueString(),
		Log:     plan.Log.ValueString(),
		Comment: plan.Comment.ValueString(),
	}

	tflog.Debug(ctx, "Creating firewall rule", map[string]any{"scope": plan.Scope.ValueString(), "action": createReq.Action})

	if err := r.client.CreateFirewallRule(ctx, pathPrefix, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating firewall rule", err.Error())
		return
	}

	// read back to find out what position was assigned
	rules, err := r.client.GetFirewallRules(ctx, pathPrefix)
	if err != nil {
		resp.Diagnostics.AddError("Error reading firewall rules after create", err.Error())
		return
	}

	// find our rule — take the last match on action+type
	pos := -1
	for _, rule := range rules {
		if rule.Type == plan.Type.ValueString() && rule.Action == plan.Action.ValueString() {
			pos = rule.Pos
		}
	}

	if pos < 0 {
		resp.Diagnostics.AddError("Error finding created rule", "Could not find the newly created firewall rule")
		return
	}

	plan.Pos = types.Int64Value(int64(pos))
	plan.ID = types.StringValue(fmt.Sprintf("%s/%d", plan.Scope.ValueString(), pos))
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FirewallRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FirewallRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FirewallRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := scopeToPath(plan.Scope.ValueString())
	pos := int(plan.Pos.ValueInt64())
	enableInt := boolToInt(plan.Enabled.ValueBool())
	action := plan.Action.ValueString()
	ruleType := plan.Type.ValueString()
	macro := plan.Macro.ValueString()
	proto := plan.Proto.ValueString()
	source := plan.Source.ValueString()
	dest := plan.Dest.ValueString()
	dport := plan.DPort.ValueString()
	sport := plan.Sport.ValueString()
	iface := plan.IFace.ValueString()
	log := plan.Log.ValueString()
	comment := plan.Comment.ValueString()

	updateReq := &models.FirewallRuleUpdateRequest{
		Type:    &ruleType,
		Action:  &action,
		Enable:  &enableInt,
		Macro:   &macro,
		Proto:   &proto,
		Source:  &source,
		Dest:    &dest,
		DPort:   &dport,
		Sport:   &sport,
		IFace:   &iface,
		Log:     &log,
		Comment: &comment,
	}

	if err := r.client.UpdateFirewallRule(ctx, pathPrefix, pos, updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating firewall rule", err.Error())
		return
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FirewallRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := scopeToPath(state.Scope.ValueString())
	if err := r.client.DeleteFirewallRule(ctx, pathPrefix, int(state.Pos.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Error deleting firewall rule", err.Error())
	}
}

func (r *FirewallRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// format: scope/pos  e.g. "cluster/0" or "vm/pve1/100/5"
	parts := strings.Split(req.ID, "/")
	if len(parts) < 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Format: scope/pos (e.g. 'cluster/0', 'node/pve1/3', 'vm/pve1/100/0')")
		return
	}

	posStr := parts[len(parts)-1]
	pos, err := strconv.Atoi(posStr)
	if err != nil {
		resp.Diagnostics.AddError("Invalid position", err.Error())
		return
	}

	scope := strings.Join(parts[:len(parts)-1], "/")
	state := FirewallRuleResourceModel{
		ID:    types.StringValue(req.ID),
		Scope: types.StringValue(scope),
		Pos:   types.Int64Value(int64(pos)),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FirewallRuleResource) readIntoModel(ctx context.Context, model *FirewallRuleResourceModel, diagnostics *diag.Diagnostics) {
	pathPrefix := scopeToPath(model.Scope.ValueString())
	rule, err := r.client.GetFirewallRule(ctx, pathPrefix, int(model.Pos.ValueInt64()))
	if err != nil {
		diagnostics.AddError("Error reading firewall rule", err.Error())
		return
	}

	model.Type = types.StringValue(rule.Type)
	model.Action = types.StringValue(rule.Action)
	if rule.Enable != nil {
		model.Enabled = types.BoolValue(*rule.Enable == 1)
	} else {
		model.Enabled = types.BoolValue(true)
	}
	model.Macro = types.StringValue(rule.Macro)
	model.Proto = types.StringValue(rule.Proto)
	model.Source = types.StringValue(rule.Source)
	model.Dest = types.StringValue(rule.Dest)
	model.DPort = types.StringValue(rule.DPort)
	model.Sport = types.StringValue(rule.Sport)
	model.IFace = types.StringValue(rule.IFace)
	model.Log = types.StringValue(rule.Log)
	model.Comment = types.StringValue(rule.Comment)
}

// scopeToPath turns a scope string into the right API path prefix.
func scopeToPath(scope string) string {
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
