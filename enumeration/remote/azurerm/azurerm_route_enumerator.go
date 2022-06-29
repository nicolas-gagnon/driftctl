package azurerm

import (
	"github.com/snyk/driftctl/enumeration/remote/azurerm/repository"
	remoteerror "github.com/snyk/driftctl/enumeration/remote/error"
	"github.com/snyk/driftctl/enumeration/resource"
	"github.com/snyk/driftctl/enumeration/resource/azurerm"
)

type AzurermRouteEnumerator struct {
	repository repository.NetworkRepository
	factory    resource.ResourceFactory
}

func NewAzurermRouteEnumerator(repo repository.NetworkRepository, factory resource.ResourceFactory) *AzurermRouteEnumerator {
	return &AzurermRouteEnumerator{
		repository: repo,
		factory:    factory,
	}
}

func (e *AzurermRouteEnumerator) SupportedType() resource.ResourceType {
	return azurerm.AzureRouteResourceType
}

func (e *AzurermRouteEnumerator) Enumerate() ([]*resource.Resource, error) {
	resources, err := e.repository.ListAllRouteTables()
	if err != nil {
		return nil, remoteerror.NewResourceListingErrorWithType(err, string(e.SupportedType()), azurerm.AzureRouteTableResourceType)
	}

	results := make([]*resource.Resource, 0, len(resources))

	for _, res := range resources {
		for _, route := range res.Properties.Routes {
			results = append(
				results,
				e.factory.CreateAbstractResource(
					string(e.SupportedType()),
					*route.ID,
					map[string]interface{}{
						"name":             *route.Name,
						"route_table_name": *res.Name,
					},
				),
			)
		}

	}

	return results, err
}
