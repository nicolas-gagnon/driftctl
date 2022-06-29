package aws

import (
	tf "github.com/snyk/driftctl/enumeration/terraform"

	"testing"

	"github.com/snyk/driftctl/enumeration/resource"
	testresource "github.com/snyk/driftctl/test/resource"
	"github.com/stretchr/testify/assert"
)

func TestAWS_Metadata_Flags(t *testing.T) {
	testcases := map[string][]resource.Flags{
		AwsAmiResourceType:                             {resource.FlagDeepMode},
		AwsApiGatewayAccountResourceType:               {},
		AwsApiGatewayApiKeyResourceType:                {},
		AwsApiGatewayAuthorizerResourceType:            {},
		AwsApiGatewayBasePathMappingResourceType:       {},
		AwsApiGatewayDeploymentResourceType:            {},
		AwsApiGatewayDomainNameResourceType:            {},
		AwsApiGatewayGatewayResponseResourceType:       {},
		AwsApiGatewayIntegrationResourceType:           {},
		AwsApiGatewayIntegrationResponseResourceType:   {},
		AwsApiGatewayMethodResourceType:                {},
		AwsApiGatewayMethodResponseResourceType:        {},
		AwsApiGatewayMethodSettingsResourceType:        {},
		AwsApiGatewayModelResourceType:                 {},
		AwsApiGatewayRequestValidatorResourceType:      {},
		AwsApiGatewayResourceResourceType:              {},
		AwsApiGatewayRestApiResourceType:               {},
		AwsApiGatewayRestApiPolicyResourceType:         {},
		AwsApiGatewayStageResourceType:                 {},
		AwsApiGatewayVpcLinkResourceType:               {},
		AwsApiGatewayV2ApiResourceType:                 {},
		AwsApiGatewayV2RouteResourceType:               {},
		AwsApiGatewayV2DeploymentResourceType:          {},
		AwsApiGatewayV2VpcLinkResourceType:             {},
		AwsApiGatewayV2AuthorizerResourceType:          {},
		AwsApiGatewayV2RouteResponseResourceType:       {},
		AwsApiGatewayV2DomainNameResourceType:          {},
		AwsApiGatewayV2ModelResourceType:               {},
		AwsApiGatewayV2StageResourceType:               {},
		AwsApiGatewayV2MappingResourceType:             {},
		AwsApiGatewayV2IntegrationResourceType:         {},
		AwsApiGatewayV2IntegrationResponseResourceType: {},
		AwsAppAutoscalingPolicyResourceType:            {resource.FlagDeepMode},
		AwsAppAutoscalingScheduledActionResourceType:   {},
		AwsAppAutoscalingTargetResourceType:            {resource.FlagDeepMode},
		AwsCloudformationStackResourceType:             {resource.FlagDeepMode},
		AwsCloudfrontDistributionResourceType:          {resource.FlagDeepMode},
		AwsDbInstanceResourceType:                      {resource.FlagDeepMode},
		AwsDbSubnetGroupResourceType:                   {resource.FlagDeepMode},
		AwsDefaultNetworkACLResourceType:               {resource.FlagDeepMode},
		AwsDefaultRouteTableResourceType:               {resource.FlagDeepMode},
		AwsDefaultSecurityGroupResourceType:            {resource.FlagDeepMode},
		AwsDefaultSubnetResourceType:                   {resource.FlagDeepMode},
		AwsDefaultVpcResourceType:                      {resource.FlagDeepMode},
		AwsDynamodbTableResourceType:                   {resource.FlagDeepMode},
		AwsEbsEncryptionByDefaultResourceType:          {resource.FlagDeepMode},
		AwsEbsSnapshotResourceType:                     {resource.FlagDeepMode},
		AwsEbsVolumeResourceType:                       {resource.FlagDeepMode},
		AwsEcrRepositoryResourceType:                   {resource.FlagDeepMode},
		AwsEipResourceType:                             {resource.FlagDeepMode},
		AwsEipAssociationResourceType:                  {resource.FlagDeepMode},
		AwsElastiCacheClusterResourceType:              {},
		AwsIamAccessKeyResourceType:                    {resource.FlagDeepMode},
		AwsIamPolicyResourceType:                       {resource.FlagDeepMode},
		AwsIamPolicyAttachmentResourceType:             {resource.FlagDeepMode},
		AwsIamRoleResourceType:                         {resource.FlagDeepMode},
		AwsIamRolePolicyResourceType:                   {resource.FlagDeepMode},
		AwsIamRolePolicyAttachmentResourceType:         {resource.FlagDeepMode},
		AwsIamUserResourceType:                         {resource.FlagDeepMode},
		AwsIamUserPolicyResourceType:                   {resource.FlagDeepMode},
		AwsIamUserPolicyAttachmentResourceType:         {resource.FlagDeepMode},
		AwsIamGroupPolicyResourceType:                  {},
		AwsIamGroupPolicyAttachmentResourceType:        {},
		AwsInstanceResourceType:                        {resource.FlagDeepMode},
		AwsInternetGatewayResourceType:                 {resource.FlagDeepMode},
		AwsKeyPairResourceType:                         {resource.FlagDeepMode},
		AwsKmsAliasResourceType:                        {resource.FlagDeepMode},
		AwsKmsKeyResourceType:                          {resource.FlagDeepMode},
		AwsLambdaEventSourceMappingResourceType:        {resource.FlagDeepMode},
		AwsLambdaFunctionResourceType:                  {resource.FlagDeepMode},
		AwsNatGatewayResourceType:                      {resource.FlagDeepMode},
		AwsNetworkACLResourceType:                      {resource.FlagDeepMode},
		AwsRDSClusterResourceType:                      {resource.FlagDeepMode},
		AwsRDSClusterInstanceResourceType:              {},
		AwsRouteResourceType:                           {resource.FlagDeepMode},
		AwsRoute53HealthCheckResourceType:              {resource.FlagDeepMode},
		AwsRoute53RecordResourceType:                   {resource.FlagDeepMode},
		AwsRoute53ZoneResourceType:                     {resource.FlagDeepMode},
		AwsRouteTableResourceType:                      {resource.FlagDeepMode},
		AwsRouteTableAssociationResourceType:           {resource.FlagDeepMode},
		AwsS3BucketResourceType:                        {resource.FlagDeepMode},
		AwsS3BucketAnalyticsConfigurationResourceType:  {resource.FlagDeepMode},
		AwsS3BucketInventoryResourceType:               {resource.FlagDeepMode},
		AwsS3BucketMetricResourceType:                  {resource.FlagDeepMode},
		AwsS3BucketNotificationResourceType:            {resource.FlagDeepMode},
		AwsS3BucketPolicyResourceType:                  {resource.FlagDeepMode},
		AwsS3BucketPublicAccessBlockResourceType:       {},
		AwsSecurityGroupResourceType:                   {resource.FlagDeepMode},
		AwsSnsTopicResourceType:                        {resource.FlagDeepMode},
		AwsSnsTopicPolicyResourceType:                  {resource.FlagDeepMode},
		AwsSnsTopicSubscriptionResourceType:            {resource.FlagDeepMode},
		AwsSqsQueueResourceType:                        {resource.FlagDeepMode},
		AwsSqsQueuePolicyResourceType:                  {resource.FlagDeepMode},
		AwsSubnetResourceType:                          {resource.FlagDeepMode},
		AwsVpcResourceType:                             {resource.FlagDeepMode},
		AwsSecurityGroupRuleResourceType:               {resource.FlagDeepMode},
		AwsNetworkACLRuleResourceType:                  {resource.FlagDeepMode},
		AwsLaunchTemplateResourceType:                  {resource.FlagDeepMode},
		AwsLaunchConfigurationResourceType:             {},
		AwsLoadBalancerResourceType:                    {},
		AwsApplicationLoadBalancerResourceType:         {},
		AwsClassicLoadBalancerResourceType:             {},
		AwsLoadBalancerListenerResourceType:            {},
		AwsApplicationLoadBalancerListenerResourceType: {},
		AwsIamGroupResourceType:                        {},
		AwsEcrRepositoryPolicyResourceType:             {},
	}

	schemaRepository := testresource.InitFakeSchemaRepository(tf.AWS, "3.19.0")
	InitResourcesMetadata(schemaRepository)

	for ty, flags := range testcases {
		t.Run(ty, func(tt *testing.T) {
			sch, exist := schemaRepository.GetSchema(ty)
			assert.True(tt, exist)

			if len(flags) == 0 {
				assert.Equal(tt, resource.Flags(0x0), sch.Flags, "should not have any flag")
				return
			}

			for _, flag := range flags {
				assert.Truef(tt, sch.Flags.HasFlag(flag), "should have given flag %d", flag)
			}
		})
	}
}
