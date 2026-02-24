package oci_image

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &OCIImageResource{}
var _ resource.ResourceWithConfigure = &OCIImageResource{}
var _ resource.ResourceWithImportState = &OCIImageResource{}

type OCIImageResource struct {
	client *client.Client
}

type OCIImageResourceModel struct {
	ID          types.String `tfsdk:"id"`
	NodeName    types.String `tfsdk:"node_name"`
	DatastoreID types.String `tfsdk:"datastore_id"`
	URL         types.String `tfsdk:"url"`
	FileName    types.String `tfsdk:"file_name"`
	PullMethod  types.String `tfsdk:"pull_method"`
	Size        types.String `tfsdk:"size"`
}

func NewResource() resource.Resource {
	return &OCIImageResource{}
}

func (r *OCIImageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_oci_image"
}

func (r *OCIImageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Downloads and manages an OCI container image in Proxmox VE storage (PVE 8+).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The volume ID of the downloaded image.",
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
			"datastore_id": schema.StringAttribute{
				Description: "The storage ID on which to store the image.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Description: "The OCI image URL (e.g. 'docker.io/library/ubuntu:22.04').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_name": schema.StringAttribute{
				Description: "Override the target filename on the storage.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pull_method": schema.StringAttribute{
				Description: "The pull method: 'http' or 'oci'. Defaults to 'oci'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("oci"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size": schema.StringAttribute{
				Description: "The image size as reported by Proxmox.",
				Computed:    true,
			},
		},
	}
}

func (r *OCIImageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OCIImageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OCIImageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	storage := plan.DatastoreID.ValueString()

	// generate a filename from the OCI reference when the user didnt provide one
	fileName := plan.FileName.ValueString()
	if fileName == "" {
		fileName = ociURLToFilename(plan.URL.ValueString())
	}

	downloadReq := &models.DownloadURLRequest{
		URL:      plan.URL.ValueString(),
		Filename: fileName,
		Content:  "import", // OCI/disk images always use "import" as the content type
	}

	tflog.Debug(ctx, "Downloading OCI image to Proxmox storage", map[string]any{
		"node":        node,
		"storage":     storage,
		"url":         downloadReq.URL,
		"file_name":   fileName,
		"pull_method": plan.PullMethod.ValueString(),
	})

	upid, err := r.client.DownloadFile(ctx, node, storage, downloadReq)
	if err != nil {
		resp.Diagnostics.AddError("Error initiating OCI image download", err.Error())
		return
	}

	if err := r.client.WaitForUPID(ctx, upid); err != nil {
		resp.Diagnostics.AddError("Error waiting for OCI image download task", err.Error())
		return
	}

	plan.FileName = types.StringValue(fileName)

	diags := r.readIntoModel(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OCIImageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OCIImageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := r.readIntoModel(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OCIImageResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// all attributes are ForceNew so Update is never called
}

func (r *OCIImageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OCIImageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	volid := state.ID.ValueString()
	if volid == "" {
		return
	}

	node := state.NodeName.ValueString()
	storage := state.DatastoreID.ValueString()

	tflog.Debug(ctx, "Deleting OCI image from Proxmox storage", map[string]any{
		"node":    node,
		"storage": storage,
		"volid":   volid,
	})

	upid, err := r.client.DeleteStorageContent(ctx, node, storage, volid)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return
		}
		resp.Diagnostics.AddError("Error deleting OCI image", err.Error())
		return
	}

	if upid != "" {
		if err := r.client.WaitForUPID(ctx, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for OCI image delete task", err.Error())
		}
	}
}

func (r *OCIImageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// format: {node_name}/{datastore_id}/{volid}
	parts := strings.SplitN(req.ID, "/", 3)
	if len(parts) != 3 {
		resp.Diagnostics.AddError("Invalid import ID",
			"Expected format: <node_name>/<datastore_id>/<volid>")
		return
	}

	state := OCIImageResourceModel{
		ID:          types.StringValue(parts[2]),
		NodeName:    types.StringValue(parts[0]),
		DatastoreID: types.StringValue(parts[1]),
	}

	diags := r.readIntoModel(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// readIntoModel looks up the image in storage content and fills in computed fields.
func (r *OCIImageResource) readIntoModel(ctx context.Context, model *OCIImageResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	node := model.NodeName.ValueString()
	storage := model.DatastoreID.ValueString()

	contents, err := r.client.GetStorageContent(ctx, node, storage)
	if err != nil {
		diags.AddError("Error reading storage content", err.Error())
		return diags
	}

	targetVolID := model.ID.ValueString()
	targetFileName := model.FileName.ValueString()

	for _, item := range contents {
		matched := false
		if targetVolID != "" && item.VolID == targetVolID {
			matched = true
		} else if targetFileName != "" && strings.HasSuffix(item.VolID, "/"+targetFileName) {
			matched = true
		}

		if matched {
			model.ID = types.StringValue(item.VolID)
			model.Size = types.StringValue(fmt.Sprintf("%d", item.Size))
			return diags
		}
	}

	tflog.Warn(ctx, "OCI image not found in storage, removing from state", map[string]any{
		"volid":     targetVolID,
		"file_name": targetFileName,
	})
	model.ID = types.StringValue("")
	return diags
}

// ociURLToFilename turns an OCI image reference into a safe filename.
// e.g. "docker.io/library/ubuntu:22.04" → "ubuntu_22.04.raw"
func ociURLToFilename(ociURL string) string {
	// grab just the last path component
	parts := strings.Split(ociURL, "/")
	name := parts[len(parts)-1]
	// swap the tag separator colon for underscore so its filesystem-safe
	name = strings.ReplaceAll(name, ":", "_")
	if !strings.Contains(name, ".") {
		name += ".raw"
	}
	return name
}
