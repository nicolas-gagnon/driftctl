package aws

import (
	"github.com/snyk/driftctl/enumeration/remote/aws/repository"
	remoteerror "github.com/snyk/driftctl/enumeration/remote/error"
	"github.com/snyk/driftctl/enumeration/resource"
	"github.com/snyk/driftctl/enumeration/resource/aws"
)

type EC2SubnetEnumerator struct {
	repository repository.EC2Repository
	factory    resource.ResourceFactory
}

func NewEC2SubnetEnumerator(repo repository.EC2Repository, factory resource.ResourceFactory) *EC2SubnetEnumerator {
	return &EC2SubnetEnumerator{
		repository: repo,
		factory:    factory,
	}
}

func (e *EC2SubnetEnumerator) SupportedType() resource.ResourceType {
	return aws.AwsSubnetResourceType
}

func (e *EC2SubnetEnumerator) Enumerate() ([]*resource.Resource, error) {
	subnets, _, err := e.repository.ListAllSubnets()
	if err != nil {
		return nil, remoteerror.NewResourceListingError(err, string(e.SupportedType()))
	}

	results := make([]*resource.Resource, 0, len(subnets))

	for _, subnet := range subnets {
		results = append(
			results,
			e.factory.CreateAbstractResource(
				string(e.SupportedType()),
				*subnet.SubnetId,
				map[string]interface{}{},
			),
		)
	}

	return results, err
}
