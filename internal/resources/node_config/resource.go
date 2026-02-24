package node_config

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &NodeConfigResource{}
var _ resource.ResourceWithConfigure = &NodeConfigResource{}
var _ resource.ResourceWithImportState = &NodeConfigResource{}

type NodeConfigResource struct {
	client *client.Client
}

type NodeConfigResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	NodeName            types.String `tfsdk:"node_name"`
	Description         types.String `tfsdk:"description"`
	Wakeonlan           types.String `tfsdk:"wakeonlan"`
	StartallOnbootDelay types.Int64  `tfsdk:"startall_onboot_delay"`
}

func NewResource() resource.Resource {
	return &NodeConfigResource{}
}

func (r *NodeConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_config"
}

func (r *NodeConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the configuration of a Proxmox VE node.",
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
			"description": schema.StringAttribute{
				Description: "A description for the node.",
				Optional:    true,
				Computed:    true,
			},
			"wakeonlan": schema.StringAttribute{
				Description: "The Wake-on-LAN MAC address for this node.",
				Optional:    true,
				Computed:    true,
			},
			"startall_onboot_delay": schema.Int64Attribute{
				Description: "Initial delay in seconds, before starting all the Virtual Guests with on-boot enabled.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *NodeConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NodeConfigResource) applyAndRead(ctx context.Context, model *NodeConfigResourceModel) error {
	node := model.NodeName.ValueString()

	updateReq := &models.NodeConfigUpdateRequest{
		Description:         model.Description.ValueString(),
		Wakeonlan:           model.Wakeonlan.ValueString(),
		StartallOnbootDelay: int(model.StartallOnbootDelay.ValueInt64()),
	}

	if err := r.client.UpdateNodeConfig(ctx, node, updateReq); err != nil {
		return fmt.Errorf("error updating node config: %w", err)
	}

	return r.readState(ctx, model)
}

func (r *NodeConfigResource) readState(ctx context.Context, model *NodeConfigResourceModel) error {
	cfg, err := r.client.GetNodeConfig(ctx, model.NodeName.ValueString())
	if err != nil {
		return fmt.Errorf("error reading node config: %w", err)
	}
	model.Description = types.StringValue(cfg.Description)
	model.Wakeonlan = types.StringValue(cfg.Wakeonlan)
	model.StartallOnbootDelay = types.Int64Value(int64(cfg.StartallOnbootDelay))
	return nil
}

func (r *NodeConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NodeConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(plan.NodeName.ValueString())

	tflog.Debug(ctx, "Creating node config", map[string]any{"node": plan.NodeName.ValueString()})

	if err := r.applyAndRead(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error creating node config", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NodeConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NodeConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readState(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error reading node config", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NodeConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NodeConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating node config", map[string]any{"node": plan.NodeName.ValueString()})

	if err := r.applyAndRead(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error updating node config", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NodeConfigResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// node config cant be deleted from proxmox — just remove from state
}

func (r *NodeConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	node := req.ID
	state := NodeConfigResourceModel{
		ID:       types.StringValue(node),
		NodeName: types.StringValue(node),
	}

	if err := r.readState(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing node config", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
