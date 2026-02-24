package sdn_controller

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

var _ resource.Resource = &SDNControllerResource{}
var _ resource.ResourceWithConfigure = &SDNControllerResource{}
var _ resource.ResourceWithImportState = &SDNControllerResource{}

type SDNControllerResource struct {
	client *client.Client
}

type SDNControllerResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Controller types.String `tfsdk:"controller"`
	Type       types.String `tfsdk:"type"`
	ASN        types.Int64  `tfsdk:"asn"`
	EBGP       types.Bool   `tfsdk:"ebgp"`
	Node       types.String `tfsdk:"node"`
	Peers      types.String `tfsdk:"peers"`
}

func NewResource() resource.Resource {
	return &SDNControllerResource{}
}

func (r *SDNControllerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_sdn_controller"
}

func (r *SDNControllerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE SDN controller.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"controller": schema.StringAttribute{
				Description: "The SDN controller name (identifier).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The SDN controller type (evpn, bgp, isis, or simple).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"asn": schema.Int64Attribute{
				Description: "The BGP/EVPN autonomous system number.",
				Optional:    true,
				Computed:    true,
			},
			"ebgp": schema.BoolAttribute{
				Description: "Enable eBGP (inter-AS routing). Only applicable for EVPN controllers.",
				Optional:    true,
				Computed:    true,
			},
			"node": schema.StringAttribute{
				Description: "The cluster node name (for ISIS controllers).",
				Optional:    true,
				Computed:    true,
			},
			"peers": schema.StringAttribute{
				Description: "Comma-separated list of BGP peer addresses.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *SDNControllerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SDNControllerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SDNControllerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating SDN controller", map[string]any{
		"controller": plan.Controller.ValueString(),
		"type":       plan.Type.ValueString(),
	})

	if err := r.client.CreateSDNController(ctx, &models.SDNControllerCreateRequest{
		Controller: plan.Controller.ValueString(),
		Type:       plan.Type.ValueString(),
		ASN:        int(plan.ASN.ValueInt64()),
		EBGP:       plan.EBGP.ValueBool(),
		Node:       plan.Node.ValueString(),
		Peers:      plan.Peers.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error creating SDN controller", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	plan.ID = plan.Controller
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNControllerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SDNControllerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if isNotFound := r.readIntoModelWithNotFound(ctx, &state, &resp.Diagnostics); isNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNControllerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SDNControllerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateSDNController(ctx, plan.Controller.ValueString(), &models.SDNControllerUpdateRequest{
		ASN:   int(plan.ASN.ValueInt64()),
		EBGP:  plan.EBGP.ValueBool(),
		Node:  plan.Node.ValueString(),
		Peers: plan.Peers.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error updating SDN controller", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNControllerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SDNControllerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSDNController(ctx, state.Controller.ValueString()); err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return
		}
		resp.Diagnostics.AddError("Error deleting SDN controller", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)
}

func (r *SDNControllerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := SDNControllerResourceModel{
		ID:         types.StringValue(req.ID),
		Controller: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// readIntoModelWithNotFound fetches controller state; returns true if the resource doesnt exist.
func (r *SDNControllerResource) readIntoModelWithNotFound(ctx context.Context, model *SDNControllerResourceModel, diags interface{ AddError(string, string) }) bool {
	ctrl, err := r.client.GetSDNController(ctx, model.Controller.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return true
		}
		diags.AddError("Error reading SDN controller", err.Error())
		return false
	}
	populateModel(model, ctrl)
	return false
}

// readIntoModel fetches controller state and adds diagnostic errors on failure.
func (r *SDNControllerResource) readIntoModel(ctx context.Context, model *SDNControllerResourceModel, diags interface{ AddError(string, string) }) {
	ctrl, err := r.client.GetSDNController(ctx, model.Controller.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diags.AddError("SDN controller not found", "The SDN controller no longer exists.")
			return
		}
		diags.AddError("Error reading SDN controller", err.Error())
		return
	}
	populateModel(model, ctrl)
}

func populateModel(model *SDNControllerResourceModel, ctrl *models.SDNController) {
	model.Type = types.StringValue(ctrl.Type)
	model.ASN = types.Int64Value(int64(ctrl.ASN))
	model.EBGP = types.BoolValue(ctrl.EBGP)
	model.Node = types.StringValue(ctrl.Node)
	model.Peers = types.StringValue(ctrl.Peers)
}
