package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// DevModel describes the data source data model.
type DevModel struct {
	Id        types.String    `tfsdk:"id"`
	Name      types.String    `tfsdk:"name"`
	Engineers []EngineerModel `tfsdk:"engineers"`
}
