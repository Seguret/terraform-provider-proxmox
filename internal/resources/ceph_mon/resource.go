package ceph_mon

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

var _ resource.Resource = &CephMONResource{}
var _ resource.ResourceWithConfigure = &CephMONResource{}
var _ resource.ResourceWithImportState = &CephMONResource{}

type CephMONResource struct {
	client *client.Client
}

type CephMONResourceModel struct {
	ID       types.String `tfsdk:"id"`
	NodeName types.String `tfsdk:"node_name"`
	MonName  types.String `tfsdk:"mon_name"`
	Host     types.String `tfsdk:"host"`
	Addr     types.String `tfsdk:"addr"`
}

func NewResource() resource.Resource {
	return &CephMONResource{}
}

func (r *CephMONResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_ceph_mon"
}

func (r *CephMONResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Ceph monitor (MON) on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The Proxmox node on which to create the MON.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mon_name": schema.StringAttribute{
				Description: "The MON name (typically the node name).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"host": schema.StringAttribute{
				Description: "The MON host address.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"addr": schema.StringAttribute{
				Description: "The MON bind address.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *CephMONResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CephMONResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CephMONResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Ceph MON", map[string]any{
		"node":     plan.NodeName.ValueString(),
		"mon_name": plan.MonName.ValueString(),
	})

	if err := r.client.CreateCephMON(ctx, plan.NodeName.ValueString(), plan.MonName.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error creating Ceph MON", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.NodeName.ValueString() + "/" + plan.MonName.ValueString())

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading Ceph MON after create", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CephMONResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CephMONResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Ceph MON", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CephMONResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// all fields are ForceNew so Update is never actually called
}

func (r *CephMONResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CephMONResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteCephMON(ctx, state.NodeName.ValueString(), state.MonName.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting Ceph MON", err.Error())
	}
}

func (r *CephMONResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid import ID",
			"Expected format: {node_name}/{mon_name}")
		return
	}

	state := CephMONResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(parts[0]),
		MonName:  types.StringValue(parts[1]),
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing Ceph MON", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CephMONResource) readIntoModel(ctx context.Context, model *CephMONResourceModel) error {
	mons, err := r.client.GetNodeCephMON(ctx, model.NodeName.ValueString())
	if err != nil {
		return err
	}

	for _, m := range mons {
		if m.Name == model.MonName.ValueString() {
			model.Host = types.StringValue(m.Host)
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
