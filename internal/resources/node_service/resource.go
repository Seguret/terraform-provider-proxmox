package node_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ resource.Resource = &NodeServiceResource{}
var _ resource.ResourceWithConfigure = &NodeServiceResource{}
var _ resource.ResourceWithImportState = &NodeServiceResource{}

// NodeServiceResource manages the run state of a systemd service on a Proxmox node.
type NodeServiceResource struct {
	client *client.Client
}

type NodeServiceResourceModel struct {
	ID       types.String `tfsdk:"id"`
	NodeName types.String `tfsdk:"node_name"`
	Service  types.String `tfsdk:"service"`
	// State is the desired run state: "started", "stopped", or "restarted".
	State types.String `tfsdk:"state"`
}

func NewResource() resource.Resource {
	return &NodeServiceResource{}
}

func (r *NodeServiceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_service"
}

func (r *NodeServiceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the desired run state of a system service on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The resource identifier (node_name:service).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox VE node.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service": schema.StringAttribute{
				Description: "The service name (e.g. 'pveproxy', 'pvedaemon').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Description: "The desired service state. Valid values: 'started', 'stopped', 'restarted'.",
				Required:    true,
			},
		},
	}
}

func (r *NodeServiceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// applyState puts the service into the requested state and waits for the task to finish.
func (r *NodeServiceResource) applyState(ctx context.Context, node, service, desiredState string) error {
	var upid string
	var err error

	switch desiredState {
	case "started":
		upid, err = r.client.StartNodeService(ctx, node, service)
	case "stopped":
		upid, err = r.client.StopNodeService(ctx, node, service)
	case "restarted":
		upid, err = r.client.RestartNodeService(ctx, node, service)
	default:
		return fmt.Errorf("unknown service state %q: must be 'started', 'stopped', or 'restarted'", desiredState)
	}

	if err != nil {
		return fmt.Errorf("error applying service state %q: %w", desiredState, err)
	}

	if upid != "" {
		if err := r.client.WaitForUPID(ctx, upid); err != nil {
			return fmt.Errorf("error waiting for service task: %w", err)
		}
	}

	return nil
}

// readIntoModel reads the current service state from the API and populates the model.
// func (r *NodeServiceResource) readIntoMcfodel(ctx context.Context, model *NodeServiceResourceModel) error {
// 	svc, err := r.client.GetNodeService(ctx, model.NodeName.ValueString(), model.Service.ValueString())
// 	if err != nil {
// 		return fmt.Errorf("error reading service state: %w", err)
// 	}
// 	// record the live state from proxmox so terraform can detect drift
// 	model.State = types.StringValue(svc.State)
// 	return nil
// }

func (r *NodeServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NodeServiceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	service := plan.Service.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", node, service))

	tflog.Debug(ctx, "Creating node service state", map[string]any{
		"node":    node,
		"service": service,
		"state":   plan.State.ValueString(),
	})

	if err := r.applyState(ctx, node, service, plan.State.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error managing node service", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NodeServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NodeServiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.GetNodeService(ctx, state.NodeName.ValueString(), state.Service.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading node service", err.Error())
		return
	}

	// dont overwrite model.State with the live state - this field is intentional/desired state,
	// not an observation. drift is only visible when the user changes their config.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NodeServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NodeServiceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	service := plan.Service.ValueString()

	tflog.Debug(ctx, "Updating node service state", map[string]any{
		"node":    node,
		"service": service,
		"state":   plan.State.ValueString(),
	})

	if err := r.applyState(ctx, node, service, plan.State.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error updating node service", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NodeServiceResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// services cant be deleted, just removed from state
}

func (r *NodeServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// expected format is "node_name:service"
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Format must be '<node_name>:<service>' (e.g. 'pve:pveproxy')",
		)
		return
	}

	node := parts[0]
	service := parts[1]

	svc, err := r.client.GetNodeService(ctx, node, service)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			resp.Diagnostics.AddError("Node service not found", fmt.Sprintf("Service '%s' not found on node '%s'", service, node))
			return
		}
		resp.Diagnostics.AddError("Error importing node service", err.Error())
		return
	}

	state := NodeServiceResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(node),
		Service:  types.StringValue(service),
		State:    types.StringValue(svc.State),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
