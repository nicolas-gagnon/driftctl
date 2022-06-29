package aws

import (
	"github.com/snyk/driftctl/enumeration/remote/aws/repository"
	remoteerror "github.com/snyk/driftctl/enumeration/remote/error"
	"github.com/snyk/driftctl/enumeration/resource"
	"github.com/snyk/driftctl/enumeration/resource/aws"
)

type CloudfrontDistributionEnumerator struct {
	repository repository.CloudfrontRepository
	factory    resource.ResourceFactory
}

func NewCloudfrontDistributionEnumerator(repo repository.CloudfrontRepository, factory resource.ResourceFactory) *CloudfrontDistributionEnumerator {
	return &CloudfrontDistributionEnumerator{
		repository: repo,
		factory:    factory,
	}
}

func (e *CloudfrontDistributionEnumerator) SupportedType() resource.ResourceType {
	return aws.AwsCloudfrontDistributionResourceType
}

func (e *CloudfrontDistributionEnumerator) Enumerate() ([]*resource.Resource, error) {
	distributions, err := e.repository.ListAllDistributions()
	if err != nil {
		return nil, remoteerror.NewResourceListingError(err, string(e.SupportedType()))
	}

	results := make([]*resource.Resource, 0, len(distributions))

	for _, distribution := range distributions {
		results = append(
			results,
			e.factory.CreateAbstractResource(
				string(e.SupportedType()),
				*distribution.Id,
				map[string]interface{}{},
			),
		)
	}

	return results, err
}
