package firewall_alias

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
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &FirewallAliasResource{}
var _ resource.ResourceWithConfigure = &FirewallAliasResource{}
var _ resource.ResourceWithImportState = &FirewallAliasResource{}

type FirewallAliasResource struct {
	client *client.Client
}

type FirewallAliasResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Scope   types.String `tfsdk:"scope"`
	Name    types.String `tfsdk:"name"`
	CIDR    types.String `tfsdk:"cidr"`
	Comment types.String `tfsdk:"comment"`
}

func NewResource() resource.Resource {
	return &FirewallAliasResource{}
}

func (r *FirewallAliasResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_firewall_alias"
}

func (r *FirewallAliasResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE firewall alias (named IP/CIDR).",
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
				Description: "The alias name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cidr": schema.StringAttribute{
				Description: "The IP or CIDR that the alias represents.",
				Required:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Alias description.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *FirewallAliasResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FirewallAliasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallAliasResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := scopeToPath(plan.Scope.ValueString())

	tflog.Debug(ctx, "Creating firewall alias", map[string]any{
		"scope": plan.Scope.ValueString(), "name": plan.Name.ValueString(),
	})

	createReq := &models.FirewallAliasCreateRequest{
		Name:    plan.Name.ValueString(),
		CIDR:    plan.CIDR.ValueString(),
		Comment: plan.Comment.ValueString(),
	}

	if err := r.client.CreateFirewallAlias(ctx, pathPrefix, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating firewall alias", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.Scope.ValueString(), plan.Name.ValueString()))

	alias, err := r.client.GetFirewallAlias(ctx, pathPrefix, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading firewall alias", err.Error())
		return
	}
	plan.CIDR = types.StringValue(alias.CIDR)
	plan.Comment = types.StringValue(alias.Comment)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallAliasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FirewallAliasResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := scopeToPath(state.Scope.ValueString())
	alias, err := r.client.GetFirewallAlias(ctx, pathPrefix, state.Name.ValueString())
	if err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading firewall alias", err.Error())
		return
	}

	state.CIDR = types.StringValue(alias.CIDR)
	state.Comment = types.StringValue(alias.Comment)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FirewallAliasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FirewallAliasResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := scopeToPath(plan.Scope.ValueString())

	updateReq := &models.FirewallAliasUpdateRequest{
		CIDR:    plan.CIDR.ValueString(),
		Comment: plan.Comment.ValueString(),
	}

	if err := r.client.UpdateFirewallAlias(ctx, pathPrefix, plan.Name.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating firewall alias", err.Error())
		return
	}

	alias, err := r.client.GetFirewallAlias(ctx, pathPrefix, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading firewall alias", err.Error())
		return
	}
	plan.CIDR = types.StringValue(alias.CIDR)
	plan.Comment = types.StringValue(alias.Comment)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallAliasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FirewallAliasResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathPrefix := scopeToPath(state.Scope.ValueString())
	if err := r.client.DeleteFirewallAlias(ctx, pathPrefix, state.Name.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting firewall alias", err.Error())
	}
}

func (r *FirewallAliasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	lastSlash := strings.LastIndex(req.ID, "/")
	if lastSlash < 0 {
		resp.Diagnostics.AddError("Invalid import ID", "Format: scope/name")
		return
	}
	scope := req.ID[:lastSlash]
	name := req.ID[lastSlash+1:]

	pathPrefix := scopeToPath(scope)
	alias, err := r.client.GetFirewallAlias(ctx, pathPrefix, name)
	if err != nil {
		resp.Diagnostics.AddError("Error importing firewall alias", err.Error())
		return
	}

	state := FirewallAliasResourceModel{
		ID:      types.StringValue(req.ID),
		Scope:   types.StringValue(scope),
		Name:    types.StringValue(name),
		CIDR:    types.StringValue(alias.CIDR),
		Comment: types.StringValue(alias.Comment),
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
