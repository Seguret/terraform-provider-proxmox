package hardware_mapping_pci

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

var _ resource.Resource = &PCIHardwareMappingResource{}
var _ resource.ResourceWithConfigure = &PCIHardwareMappingResource{}
var _ resource.ResourceWithImportState = &PCIHardwareMappingResource{}

type PCIHardwareMappingResource struct {
	client *client.Client
}

type PCIHardwareMappingResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Comment types.String `tfsdk:"comment"`
	Map     types.List   `tfsdk:"map"`
	MDevs   types.String `tfsdk:"mdevs"`
}

func NewResource() resource.Resource {
	return &PCIHardwareMappingResource{}
}

func (r *PCIHardwareMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_hardware_mapping_pci"
}

func (r *PCIHardwareMappingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE cluster-level PCI hardware mapping.",
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
				Description: "List of per-node PCI device entries. Each entry uses Proxmox format: " +
					"'node=<node>,path=<path>,id=<vendor>:<device>[,iommu-group=<n>]'.",
				Required:    true,
				ElementType: types.StringType,
			},
			"mdevs": schema.StringAttribute{
				Description: "A list of MDev types.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *PCIHardwareMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PCIHardwareMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PCIHardwareMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapEntries := listToStringSlice(ctx, plan.Map, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating PCI hardware mapping", map[string]any{"name": plan.Name.ValueString()})

	if err := r.client.CreatePCIHardwareMapping(ctx, &models.PCIHardwareMappingCreateRequest{
		ID:      plan.Name.ValueString(),
		Comment: plan.Comment.ValueString(),
		Map:     mapEntries,
		MDevs:   plan.MDevs.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error creating PCI hardware mapping", err.Error())
		return
	}

	plan.ID = plan.Name
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PCIHardwareMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PCIHardwareMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PCIHardwareMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PCIHardwareMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapEntries := listToStringSlice(ctx, plan.Map, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdatePCIHardwareMapping(ctx, plan.Name.ValueString(), &models.PCIHardwareMappingUpdateRequest{
		Comment: plan.Comment.ValueString(),
		Map:     mapEntries,
		MDevs:   plan.MDevs.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error updating PCI hardware mapping", err.Error())
		return
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PCIHardwareMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PCIHardwareMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeletePCIHardwareMapping(ctx, state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting PCI hardware mapping", err.Error())
	}
}

func (r *PCIHardwareMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := PCIHardwareMappingResourceModel{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PCIHardwareMappingResource) readIntoModel(ctx context.Context, model *PCIHardwareMappingResourceModel, diagnostics interface{ AddError(string, string) }) {
	mapping, err := r.client.GetPCIHardwareMapping(ctx, model.Name.ValueString())
	if err != nil {
		diagnostics.AddError("Error reading PCI hardware mapping", err.Error())
		return
	}

	model.Comment = types.StringValue(mapping.Comment)
	model.MDevs = types.StringValue(mapping.MDevs)

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
