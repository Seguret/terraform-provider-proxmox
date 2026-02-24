package firewall_security_group_rule

import (
	"context"
	"fmt"
	"strconv"

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

var _ resource.Resource = &FirewallSecurityGroupRuleResource{}
var _ resource.ResourceWithConfigure = &FirewallSecurityGroupRuleResource{}
var _ resource.ResourceWithImportState = &FirewallSecurityGroupRuleResource{}

type FirewallSecurityGroupRuleResource struct {
	client *client.Client
}

type FirewallSecurityGroupRuleResourceModel struct {
	ID            types.String `tfsdk:"id"`
	SecurityGroup types.String `tfsdk:"security_group"`
	Pos           types.Int64  `tfsdk:"pos"`
	Type          types.String `tfsdk:"type"`
	Action        types.String `tfsdk:"action"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Macro         types.String `tfsdk:"macro"`
	Proto         types.String `tfsdk:"proto"`
	Source        types.String `tfsdk:"source"`
	Dest          types.String `tfsdk:"dest"`
	DPort         types.String `tfsdk:"dport"`
	Sport         types.String `tfsdk:"sport"`
	IFace         types.String `tfsdk:"iface"`
	Log           types.String `tfsdk:"log"`
	Comment       types.String `tfsdk:"comment"`
}

func NewResource() resource.Resource {
	return &FirewallSecurityGroupRuleResource{}
}

func (r *FirewallSecurityGroupRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_firewall_security_group_rule"
}

func (r *FirewallSecurityGroupRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a rule inside a Proxmox VE cluster firewall security group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"security_group": schema.StringAttribute{
				Description: "The security group name this rule belongs to.",
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
			"log":     schema.StringAttribute{Description: "Log level.", Optional: true, Computed: true},
			"comment": schema.StringAttribute{Description: "Rule comment.", Optional: true, Computed: true},
		},
	}
}

func (r *FirewallSecurityGroupRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FirewallSecurityGroupRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallSecurityGroupRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := sgPathPrefix(plan.SecurityGroup.ValueString())
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

	tflog.Debug(ctx, "Creating firewall security group rule", map[string]any{
		"security_group": plan.SecurityGroup.ValueString(),
		"action":         createReq.Action,
	})

	if err := r.client.CreateFirewallRule(ctx, pathPrefix, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating firewall security group rule", err.Error())
		return
	}

	// read back to find out what position was assigned
	rules, err := r.client.GetFirewallRules(ctx, pathPrefix)
	if err != nil {
		resp.Diagnostics.AddError("Error reading firewall rules after create", err.Error())
		return
	}

	pos := -1
	for _, rule := range rules {
		if rule.Type == plan.Type.ValueString() && rule.Action == plan.Action.ValueString() {
			pos = rule.Pos
		}
	}

	if pos < 0 {
		resp.Diagnostics.AddError("Error finding created rule", "Could not find the newly created security group rule")
		return
	}

	plan.Pos = types.Int64Value(int64(pos))
	plan.ID = types.StringValue(fmt.Sprintf("%s/%d", plan.SecurityGroup.ValueString(), pos))
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallSecurityGroupRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FirewallSecurityGroupRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FirewallSecurityGroupRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FirewallSecurityGroupRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := sgPathPrefix(plan.SecurityGroup.ValueString())
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
		resp.Diagnostics.AddError("Error updating firewall security group rule", err.Error())
		return
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallSecurityGroupRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FirewallSecurityGroupRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := sgPathPrefix(state.SecurityGroup.ValueString())
	if err := r.client.DeleteFirewallRule(ctx, pathPrefix, int(state.Pos.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Error deleting firewall security group rule", err.Error())
	}
}

func (r *FirewallSecurityGroupRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// format: "<group>/<pos>"  e.g. "my-sg/3"
	// group names dont have "/" so the last slash always splits group from pos
	lastSlash := -1
	for i := len(req.ID) - 1; i >= 0; i-- {
		if req.ID[i] == '/' {
			lastSlash = i
			break
		}
	}
	if lastSlash < 0 {
		resp.Diagnostics.AddError("Invalid import ID", "Format: <security_group>/<pos>")
		return
	}

	groupName := req.ID[:lastSlash]
	posStr := req.ID[lastSlash+1:]
	pos, err := strconv.Atoi(posStr)
	if err != nil {
		resp.Diagnostics.AddError("Invalid position", err.Error())
		return
	}

	state := FirewallSecurityGroupRuleResourceModel{
		ID:            types.StringValue(req.ID),
		SecurityGroup: types.StringValue(groupName),
		Pos:           types.Int64Value(int64(pos)),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FirewallSecurityGroupRuleResource) readIntoModel(ctx context.Context, model *FirewallSecurityGroupRuleResourceModel, diagnostics *diag.Diagnostics) {
	pathPrefix := sgPathPrefix(model.SecurityGroup.ValueString())
	rule, err := r.client.GetFirewallRule(ctx, pathPrefix, int(model.Pos.ValueInt64()))
	if err != nil {
		diagnostics.AddError("Error reading firewall security group rule", err.Error())
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

// sgPathPrefix builds the API path for rules inside a security group.
func sgPathPrefix(group string) string {
	return fmt.Sprintf("/cluster/firewall/groups/%s", group)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
