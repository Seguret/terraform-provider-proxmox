package notification_matcher

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

var _ resource.Resource = &NotificationMatcherResource{}
var _ resource.ResourceWithConfigure = &NotificationMatcherResource{}
var _ resource.ResourceWithImportState = &NotificationMatcherResource{}

type NotificationMatcherResource struct {
	client *client.Client
}

type NotificationMatcherResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	MatchSeverity types.List   `tfsdk:"match_severity"`
	MatchCalendar types.List   `tfsdk:"match_calendar"`
	MatchField    types.List   `tfsdk:"match_field"`
	Target        types.List   `tfsdk:"target"`
	Mode          types.String `tfsdk:"mode"`
	Comment       types.String `tfsdk:"comment"`
	Disable       types.Bool   `tfsdk:"disable"`
}

func NewResource() resource.Resource {
	return &NotificationMatcherResource{}
}

func (r *NotificationMatcherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_notification_matcher"
}

func (r *NotificationMatcherResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE notification matcher.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the notification matcher.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"match_severity": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "List of severity levels to match.",
				Optional:    true,
			},
			"match_calendar": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "List of calendar entries to match.",
				Optional:    true,
			},
			"match_field": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "List of field match expressions.",
				Optional:    true,
			},
			"target": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "List of notification endpoint targets.",
				Optional:    true,
			},
			"mode": schema.StringAttribute{
				Description: "The matcher mode.",
				Optional:    true,
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Comment for the notification matcher.",
				Optional:    true,
				Computed:    true,
			},
			"disable": schema.BoolAttribute{
				Description: "Whether the notification matcher is disabled.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *NotificationMatcherResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationMatcherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NotificationMatcherResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating notification matcher", map[string]any{"name": plan.Name.ValueString()})

	createReq := &models.NotificationMatcherCreateRequest{
		Name:          plan.Name.ValueString(),
		MatchSeverity: listToStrings(plan.MatchSeverity),
		MatchCalendar: listToStrings(plan.MatchCalendar),
		MatchField:    listToStrings(plan.MatchField),
		Target:        listToStrings(plan.Target),
		Mode:          plan.Mode.ValueString(),
		Comment:       plan.Comment.ValueString(),
		Disable:       plan.Disable.ValueBool(),
	}

	if err := r.client.CreateNotificationMatcher(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating notification matcher", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Name.ValueString())

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading notification matcher", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NotificationMatcherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NotificationMatcherResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading notification matcher", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NotificationMatcherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NotificationMatcherResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	disable := plan.Disable.ValueBool()
	updateReq := &models.NotificationMatcherUpdateRequest{
		MatchSeverity: listToStrings(plan.MatchSeverity),
		MatchCalendar: listToStrings(plan.MatchCalendar),
		MatchField:    listToStrings(plan.MatchField),
		Target:        listToStrings(plan.Target),
		Mode:          plan.Mode.ValueString(),
		Comment:       plan.Comment.ValueString(),
		Disable:       &disable,
	}

	if err := r.client.UpdateNotificationMatcher(ctx, plan.Name.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating notification matcher", err.Error())
		return
	}

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading notification matcher", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NotificationMatcherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NotificationMatcherResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteNotificationMatcher(ctx, state.Name.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting notification matcher", err.Error())
	}
}

func (r *NotificationMatcherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := NotificationMatcherResourceModel{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing notification matcher", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NotificationMatcherResource) readIntoModel(ctx context.Context, model *NotificationMatcherResourceModel) error {
	m, err := r.client.GetNotificationMatcher(ctx, model.Name.ValueString())
	if err != nil {
		return err
	}
	model.Name = types.StringValue(m.Name)
	model.MatchSeverity = stringsToList(m.MatchSeverity)
	model.MatchCalendar = stringsToList(m.MatchCalendar)
	model.MatchField = stringsToList(m.MatchField)
	model.Target = stringsToList(m.Target)
	model.Mode = types.StringValue(m.Mode)
	model.Comment = types.StringValue(m.Comment)
	model.Disable = types.BoolValue(m.Disable)
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
