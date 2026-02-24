package firewall_ipset

import (
	"context"
	"fmt"
	"strconv"
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

var _ resource.Resource = &FirewallIPSetResource{}
var _ resource.ResourceWithConfigure = &FirewallIPSetResource{}
var _ resource.ResourceWithImportState = &FirewallIPSetResource{}

type FirewallIPSetResource struct {
	client *client.Client
}

// CIDREntryModel is a single CIDR entry inside an IP set.
type CIDREntryModel struct {
	CIDR    types.String `tfsdk:"cidr"`
	Comment types.String `tfsdk:"comment"`
	NoMatch types.Bool   `tfsdk:"no_match"`
}

type FirewallIPSetResourceModel struct {
	ID      types.String     `tfsdk:"id"`
	Scope   types.String     `tfsdk:"scope"`
	Name    types.String     `tfsdk:"name"`
	Comment types.String     `tfsdk:"comment"`
	CIDRs   []CIDREntryModel `tfsdk:"cidrs"`
}

func NewResource() resource.Resource {
	return &FirewallIPSetResource{}
}

func (r *FirewallIPSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_firewall_ipset"
}

func (r *FirewallIPSetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE firewall IP set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scope": schema.StringAttribute{
				Description: "The firewall scope: 'cluster', 'node/<node>', 'vm/<node>/<vmid>', 'ct/<node>/<vmid>'.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The IP set name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Description: "IP set description.",
				Optional:    true,
				Computed:    true,
			},
			"cidrs": schema.ListNestedAttribute{
				Description: "CIDR entries in the IP set.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cidr": schema.StringAttribute{
							Description: "The CIDR or IP address.",
							Required:    true,
						},
						"comment": schema.StringAttribute{
							Description: "Entry comment.",
							Optional:    true,
							Computed:    true,
						},
						"no_match": schema.BoolAttribute{
							Description: "Whether to negate the match (prefix with '!').",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
					},
				},
			},
		},
	}
}

func (r *FirewallIPSetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FirewallIPSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallIPSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := scopeToPath(plan.Scope.ValueString())

	tflog.Debug(ctx, "Creating firewall IP set", map[string]any{"scope": plan.Scope.ValueString(), "name": plan.Name.ValueString()})

	if err := r.client.CreateFirewallIPSet(ctx, pathPrefix, plan.Name.ValueString(), plan.Comment.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error creating firewall IP set", err.Error())
		return
	}

	// add each CIDR entry to the set
	for _, entry := range plan.CIDRs {
		nomatch := 0
		if entry.NoMatch.ValueBool() {
			nomatch = 1
		}
		e := &models.FirewallIPSetEntry{
			CIDR:    entry.CIDR.ValueString(),
			Comment: entry.Comment.ValueString(),
			NoMatch: &nomatch,
		}
		if err := r.client.CreateFirewallIPSetEntry(ctx, pathPrefix, plan.Name.ValueString(), e); err != nil {
			resp.Diagnostics.AddError("Error adding CIDR to IP set", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.Scope.ValueString(), plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallIPSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FirewallIPSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := scopeToPath(state.Scope.ValueString())
	entries, err := r.client.GetFirewallIPSetEntries(ctx, pathPrefix, state.Name.ValueString())
	if err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading firewall IP set", err.Error())
		return
	}

	state.CIDRs = nil
	for _, e := range entries {
		nomatch := e.NoMatch != nil && *e.NoMatch == 1
		state.CIDRs = append(state.CIDRs, CIDREntryModel{
			CIDR:    types.StringValue(e.CIDR),
			Comment: types.StringValue(e.Comment),
			NoMatch: types.BoolValue(nomatch),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FirewallIPSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FirewallIPSetResourceModel
	var state FirewallIPSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := scopeToPath(plan.Scope.ValueString())

	// clear all entries then re-add the desired set
	existingEntries, err := r.client.GetFirewallIPSetEntries(ctx, pathPrefix, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading existing IP set entries", err.Error())
		return
	}
	for _, e := range existingEntries {
		if err := r.client.DeleteFirewallIPSetEntry(ctx, pathPrefix, plan.Name.ValueString(), e.CIDR); err != nil {
			resp.Diagnostics.AddError("Error removing IP set entry", err.Error())
			return
		}
	}

	for _, entry := range plan.CIDRs {
		nomatch := 0
		if entry.NoMatch.ValueBool() {
			nomatch = 1
		}
		e := &models.FirewallIPSetEntry{
			CIDR:    entry.CIDR.ValueString(),
			Comment: entry.Comment.ValueString(),
			NoMatch: &nomatch,
		}
		if err := r.client.CreateFirewallIPSetEntry(ctx, pathPrefix, plan.Name.ValueString(), e); err != nil {
			resp.Diagnostics.AddError("Error adding CIDR to IP set", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallIPSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FirewallIPSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := scopeToPath(state.Scope.ValueString())

	// entries must be removed before the set itself can be deleted
	entries, err := r.client.GetFirewallIPSetEntries(ctx, pathPrefix, state.Name.ValueString())
	if err == nil {
		for _, e := range entries {
			_ = r.client.DeleteFirewallIPSetEntry(ctx, pathPrefix, state.Name.ValueString(), e.CIDR)
		}
	}

	if err := r.client.DeleteFirewallIPSet(ctx, pathPrefix, state.Name.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting firewall IP set", err.Error())
	}
}

func (r *FirewallIPSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// format: scope/name (e.g. "cluster/myipset" or "vm/node1/100/myipset")
	lastSlash := strings.LastIndex(req.ID, "/")
	if lastSlash < 0 {
		resp.Diagnostics.AddError("Invalid import ID", "Format: scope/name")
		return
	}
	scope := req.ID[:lastSlash]
	name := req.ID[lastSlash+1:]

	state := FirewallIPSetResourceModel{
		ID:    types.StringValue(req.ID),
		Scope: types.StringValue(scope),
		Name:  types.StringValue(name),
	}

	pathPrefix := scopeToPath(scope)
	entries, err := r.client.GetFirewallIPSetEntries(ctx, pathPrefix, name)
	if err != nil {
		resp.Diagnostics.AddError("Error importing firewall IP set", err.Error())
		return
	}

	for _, e := range entries {
		nomatch := e.NoMatch != nil && *e.NoMatch == 1
		state.CIDRs = append(state.CIDRs, CIDREntryModel{
			CIDR:    types.StringValue(e.CIDR),
			Comment: types.StringValue(e.Comment),
			NoMatch: types.BoolValue(nomatch),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}


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

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
