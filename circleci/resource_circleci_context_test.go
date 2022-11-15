package circleci

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/healx/circleci-cli/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	client "github.com/healx/terraform-provider-circleci/circleci/client"
	"github.com/healx/terraform-provider-circleci/circleci/template"
)

func TestAccCircleCIContext_basic(t *testing.T) {
	context := &api.Context{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccOrgProviders,
		CheckDestroy: testAccCheckCircleCIContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: template.ParseRandName(testAccCircleCIContext_basic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCircleCIContextExists("circleci_context.foo", context),
					testAccCheckCircleCIContextAttributes_basic(context),
					resource.TestMatchResourceAttr("circleci_context.foo", "name", regexp.MustCompile("^terraform-test")),
				),
			},
		},
	})
}

func TestAccCircleCIContext_update(t *testing.T) {
	context := &api.Context{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccOrgProviders,
		CheckDestroy: testAccCheckCircleCIContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: template.ParseRandName(testAccCircleCIContext_basic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCircleCIContextExists("circleci_context.foo", context),
					testAccCheckCircleCIContextAttributes_basic(context),
					resource.TestMatchResourceAttr("circleci_context.foo", "name", regexp.MustCompile("^terraform-test")),
				),
			},
			{
				Config: template.ParseRandName(testAccCircleCIContext_update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCircleCIContextExists("circleci_context.foo", context),
					testAccCheckCircleCIContextAttributes_update(context),
					resource.TestMatchResourceAttr("circleci_context.foo", "name", regexp.MustCompile("^terraform-test-updated")),
				),
			},
		},
	})
}

func TestAccCircleCIContext_import(t *testing.T) {
	context := &api.Context{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccOrgProviders,
		CheckDestroy: testAccCheckCircleCIContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: template.ParseRandName(testAccCircleCIContext_basic),
				Check:  testAccCheckCircleCIContextExists("circleci_context.foo", context),
			},
			{
				ResourceName: "circleci_context.foo",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					org, err := testAccOrgProvider.Meta().(*client.Client).Organization("")
					if err != nil {
						return "", err
					}

					return fmt.Sprintf(
						"%s/%s",
						org,
						context.ID,
					), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCircleCIContext_import_name(t *testing.T) {
	context := &api.Context{}
	cfg, randName := template.ParseRandNameAndReturn(testAccCircleCIContext_basic)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccOrgProviders,
		CheckDestroy: testAccCheckCircleCIContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check:  testAccCheckCircleCIContextExists("circleci_context.foo", context),
			},
			{
				ResourceName: "circleci_context.foo",
				ImportStateId: fmt.Sprintf(
					"%s/%s",
					os.Getenv("TEST_CIRCLECI_ORGANIZATION"),
					fmt.Sprintf("terraform-test-%s", randName),
				),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCircleCIContextExists(addr string, context *api.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccOrgProvider.Meta().(*client.Client)

		resource, ok := s.RootModule().Resources[addr]
		if !ok {
			return fmt.Errorf("Not found: %s", addr)
		}
		if resource.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		ctx, err := c.GetContext(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting context: %w", err)
		}

		*context = *ctx

		return nil
	}
}

func testAccCheckCircleCIContextDestroy(s *terraform.State) error {
	c := testAccOrgProvider.Meta().(*client.Client)

	for _, resource := range s.RootModule().Resources {
		if resource.Type != "circleci_context" {
			continue
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		_, err := c.GetContext(resource.Primary.ID)
		if err == nil {
			return fmt.Errorf("Context %s still exists: %w", resource.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckCircleCIContextAttributes_basic(context *api.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !regexp.MustCompile("^terraform-test").MatchString(context.Name) {
			return fmt.Errorf("Unexpected context name: %s", context.Name)
		}

		return nil
	}
}

func testAccCheckCircleCIContextAttributes_update(context *api.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !regexp.MustCompile("^terraform-test-updated").MatchString(context.Name) {
			return fmt.Errorf("Unexpected context name: %s", context.Name)
		}

		return nil
	}
}

const testAccCircleCIContext_basic = `
resource "circleci_context" "foo" {
	name = "terraform-test-{{.randName}}"
}
`

const testAccCircleCIContext_update = `
resource "circleci_context" "foo" {
	name = "terraform-test-updated-{{.randName}}"
}
`
