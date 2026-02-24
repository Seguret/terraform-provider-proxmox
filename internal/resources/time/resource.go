package time

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ resource.Resource = &TimeResource{}
var _ resource.ResourceWithConfigure = &TimeResource{}
var _ resource.ResourceWithImportState = &TimeResource{}

type TimeResource struct {
	client *client.Client
}

type TimeResourceModel struct {
	ID       types.String `tfsdk:"id"`
	NodeName types.String `tfsdk:"node_name"`
	Timezone types.String `tfsdk:"timezone"`
}

func NewResource() resource.Resource {
	return &TimeResource{}
}

func (r *TimeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_time"
}

func (r *TimeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the timezone configuration of a Proxmox VE node.",
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
			"timezone": schema.StringAttribute{
				Description: "The timezone (e.g. 'Europe/Rome', 'UTC').",
				Required:    true,
			},
		},
	}
}

func (r *TimeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TimeResource) readState(ctx context.Context, model *TimeResourceModel) error {
	data, err := r.client.GetNodeTime(ctx, model.NodeName.ValueString())
	if err != nil {
		return fmt.Errorf("error reading node time: %w", err)
	}
	if tz, ok := data["timezone"].(string); ok {
		model.Timezone = types.StringValue(tz)
	}
	return nil
}

func (r *TimeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TimeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(plan.NodeName.ValueString())

	if err := r.client.SetNodeTimezone(ctx, plan.NodeName.ValueString(), plan.Timezone.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error setting timezone", err.Error())
		return
	}

	if err := r.readState(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading timezone after create", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *TimeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TimeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readState(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error reading timezone", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TimeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TimeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.SetNodeTimezone(ctx, plan.NodeName.ValueString(), plan.Timezone.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error updating timezone", err.Error())
		return
	}

	if err := r.readState(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading timezone after update", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *TimeResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// timezone config cant be deleted, just remove it from state
}

func (r *TimeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := TimeResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(req.ID),
	}

	if err := r.readState(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing timezone", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
