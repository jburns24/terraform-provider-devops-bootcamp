package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// EngineerModel describes the data source data model.
type EngineerModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Email       types.String `tfsdk:"email"`
	LastUpdated types.String `tfsdk:"last_updated"`
}
