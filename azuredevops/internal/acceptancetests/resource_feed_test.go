//go:build (all || core || data_sources || data_feed) && (!data_sources || !exclude_feed)
// +build all core data_sources data_feed
// +build !data_sources !exclude_feed

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccFeed_basic(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfNode := "azuredevops_feed.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclFeedBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckNoResourceAttr(tfNode, "project"),
				),
			},
		},
	})
}

func TestAccFeed_with_Project(t *testing.T) {
	name := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()

	tfNode := "azuredevops_feed.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclFeedWithProject(projectName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
				),
			},
		},
	})
}

func TestAccFeed_softDeleteRecovery(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfNode := "azuredevops_feed.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:  hclFeedBasic(name),
				Destroy: true,
			},
			{
				Config: hclFeedRestore(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckNoResourceAttr(tfNode, "project"),
				),
			},
		},
	})
}

func hclFeedBasic(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_feed" "test" {
  name = "%s"
}`, name)
}

func hclFeedWithProject(projectName, name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "%[1]s-description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_feed" "test" {
  name       = "%[1]s"
  project_id = azuredevops_project.test.id
}`, projectName, name)
}

func hclFeedRestore(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_feed" "test" {
  name = "%s"
  features {
    restore = true
  }
}
`, name)
}
