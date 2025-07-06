// Copyright 2025 Jamf Software LLC.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Jamf-Concepts/terraform-provider-jamfcompliancebenchmarkengine/internal/client"
	"github.com/Jamf-Concepts/terraform-provider-jamfcompliancebenchmarkengine/internal/resources/baselines"
	"github.com/Jamf-Concepts/terraform-provider-jamfcompliancebenchmarkengine/internal/resources/benchmark"
	"github.com/Jamf-Concepts/terraform-provider-jamfcompliancebenchmarkengine/internal/resources/rules"
)

// providerModel describes the provider data model for configuration.
type providerModel struct {
	Region       types.String `tfsdk:"region"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

// jamfComplianceBenchmarkEngineProvider implements the Terraform provider for Jamf Compliance Benchmark Engine.
type jamfComplianceBenchmarkEngineProvider struct {
	client *client.Client
}

// Ensure jamfComplianceBenchmarkEngineProvider satisfies the provider.Provider interface.
var _ provider.Provider = &jamfComplianceBenchmarkEngineProvider{}

// Metadata sets the provider type name for the Terraform provider.
func (p *jamfComplianceBenchmarkEngineProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jamfcompliancebenchmarkengine"
}

// Schema sets the Terraform schema for the provider.
func (p *jamfComplianceBenchmarkEngineProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provider for Jamf Compliance Benchmark Engine. https://learn.jamf.com/en-US/bundle/jamf-compliance-benchmarks-configuration-guide/page/Compliance_Benchmarks.html Configure region, client_id, and client_secret.",
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required:    true,
				Description: "The Jamf region to use (us, eu, apac)",
				Validators: []validator.String{
					RegionValidator{},
				},
			},
			"client_id": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The OAuth client ID for authentication.",
			},
			"client_secret": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The OAuth client secret for authentication.",
			},
		},
	}
}

// Configure sets up the API client for the provider from the provider configuration.
func (p *jamfComplianceBenchmarkEngineProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := config.Region.ValueString()
	clientID := config.ClientID.ValueString()
	clientSecret := config.ClientSecret.ValueString()

	if region == "" || clientID == "" || clientSecret == "" {
		resp.Diagnostics.AddError(
			"Missing Required Provider Configuration",
			"All of region, client_id, and client_secret must be set.",
		)
		return
	}

	apiClient := client.NewClient(region, clientID, clientSecret)
	p.client = apiClient
	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

// Resources returns the list of resource constructors for the provider.
func (p *jamfComplianceBenchmarkEngineProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		benchmark.NewBenchmarkResource,
	}
}

// DataSources returns the list of data source constructors for the provider.
func (p *jamfComplianceBenchmarkEngineProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		baselines.NewBaselinesDataSource,
		rules.NewRulesDataSource,
		benchmark.NewBenchmarkDataSource,
	}
}

// New creates a new instance of the Jamf Compliance Benchmark Engine provider.
func New() provider.Provider {
	return &jamfComplianceBenchmarkEngineProvider{}
}

// RegionValidator ensures region is one of the allowed values for the provider.
type RegionValidator struct{}

// ValidateString validates that the region is one of 'us', 'eu', or 'apac'.
func (v RegionValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()
	switch value {
	case "us", "eu", "apac":
		return
	default:
		resp.Diagnostics.AddError(
			"Invalid Region",
			fmt.Sprintf("Region must be one of 'us', 'eu', or 'apac', got: %s", value),
		)
	}
}

// Description returns a description of the region validator.
func (v RegionValidator) Description(_ context.Context) string {
	return "Validates that the region is one of 'us', 'eu', or 'apac'"
}

// MarkdownDescription returns a markdown description of the region validator.
func (v RegionValidator) MarkdownDescription(_ context.Context) string {
	return "Validates that the region is one of `us`, `eu`, or `apac`"
}
