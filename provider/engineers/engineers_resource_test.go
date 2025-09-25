package engineers_test

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEngineerResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Create and Read testing
            {
                Config: providerConfig + `
resource "dob_engineer" "test" {
    name = "Test User 123"
    email = "testuser123@liatrio.com"
}
`,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("dob_engineer.test", "name", "Test User 123"),
                    resource.TestCheckResourceAttr("dob_engineer.test", "email", "testuser123@liatrio.com"),
                    resource.TestCheckResourceAttrSet("dob_engineer.test", "id"),
                ),
            },
            // Update and Read testing
            {
                Config: providerConfig + `
resource "dob_engineer" "test" {
    name = "Test User 123"
    email = "testuser123@liatrio.com"
}
`,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("dob_engineer.test", "name", "Test User 123"),
                    resource.TestCheckResourceAttr("dob_engineer.test", "email", "testuser123@liatrio.com"),
					resource.TestCheckResourceAttrSet("dob_engineer.test", "id"),
                ),
            },
            // Delete testing automatically occurs in TestCase
        },
    })
}
