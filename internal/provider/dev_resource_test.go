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
resource "devops-bootcamp_dev" "test" {
	name  = "updatedBobby"
	engineers = [ {id = devops-bootcamp_engineer.test_engineer.id} ]
}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devops-bootcamp_dev.test", "name", "updatedBobby"),
					resource.TestCheckResourceAttrSet("devops-bootcamp_dev.test", "id"),
					resource.TestCheckResourceAttr("devops-bootcamp_dev.test", "engineers.#", "1"),
					resource.TestCheckResourceAttrPair("devops-bootcamp_dev.test", "engineers.0.id", "devops-bootcamp_engineer.test_engineer", "id"),
					resource.TestCheckResourceAttrPair("devops-bootcamp_dev.test", "engineers.0.name", "devops-bootcamp_engineer.test_engineer", "name"),
				),
			},
			// Add a second engineer
			{

				Config: providerConfig + `
resource "devops-bootcamp_engineer" "test_engineer" {
	name  = "Bobby"
	email = "bobby@bobby.com"
}
resource "devops-bootcamp_engineer" "test_engineer2" {
	name  = "BobbysBrother"
	email = "bobbysBrother@bobby.com"
}
resource "devops-bootcamp_dev" "test" {
	name  = "updatedBobby"
	engineers = [ {id = devops-bootcamp_engineer.test_engineer.id}, {id = devops-bootcamp_engineer.test_engineer2.id} ]
}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devops-bootcamp_dev.test", "name", "updatedBobby"),
					resource.TestCheckResourceAttrSet("devops-bootcamp_dev.test", "id"),
					resource.TestCheckResourceAttr("devops-bootcamp_dev.test", "engineers.#", "2"),
					resource.TestCheckResourceAttrPair("devops-bootcamp_dev.test", "engineers.0.id", "devops-bootcamp_engineer.test_engineer", "id"),
					resource.TestCheckResourceAttrPair("devops-bootcamp_dev.test", "engineers.0.name", "devops-bootcamp_engineer.test_engineer", "name"),
					resource.TestCheckResourceAttrPair("devops-bootcamp_dev.test", "engineers.1.id", "devops-bootcamp_engineer.test_engineer2", "id"),
					resource.TestCheckResourceAttrPair("devops-bootcamp_dev.test", "engineers.1.name", "devops-bootcamp_engineer.test_engineer2", "name"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
