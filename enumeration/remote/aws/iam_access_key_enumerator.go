package aws

import (
	"github.com/snyk/driftctl/enumeration/remote/aws/repository"
	remoteerror "github.com/snyk/driftctl/enumeration/remote/error"
	"github.com/snyk/driftctl/enumeration/resource"
	resourceaws "github.com/snyk/driftctl/enumeration/resource/aws"
)

type IamAccessKeyEnumerator struct {
	repository repository.IAMRepository
	factory    resource.ResourceFactory
}

func NewIamAccessKeyEnumerator(repository repository.IAMRepository, factory resource.ResourceFactory) *IamAccessKeyEnumerator {
	return &IamAccessKeyEnumerator{
		repository,
		factory,
	}
}

func (e *IamAccessKeyEnumerator) SupportedType() resource.ResourceType {
	return resourceaws.AwsIamAccessKeyResourceType
}

func (e *IamAccessKeyEnumerator) Enumerate() ([]*resource.Resource, error) {
	users, err := e.repository.ListAllUsers()
	if err != nil {
		return nil, remoteerror.NewResourceListingErrorWithType(err, string(e.SupportedType()), resourceaws.AwsIamUserResourceType)
	}

	keys, err := e.repository.ListAllAccessKeys(users)
	if err != nil {
		return nil, remoteerror.NewResourceListingError(err, string(e.SupportedType()))
	}

	results := make([]*resource.Resource, 0)
	for _, key := range keys {
		results = append(
			results,
			e.factory.CreateAbstractResource(
				string(e.SupportedType()),
				*key.AccessKeyId,
				map[string]interface{}{
					"user": *key.UserName,
				},
			),
		)
	}

	return results, nil
}
