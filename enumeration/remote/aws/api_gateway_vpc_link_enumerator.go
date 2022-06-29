package aws

import (
	"github.com/snyk/driftctl/enumeration/remote/aws/repository"
	remoteerror "github.com/snyk/driftctl/enumeration/remote/error"
	"github.com/snyk/driftctl/enumeration/resource"
	"github.com/snyk/driftctl/enumeration/resource/aws"
)

type ApiGatewayVpcLinkEnumerator struct {
	repository repository.ApiGatewayRepository
	factory    resource.ResourceFactory
}

func NewApiGatewayVpcLinkEnumerator(repo repository.ApiGatewayRepository, factory resource.ResourceFactory) *ApiGatewayVpcLinkEnumerator {
	return &ApiGatewayVpcLinkEnumerator{
		repository: repo,
		factory:    factory,
	}
}

func (e *ApiGatewayVpcLinkEnumerator) SupportedType() resource.ResourceType {
	return aws.AwsApiGatewayVpcLinkResourceType
}

func (e *ApiGatewayVpcLinkEnumerator) Enumerate() ([]*resource.Resource, error) {
	vpcLinks, err := e.repository.ListAllVpcLinks()
	if err != nil {
		return nil, remoteerror.NewResourceListingError(err, string(e.SupportedType()))
	}

	results := make([]*resource.Resource, 0, len(vpcLinks))

	for _, vpcLink := range vpcLinks {
		results = append(
			results,
			e.factory.CreateAbstractResource(
				string(e.SupportedType()),
				*vpcLink.Id,
				map[string]interface{}{},
			),
		)
	}

	return results, err
}
