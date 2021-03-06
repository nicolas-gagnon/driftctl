package aws

import "github.com/snyk/driftctl/pkg/resource"

const AwsEbsEncryptionByDefaultResourceType = "aws_ebs_encryption_by_default"

func initAwsEbsEncryptionByDefaultMetaData(resourceSchemaRepository resource.SchemaRepositoryInterface) {
	resourceSchemaRepository.SetFlags(AwsEbsEncryptionByDefaultResourceType, resource.FlagDeepMode)
}
