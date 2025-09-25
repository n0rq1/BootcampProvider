package devops_test

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDevOpsResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Create and Read testing
            {
                Config: providerConfig + `
resource "dob_engineer" "e1" {
    name  = "Test Engineer 123"
    email = "testuser1@liatrio.com"
}

resource "dob_engineer" "e2" {
    name  = "Test Engineer 123"
    email = "testuser2@liatrio.com"
}

resource "dob_dev" "test" {
	name = "Test Dev #123"
    engineers = [dob_engineer.e1.id]
}

resource "dob_ops" "test" {
	name = "Test Ops #123"
    engineers = [dob_engineer.e2.id]
}

resource "dob_devops" "test" {
    devs = [dob_dev.test.id]
    ops = [dob_ops.test.id]
}
`,
                Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dob_devops.test", "id"),
                ),
            },
        },
    })
}
