package remote

import (
	"testing"

	"github.com/snyk/driftctl/enumeration"
	"github.com/snyk/driftctl/enumeration/remote/alerts"
	aws2 "github.com/snyk/driftctl/enumeration/remote/aws"
	repository2 "github.com/snyk/driftctl/enumeration/remote/aws/repository"
	"github.com/snyk/driftctl/enumeration/remote/cache"
	common2 "github.com/snyk/driftctl/enumeration/remote/common"
	remoteerr "github.com/snyk/driftctl/enumeration/remote/error"
	terraform3 "github.com/snyk/driftctl/enumeration/terraform"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
	"github.com/snyk/driftctl/mocks"

	testresource "github.com/snyk/driftctl/test/resource"
	terraform2 "github.com/snyk/driftctl/test/terraform"
	"github.com/stretchr/testify/mock"

	"github.com/snyk/driftctl/enumeration/resource"
	resourceaws "github.com/snyk/driftctl/enumeration/resource/aws"

	"github.com/snyk/driftctl/test"
	"github.com/snyk/driftctl/test/goldenfile"
	"github.com/stretchr/testify/assert"
)

func TestScanLambdaFunction(t *testing.T) {

	tests := []struct {
		test    string
		dirName string
		mocks   func(*repository2.MockLambdaRepository, *mocks.AlerterInterface)
		err     error
	}{
		{
			test:    "no lambda functions",
			dirName: "aws_lambda_function_empty",
			mocks: func(repo *repository2.MockLambdaRepository, alerter *mocks.AlerterInterface) {
				repo.On("ListAllLambdaFunctions").Return([]*lambda.FunctionConfiguration{}, nil)
			},
			err: nil,
		},
		{
			test:    "with lambda functions",
			dirName: "aws_lambda_function_multiple",
			mocks: func(repo *repository2.MockLambdaRepository, alerter *mocks.AlerterInterface) {
				repo.On("ListAllLambdaFunctions").Return([]*lambda.FunctionConfiguration{
					{
						FunctionName: awssdk.String("foo"),
					},
					{
						FunctionName: awssdk.String("bar"),
					},
				}, nil)
			},
			err: nil,
		},
		{
			test:    "One lambda with signing",
			dirName: "aws_lambda_function_signed",
			mocks: func(repo *repository2.MockLambdaRepository, alerter *mocks.AlerterInterface) {
				repo.On("ListAllLambdaFunctions").Return([]*lambda.FunctionConfiguration{
					{
						FunctionName: awssdk.String("foo"),
					},
				}, nil)
			},
			err: nil,
		},
		{
			test:    "cannot list lambda functions",
			dirName: "aws_lambda_function_empty",
			mocks: func(repo *repository2.MockLambdaRepository, alerter *mocks.AlerterInterface) {
				awsError := awserr.NewRequestFailure(awserr.New("AccessDeniedException", "", errors.New("")), 403, "")
				repo.On("ListAllLambdaFunctions").Return([]*lambda.FunctionConfiguration{}, awsError)

				alerter.On("SendAlert", resourceaws.AwsLambdaFunctionResourceType, alerts.NewRemoteAccessDeniedAlert(common2.RemoteAWSTerraform, remoteerr.NewResourceListingErrorWithType(awsError, resourceaws.AwsLambdaFunctionResourceType, resourceaws.AwsLambdaFunctionResourceType), alerts.EnumerationPhase)).Return()
			},
			err: nil,
		},
	}

	schemaRepository := testresource.InitFakeSchemaRepository("aws", "3.19.0")
	resourceaws.InitResourcesMetadata(schemaRepository)
	factory := terraform3.NewTerraformResourceFactory(schemaRepository)
	deserializer := resource.NewDeserializer(factory)

	for _, c := range tests {
		t.Run(c.test, func(tt *testing.T) {
			shouldUpdate := c.dirName == *goldenfile.Update

			session := session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))

			scanOptions := ScannerOptions{Deep: true}
			providerLibrary := terraform3.NewProviderLibrary()
			remoteLibrary := common2.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			fakeRepo := &repository2.MockLambdaRepository{}
			c.mocks(fakeRepo, alerter)

			var repo repository2.LambdaRepository = fakeRepo
			providerVersion := "3.19.0"
			realProvider, err := terraform2.InitTestAwsProvider(providerLibrary, providerVersion)
			if err != nil {
				t.Fatal(err)
			}
			provider := terraform2.NewFakeTerraformProvider(realProvider)
			provider.WithResponse(c.dirName)

			// Replace mock by real resources if we are in update mode
			if shouldUpdate {
				err := realProvider.Init()
				if err != nil {
					t.Fatal(err)
				}
				provider.ShouldUpdate()
				repo = repository2.NewLambdaRepository(session, cache.New(0))
			}

			remoteLibrary.AddEnumerator(aws2.NewLambdaFunctionEnumerator(repo, factory))
			remoteLibrary.AddDetailsFetcher(resourceaws.AwsLambdaFunctionResourceType, common2.NewGenericDetailsFetcher(resourceaws.AwsLambdaFunctionResourceType, provider, deserializer))

			testFilter := &enumeration.MockFilter{}
			testFilter.On("IsTypeIgnored", mock.Anything).Return(false)

			s := NewScanner(remoteLibrary, alerter, scanOptions, testFilter)
			got, err := s.Resources()
			assert.Equal(tt, c.err, err)
			if err != nil {
				return
			}
			test.TestAgainstGoldenFile(got, resourceaws.AwsLambdaFunctionResourceType, c.dirName, provider, deserializer, shouldUpdate, tt)
			alerter.AssertExpectations(tt)
			fakeRepo.AssertExpectations(tt)
		})
	}
}

func TestScanLambdaEventSourceMapping(t *testing.T) {

	tests := []struct {
		test    string
		dirName string
		mocks   func(*repository2.MockLambdaRepository, *mocks.AlerterInterface)
		err     error
	}{
		{
			test:    "no EventSourceMapping",
			dirName: "aws_lambda_source_mapping_empty",
			mocks: func(repo *repository2.MockLambdaRepository, alerter *mocks.AlerterInterface) {
				repo.On("ListAllLambdaEventSourceMappings").Return([]*lambda.EventSourceMappingConfiguration{}, nil)
			},
			err: nil,
		},
		{
			test:    "with 2 sqs EventSourceMapping",
			dirName: "aws_lambda_source_mapping_sqs_multiple",
			mocks: func(repo *repository2.MockLambdaRepository, alerter *mocks.AlerterInterface) {
				repo.On("ListAllLambdaEventSourceMappings").Return([]*lambda.EventSourceMappingConfiguration{
					{
						UUID: awssdk.String("13ff66f8-37eb-4ad6-a0a8-594fea72df4f"),
					},
					{
						UUID: awssdk.String("4ad7e2b3-79e9-4713-9d9d-5af2c01d9058"),
					},
				}, nil)
			},
			err: nil,
		},
		{
			test:    "with dynamo EventSourceMapping",
			dirName: "aws_lambda_source_mapping_dynamo_multiple",
			mocks: func(repo *repository2.MockLambdaRepository, alerter *mocks.AlerterInterface) {
				repo.On("ListAllLambdaEventSourceMappings").Return([]*lambda.EventSourceMappingConfiguration{
					{
						UUID: awssdk.String("1aa9c4a0-060b-41c1-a9ae-dc304ebcdb00"),
					},
				}, nil)
			},
			err: nil,
		},
		{
			test:    "cannot list lambda functions",
			dirName: "aws_lambda_function_empty",
			mocks: func(repo *repository2.MockLambdaRepository, alerter *mocks.AlerterInterface) {
				awsError := awserr.NewRequestFailure(awserr.New("AccessDeniedException", "", errors.New("")), 403, "")
				repo.On("ListAllLambdaEventSourceMappings").Return([]*lambda.EventSourceMappingConfiguration{}, awsError)

				alerter.On("SendAlert", resourceaws.AwsLambdaEventSourceMappingResourceType, alerts.NewRemoteAccessDeniedAlert(common2.RemoteAWSTerraform, remoteerr.NewResourceListingErrorWithType(awsError, resourceaws.AwsLambdaEventSourceMappingResourceType, resourceaws.AwsLambdaEventSourceMappingResourceType), alerts.EnumerationPhase)).Return()
			},
			err: nil,
		},
	}

	schemaRepository := testresource.InitFakeSchemaRepository("aws", "3.19.0")
	resourceaws.InitResourcesMetadata(schemaRepository)
	factory := terraform3.NewTerraformResourceFactory(schemaRepository)
	deserializer := resource.NewDeserializer(factory)

	for _, c := range tests {
		t.Run(c.test, func(tt *testing.T) {
			shouldUpdate := c.dirName == *goldenfile.Update

			session := session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))

			scanOptions := ScannerOptions{Deep: true}
			providerLibrary := terraform3.NewProviderLibrary()
			remoteLibrary := common2.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			fakeRepo := &repository2.MockLambdaRepository{}
			c.mocks(fakeRepo, alerter)

			var repo repository2.LambdaRepository = fakeRepo
			providerVersion := "3.19.0"
			realProvider, err := terraform2.InitTestAwsProvider(providerLibrary, providerVersion)
			if err != nil {
				t.Fatal(err)
			}
			provider := terraform2.NewFakeTerraformProvider(realProvider)
			provider.WithResponse(c.dirName)

			// Replace mock by real resources if we are in update mode
			if shouldUpdate {
				err := realProvider.Init()
				if err != nil {
					t.Fatal(err)
				}
				provider.ShouldUpdate()
				repo = repository2.NewLambdaRepository(session, cache.New(0))
			}

			remoteLibrary.AddEnumerator(aws2.NewLambdaEventSourceMappingEnumerator(repo, factory))
			remoteLibrary.AddDetailsFetcher(resourceaws.AwsLambdaEventSourceMappingResourceType, common2.NewGenericDetailsFetcher(resourceaws.AwsLambdaEventSourceMappingResourceType, provider, deserializer))

			testFilter := &enumeration.MockFilter{}
			testFilter.On("IsTypeIgnored", mock.Anything).Return(false)

			s := NewScanner(remoteLibrary, alerter, scanOptions, testFilter)
			got, err := s.Resources()
			assert.Equal(tt, c.err, err)
			if err != nil {
				return
			}
			test.TestAgainstGoldenFile(got, resourceaws.AwsLambdaEventSourceMappingResourceType, c.dirName, provider, deserializer, shouldUpdate, tt)
			alerter.AssertExpectations(tt)
			fakeRepo.AssertExpectations(tt)
		})
	}
}
