package firewall_security_group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &FirewallSecurityGroupResource{}
var _ resource.ResourceWithConfigure = &FirewallSecurityGroupResource{}
var _ resource.ResourceWithImportState = &FirewallSecurityGroupResource{}

type FirewallSecurityGroupResource struct {
	client *client.Client
}

type FirewallSecurityGroupResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Comment types.String `tfsdk:"comment"`
}

func NewResource() resource.Resource {
	return &FirewallSecurityGroupResource{}
}

func (r *FirewallSecurityGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_firewall_security_group"
}

func (r *FirewallSecurityGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE cluster firewall security group. Rules within the group are managed by a separate resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The security group name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Description: "A comment for the security group.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *FirewallSecurityGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FirewallSecurityGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallSecurityGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &models.FirewallSecurityGroupCreateRequest{
		Group:   plan.Name.ValueString(),
		Comment: plan.Comment.ValueString(),
	}

	tflog.Debug(ctx, "Creating firewall security group", map[string]any{"name": createReq.Group})

	if err := r.client.CreateFirewallSecurityGroup(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating firewall security group", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Name.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallSecurityGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FirewallSecurityGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FirewallSecurityGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FirewallSecurityGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &models.FirewallSecurityGroupUpdateRequest{
		Comment: plan.Comment.ValueString(),
	}

	if err := r.client.UpdateFirewallSecurityGroup(ctx, plan.Name.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating firewall security group", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Name.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallSecurityGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FirewallSecurityGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteFirewallSecurityGroup(ctx, state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting firewall security group", err.Error())
	}
}

func (r *FirewallSecurityGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := FirewallSecurityGroupResourceModel{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FirewallSecurityGroupResource) readIntoModel(ctx context.Context, model *FirewallSecurityGroupResourceModel, diagnostics *diag.Diagnostics) {
	group, err := r.client.GetFirewallSecurityGroup(ctx, model.Name.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok {
			if apiErr.IsNotFound() {
				diagnostics.AddWarning("Firewall security group not found", "The security group no longer exists.")
				return
			}
		}
		diagnostics.AddError("Error reading firewall security group", err.Error())
		return
	}

	model.Comment = types.StringValue(group.Comment)
}
