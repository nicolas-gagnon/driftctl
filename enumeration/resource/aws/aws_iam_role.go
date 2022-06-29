package aws

import (
	"github.com/snyk/driftctl/enumeration/resource"
)

const AwsIamRoleResourceType = "aws_iam_role"

func initAwsIAMRoleMetaData(resourceSchemaRepository resource.SchemaRepositoryInterface) {
	resourceSchemaRepository.UpdateSchema(AwsIamRoleResourceType, map[string]func(attributeSchema *resource.AttributeSchema){
		"assume_role_policy": func(attributeSchema *resource.AttributeSchema) {
			attributeSchema.JsonString = true
		},
	})
	resourceSchemaRepository.SetFlags(AwsIamRoleResourceType, resource.FlagDeepMode)
}
