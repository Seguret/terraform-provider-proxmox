package backups

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &BackupsDataSource{}
var _ datasource.DataSourceWithConfigure = &BackupsDataSource{}

type BackupsDataSource struct {
	client *client.Client
}

type BackupsDataSourceModel struct {
	ID        types.String   `tfsdk:"id"`
	IDs       []types.String `tfsdk:"ids"`
	Storages  []types.String `tfsdk:"storages"`
	Schedules []types.String `tfsdk:"schedules"`
}

func NewDataSource() datasource.DataSource {
	return &BackupsDataSource{}
}

func (d *BackupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_backups"
}

func (d *BackupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves all Proxmox VE backup schedule jobs.",
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Computed: true},
			"ids":       schema.ListAttribute{Description: "Backup job IDs.", Computed: true, ElementType: types.StringType},
			"storages":  schema.ListAttribute{Description: "Target storage for each job.", Computed: true, ElementType: types.StringType},
			"schedules": schema.ListAttribute{Description: "Schedule for each job.", Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *BackupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	d.client = c
}

func (d *BackupsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading backup jobs list")

	jobs, err := d.client.GetBackupJobs(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading backup jobs", err.Error())
		return
	}

	state := BackupsDataSourceModel{
		ID:        types.StringValue("backups"),
		IDs:       make([]types.String, len(jobs)),
		Storages:  make([]types.String, len(jobs)),
		Schedules: make([]types.String, len(jobs)),
	}
	for i, j := range jobs {
		state.IDs[i] = types.StringValue(j.ID)
		state.Storages[i] = types.StringValue(j.Storage)
		state.Schedules[i] = types.StringValue(j.Schedule)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
