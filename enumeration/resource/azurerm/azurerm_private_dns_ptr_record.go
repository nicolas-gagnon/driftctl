package azurerm

import (
	"github.com/snyk/driftctl/enumeration/resource"
)

const AzurePrivateDNSPTRRecordResourceType = "azurerm_private_dns_ptr_record"

func initAzurePrivateDNSPTRRecordMetaData(resourceSchemaRepository resource.SchemaRepositoryInterface) {
	resourceSchemaRepository.SetHumanReadableAttributesFunc(AzurePrivateDNSPTRRecordResourceType, func(res *resource.Resource) map[string]string {
		val := res.Attrs
		attrs := make(map[string]string)
		if name := val.GetString("name"); name != nil && *name != "" {
			attrs["Name"] = *name
		}
		if zone := val.GetString("zone_name"); zone != nil && *zone != "" {
			attrs["Zone"] = *zone
		}
		return attrs
	})
	resourceSchemaRepository.SetFlags(AzurePrivateDNSPTRRecordResourceType, resource.FlagDeepMode)
}
