// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDevDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccDevDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.devops-bootcamp_dev.test", "name", "Ryan"),
					resource.TestCheckResourceAttrSet("data.devops-bootcamp_dev.test", "id"),
				),
			},
		},
	})
}

const testAccDevDataSourceConfig = providerConfig + `

	resource "devops-bootcamp_dev" "test" {
		name  = "Ryan"
	}

	data "devops-bootcamp_dev" "test" {
	  name = devops-bootcamp_dev.test.name
	}

`

// const testAccDevDataSourceConfig = providerConfig + `
// 	data "devops-bootcamp_dev" "test" {
// 	  name = "dev_bengal"
// 	}

// `
