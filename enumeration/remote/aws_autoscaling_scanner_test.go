package remote

import (
	"testing"

	"github.com/snyk/driftctl/enumeration"
	"github.com/snyk/driftctl/enumeration/remote/alerts"
	"github.com/snyk/driftctl/enumeration/remote/aws"
	repository2 "github.com/snyk/driftctl/enumeration/remote/aws/repository"
	common2 "github.com/snyk/driftctl/enumeration/remote/common"
	remoteerr "github.com/snyk/driftctl/enumeration/remote/error"
	"github.com/snyk/driftctl/enumeration/terraform"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
	"github.com/snyk/driftctl/enumeration/resource"
	resourceaws "github.com/snyk/driftctl/enumeration/resource/aws"
	"github.com/snyk/driftctl/mocks"

	testresource "github.com/snyk/driftctl/test/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAutoscaling_LaunchConfiguration(t *testing.T) {
	tests := []struct {
		test           string
		mocks          func(*repository2.MockAutoScalingRepository, *mocks.AlerterInterface)
		assertExpected func(*testing.T, []*resource.Resource)
		wantErr        error
	}{
		{
			test: "no launch configuration",
			mocks: func(repository *repository2.MockAutoScalingRepository, alerter *mocks.AlerterInterface) {
				repository.On("DescribeLaunchConfigurations").Return([]*autoscaling.LaunchConfiguration{}, nil)
			},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 0)
			},
		},
		{
			test: "multiple launch configurations",
			mocks: func(repository *repository2.MockAutoScalingRepository, alerter *mocks.AlerterInterface) {
				repository.On("DescribeLaunchConfigurations").Return([]*autoscaling.LaunchConfiguration{
					{LaunchConfigurationName: awssdk.String("web_config_1")},
					{LaunchConfigurationName: awssdk.String("web_config_2")},
				}, nil)
			},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 2)

				assert.Equal(t, "web_config_1", got[0].ResourceId())
				assert.Equal(t, resourceaws.AwsLaunchConfigurationResourceType, got[0].ResourceType())

				assert.Equal(t, "web_config_2", got[1].ResourceId())
				assert.Equal(t, resourceaws.AwsLaunchConfigurationResourceType, got[1].ResourceType())
			},
		},
		{
			test: "cannot list launch configurations",
			mocks: func(repository *repository2.MockAutoScalingRepository, alerter *mocks.AlerterInterface) {
				awsError := awserr.NewRequestFailure(awserr.New("AccessDeniedException", "", errors.New("")), 403, "")
				repository.On("DescribeLaunchConfigurations").Return(nil, awsError)

				alerter.On("SendAlert", resourceaws.AwsLaunchConfigurationResourceType, alerts.NewRemoteAccessDeniedAlert(common2.RemoteAWSTerraform, remoteerr.NewResourceListingErrorWithType(awsError, resourceaws.AwsLaunchConfigurationResourceType, resourceaws.AwsLaunchConfigurationResourceType), alerts.EnumerationPhase)).Return()
			},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 0)
			},
			wantErr: nil,
		},
	}

	schemaRepository := testresource.InitFakeSchemaRepository("aws", "3.19.0")
	resourceaws.InitResourcesMetadata(schemaRepository)
	factory := terraform.NewTerraformResourceFactory(schemaRepository)

	for _, c := range tests {
		t.Run(c.test, func(tt *testing.T) {

			scanOptions := ScannerOptions{}
			remoteLibrary := common2.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			fakeRepo := &repository2.MockAutoScalingRepository{}
			c.mocks(fakeRepo, alerter)

			var repo repository2.AutoScalingRepository = fakeRepo

			remoteLibrary.AddEnumerator(aws.NewLaunchConfigurationEnumerator(repo, factory))

			testFilter := &enumeration.MockFilter{}
			testFilter.On("IsTypeIgnored", mock.Anything).Return(false)

			s := NewScanner(remoteLibrary, alerter, scanOptions, testFilter)
			got, err := s.Resources()
			assert.Equal(tt, c.wantErr, err)
			if err != nil {
				return
			}

			c.assertExpected(tt, got)
			alerter.AssertExpectations(tt)
			fakeRepo.AssertExpectations(tt)
		})
	}
}
