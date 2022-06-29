package aws

import (
	"github.com/snyk/driftctl/enumeration/remote/aws/repository"
	remoteerror "github.com/snyk/driftctl/enumeration/remote/error"
	"github.com/snyk/driftctl/enumeration/resource"
	"github.com/snyk/driftctl/enumeration/resource/aws"
)

type IamUserPolicyEnumerator struct {
	repository repository.IAMRepository
	factory    resource.ResourceFactory
}

func NewIamUserPolicyEnumerator(repo repository.IAMRepository, factory resource.ResourceFactory) *IamUserPolicyEnumerator {
	return &IamUserPolicyEnumerator{
		repository: repo,
		factory:    factory,
	}
}

func (e *IamUserPolicyEnumerator) SupportedType() resource.ResourceType {
	return aws.AwsIamUserPolicyResourceType
}

func (e *IamUserPolicyEnumerator) Enumerate() ([]*resource.Resource, error) {
	users, err := e.repository.ListAllUsers()
	if err != nil {
		return nil, remoteerror.NewResourceListingErrorWithType(err, string(e.SupportedType()), aws.AwsIamUserResourceType)
	}
	userPolicies, err := e.repository.ListAllUserPolicies(users)
	if err != nil {
		return nil, remoteerror.NewResourceListingError(err, string(e.SupportedType()))
	}

	results := make([]*resource.Resource, 0, len(userPolicies))

	for _, userPolicy := range userPolicies {
		results = append(
			results,
			e.factory.CreateAbstractResource(
				string(e.SupportedType()),
				userPolicy,
				map[string]interface{}{},
			),
		)
	}

	return results, err
}
