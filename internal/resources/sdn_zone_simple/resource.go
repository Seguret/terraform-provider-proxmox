package sdn_zone_simple

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

var _ resource.Resource = &SDNZoneSimpleResource{}
var _ resource.ResourceWithConfigure = &SDNZoneSimpleResource{}
var _ resource.ResourceWithImportState = &SDNZoneSimpleResource{}

type SDNZoneSimpleResource struct {
	client *client.Client
}

type SDNZoneSimpleModel struct {
	ID         types.String `tfsdk:"id"`
	Zone       types.String `tfsdk:"zone"`
	Nodes      types.String `tfsdk:"nodes"`
	IPAM       types.String `tfsdk:"ipam"`
	DNS        types.String `tfsdk:"dns"`
	DNSZone    types.String `tfsdk:"dns_zone"`
	ReverseDNS types.String `tfsdk:"reverse_dns"`
}

func NewResource() resource.Resource {
	return &SDNZoneSimpleResource{}
}

func (r *SDNZoneSimpleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_sdn_zone_simple"
}

func (r *SDNZoneSimpleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE SDN zone of type 'simple'.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone": schema.StringAttribute{
				Description: "The SDN zone identifier.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"nodes": schema.StringAttribute{
				Description: "Comma-separated list of nodes where the zone is deployed.",
				Optional:    true,
			},
			"ipam": schema.StringAttribute{
				Description: "IPAM plugin name.",
				Optional:    true,
			},
			"dns": schema.StringAttribute{
				Description: "DNS plugin name.",
				Optional:    true,
			},
			"dns_zone": schema.StringAttribute{
				Description: "DNS domain.",
				Optional:    true,
			},
			"reverse_dns": schema.StringAttribute{
				Description: "Reverse DNS plugin name.",
				Optional:    true,
			},
		},
	}
}

func (r *SDNZoneSimpleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SDNZoneSimpleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SDNZoneSimpleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating SDN simple zone", map[string]any{"zone": plan.Zone.ValueString()})

	if err := r.client.CreateSDNZone(ctx, &models.SDNZoneCreateRequest{
		Zone:       plan.Zone.ValueString(),
		Type:       "simple",
		Nodes:      plan.Nodes.ValueString(),
		IPAM:       plan.IPAM.ValueString(),
		DNS:        plan.DNS.ValueString(),
		DNSZone:    plan.DNSZone.ValueString(),
		ReverseDNS: plan.ReverseDNS.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error creating SDN simple zone", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	plan.ID = plan.Zone
	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading SDN simple zone after create", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNZoneSimpleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SDNZoneSimpleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading SDN simple zone", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNZoneSimpleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SDNZoneSimpleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateSDNZone(ctx, plan.Zone.ValueString(), &models.SDNZoneUpdateRequest{
		Nodes:      plan.Nodes.ValueString(),
		IPAM:       plan.IPAM.ValueString(),
		DNS:        plan.DNS.ValueString(),
		DNSZone:    plan.DNSZone.ValueString(),
		ReverseDNS: plan.ReverseDNS.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error updating SDN simple zone", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading SDN simple zone after update", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNZoneSimpleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SDNZoneSimpleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSDNZone(ctx, state.Zone.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting SDN simple zone", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)
}

func (r *SDNZoneSimpleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := SDNZoneSimpleModel{
		ID:   types.StringValue(req.ID),
		Zone: types.StringValue(req.ID),
	}
	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing SDN simple zone", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNZoneSimpleResource) readIntoModel(ctx context.Context, model *SDNZoneSimpleModel) error {
	zone, err := r.client.GetSDNZone(ctx, model.Zone.ValueString())
	if err != nil {
		return err
	}
	if zone.Type != "simple" {
		return &client.APIError{StatusCode: 404, Status: "404 Not Found"}
	}

	model.Nodes = types.StringValue(zone.Nodes)
	model.IPAM = types.StringValue(zone.IPAM)
	model.DNS = types.StringValue(zone.DNS)
	model.DNSZone = types.StringValue(zone.DNSZone)
	model.ReverseDNS = types.StringValue(zone.ReverseDNS)
	return nil
}
