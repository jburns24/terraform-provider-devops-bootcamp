// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccEngineerDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.devops-bootcamp_engineer.test", "name", "Ryan"),
					resource.TestCheckResourceAttrSet("data.devops-bootcamp_engineer.test", "id"),
				),
			},
		},
	})
}

const testAccEngineerDataSourceConfig = providerConfig + `

	resource "devops-bootcamp_engineer" "test" {
		name  = "Ryan"
		email = "Ryan@gmail.com"
	}

	data "devops-bootcamp_engineer" "test" {
	  name = devops-bootcamp_engineer.test.name
	}

`
