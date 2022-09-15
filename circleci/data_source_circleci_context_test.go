package circleci

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/healx/terraform-provider-circleci/circleci/template"
)

func TestAccCircleCIContextDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccOrgProviders,
		Steps: []resource.TestStep{
			{
				Config: template.ParseRandName(testAccCircleCIContextDataSource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.circleci_context.foo", "name", regexp.MustCompile("^terraform-test")),
				),
			},
		},
	})
}

const testAccCircleCIContextDataSource = `
resource "circleci_context" "foo" {
  name = "terraform-test-{{.randName}}"
}

data "circleci_context" "foo" {
  name = circleci_context.foo.name
}
`
