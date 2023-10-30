package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDevResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "devops-bootcamp_dev" "test" {
	name  = "Bobby"
	engineers = []
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name
					resource.TestCheckResourceAttr("devops-bootcamp_dev.test", "name", "Bobby"),
					// Verify email
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devops-bootcamp_dev.test", "id")),
			},
			// Update and Read testing
			{

				Config: providerConfig + `
resource "devops-bootcamp_dev" "test" {
	name  = "updatedBobby"
	engineers = []
}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devops-bootcamp_dev.test", "name", "updatedBobby"),
					resource.TestCheckResourceAttrSet("devops-bootcamp_dev.test", "id"),
				),
			},
			// Add a new engineer
			{

				Config: providerConfig + `
resource "devops-bootcamp_engineer" "test_engineer" {
	name  = "Bobby"
	email = "bobby@bobby.com"
}
resource "devops-bootcamp_dev" "test_dev" {
	name  = "Bangle"
	engineers = [ {id = devops-bootcamp_engineer.test_engineer.id} ]
}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devops-bootcamp_dev.test_dev", "name", "Bangle"),
					resource.TestCheckResourceAttrSet("devops-bootcamp_dev.test_dev", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
