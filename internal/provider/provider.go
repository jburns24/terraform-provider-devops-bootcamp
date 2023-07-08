// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure DevOpsBootcampProvider satisfies various provider interfaces.
var _ provider.Provider = &DevOpsBootcampProvider{}

// DevOpsBootcampProvider defines the provider implementation.
type DevOpsBootcampProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// DevOpsBootcampProviderModel describes the provider data model.
type DevOpsBootcampProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *DevOpsBootcampProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "devops-bootcamp"
	resp.Version = p.version
}

func (p *DevOpsBootcampProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
		},
	}
}

func (p *DevOpsBootcampProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data DevOpsBootcampProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	var host string
	if !data.Endpoint.IsNull() {
		host = data.Endpoint.ValueString()
	}

	// Example client configuration for data sources and resources
	client, err := NewClient(&host)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to API Client",
			"An unexpected error occurred when creating the API client. "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *DevOpsBootcampProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewEngineerResource,
	}
}

func (p *DevOpsBootcampProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEngineerDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DevOpsBootcampProvider{
			version: version,
		}
	}
}
