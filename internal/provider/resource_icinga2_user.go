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
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func User() resource.Resource {
	return &userResource{}
}

type userResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Name        types.String `tfsdk:"name"`
	Email       types.String `tfsdk:"email"`
	Vars        types.Map    `tfsdk:"vars"`
}

type userResource struct {
	client *iapi.Server
}

func (r *userResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Description: "Username",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"vars": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *userResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	email := plan.Email.ValueString()

	vars := make(map[string]string)
	if !plan.Vars.IsNull() && !plan.Vars.IsUnknown() {
		for key, value := range plan.Vars.Elements() {
			if strVal, ok := value.(types.String); ok {
				vars[key] = strVal.ValueString()
			} else {
				resp.Diagnostics.AddError(
					"Error creating User",
					"Variable not a string",
				)
			}
		}
	}

	users, err := r.client.CreateUser(ctx, name, email, vars)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating User",
			"Could not create user unexpected error: "+err.Error(),
		)
		return
	}

	for _, user := range users {
		if user.Name == name {
			plan.ID = types.StringValue(user.Name)
			plan.Name = types.StringValue(user.Name)
			plan.Email = types.StringValue(user.Attrs.Email)
		}
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name.ValueString()

	users, err := r.client.GetUser(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "No objects found.") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading User",
			"Could not read user "+name+": "+err.Error(),
		)
		return
	}

	found := false
	for _, user := range users {
		if user.Name == name {
			found = true
			state.ID = types.StringValue(user.Name)
			state.Name = types.StringValue(user.Name)
			state.Email = types.StringValue(user.Attrs.Email)
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

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state userResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID
	if plan.Email.IsNull() || plan.Email.IsUnknown() {
		plan.Email = state.Email
	}

	vars := make(map[string]string)
	if !plan.Vars.IsNull() && !plan.Vars.IsUnknown() {
		for key, value := range plan.Vars.Elements() {
			if strVal, ok := value.(types.String); ok {
				vars[key] = strVal.ValueString()
			}
		}
	}

	attrs := iapi.UserAttrs{
		Email: plan.Email.ValueString(),
		Vars:  vars,
	}

	_, err := r.client.UpdateUser(ctx, plan.Name.ValueString(), attrs)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating User",
			"Could not update user unexpected error: "+err.Error(),
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

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name.ValueString()

	err := r.client.DeleteUser(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "No objects found.") {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting User",
			"Could not delete user, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
