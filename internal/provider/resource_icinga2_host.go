package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/lrsmith/go-icinga2-api/iapi"
)

var (
	_ resource.Resource                = &hostResource{}
	_ resource.ResourceWithConfigure   = &hostResource{}
	_ resource.ResourceWithImportState = &hostResource{}
)

func Host() resource.Resource {
	return &hostResource{}
}

type hostResourceModel struct {
	ID           types.String `tfsdk:"id"`
	LastUpdated  types.String `tfsdk:"last_updated"`
	Hostname     types.String `tfsdk:"hostname"`
	Address      types.String `tfsdk:"address"`
	CheckCommand types.String `tfsdk:"check_command"`
	Groups       types.List   `tfsdk:"groups"`
	Vars         types.Map    `tfsdk:"vars"`
	Templates    types.List   `tfsdk:"templates"`
}

type hostResource struct {
	client *iapi.Server
}

func (r *hostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (r *hostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"hostname": schema.StringAttribute{
				Required:    true,
				Description: "Hostname",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"address": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"check_command": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"groups": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"vars": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"templates": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *hostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *hostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan hostResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var groups []string
	if !plan.Groups.IsNull() && !plan.Groups.IsUnknown() {
		for _, group := range plan.Groups.Elements() {
			if strVal, ok := group.(types.String); ok {
				groups = append(groups, strVal.ValueString())
			} else {
				resp.Diagnostics.AddError(
					"Error creating Host",
					"Group not a string",
				)
			}
		}
	}

	vars := make(map[string]interface{})
	if !plan.Vars.IsNull() && !plan.Vars.IsUnknown() {
		for key, value := range plan.Vars.Elements() {
			if strVal, ok := value.(types.String); ok {
				vars[key] = strVal.ValueString()
			} else {
				resp.Diagnostics.AddError(
					"Error creating Host",
					"Variable not a string",
				)
			}
		}
	}

	var templates []string
	if !plan.Templates.IsNull() && !plan.Templates.IsUnknown() {
		for _, template := range plan.Templates.Elements() {
			if strVal, ok := template.(types.String); ok {
				templates = append(templates, strVal.ValueString())
			} else {
				resp.Diagnostics.AddError(
					"Error creating Host",
					"Template not a string",
				)
			}
		}
	}

	hosts, err := r.client.CreateHost(plan.Hostname.ValueString(), plan.Address.ValueString(), "", plan.CheckCommand.ValueString(), vars, templates, groups, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Host",
			"Could not create host unexpected error: "+err.Error(),
		)
		return
	}

	for _, host := range hosts {
		if host.Name == plan.Hostname.ValueString() {
			plan.ID = types.StringValue(host.Name)
		}
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *hostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state hostResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hosts, err := r.client.GetHost(state.Hostname.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Host",
			"Could not read host "+state.Hostname.ValueString()+": "+err.Error(),
		)
		return
	}

	for _, host := range hosts {
		if host.Name == state.Hostname.ValueString() {
			state.ID = types.StringValue(host.Name)
			state.Hostname = types.StringValue(host.Name)
			state.Address = types.StringValue(host.Attrs.Address)
			state.CheckCommand = types.StringValue(host.Attrs.CheckCommand)

			// Note: We might need to map vars back to state correctly for lists/maps. For simplicity keeping it string mapped to attributes if they existed directly.
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *hostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"Updates are currently not supported for host resources",
	)
}

func (r *hostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state hostResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteHost(state.Hostname.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Host",
			"Could not delete host, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *hostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	hosts, err := r.client.GetHost(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Host",
			"Could not read host "+req.ID+": "+err.Error(),
		)
		return
	}

	for _, host := range hosts {
		if host.Name == req.ID {
			resource.ImportStatePassthroughID(ctx, path.Root("hostname"), req, resp)
			return
		}
	}

	resp.Diagnostics.AddError(
		"Error Importing Host",
		"Host "+req.ID+" does not exist",
	)
}
