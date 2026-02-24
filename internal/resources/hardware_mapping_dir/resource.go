package hardware_mapping_dir

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

var _ resource.Resource = &DirHardwareMappingResource{}
var _ resource.ResourceWithConfigure = &DirHardwareMappingResource{}
var _ resource.ResourceWithImportState = &DirHardwareMappingResource{}

// mapEntryAttrTypes holds the attr types for a single map entry object.
var mapEntryAttrTypes = map[string]attr.Type{
	"node": types.StringType,
	"path": types.StringType,
}

type DirHardwareMappingResource struct {
	client *client.Client
}

type DirHardwareMappingResourceModel struct {
	ID        types.String `tfsdk:"id"`
	MappingID types.String `tfsdk:"mapping_id"`
	Comment   types.String `tfsdk:"comment"`
	Map       types.List   `tfsdk:"map"`
}

// mapEntryModel is one entry in the map list.
type mapEntryModel struct {
	Node types.String `tfsdk:"node"`
	Path types.String `tfsdk:"path"`
}

func NewResource() resource.Resource {
	return &DirHardwareMappingResource{}
}

func (r *DirHardwareMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_hardware_mapping_dir"
}

func (r *DirHardwareMappingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE cluster-level directory hardware mapping.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"mapping_id": schema.StringAttribute{
				Description:   "The hardware mapping identifier.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"comment": schema.StringAttribute{
				Description: "A human-readable description.",
				Optional:    true,
				Computed:    true,
			},
			"map": schema.ListNestedAttribute{
				Description: "List of per-node directory path entries.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"node": schema.StringAttribute{
							Description: "The Proxmox VE node name.",
							Required:    true,
						},
						"path": schema.StringAttribute{
							Description: "The directory path on the node.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *DirHardwareMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DirHardwareMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DirHardwareMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapStrings := modelMapToStrings(ctx, plan.Map, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating directory hardware mapping", map[string]any{"id": plan.MappingID.ValueString()})

	if err := r.client.CreateDirHardwareMapping(ctx, &models.DirHardwareMappingCreateRequest{
		ID:      plan.MappingID.ValueString(),
		Comment: plan.Comment.ValueString(),
		Map:     mapStrings,
	}); err != nil {
		resp.Diagnostics.AddError("Error creating directory hardware mapping", err.Error())
		return
	}

	plan.ID = plan.MappingID
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DirHardwareMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DirHardwareMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DirHardwareMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DirHardwareMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapStrings := modelMapToStrings(ctx, plan.Map, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating directory hardware mapping", map[string]any{"id": plan.MappingID.ValueString()})

	if err := r.client.UpdateDirHardwareMapping(ctx, plan.MappingID.ValueString(), &models.DirHardwareMappingUpdateRequest{
		Comment: plan.Comment.ValueString(),
		Map:     mapStrings,
	}); err != nil {
		resp.Diagnostics.AddError("Error updating directory hardware mapping", err.Error())
		return
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DirHardwareMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DirHardwareMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting directory hardware mapping", map[string]any{"id": state.MappingID.ValueString()})

	if err := r.client.DeleteDirHardwareMapping(ctx, state.MappingID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting directory hardware mapping", err.Error())
	}
}

func (r *DirHardwareMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := DirHardwareMappingResourceModel{
		ID:        types.StringValue(req.ID),
		MappingID: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DirHardwareMappingResource) readIntoModel(ctx context.Context, model *DirHardwareMappingResourceModel, diagnostics *diag.Diagnostics) {
	mapping, err := r.client.GetDirHardwareMapping(ctx, model.MappingID.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diagnostics.AddWarning("Directory hardware mapping not found",
				fmt.Sprintf("Mapping '%s' no longer exists.", model.MappingID.ValueString()))
			return
		}
		diagnostics.AddError("Error reading directory hardware mapping", err.Error())
		return
	}

	model.Comment = types.StringValue(mapping.Comment)

	entries := make([]attr.Value, len(mapping.Map))
	for i, e := range mapping.Map {
		obj, diags := types.ObjectValue(mapEntryAttrTypes, map[string]attr.Value{
			"node": types.StringValue(e.Node),
			"path": types.StringValue(e.Path),
		})
		if diags.HasError() {
			diagnostics.AddError("Error building map entry", "Failed to construct map entry object")
			return
		}
		entries[i] = obj
	}

	listVal, diags := types.ListValue(types.ObjectType{AttrTypes: mapEntryAttrTypes}, entries)
	if diags.HasError() {
		diagnostics.AddError("Error building map list", "Failed to convert map entries to list")
		return
	}
	model.Map = listVal
}

// modelMapToStrings turns the list-of-objects map attribute into the "node=X,path=Y" strings
// that the proxmox API expects.
func modelMapToStrings(ctx context.Context, list types.List, diagnostics *diag.Diagnostics) []string {
	var entries []mapEntryModel
	if diags := list.ElementsAs(ctx, &entries, false); diags.HasError() {
		diagnostics.Append(diags...)
		return nil
	}

	result := make([]string, len(entries))
	for i, e := range entries {
		result[i] = fmt.Sprintf("node=%s,path=%s", e.Node.ValueString(), e.Path.ValueString())
	}
	return result
}
