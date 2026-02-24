package file

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &FileResource{}
var _ resource.ResourceWithConfigure = &FileResource{}
var _ resource.ResourceWithImportState = &FileResource{}

type FileResource struct {
	client *client.Client
}

type FileResourceModel struct {
	ID                types.String `tfsdk:"id"`
	NodeName          types.String `tfsdk:"node_name"`
	DatastoreID       types.String `tfsdk:"datastore_id"`
	SourceFile        types.String `tfsdk:"source_file"`
	ContentType       types.String `tfsdk:"content_type"`
	FileName          types.String `tfsdk:"file_name"`
	Checksum          types.String `tfsdk:"checksum"`
	ChecksumAlgorithm types.String `tfsdk:"checksum_algorithm"`
	Overwrite         types.Bool   `tfsdk:"overwrite"`
	UploadTimeout     types.Int64  `tfsdk:"upload_timeout"`
	Size              types.String `tfsdk:"size"`
}

func NewResource() resource.Resource {
	return &FileResource{}
}

func (r *FileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_file"
}

func (r *FileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a file uploaded to Proxmox VE storage (ISOs, snippets, container templates, etc.).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The volume ID of the file (e.g. 'local:iso/ubuntu.iso').",
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
				Description: "The storage ID on which to store the file (e.g. 'local').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_file": schema.StringAttribute{
				Description: "The source URL to download the file from.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content_type": schema.StringAttribute{
				Description: "The content type: 'iso', 'vztmpl', 'snippets', 'import'. Inferred from URL if omitted.",
				Optional:    true,
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
			"checksum": schema.StringAttribute{
				Description: "Expected checksum of the file for verification.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"checksum_algorithm": schema.StringAttribute{
				Description: "The checksum algorithm (md5, sha1, sha256, sha512).",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"overwrite": schema.BoolAttribute{
				Description: "Whether to overwrite an existing file with the same name. Defaults to true.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"upload_timeout": schema.Int64Attribute{
				Description: "The upload timeout in seconds. Defaults to 1800.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1800),
			},
			"size": schema.StringAttribute{
				Description: "The file size as reported by Proxmox.",
				Computed:    true,
			},
		},
	}
}

func (r *FileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	storage := plan.DatastoreID.ValueString()

	// derive filename from the URL when not explicitly set
	fileName := plan.FileName.ValueString()
	if fileName == "" {
		parts := strings.Split(plan.SourceFile.ValueString(), "/")
		fileName = parts[len(parts)-1]
		// strip query string if present
		if idx := strings.Index(fileName, "?"); idx >= 0 {
			fileName = fileName[:idx]
		}
	}

	// guess content type from extension if the user didnt set it
	contentType := plan.ContentType.ValueString()
	if contentType == "" {
		contentType = inferContentType(fileName)
	}

	downloadReq := &models.DownloadURLRequest{
		URL:               plan.SourceFile.ValueString(),
		Filename:          fileName,
		Content:           contentType,
		Checksum:          plan.Checksum.ValueString(),
		ChecksumAlgorithm: plan.ChecksumAlgorithm.ValueString(),
	}

	tflog.Debug(ctx, "Downloading file to Proxmox storage", map[string]any{
		"node":         node,
		"storage":      storage,
		"url":          downloadReq.URL,
		"file_name":    fileName,
		"content_type": contentType,
	})

	upid, err := r.client.DownloadFile(ctx, node, storage, downloadReq)
	if err != nil {
		resp.Diagnostics.AddError("Error initiating file download", err.Error())
		return
	}

	if err := r.client.WaitForUPID(ctx, upid); err != nil {
		resp.Diagnostics.AddError("Error waiting for file download task", err.Error())
		return
	}

	// set resolved filename and content type before reading back from storage
	plan.FileName = types.StringValue(fileName)
	plan.ContentType = types.StringValue(contentType)

	diags := r.readIntoModel(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := r.readIntoModel(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// empty ID means the file wasnt found — remove from state so terraform plans a re-create
	if state.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FileResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// all structural fields are ForceNew so Update is never called
	// overwrite and upload_timeout are plan-only — no remote update needed
}

func (r *FileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FileResourceModel
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

	tflog.Debug(ctx, "Deleting file from Proxmox storage", map[string]any{
		"node":    node,
		"storage": storage,
		"volid":   volid,
	})

	upid, err := r.client.DeleteStorageContent(ctx, node, storage, volid)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return
		}
		resp.Diagnostics.AddError("Error deleting storage content", err.Error())
		return
	}

	if upid != "" {
		if err := r.client.WaitForUPID(ctx, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for delete task", err.Error())
		}
	}
}

func (r *FileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// format: {node_name}/{datastore_id}/{volid}
	// volids can contain slashes (e.g. local:iso/ubuntu.iso) so only split on the first two "/"
	parts := strings.SplitN(req.ID, "/", 3)
	if len(parts) != 3 {
		resp.Diagnostics.AddError("Invalid import ID",
			"Expected format: <node_name>/<datastore_id>/<volid>")
		return
	}

	state := FileResourceModel{
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

// readIntoModel looks up the file in storage content and fills in computed fields.
// Tries to match by volid first, then falls back to filename suffix.
// Clears the ID if nothing matches so the caller can detect absence.
func (r *FileResource) readIntoModel(ctx context.Context, model *FileResourceModel) diag.Diagnostics {
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
			if model.ContentType.IsNull() || model.ContentType.IsUnknown() || model.ContentType.ValueString() == "" {
				model.ContentType = types.StringValue(item.Content)
			}
			return diags
		}
	}

	// not found — clear the ID so the caller knows and can remove from state
	tflog.Warn(ctx, "File not found in storage, will remove from state", map[string]any{
		"volid":     targetVolID,
		"file_name": targetFileName,
		"node":      node,
		"storage":   storage,
	})
	model.ID = types.StringValue("")
	return diags
}

// inferContentType guesses proxmox content type from the file extension.
func inferContentType(filename string) string {
	lower := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(lower, ".iso"):
		return "iso"
	case strings.HasSuffix(lower, ".tar.gz"),
		strings.HasSuffix(lower, ".tar.xz"),
		strings.HasSuffix(lower, ".tar.zst"):
		return "vztmpl"
	default:
		return "iso"
	}
}

