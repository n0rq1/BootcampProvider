// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"terraform-provider-devops/provider/client"
	"terraform-provider-devops/provider/devops"
	"terraform-provider-devops/provider/devs"
	"terraform-provider-devops/provider/engineers"
	"terraform-provider-devops/provider/ops"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &DOBProvider{}

// DOBProvider defines the provider implementation.
type DOBProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// DOBProviderModel describes the provider data model.
type DOBProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *DOBProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dob"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *DOBProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
		},
	}
}

func (p *DOBProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config DOBProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if config.Endpoint.IsNull() { /* ... */ }

	// Initialize custom API client for data sources and resources
	var endpointPtr *string
	if !config.Endpoint.IsNull() && !config.Endpoint.IsUnknown() {
		v := config.Endpoint.ValueString()
		endpointPtr = &v
	}

	c, err := client.NewClient(endpointPtr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create API client",
			err.Error(),
		)
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

// Resources defines the resources implemented in the provider.
func (p *DOBProvider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        engineers.NewEngineerResource,
        devs.NewDevResource,
		ops.NewOpsResource,
		devops.NewDevOpsResource,
    }
}

func (p *DOBProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{
        engineers.NewEngineersDataSource,
        devs.NewDevDataSource,
		ops.NewOpsDataSource,
		devops.NewDevopsDataSource,
    }
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DOBProvider{
			version: version,
		}
	}
}
