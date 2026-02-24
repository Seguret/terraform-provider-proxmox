package node_subscription

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
)

var _ resource.Resource = &NodeSubscriptionResource{}
var _ resource.ResourceWithConfigure = &NodeSubscriptionResource{}
var _ resource.ResourceWithImportState = &NodeSubscriptionResource{}

type NodeSubscriptionResource struct {
	client *client.Client
}

type NodeSubscriptionResourceModel struct {
	ID          types.String `tfsdk:"id"`
	NodeName    types.String `tfsdk:"node_name"`
	Key         types.String `tfsdk:"key"`
	Status      types.String `tfsdk:"status"`
	ProductName types.String `tfsdk:"product_name"`
	RegDate     types.String `tfsdk:"reg_date"`
	NextDueDate types.String `tfsdk:"next_due_date"`
}

func NewResource() resource.Resource {
	return &NodeSubscriptionResource{}
}

func (r *NodeSubscriptionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_subscription"
}

func (r *NodeSubscriptionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the subscription key for a Proxmox VE node.",
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
			"key": schema.StringAttribute{
				Description: "The Proxmox VE subscription key.",
				Required:    true,
				Sensitive:   true,
			},
			"status": schema.StringAttribute{
				Description: "The subscription status (e.g. Active, Invalid).",
				Computed:    true,
			},
			"product_name": schema.StringAttribute{
				Description: "The subscription product name.",
				Computed:    true,
			},
			"reg_date": schema.StringAttribute{
				Description: "The subscription registration date.",
				Computed:    true,
			},
			"next_due_date": schema.StringAttribute{
				Description: "The next renewal due date.",
				Computed:    true,
			},
		},
	}
}

func (r *NodeSubscriptionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NodeSubscriptionResource) applyAndRead(ctx context.Context, model *NodeSubscriptionResourceModel) error {
	node := model.NodeName.ValueString()
	key := model.Key.ValueString()

	if err := r.client.SetNodeSubscription(ctx, node, key); err != nil {
		return fmt.Errorf("error setting node subscription key: %w", err)
	}

	return r.readState(ctx, model)
}

func (r *NodeSubscriptionResource) readState(ctx context.Context, model *NodeSubscriptionResourceModel) error {
	sub, err := r.client.GetNodeSubscription(ctx, model.NodeName.ValueString())
	if err != nil {
		return fmt.Errorf("error reading node subscription: %w", err)
	}
	model.Status = types.StringValue(sub.Status)
	model.ProductName = types.StringValue(sub.ProductName)
	model.RegDate = types.StringValue(sub.RegDate)
	model.NextDueDate = types.StringValue(sub.NextDueDate)
	// key is never returned by the GET endpoint — dont overwrite what we have in state
	return nil
}

func (r *NodeSubscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NodeSubscriptionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(plan.NodeName.ValueString())

	tflog.Debug(ctx, "Setting node subscription key", map[string]any{"node": plan.NodeName.ValueString()})

	if err := r.applyAndRead(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error creating node subscription", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NodeSubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NodeSubscriptionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readState(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error reading node subscription", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NodeSubscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NodeSubscriptionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating node subscription key", map[string]any{"node": plan.NodeName.ValueString()})

	if err := r.applyAndRead(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error updating node subscription", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NodeSubscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NodeSubscriptionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting node subscription", map[string]any{"node": state.NodeName.ValueString()})

	if err := r.client.DeleteNodeSubscription(ctx, state.NodeName.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting node subscription", err.Error())
	}
}

func (r *NodeSubscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	node := req.ID
	state := NodeSubscriptionResourceModel{
		ID:       types.StringValue(node),
		NodeName: types.StringValue(node),
		// key isnt in the GET response, user will need to re-apply with it after import
		Key: types.StringValue(""),
	}

	if err := r.readState(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing node subscription", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
