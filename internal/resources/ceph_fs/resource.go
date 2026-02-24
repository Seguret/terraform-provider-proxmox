package ceph_fs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ resource.Resource = &CephFSResource{}
var _ resource.ResourceWithConfigure = &CephFSResource{}
var _ resource.ResourceWithImportState = &CephFSResource{}

type CephFSResource struct {
	client *client.Client
}

type CephFSResourceModel struct {
	ID           types.String `tfsdk:"id"`
	NodeName     types.String `tfsdk:"node_name"`
	Name         types.String `tfsdk:"name"`
	PGNum        types.Int64  `tfsdk:"pg_num"`
	MetadataPool types.String `tfsdk:"metadata_pool"`
	DataPool     types.String `tfsdk:"data_pool"`
}

func NewResource() resource.Resource {
	return &CephFSResource{}
}

func (r *CephFSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_ceph_fs"
}

func (r *CephFSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a CephFS filesystem on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The Proxmox node on which to create the CephFS.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The CephFS filesystem name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pg_num": schema.Int64Attribute{
				Description: "Number of placement groups for the filesystem pools.",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"metadata_pool": schema.StringAttribute{
				Description: "The metadata pool name.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"data_pool": schema.StringAttribute{
				Description: "The data pool name.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *CephFSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CephFSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CephFSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var pgNum *int
	if !plan.PGNum.IsNull() && !plan.PGNum.IsUnknown() {
		v := int(plan.PGNum.ValueInt64())
		pgNum = &v
	}

	tflog.Debug(ctx, "Creating CephFS", map[string]any{
		"node": plan.NodeName.ValueString(),
		"name": plan.Name.ValueString(),
	})

	if err := r.client.CreateCephFS(ctx, plan.NodeName.ValueString(), plan.Name.ValueString(), pgNum); err != nil {
		resp.Diagnostics.AddError("Error creating CephFS", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.NodeName.ValueString() + "/" + plan.Name.ValueString())

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading CephFS after create", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CephFSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CephFSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading CephFS", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CephFSResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// all fields are ForceNew so Update is never actually called
}

func (r *CephFSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CephFSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCephFS(ctx, state.NodeName.ValueString(), state.Name.ValueString())
	if err == nil {
		return
	}

	// 404 or 405 MethodNotAllowed is fine — older PVE versions dont support this endpoint
	if apiErr, ok := err.(*client.APIError); ok {
		if apiErr.StatusCode == 404 || apiErr.StatusCode == 405 {
			return
		}
	}

	resp.Diagnostics.AddError("Error deleting CephFS", err.Error())
}

func (r *CephFSResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid import ID",
			"Expected format: {node_name}/{name}")
		return
	}

	state := CephFSResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(parts[0]),
		Name:     types.StringValue(parts[1]),
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing CephFS", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CephFSResource) readIntoModel(ctx context.Context, model *CephFSResourceModel) error {
	fsList, err := r.client.GetCephFSList(ctx, model.NodeName.ValueString())
	if err != nil {
		return err
	}

	for _, fs := range fsList {
		if fs.Name == model.Name.ValueString() {
			model.MetadataPool = types.StringValue(fs.MetadataPool)
			model.DataPool = types.StringValue(fs.DataPool)
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
