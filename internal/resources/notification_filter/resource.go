package notification_filter

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &NotificationFilterResource{}
var _ resource.ResourceWithConfigure = &NotificationFilterResource{}
var _ resource.ResourceWithImportState = &NotificationFilterResource{}

type NotificationFilterResource struct {
	client *client.Client
}

type NotificationFilterResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	MinSeverity types.String `tfsdk:"min_severity"`
	MaxSeverity types.String `tfsdk:"max_severity"`
	Mode        types.String `tfsdk:"mode"`
	Rules       types.List   `tfsdk:"rules"`
	Comment     types.String `tfsdk:"comment"`
	Disable     types.Bool   `tfsdk:"disable"`
}

func NewResource() resource.Resource {
	return &NotificationFilterResource{}
}

func (r *NotificationFilterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_notification_filter"
}

func (r *NotificationFilterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE notification filter.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the notification filter.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"min_severity": schema.StringAttribute{
				Description: "Minimum severity level for the filter.",
				Optional:    true,
				Computed:    true,
			},
			"max_severity": schema.StringAttribute{
				Description: "Maximum severity level for the filter.",
				Optional:    true,
				Computed:    true,
			},
			"mode": schema.StringAttribute{
				Description: "The filter mode.",
				Optional:    true,
				Computed:    true,
			},
			"rules": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "List of filter rules.",
				Optional:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Comment for the notification filter.",
				Optional:    true,
				Computed:    true,
			},
			"disable": schema.BoolAttribute{
				Description: "Whether the notification filter is disabled.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *NotificationFilterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationFilterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NotificationFilterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating notification filter", map[string]any{"name": plan.Name.ValueString()})

	createReq := &models.NotificationFilterCreateRequest{
		Name:        plan.Name.ValueString(),
		MinSeverity: plan.MinSeverity.ValueString(),
		MaxSeverity: plan.MaxSeverity.ValueString(),
		Mode:        plan.Mode.ValueString(),
		Rules:       listToStrings(plan.Rules),
		Comment:     plan.Comment.ValueString(),
		Disable:     plan.Disable.ValueBool(),
	}

	if err := r.client.CreateNotificationFilter(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating notification filter", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Name.ValueString())

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading notification filter", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NotificationFilterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NotificationFilterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading notification filter", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NotificationFilterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NotificationFilterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	disable := plan.Disable.ValueBool()
	updateReq := &models.NotificationFilterUpdateRequest{
		MinSeverity: plan.MinSeverity.ValueString(),
		MaxSeverity: plan.MaxSeverity.ValueString(),
		Mode:        plan.Mode.ValueString(),
		Rules:       listToStrings(plan.Rules),
		Comment:     plan.Comment.ValueString(),
		Disable:     &disable,
	}

	if err := r.client.UpdateNotificationFilter(ctx, plan.Name.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating notification filter", err.Error())
		return
	}

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading notification filter", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NotificationFilterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NotificationFilterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteNotificationFilter(ctx, state.Name.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting notification filter", err.Error())
	}
}

func (r *NotificationFilterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := NotificationFilterResourceModel{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing notification filter", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NotificationFilterResource) readIntoModel(ctx context.Context, model *NotificationFilterResourceModel) error {
	f, err := r.client.GetNotificationFilter(ctx, model.Name.ValueString())
	if err != nil {
		return err
	}
	model.Name = types.StringValue(f.Name)
	model.MinSeverity = types.StringValue(f.MinSeverity)
	model.MaxSeverity = types.StringValue(f.MaxSeverity)
	model.Mode = types.StringValue(f.Mode)
	model.Rules = stringsToList(f.Rules)
	model.Comment = types.StringValue(f.Comment)
	model.Disable = types.BoolValue(f.Disable)
	return nil
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}

func stringsToList(strs []string) types.List {
	if len(strs) == 0 {
		return types.ListValueMust(types.StringType, []attr.Value{})
	}
	vals := make([]attr.Value, len(strs))
	for i, s := range strs {
		vals[i] = types.StringValue(s)
	}
	return types.ListValueMust(types.StringType, vals)
}

func listToStrings(list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	var strs []string
	for _, v := range list.Elements() {
		strs = append(strs, v.(types.String).ValueString())
	}
	return strs
}
