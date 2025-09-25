package ops_test

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOpsResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Create and Read testing
            {
                Config: providerConfig + `
resource "dob_engineer" "e1" {
    name  = "Test Engineer 1"
    email = "testuser1@liatrio.com"
}

resource "dob_ops" "test" {
    name = "Test User 123"
    engineers = [dob_engineer.e1.id]
}
`,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("dob_ops.test", "name", "Test User 123"),
                    resource.TestCheckResourceAttr("dob_ops.test", "engineers.#", "1"),
                    resource.TestCheckResourceAttrSet("dob_ops.test", "engineers.0"),
                    resource.TestCheckResourceAttrSet("dob_ops.test", "id"),
                ),
            },
            // Update and Read testing
            {
                Config: providerConfig + `
resource "dob_engineer" "e1" {
    name  = "Test Engineer 1"
    email = "testuser1@liatrio.com"
}

resource "dob_engineer" "e2" {
    name  = "Test Engineer 2"
    email = "testuser2@liatrio.com"
}

resource "dob_ops" "test" {
    name = "Test User 123"
    engineers = [dob_engineer.e1.id, dob_engineer.e2.id]
}
`,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("dob_ops.test", "name", "Test User 123"),
                    resource.TestCheckResourceAttr("dob_ops.test", "engineers.#", "2"),
                    resource.TestCheckResourceAttrSet("dob_ops.test", "engineers.0"),
                    resource.TestCheckResourceAttrSet("dob_ops.test", "engineers.1"),
					resource.TestCheckResourceAttrSet("dob_ops.test", "id"),
                ),
            },
            // Delete testing automatically occurs in TestCase
        },
    })
}
