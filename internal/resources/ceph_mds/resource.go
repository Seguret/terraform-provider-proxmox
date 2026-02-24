package ceph_mds

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

var _ resource.Resource = &CephMDSResource{}
var _ resource.ResourceWithConfigure = &CephMDSResource{}
var _ resource.ResourceWithImportState = &CephMDSResource{}

type CephMDSResource struct {
	client *client.Client
}

type CephMDSResourceModel struct {
	ID       types.String `tfsdk:"id"`
	NodeName types.String `tfsdk:"node_name"`
	Name     types.String `tfsdk:"name"`
	State    types.String `tfsdk:"state"`
	Addr     types.String `tfsdk:"addr"`
}

func NewResource() resource.Resource {
	return &CephMDSResource{}
}

func (r *CephMDSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_ceph_mds"
}

func (r *CephMDSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Ceph MDS daemon on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The Proxmox node on which to create the MDS.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The MDS daemon name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Description: "The MDS daemon state.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"addr": schema.StringAttribute{
				Description: "The MDS daemon address.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *CephMDSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CephMDSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CephMDSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Ceph MDS", map[string]any{
		"node": plan.NodeName.ValueString(),
		"name": plan.Name.ValueString(),
	})

	if err := r.client.CreateCephMDS(ctx, plan.NodeName.ValueString(), plan.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error creating Ceph MDS", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.NodeName.ValueString() + "/" + plan.Name.ValueString())

	if err := r.readIntoModel(ctx, &plan); err != nil {
		// MDS state isnt always visible right away — just zero out these fields
		plan.State = types.StringValue("")
		plan.Addr = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CephMDSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CephMDSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Ceph MDS", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CephMDSResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// all fields are ForceNew so Update is never actually called
}

func (r *CephMDSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CephMDSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteCephMDS(ctx, state.NodeName.ValueString(), state.Name.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting Ceph MDS", err.Error())
	}
}

func (r *CephMDSResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid import ID",
			"Expected format: {node_name}/{name}")
		return
	}

	state := CephMDSResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(parts[0]),
		Name:     types.StringValue(parts[1]),
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing Ceph MDS", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CephMDSResource) readIntoModel(ctx context.Context, model *CephMDSResourceModel) error {
	mdsList, err := r.client.GetCephMDSList(ctx, model.NodeName.ValueString())
	if err != nil {
		return err
	}

	for _, m := range mdsList {
		if m.Name == model.Name.ValueString() {
			model.State = types.StringValue(m.State)
			model.Addr = types.StringValue(m.Addr)
			return nil
		}
	}

	return &client.APIError{StatusCode: 404, Status: "404 Not Found"}
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
