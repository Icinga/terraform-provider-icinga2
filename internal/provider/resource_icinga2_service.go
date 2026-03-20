package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/lrsmith/go-icinga2-api/iapi"
)

var (
	_ resource.Resource              = &serviceResource{}
	_ resource.ResourceWithConfigure = &serviceResource{}
)

func Service() resource.Resource {
	return &serviceResource{}
}

type serviceResourceModel struct {
	ID           types.String `tfsdk:"id"`
	LastUpdated  types.String `tfsdk:"last_updated"`
	Name         types.String `tfsdk:"name"`
	Hostname     types.String `tfsdk:"hostname"`
	CheckCommand types.String `tfsdk:"check_command"`
	Vars         types.Map    `tfsdk:"vars"`
	Templates    types.List   `tfsdk:"templates"`
}

type serviceResource struct {
	client *iapi.Server
}

func (r *serviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (r *serviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "ServiceName",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hostname": schema.StringAttribute{
				Required:    true,
				Description: "Hostname",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"check_command": schema.StringAttribute{
				Required:    true,
				Description: "CheckCommand",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vars": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"templates": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Templates",
			},
		},
	}
}

func (r *serviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*iapi.Server)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *iapi.Server, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *serviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serviceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var templates []string
	if !plan.Templates.IsNull() && !plan.Templates.IsUnknown() {
		for _, template := range plan.Templates.Elements() {
			templates = append(templates, template.(types.String).ValueString())
		}
	}

	vars := make(map[string]string)
	if !plan.Vars.IsNull() && !plan.Vars.IsUnknown() {
		for key, value := range plan.Vars.Elements() {
			vars[key] = value.(types.String).ValueString()
		}
	}

	hostname := plan.Hostname.ValueString()
	name := plan.Name.ValueString()

	services, err := r.client.CreateService(name, hostname, plan.CheckCommand.ValueString(), vars, templates)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Service",
			"Could not create service unexpected error: "+err.Error(),
		)
		return
	}

	for _, service := range services {
		if service.Name == hostname+"!"+name {
			plan.ID = types.StringValue(hostname + "!" + name)
		}
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *serviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serviceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hostname := state.Hostname.ValueString()
	name := state.Name.ValueString()

	services, err := r.client.GetService(name, hostname)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Service",
			"Could not read service "+name+": "+err.Error(),
		)
		return
	}

	for _, service := range services {
		if service.Name == hostname+"!"+name {
			state.ID = types.StringValue(hostname + "!" + name)
			state.Hostname = types.StringValue(hostname)
			state.CheckCommand = types.StringValue(service.Attrs.CheckCommand)
			// keeping vars simple
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *serviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"Updates are currently not supported for service resources",
	)
}

func (r *serviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serviceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hostname := state.Hostname.ValueString()
	name := state.Name.ValueString()

	err := r.client.DeleteService(name, hostname)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Service",
			"Could not delete service, unexpected error: "+err.Error(),
		)
		return
	}
}
