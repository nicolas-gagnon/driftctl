package aws

import (
	"github.com/snyk/driftctl/enumeration/remote/aws/repository"
	remoteerror "github.com/snyk/driftctl/enumeration/remote/error"
	"github.com/snyk/driftctl/enumeration/resource"
	resourceaws "github.com/snyk/driftctl/enumeration/resource/aws"

	"github.com/aws/aws-sdk-go/aws"
)

type VPCDefaultSecurityGroupEnumerator struct {
	repository repository.EC2Repository
	factory    resource.ResourceFactory
}

func NewVPCDefaultSecurityGroupEnumerator(repo repository.EC2Repository, factory resource.ResourceFactory) *VPCDefaultSecurityGroupEnumerator {
	return &VPCDefaultSecurityGroupEnumerator{
		repo,
		factory,
	}
}

func (e *VPCDefaultSecurityGroupEnumerator) SupportedType() resource.ResourceType {
	return resourceaws.AwsDefaultSecurityGroupResourceType
}

func (e *VPCDefaultSecurityGroupEnumerator) Enumerate() ([]*resource.Resource, error) {
	_, defaultSecurityGroups, err := e.repository.ListAllSecurityGroups()
	if err != nil {
		return nil, remoteerror.NewResourceListingError(err, string(e.SupportedType()))
	}

	results := make([]*resource.Resource, 0, len(defaultSecurityGroups))

	for _, item := range defaultSecurityGroups {
		results = append(
			results,
			e.factory.CreateAbstractResource(
				string(e.SupportedType()),
				aws.StringValue(item.GroupId),
				map[string]interface{}{},
			),
		)
	}

	return results, nil
}
