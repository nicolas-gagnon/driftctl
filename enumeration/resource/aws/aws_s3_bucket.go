package aws

import (
	"github.com/snyk/driftctl/enumeration/resource"
)

const AwsS3BucketResourceType = "aws_s3_bucket"

func initAwsS3BucketMetaData(resourceSchemaRepository resource.SchemaRepositoryInterface) {
	resourceSchemaRepository.SetResolveReadAttributesFunc(AwsS3BucketResourceType, func(res *resource.Resource) map[string]string {
		return map[string]string{
			"alias": *res.Attributes().GetString("region"),
		}
	})
	resourceSchemaRepository.UpdateSchema(AwsS3BucketResourceType, map[string]func(attributeSchema *resource.AttributeSchema){
		"policy": func(attributeSchema *resource.AttributeSchema) {
			attributeSchema.JsonString = true
		},
	})
	resourceSchemaRepository.SetFlags(AwsS3BucketResourceType, resource.FlagDeepMode)
}
