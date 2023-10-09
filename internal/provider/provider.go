package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &BlueChiProvider{}

type BlueChiProvider struct {
	version string
}

type BlueChiProviderModel struct {
	UseMock types.Bool `tfsdk:"use_mock"`
}

func (p *BlueChiProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bluechi"
	resp.Version = p.version
}

func (p *BlueChiProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"use_mock": schema.BoolAttribute{
				Optional:    true,
				Description: "Flag to indicate if a mock client should be used",
			},
		},
	}
}

func (p *BlueChiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data BlueChiProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = data.UseMock
	resp.ResourceData = data.UseMock
}

func (p *BlueChiProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewBlueChiNodeResource,
	}
}

func (p *BlueChiProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BlueChiProvider{
			version: version,
		}
	}
}
