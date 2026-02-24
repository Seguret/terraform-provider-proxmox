package hardware_mapping_usb

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

var _ resource.Resource = &USBHardwareMappingResource{}
var _ resource.ResourceWithConfigure = &USBHardwareMappingResource{}
var _ resource.ResourceWithImportState = &USBHardwareMappingResource{}

type USBHardwareMappingResource struct {
	client *client.Client
}

type USBHardwareMappingResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Comment types.String `tfsdk:"comment"`
	Map     types.List   `tfsdk:"map"`
}

func NewResource() resource.Resource {
	return &USBHardwareMappingResource{}
}

func (r *USBHardwareMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_hardware_mapping_usb"
}

func (r *USBHardwareMappingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE cluster-level USB hardware mapping.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The hardware mapping name (identifier).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Description: "A human-readable description.",
				Optional:    true,
				Computed:    true,
			},
			"map": schema.ListAttribute{
				Description: "List of per-node USB device entries. Each entry uses Proxmox format: " +
					"'node=<node>,id=<vendor>:<product>[,path=<path>]'.",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *USBHardwareMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *USBHardwareMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan USBHardwareMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapEntries := listToStringSlice(ctx, plan.Map, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating USB hardware mapping", map[string]any{"name": plan.Name.ValueString()})

	if err := r.client.CreateUSBHardwareMapping(ctx, &models.USBHardwareMappingCreateRequest{
		ID:      plan.Name.ValueString(),
		Comment: plan.Comment.ValueString(),
		Map:     mapEntries,
	}); err != nil {
		resp.Diagnostics.AddError("Error creating USB hardware mapping", err.Error())
		return
	}

	plan.ID = plan.Name
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *USBHardwareMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state USBHardwareMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *USBHardwareMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan USBHardwareMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapEntries := listToStringSlice(ctx, plan.Map, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateUSBHardwareMapping(ctx, plan.Name.ValueString(), &models.USBHardwareMappingUpdateRequest{
		Comment: plan.Comment.ValueString(),
		Map:     mapEntries,
	}); err != nil {
		resp.Diagnostics.AddError("Error updating USB hardware mapping", err.Error())
		return
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *USBHardwareMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state USBHardwareMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteUSBHardwareMapping(ctx, state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting USB hardware mapping", err.Error())
	}
}

func (r *USBHardwareMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := USBHardwareMappingResourceModel{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *USBHardwareMappingResource) readIntoModel(ctx context.Context, model *USBHardwareMappingResourceModel, diagnostics interface{ AddError(string, string) }) {
	mapping, err := r.client.GetUSBHardwareMapping(ctx, model.Name.ValueString())
	if err != nil {
		diagnostics.AddError("Error reading USB hardware mapping", err.Error())
		return
	}

	model.Comment = types.StringValue(mapping.Comment)

	mapVals := make([]types.String, len(mapping.Map))
	for i, e := range mapping.Map {
		mapVals[i] = types.StringValue(e)
	}
	listVal, diags := types.ListValueFrom(ctx, types.StringType, mapVals)
	if diags.HasError() {
		diagnostics.AddError("Error building map list", "Failed to convert map entries to list")
		return
	}
	model.Map = listVal
}

func listToStringSlice(ctx context.Context, list types.List, diagnostics interface{ AddError(string, string) }) []string {
	var elems []types.String
	if diags := list.ElementsAs(ctx, &elems, false); diags.HasError() {
		diagnostics.AddError("Error reading map list", "Failed to convert list elements to strings")
		return nil
	}
	result := make([]string, len(elems))
	for i, e := range elems {
		result[i] = e.ValueString()
	}
	return result
}
