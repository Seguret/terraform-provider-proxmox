package hosts

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &HostsResource{}
var _ resource.ResourceWithConfigure = &HostsResource{}
var _ resource.ResourceWithImportState = &HostsResource{}

type HostsResource struct {
	client *client.Client
}

type HostsResourceModel struct {
	ID       types.String `tfsdk:"id"`
	NodeName types.String `tfsdk:"node_name"`
	Entries  types.String `tfsdk:"entries"`
	Digest   types.String `tfsdk:"digest"`
}

func NewResource() resource.Resource {
	return &HostsResource{}
}

func (r *HostsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_hosts"
}

func (r *HostsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the /etc/hosts file of a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The node name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"entries": schema.StringAttribute{
				Description: "The full content of /etc/hosts.",
				Required:    true,
			},
			"digest": schema.StringAttribute{
				Description: "The digest of the hosts file (for conflict detection).",
				Computed:    true,
			},
		},
	}
}

func (r *HostsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *HostsResource) applyAndRead(ctx context.Context, model *HostsResourceModel) error {
	node := model.NodeName.ValueString()

	updateReq := &models.NodeHostsUpdateRequest{
		Data:   model.Entries.ValueString(),
		Digest: model.Digest.ValueString(),
	}

	if err := r.client.UpdateNodeHosts(ctx, node, updateReq); err != nil {
		return fmt.Errorf("error writing hosts: %w", err)
	}

	hosts, err := r.client.GetNodeHosts(ctx, node)
	if err != nil {
		return fmt.Errorf("error reading hosts: %w", err)
	}

	model.Entries = types.StringValue(hosts.Data)
	model.Digest = types.StringValue(hosts.Digest)
	return nil
}

func (r *HostsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan HostsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(plan.NodeName.ValueString())

	if err := r.applyAndRead(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error creating hosts config", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *HostsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state HostsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hosts, err := r.client.GetNodeHosts(ctx, state.NodeName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading hosts config", err.Error())
		return
	}

	state.Entries = types.StringValue(hosts.Data)
	state.Digest = types.StringValue(hosts.Digest)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *HostsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan HostsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.applyAndRead(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error updating hosts config", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *HostsResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// /etc/hosts cant be deleted — just remove from state
}

func (r *HostsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	node := req.ID
	state := HostsResourceModel{
		ID:       types.StringValue(node),
		NodeName: types.StringValue(node),
	}

	hosts, err := r.client.GetNodeHosts(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error importing hosts config", err.Error())
		return
	}

	state.Entries = types.StringValue(hosts.Data)
	state.Digest = types.StringValue(hosts.Digest)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
