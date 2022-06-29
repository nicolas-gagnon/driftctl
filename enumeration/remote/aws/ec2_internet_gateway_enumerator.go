package aws

import (
	"github.com/snyk/driftctl/enumeration/remote/aws/repository"
	remoteerror "github.com/snyk/driftctl/enumeration/remote/error"
	"github.com/snyk/driftctl/enumeration/resource"
	"github.com/snyk/driftctl/enumeration/resource/aws"
)

type EC2InternetGatewayEnumerator struct {
	repository repository.EC2Repository
	factory    resource.ResourceFactory
}

func NewEC2InternetGatewayEnumerator(repo repository.EC2Repository, factory resource.ResourceFactory) *EC2InternetGatewayEnumerator {
	return &EC2InternetGatewayEnumerator{
		repository: repo,
		factory:    factory,
	}
}

func (e *EC2InternetGatewayEnumerator) SupportedType() resource.ResourceType {
	return aws.AwsInternetGatewayResourceType
}

func (e *EC2InternetGatewayEnumerator) Enumerate() ([]*resource.Resource, error) {
	internetGateways, err := e.repository.ListAllInternetGateways()
	if err != nil {
		return nil, remoteerror.NewResourceListingError(err, string(e.SupportedType()))
	}

	results := make([]*resource.Resource, 0, len(internetGateways))

	for _, internetGateway := range internetGateways {
		data := map[string]interface{}{}
		if len(internetGateway.Attachments) > 0 && internetGateway.Attachments[0].VpcId != nil {
			data["vpc_id"] = *internetGateway.Attachments[0].VpcId
		}
		results = append(
			results,
			e.factory.CreateAbstractResource(
				string(e.SupportedType()),
				*internetGateway.InternetGatewayId,
				data,
			),
		)
	}

	return results, err
}
