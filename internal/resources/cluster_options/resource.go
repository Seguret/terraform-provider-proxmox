package cluster_options

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &ClusterOptionsResource{}
var _ resource.ResourceWithConfigure = &ClusterOptionsResource{}
var _ resource.ResourceWithImportState = &ClusterOptionsResource{}

type ClusterOptionsResource struct {
	client *client.Client
}

type ClusterOptionsResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Keyboard         types.String `tfsdk:"keyboard"`
	Language         types.String `tfsdk:"language"`
	EmailFrom        types.String `tfsdk:"email_from"`
	HTTPProxy        types.String `tfsdk:"http_proxy"`
	MaxWorkers       types.Int64  `tfsdk:"max_workers"`
	MigrationUnsecure types.Bool  `tfsdk:"migration_unsecure"`
	MigrationType    types.String `tfsdk:"migration_type"`
	HAShutdownPolicy types.String `tfsdk:"ha_shutdown_policy"`
}

func NewResource() resource.Resource {
	return &ClusterOptionsResource{}
}

func (r *ClusterOptionsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_cluster_options"
}

func (r *ClusterOptionsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages cluster-wide Proxmox VE options.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"keyboard": schema.StringAttribute{
				Description: "Default keyboard layout (e.g., 'en-us', 'de', 'fr').",
				Optional:    true,
				Computed:    true,
			},
			"language": schema.StringAttribute{
				Description: "Default web UI language.",
				Optional:    true,
				Computed:    true,
			},
			"email_from": schema.StringAttribute{
				Description: "Email address used as sender for notifications.",
				Optional:    true,
				Computed:    true,
			},
			"http_proxy": schema.StringAttribute{
				Description: "HTTP proxy URL for the cluster (used for apt, etc.).",
				Optional:    true,
				Computed:    true,
			},
			"max_workers": schema.Int64Attribute{
				Description: "Maximum number of workers for bulk operations.",
				Optional:    true,
				Computed:    true,
			},
			"migration_unsecure": schema.BoolAttribute{
				Description: "Whether to allow unsecured migrations (non-TLS).",
				Optional:    true,
				Computed:    true,
			},
			"migration_type": schema.StringAttribute{
				Description: "Migration type: 'secure' or 'insecure' or 'websocket'.",
				Optional:    true,
				Computed:    true,
			},
			"ha_shutdown_policy": schema.StringAttribute{
				Description: "HA shutdown policy: 'freeze', 'failover', 'conditional', or 'migrate'.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *ClusterOptionsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ClusterOptionsResource) applyAndRead(ctx context.Context, model *ClusterOptionsResourceModel) error {
	migUnsecure := 0
	if model.MigrationUnsecure.ValueBool() {
		migUnsecure = 1
	}
	maxWorkers := int(model.MaxWorkers.ValueInt64())

	updateReq := &models.ClusterOptionsUpdateRequest{
		Keyboard:          model.Keyboard.ValueString(),
		Language:          model.Language.ValueString(),
		EmailFrom:         model.EmailFrom.ValueString(),
		HTTPProxy:         model.HTTPProxy.ValueString(),
		MigrationUnsecure: &migUnsecure,
		MigrationType:     model.MigrationType.ValueString(),
		HAShutdownPolicy:  model.HAShutdownPolicy.ValueString(),
	}
	if maxWorkers > 0 {
		updateReq.MaxWorkers = &maxWorkers
	}

	if err := r.client.UpdateClusterOptions(ctx, updateReq); err != nil {
		return fmt.Errorf("error setting cluster options: %w", err)
	}

	return r.readState(ctx, model)
}

func (r *ClusterOptionsResource) readState(ctx context.Context, model *ClusterOptionsResourceModel) error {
	opts, err := r.client.GetClusterOptions(ctx)
	if err != nil {
		return fmt.Errorf("error reading cluster options: %w", err)
	}

	model.Keyboard = types.StringValue(opts.Keyboard)
	model.Language = types.StringValue(opts.Language)
	model.EmailFrom = types.StringValue(opts.EmailFrom)
	model.HTTPProxy = types.StringValue(opts.HTTPProxy)
	model.MigrationType = types.StringValue(opts.MigrationType)
	model.HAShutdownPolicy = types.StringValue(opts.HAShutdownPolicy)

	if opts.MaxWorkers != nil {
		model.MaxWorkers = types.Int64Value(int64(*opts.MaxWorkers))
	} else {
		model.MaxWorkers = types.Int64Value(0)
	}

	if opts.MigrationUnsecure != nil {
		model.MigrationUnsecure = types.BoolValue(*opts.MigrationUnsecure == 1)
	} else {
		model.MigrationUnsecure = types.BoolValue(false)
	}

	return nil
}

func (r *ClusterOptionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ClusterOptionsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue("cluster")

	if err := r.applyAndRead(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error creating cluster options", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ClusterOptionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ClusterOptionsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readState(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error reading cluster options", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ClusterOptionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ClusterOptionsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.applyAndRead(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error updating cluster options", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ClusterOptionsResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// cluster options cant be deleted, only changed — just drop from state
}

func (r *ClusterOptionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := ClusterOptionsResourceModel{
		ID: types.StringValue("cluster"),
	}

	if err := r.readState(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing cluster options", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
