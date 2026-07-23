package provider

import (
	"context"
	"fmt"
	"strings"
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
	_ resource.Resource                = &checkCommandResource{}
	_ resource.ResourceWithConfigure   = &checkCommandResource{}
	_ resource.ResourceWithImportState = &checkCommandResource{}
)

func CheckCommand() resource.Resource {
	return &checkCommandResource{}
}

type checkCommandResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Name        types.String `tfsdk:"name"`
	Command     types.String `tfsdk:"command"`
	Templates   types.List   `tfsdk:"templates"`
	Arguments   types.Map    `tfsdk:"arguments"`
}

type checkCommandResource struct {
	client *iapi.Server
}

func (r *checkCommandResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_checkcommand"
}

func (r *checkCommandResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Description: "Name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"command": schema.StringAttribute{
				Required: true,
			},
			"templates": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"arguments": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *checkCommandResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *checkCommandResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan checkCommandResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	arguments := make(map[string]string)
	if !plan.Arguments.IsNull() && !plan.Arguments.IsUnknown() {
		for key, value := range plan.Arguments.Elements() {
			if strVal, ok := value.(types.String); ok {
				arguments[key] = strVal.ValueString()
			} else {
				resp.Diagnostics.AddError(
					"Error creating CheckCommand",
					"Argument not a string",
				)
			}
		}
	}

	checkcommands, err := r.client.CreateCheckcommand(ctx, plan.Name.ValueString(), plan.Command.ValueString(), arguments)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Checkcommand",
			"Could not create checkcommand unexpected error: "+err.Error(),
		)
		return
	}

	for _, checkcommand := range checkcommands {
		if checkcommand.Name == plan.Name.ValueString() {
			plan.ID = types.StringValue(checkcommand.Name)
		}
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *checkCommandResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state checkCommandResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkcommands, err := r.client.GetCheckcommand(ctx, state.Name.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "No objects found.") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Checkcommand",
			"Could not read checkcommand "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	found := false
	for _, checkcommand := range checkcommands {
		if checkcommand.Name == state.Name.ValueString() {
			found = true
			state.ID = types.StringValue(checkcommand.Name)
			if len(checkcommand.Attrs.Command) > 0 {
				state.Command = types.StringValue(checkcommand.Attrs.Command[0])
			}
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *checkCommandResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan checkCommandResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state checkCommandResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID

	arguments := make(map[string]string)
	if !plan.Arguments.IsNull() && !plan.Arguments.IsUnknown() {
		for key, value := range plan.Arguments.Elements() {
			if strVal, ok := value.(types.String); ok {
				arguments[key] = strVal.ValueString()
			}
		}
	}

	var templates []string
	if !plan.Templates.IsNull() && !plan.Templates.IsUnknown() {
		for _, template := range plan.Templates.Elements() {
			if strVal, ok := template.(types.String); ok {
				templates = append(templates, strVal.ValueString())
			}
		}
	}

	attrs := iapi.CheckcommandAttrs{
		Command:   []string{plan.Command.ValueString()},
		Arguments: arguments,
		Templates: templates,
	}

	_, err := r.client.UpdateCheckcommand(ctx, plan.Name.ValueString(), attrs)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Checkcommand",
			"Could not update checkcommand unexpected error: "+err.Error(),
		)
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *checkCommandResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state checkCommandResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheckcommand(ctx, state.Name.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "No objects found.") {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Checkcommand",
			"Could not delete checkcommand, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *checkCommandResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
