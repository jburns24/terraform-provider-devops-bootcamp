// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	devops_resource "github.com/liatrio/devops-bootcamp/examples/ch6/devops-resources"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DevResource{}
var _ resource.ResourceWithImportState = &DevResource{}

func NewDevResource() resource.Resource {
	return &DevResource{}
}

// DevResource defines the resource implementation.
type DevResource struct {
	client *Client
}

// DevResourceModel describes the resource data model.
type DevResourceModel struct {
	Id          types.String    `tfsdk:"id"`
	Name        types.String    `tfsdk:"name"`
	Engineers   []EngineerModel `tfsdk:"engineers"`
	LastUpdated types.String    `tfsdk:"last_updated"`
}

func (r *DevResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dev"
}

func (r *DevResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Developer group resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the developer group",
				Required:            true,
			},
			"engineers": schema.ListNestedAttribute{
				MarkdownDescription: "List of engineers in the developer group by id",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"email": schema.StringAttribute{
							Computed: true,
						},
						"last_updated": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Dev identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *DevResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DevResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planned *DevResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planned)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var reqObj devops_resource.Dev
	reqObj.Name = planned.Name.ValueString()
	for _, engineer := range planned.Engineers {
		reqObj.Engineers = append(reqObj.Engineers, &devops_resource.Engineer{
			Id: engineer.Id.ValueString(),
			// Name:  engineer.Name.ValueString(),
			// Email: engineer.Email.ValueString(),
		})
	}

	// Make empty list if no engineers are provided
	if reqObj.Engineers == nil {
		reqObj.Engineers = make([]*devops_resource.Engineer, 0)
	}

	dev, err := r.client.CreateDev(&reqObj)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating dev",
			"Could not create dev,unexpected error:"+err.Error(),
		)
		return
	}

	// Map the response to the planned model
	planned.Id = types.StringValue(dev.Id)
	planned.Name = types.StringValue(dev.Name)
	planned.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	planned.Engineers = []EngineerModel{}

	for _, engineer := range dev.Engineers {
		planned.Engineers = append(planned.Engineers, EngineerModel{
			Id:    types.StringValue(engineer.Id),
			Name:  types.StringValue(engineer.Name),
			Email: types.StringValue(engineer.Email),
		})
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a dev resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &planned)...)
}

func (r *DevResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *DevResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dev, err := r.client.GetDevById(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	state.Id = types.StringValue(dev.Id)
	state.Name = types.StringValue(dev.Name)

	state.Engineers = []EngineerModel{}
	for _, engineer := range dev.Engineers {
		state.Engineers = append(state.Engineers, EngineerModel{
			Id:    types.StringValue(engineer.Id),
			Name:  types.StringValue(engineer.Name),
			Email: types.StringValue(engineer.Email),
		})
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("read resource for id:%s", state.Id))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DevResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planned *DevResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planned)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var reqObj devops_resource.Dev
	reqObj.Id = planned.Id.ValueString()
	reqObj.Name = planned.Name.ValueString()
	for _, engineer := range planned.Engineers {
		reqObj.Engineers = append(reqObj.Engineers, &devops_resource.Engineer{
			Id:    engineer.Id.ValueString(),
			Name:  engineer.Name.ValueString(),
			Email: engineer.Email.ValueString(),
		})
	}

	// update the dev
	_, err := r.client.UpdateDev(&reqObj)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating dev",
			"Could not update dev,unexpected error:"+err.Error(),
		)
		return
	}

	// Fetch the updated dev from the API
	dev, err := r.client.GetDevById(planned.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	// Update the planned model with the updated dev
	planned.Id = types.StringValue(dev.Id)
	planned.Name = types.StringValue(dev.Name)
	planned.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	planned.Engineers = []EngineerModel{}
	for _, engineer := range dev.Engineers {
		planned.Engineers = append(planned.Engineers, EngineerModel{
			Id:    types.StringValue(engineer.Id),
			Name:  types.StringValue(engineer.Name),
			Email: types.StringValue(engineer.Email),
		})
	}

	tflog.Trace(ctx, "updated a dev resource")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &planned)...)
}

func (r *DevResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *DevResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dev := devops_resource.Dev{
		Id:   state.Id.ValueString(),
		Name: state.Name.ValueString(),
		// Engineers: make([]*devops_resource.Engineer, 0),
	}

	err := r.client.DeleteDev(&dev)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Dev",
			"Could not delete dev, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *DevResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
