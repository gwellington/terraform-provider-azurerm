package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PostgresqlFlexibleServerDatabaseResource struct {
}

func TestAccPostgresqlFlexibleServerDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_database", "test")
	r := PostgresqlFlexibleServerDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("charset").HasValue("UTF8"),
				check.That(data.ResourceName).Key("collation").HasValue("en_US.UTF8"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPostgresqlFlexibleServerDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_database", "test")
	r := PostgresqlFlexibleServerDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccPostgresqlFlexibleServerDatabase_charsetLowercase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_database", "test")
	r := PostgresqlFlexibleServerDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.charsetLowercase(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("charset").HasValue("LATIN1"),
				check.That(data.ResourceName).Key("collation").HasValue("en_US.latin1"),
			),
		},
		data.ImportStep(),
	})
}

func (PostgresqlFlexibleServerDatabaseResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FlexibleServerDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.FlexibleServerDatabaseClient.Get(ctx, id.ResourceGroup, id.FlexibleServerName, id.DatabaseName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.DatabaseProperties != nil), nil
}

func (r PostgresqlFlexibleServerDatabaseResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server_database" "import" {
  name      = azurerm_postgresql_flexible_server_database.test.name
  server_id = azurerm_postgresql_flexible_server_database.test.server_id
  collation = azurerm_postgresql_flexible_server_database.test.collation
  charset   = azurerm_postgresql_flexible_server_database.test.charset
}
`, r.basic(data))
}

func (PostgresqlFlexibleServerDatabaseResource) charsetLowercase(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server_database" "test" {
  name      = "acctest-fsd-%d"
  server_id = azurerm_postgresql_flexible_server.test.id
  collation = "en_US.latin1"
  charset   = "latin1"
}
`, PostgresqlFlexibleServerResource{}.basic(data), data.RandomInteger)
}

func (PostgresqlFlexibleServerDatabaseResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server_database" "test" {
  name      = "acctest-fsd-%d"
  server_id = azurerm_postgresql_flexible_server.test.id
  collation = "en_US.UTF8"
  charset   = "UTF8"
}
`, PostgresqlFlexibleServerResource{}.basic(data), data.RandomInteger)
}
