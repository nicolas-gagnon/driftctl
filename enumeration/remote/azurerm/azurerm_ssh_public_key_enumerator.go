package azurerm

import (
	"github.com/snyk/driftctl/enumeration/remote/azurerm/repository"
	remoteerror "github.com/snyk/driftctl/enumeration/remote/error"
	"github.com/snyk/driftctl/enumeration/resource"
	"github.com/snyk/driftctl/enumeration/resource/azurerm"
)

type AzurermSSHPublicKeyEnumerator struct {
	repository repository.ComputeRepository
	factory    resource.ResourceFactory
}

func NewAzurermSSHPublicKeyEnumerator(repo repository.ComputeRepository, factory resource.ResourceFactory) *AzurermSSHPublicKeyEnumerator {
	return &AzurermSSHPublicKeyEnumerator{
		repository: repo,
		factory:    factory,
	}
}

func (e *AzurermSSHPublicKeyEnumerator) SupportedType() resource.ResourceType {
	return azurerm.AzureSSHPublicKeyResourceType
}

func (e *AzurermSSHPublicKeyEnumerator) Enumerate() ([]*resource.Resource, error) {
	keys, err := e.repository.ListAllSSHPublicKeys()
	if err != nil {
		return nil, remoteerror.NewResourceListingError(err, string(e.SupportedType()))
	}

	results := make([]*resource.Resource, 0, len(keys))

	for _, res := range keys {
		results = append(
			results,
			e.factory.CreateAbstractResource(
				string(e.SupportedType()),
				*res.ID,
				map[string]interface{}{
					"name": *res.Name,
				},
			),
		)
	}

	return results, err
}
