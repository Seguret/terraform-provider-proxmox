package apt_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

// handleURIs maps known proxmox apt repository handles to their URIs.
// Used to locate an entry in the repository list after adding it.
var handleURIs = map[string]string{
	"pve-enterprise":              "https://enterprise.proxmox.com/debian/pve",
	"pve-no-subscription":         "http://download.proxmox.com/debian/pve",
	"pvetest":                     "http://download.proxmox.com/debian/pvetest",
	"ceph-quincy-enterprise":      "https://enterprise.proxmox.com/debian/ceph-quincy",
	"ceph-quincy-no-subscription": "http://download.proxmox.com/debian/ceph-quincy",
	"ceph-reef-enterprise":        "https://enterprise.proxmox.com/debian/ceph-reef",
	"ceph-reef-no-subscription":   "http://download.proxmox.com/debian/ceph-reef",
}

var _ resource.Resource = &AptRepositoryResource{}
var _ resource.ResourceWithConfigure = &AptRepositoryResource{}
var _ resource.ResourceWithImportState = &AptRepositoryResource{}

type AptRepositoryResource struct {
	client *client.Client
}

type AptRepositoryResourceModel struct {
	ID       types.String `tfsdk:"id"`
	NodeName types.String `tfsdk:"node_name"`
	Handle   types.String `tfsdk:"handle"`
	Enabled  types.Bool   `tfsdk:"enabled"`
	// stored so we can do idempotent reads/updates
	FilePath types.String `tfsdk:"file_path"`
	Index    types.Int64  `tfsdk:"index"`
}

func NewResource() resource.Resource {
	return &AptRepositoryResource{}
}

func (r *AptRepositoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_apt_repository"
}

func (r *AptRepositoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a standard Proxmox VE APT repository on a node.",
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
			"handle": schema.StringAttribute{
				Description: "The standard repository handle (e.g. 'pve-no-subscription', 'pve-enterprise').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the repository is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"file_path": schema.StringAttribute{
				Description: "The sources file path (computed after creation).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"index": schema.Int64Attribute{
				Description: "The index of the repository entry within its file (computed after creation).",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *AptRepositoryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AptRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AptRepositoryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	handle := plan.Handle.ValueString()

	tflog.Debug(ctx, "Adding APT repository", map[string]any{"node": node, "handle": handle})

	// grab current digest before adding so proxmox can detect concurrent changes
	current, err := r.client.GetAptRepositories(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading APT repositories", err.Error())
		return
	}

	// add the standard repository handle
	if err := r.client.AddAptRepository(ctx, node, &models.AptRepositoryAddRequest{
		Handle: handle,
		Digest: current.Digest,
	}); err != nil {
		resp.Diagnostics.AddError("Error adding APT repository", err.Error())
		return
	}

	// re-read to locate the newly added entry
	updated, err := r.client.GetAptRepositories(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading APT repositories after add", err.Error())
		return
	}

	filePath, index, found := findRepoByHandle(handle, updated)
	if !found {
		resp.Diagnostics.AddError("APT repository not found after add",
			fmt.Sprintf("Could not locate handle '%s' in repository list.", handle))
		return
	}

	plan.FilePath = types.StringValue(filePath)
	plan.Index = types.Int64Value(int64(index))
	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", node, handle))

	// set initial enabled state
	enabledInt := 0
	if plan.Enabled.ValueBool() {
		enabledInt = 1
	}

	// need a fresh digest before changing the entry
	updated2, err := r.client.GetAptRepositories(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading APT repositories digest", err.Error())
		return
	}

	if err := r.client.ChangeAptRepository(ctx, node, &models.AptRepositoryChangeRequest{
		Path:    filePath,
		Index:   index,
		Enabled: enabledInt,
		Digest:  updated2.Digest,
	}); err != nil {
		resp.Diagnostics.AddError("Error setting APT repository enabled state", err.Error())
		return
	}

	if err := r.readState(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading APT repository state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AptRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AptRepositoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readState(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error reading APT repository", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AptRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AptRepositoryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()

	current, err := r.client.GetAptRepositories(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading APT repositories", err.Error())
		return
	}

	enabledInt := 0
	if plan.Enabled.ValueBool() {
		enabledInt = 1
	}

	if err := r.client.ChangeAptRepository(ctx, node, &models.AptRepositoryChangeRequest{
		Path:    plan.FilePath.ValueString(),
		Index:   int(plan.Index.ValueInt64()),
		Enabled: enabledInt,
		Digest:  current.Digest,
	}); err != nil {
		resp.Diagnostics.AddError("Error updating APT repository", err.Error())
		return
	}

	if err := r.readState(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading APT repository after update", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AptRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AptRepositoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := state.NodeName.ValueString()

	current, err := r.client.GetAptRepositories(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading APT repositories for disable", err.Error())
		return
	}

	// standard repos cant be truly deleted — just disable them
	_ = r.client.ChangeAptRepository(ctx, node, &models.AptRepositoryChangeRequest{
		Path:    state.FilePath.ValueString(),
		Index:   int(state.Index.ValueInt64()),
		Enabled: 0,
		Digest:  current.Digest,
	})
}

func (r *AptRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// expected format: "<node>/<handle>"
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Format: <node>/<handle> (e.g. 'pve/pve-no-subscription')")
		return
	}

	node, handle := parts[0], parts[1]

	repos, err := r.client.GetAptRepositories(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading APT repositories", err.Error())
		return
	}

	filePath, index, found := findRepoByHandle(handle, repos)
	if !found {
		resp.Diagnostics.AddError("APT repository not found",
			fmt.Sprintf("Could not find handle '%s' on node '%s'.", handle, node))
		return
	}

	state := AptRepositoryResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(node),
		Handle:   types.StringValue(handle),
		FilePath: types.StringValue(filePath),
		Index:    types.Int64Value(int64(index)),
	}

	if err := r.readState(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error reading APT repository state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AptRepositoryResource) readState(ctx context.Context, model *AptRepositoryResourceModel) error {
	repos, err := r.client.GetAptRepositories(ctx, model.NodeName.ValueString())
	if err != nil {
		return err
	}

	filePath := model.FilePath.ValueString()
	index := int(model.Index.ValueInt64())

	for _, f := range repos.Files {
		if f.Filename != filePath {
			continue
		}
		if index >= len(f.Repositories) {
			return fmt.Errorf("repository index %d out of range in file '%s'", index, filePath)
		}
		model.Enabled = types.BoolValue(f.Repositories[index].Enabled)
		return nil
	}

	return fmt.Errorf("repository file '%s' not found", filePath)
}

// findRepoByHandle scans the repo list for an entry whose URI matches the given handle.
// Falls back to substring matching if the handle isnt in our known map.
func findRepoByHandle(handle string, repos *models.AptRepositoriesResponse) (string, int, bool) {
	uri, known := handleURIs[handle]
	for _, f := range repos.Files {
		for i, repo := range f.Repositories {
			for _, u := range repo.URIs {
				if known && strings.TrimRight(u, "/") == strings.TrimRight(uri, "/") {
					return f.Filename, i, true
				}
				// fallback: check if the URI contains the handle string
				if !known && strings.Contains(u, handle) {
					return f.Filename, i, true
				}
			}
		}
	}
	return "", 0, false
}
