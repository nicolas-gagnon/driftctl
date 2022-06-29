package azurerm

import "github.com/snyk/driftctl/enumeration/resource"

const AzurePostgresqlDatabaseResourceType = "azurerm_postgresql_database"

func initAzurePostgresqlDatabaseMetadata(resourceSchemaRepository resource.SchemaRepositoryInterface) {
	resourceSchemaRepository.SetHumanReadableAttributesFunc(AzurePostgresqlDatabaseResourceType, func(res *resource.Resource) map[string]string {
		val := res.Attrs
		attrs := make(map[string]string)
		if name := val.GetString("name"); name != nil && *name != "" {
			attrs["Name"] = *name
		}
		return attrs
	})
}
