package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEngineerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "devops-bootcamp_engineer" "test" {
	name  = "Bobby"
	email = "Bobby@gmail.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name
					resource.TestCheckResourceAttr("devops-bootcamp_engineer.test", "name", "Bobby"),
					// Verify email
					resource.TestCheckResourceAttr("devops-bootcamp_engineer.test", "email", "Bobby@gmail.com"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devops-bootcamp_engineer.test", "id")),
			},
			// Update and Read testing
			{

				Config: providerConfig + `
resource "devops-bootcamp_engineer" "test" {
	name  = "updatedBobby"
	email = "updatedBobby@gmail.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name/email updated
					resource.TestCheckResourceAttr("devops-bootcamp_engineer.test", "name", "updatedBobby"),
					resource.TestCheckResourceAttr("devops-bootcamp_engineer.test", "email", "updatedBobby@gmail.com"),
					resource.TestCheckResourceAttrSet("devops-bootcamp_engineer.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
