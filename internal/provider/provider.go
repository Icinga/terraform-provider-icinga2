package provider

import (
	"context"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/lrsmith/go-icinga2-api/iapi"
)

var (
	_ provider.Provider = &icinga2Provider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &icinga2Provider{
			version: version,
		}
	}
}

type icinga2Provider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type icinga2ProviderModel struct {
	Host                     types.String `tfsdk:"host"`
	Username                 types.String `tfsdk:"username"`
	Password                 types.String `tfsdk:"password"`
	Insecure_skip_tls_verify types.Bool   `tfsdk:"insecure_skip_tls_verify"`
}

func (p *icinga2Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "icinga2"
	resp.Version = p.version
}

func (p *icinga2Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"insecure_skip_tls_verify": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (p *icinga2Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config icinga2ProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown icinga2 API Host",
			"The provider cannot create the icinga2 API client as there is an unknown configuration value for the icinga2 API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ICINGA2_API_URL environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown icinga2 API Username",
			"The provider cannot create the icinga2 API client as there is an unknown configuration value for the icinga2 API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ICINGA2_API_USER environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown icinga2 API Password",
			"The provider cannot create the icinga2 API client as there is an unknown configuration value for the icinga2 API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ICINGA2_API_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("ICINGA2_API_URL")
	username := os.Getenv("ICINGA2_API_USER")
	password := os.Getenv("ICINGA2_API_PASSWORD")
	tlsVerify, _ := strconv.ParseBool(os.Getenv("ICINGA2_INSECURE_SKIP_TLS_VERIFY"))

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing icinga2 API Host",
			"The provider cannot create the icinga2 API client as there is a missing or empty value for the icinga2 API host. "+
				"Set the host value in the configuration or use the ICINGA2_API_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing icinga2 API Username",
			"The provider cannot create the icinga2 API client as there is a missing or empty value for the icinga2 API username. "+
				"Set the username value in the configuration or use the ICINGA2_API_USER environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing icinga2 API Password",
			"The provider cannot create the icinga2 API client as there is a missing or empty value for the icinga2 API password. "+
				"Set the password value in the configuration or use the ICINGA2_API_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := iapi.New(
		username,
		password,
		host,
		tlsVerify,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create icinga2 API Client",
			"An unexpected error occurred when creating the icinga2 API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"icinga2 Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *icinga2Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *icinga2Provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		HostGroup,
	}
}
