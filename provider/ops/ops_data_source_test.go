package ops_test

import (
	"terraform-provider-devops/provider"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Local acceptance test helpers for this package. We cannot rely on unexported
// variables from provider/*_test.go across packages, since those are not
// compiled into the provider package for import. Define local equivalents here.

const providerConfig = `
provider "dob" {
    endpoint = "http://localhost:8080"
}
`

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"dob": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func TestAccOpsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "dob_ops" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dob_ops.test", "ops.0.name", "Ops Team #1"),
					resource.TestCheckResourceAttr("data.dob_ops.test", "ops.0.id", "YVDOG"),
					resource.TestCheckResourceAttr("data.dob_ops.test", "ops.0.engineers.0", "GRESC"),
					resource.TestCheckResourceAttr("data.dob_ops.test", "ops.0.engineers.1", "5LE5Z"),
				),
			},
		},
	})

}
