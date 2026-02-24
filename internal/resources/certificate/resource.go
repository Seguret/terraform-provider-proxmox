package certificate

import (
	"context"
	"fmt"

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

var _ resource.Resource = &CertificateResource{}
var _ resource.ResourceWithConfigure = &CertificateResource{}
var _ resource.ResourceWithImportState = &CertificateResource{}

type CertificateResource struct {
	client *client.Client
}

type CertificateResourceModel struct {
	ID           types.String `tfsdk:"id"`
	NodeName     types.String `tfsdk:"node_name"`
	Certificate  types.String `tfsdk:"certificate"`
	PrivateKey   types.String `tfsdk:"private_key"`
	Force        types.Bool   `tfsdk:"force"`
	Restart      types.Bool   `tfsdk:"restart"`
	Fingerprint  types.String `tfsdk:"fingerprint"`
	Issuer       types.String `tfsdk:"issuer"`
	Subject      types.String `tfsdk:"subject"`
	NotBefore    types.Int64  `tfsdk:"not_before"`
	NotAfter     types.Int64  `tfsdk:"not_after"`
}

func NewResource() resource.Resource {
	return &CertificateResource{}
}

func (r *CertificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_certificate"
}

func (r *CertificateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a custom TLS certificate on a Proxmox VE node.",
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
			"certificate": schema.StringAttribute{
				Description: "The PEM-encoded certificate (and chain).",
				Required:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "The PEM-encoded private key.",
				Optional:    true,
				Sensitive:   true,
			},
			"force": schema.BoolAttribute{
				Description: "Whether to overwrite an existing certificate.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"restart": schema.BoolAttribute{
				Description: "Whether to restart pveproxy after uploading.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"fingerprint": schema.StringAttribute{
				Description: "The certificate fingerprint (SHA-256).",
				Computed:    true,
			},
			"issuer": schema.StringAttribute{
				Description: "The certificate issuer.",
				Computed:    true,
			},
			"subject": schema.StringAttribute{
				Description: "The certificate subject.",
				Computed:    true,
			},
			"not_before": schema.Int64Attribute{
				Description: "Certificate validity start (Unix timestamp).",
				Computed:    true,
			},
			"not_after": schema.Int64Attribute{
				Description: "Certificate validity end (Unix timestamp).",
				Computed:    true,
			},
		},
	}
}

func (r *CertificateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CertificateResource) upload(ctx context.Context, model *CertificateResourceModel) error {
	forceInt := 0
	if model.Force.ValueBool() {
		forceInt = 1
	}
	restartInt := 0
	if model.Restart.ValueBool() {
		restartInt = 1
	}

	uploadReq := &models.NodeCertificateUploadRequest{
		Certificates: model.Certificate.ValueString(),
		Key:          model.PrivateKey.ValueString(),
		Force:        &forceInt,
		Restart:      &restartInt,
	}

	cert, err := r.client.UploadNodeCertificate(ctx, model.NodeName.ValueString(), uploadReq)
	if err != nil {
		return fmt.Errorf("error uploading certificate: %w", err)
	}

	if cert != nil {
		model.Fingerprint = types.StringValue(cert.Fingerprint)
		model.Issuer = types.StringValue(cert.Issuer)
		model.Subject = types.StringValue(cert.Subject)
		model.NotBefore = types.Int64Value(cert.NotBefore)
		model.NotAfter = types.Int64Value(cert.NotAfter)
	}

	return nil
}

func (r *CertificateResource) readState(ctx context.Context, model *CertificateResourceModel) error {
	certs, err := r.client.GetNodeCertificates(ctx, model.NodeName.ValueString())
	if err != nil {
		return fmt.Errorf("error reading certificates: %w", err)
	}

	// look for the pveproxy cert specifically
	for _, cert := range certs {
		if cert.Filename == "pveproxy-ssl.pem" || cert.Filename == "/etc/pve/local/pveproxy-ssl.pem" {
			model.Fingerprint = types.StringValue(cert.Fingerprint)
			model.Issuer = types.StringValue(cert.Issuer)
			model.Subject = types.StringValue(cert.Subject)
			model.NotBefore = types.Int64Value(cert.NotBefore)
			model.NotAfter = types.Int64Value(cert.NotAfter)
			return nil
		}
	}

	// fallback to whatever is first if pveproxy cert wasnt found
	if len(certs) > 0 {
		cert := certs[0]
		model.Fingerprint = types.StringValue(cert.Fingerprint)
		model.Issuer = types.StringValue(cert.Issuer)
		model.Subject = types.StringValue(cert.Subject)
		model.NotBefore = types.Int64Value(cert.NotBefore)
		model.NotAfter = types.Int64Value(cert.NotAfter)
	}

	return nil
}

func (r *CertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(plan.NodeName.ValueString())

	tflog.Debug(ctx, "Uploading certificate", map[string]any{"node": plan.NodeName.ValueString()})

	if err := r.upload(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error uploading certificate", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readState(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error reading certificate", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.upload(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error updating certificate", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteNodeCertificate(ctx, state.NodeName.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting certificate", err.Error())
	}
}

func (r *CertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	node := req.ID
	state := CertificateResourceModel{
		ID:       types.StringValue(node),
		NodeName: types.StringValue(node),
		Force:    types.BoolValue(true),
		Restart:  types.BoolValue(true),
	}

	if err := r.readState(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing certificate", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
