package download_file

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &DownloadFileResource{}
var _ resource.ResourceWithConfigure = &DownloadFileResource{}
var _ resource.ResourceWithImportState = &DownloadFileResource{}

type DownloadFileResource struct {
	client *client.Client
}

type DownloadFileResourceModel struct {
	ID                types.String `tfsdk:"id"`
	NodeName          types.String `tfsdk:"node_name"`
	Storage           types.String `tfsdk:"storage"`
	URL               types.String `tfsdk:"url"`
	FileName          types.String `tfsdk:"file_name"`
	ContentType       types.String `tfsdk:"content_type"`
	Checksum          types.String `tfsdk:"checksum"`
	ChecksumAlgorithm types.String `tfsdk:"checksum_algorithm"`
	VerifyTLS         types.Bool   `tfsdk:"verify_tls"`
	Size              types.Int64  `tfsdk:"size"`
	FileID            types.String `tfsdk:"file_id"`
}

func NewResource() resource.Resource {
	return &DownloadFileResource{}
}

func (r *DownloadFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_download_file"
}

func (r *DownloadFileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Downloads a file from a URL to a Proxmox VE node storage.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
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
			"storage": schema.StringAttribute{
				Description: "The storage name to download to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Description: "The source URL to download from.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_name": schema.StringAttribute{
				Description: "The target filename on the storage.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content_type": schema.StringAttribute{
				Description: "The content type: 'iso' or 'vztmpl'.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"checksum": schema.StringAttribute{
				Description: "Expected checksum for verification.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"checksum_algorithm": schema.StringAttribute{
				Description: "Checksum algorithm (md5, sha1, sha256, etc.).",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"verify_tls": schema.BoolAttribute{
				Description: "Whether to verify the TLS certificate of the download URL.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"size": schema.Int64Attribute{
				Description: "The size of the downloaded file in bytes.",
				Computed:    true,
			},
			"file_id": schema.StringAttribute{
				Description: "The full volume ID of the file in storage.",
				Computed:    true,
			},
		},
	}
}

func (r *DownloadFileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DownloadFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DownloadFileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	downloadReq := &models.DownloadURLRequest{
		URL:               plan.URL.ValueString(),
		Filename:          plan.FileName.ValueString(),
		Content:           plan.ContentType.ValueString(),
		Checksum:          plan.Checksum.ValueString(),
		ChecksumAlgorithm: plan.ChecksumAlgorithm.ValueString(),
	}
	if !plan.VerifyTLS.IsNull() && !plan.VerifyTLS.IsUnknown() {
		v := boolToInt(plan.VerifyTLS.ValueBool())
		downloadReq.Verify = &v
	}

	node := plan.NodeName.ValueString()
	storage := plan.Storage.ValueString()

	tflog.Debug(ctx, "Downloading file", map[string]any{
		"node":    node,
		"storage": storage,
		"url":     downloadReq.URL,
	})

	upid, err := r.client.DownloadFile(ctx, node, storage, downloadReq)
	if err != nil {
		resp.Diagnostics.AddError("Error initiating file download", err.Error())
		return
	}

	if err := r.client.WaitForTask(ctx, node, upid); err != nil {
		resp.Diagnostics.AddError("Error waiting for download task", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s/%s", node, storage, plan.FileName.ValueString()))
	r.readStorageContent(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DownloadFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DownloadFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readStorageContent(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DownloadFileResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// all fields are ForceNew so Update is never called
}

func (r *DownloadFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DownloadFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.FileID.IsNull() || state.FileID.ValueString() == "" {
		// nothing to delete
		return
	}

	node := state.NodeName.ValueString()
	storage := state.Storage.ValueString()
	volid := state.FileID.ValueString()

	tflog.Debug(ctx, "Deleting storage content", map[string]any{"node": node, "storage": storage, "volid": volid})

	upid, err := r.client.DeleteStorageContent(ctx, node, storage, volid)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok {
			if apiErr.IsNotFound() {
				return
			}
		}
		resp.Diagnostics.AddError("Error deleting storage content", err.Error())
		return
	}

	if upid != "" {
		if err := r.client.WaitForTask(ctx, node, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for delete task", err.Error())
		}
	}
}

func (r *DownloadFileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// format: "node/storage/volid"
	parts := strings.SplitN(req.ID, "/", 3)
	if len(parts) != 3 {
		resp.Diagnostics.AddError("Invalid import ID", "Format: <node>/<storage>/<volid>")
		return
	}

	state := DownloadFileResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(parts[0]),
		Storage:  types.StringValue(parts[1]),
		FileID:   types.StringValue(parts[2]),
	}
	r.readStorageContent(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DownloadFileResource) readStorageContent(ctx context.Context, model *DownloadFileResourceModel, diagnostics *diag.Diagnostics) {
	node := model.NodeName.ValueString()
	storage := model.Storage.ValueString()

	contents, err := r.client.GetStorageContent(ctx, node, storage)
	if err != nil {
		diagnostics.AddError("Error reading storage content", err.Error())
		return
	}

	// match by file_id first since its more specific, fall back to file_name
	targetFileID := model.FileID.ValueString()
	targetFileName := model.FileName.ValueString()

	for _, item := range contents {
		matched := false
		if targetFileID != "" && item.VolID == targetFileID {
			matched = true
		} else if targetFileName != "" && strings.HasSuffix(item.VolID, "/"+targetFileName) {
			matched = true
		}

		if matched {
			model.FileID = types.StringValue(item.VolID)
			model.Size = types.Int64Value(item.Size)
			model.ContentType = types.StringValue(item.Content)
			return
		}
	}

	// file not found — leave state alone so we dont break existing plans
	tflog.Debug(ctx, "Storage content not found, may have been deleted externally", map[string]any{
		"file_id":   targetFileID,
		"file_name": targetFileName,
	})
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
