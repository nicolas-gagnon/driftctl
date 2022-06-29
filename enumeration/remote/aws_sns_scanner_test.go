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
	"github.com/aws/aws-sdk-go/service/sns"
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

func TestScanSNSTopic(t *testing.T) {

	cases := []struct {
		test    string
		dirName string
		mocks   func(*repository2.MockSNSRepository, *mocks.AlerterInterface)
		err     error
	}{
		{
			test:    "no SNS Topic",
			dirName: "aws_sns_topic_empty",
			mocks: func(client *repository2.MockSNSRepository, alerter *mocks.AlerterInterface) {
				client.On("ListAllTopics").Return([]*sns.Topic{}, nil)
			},
			err: nil,
		},
		{
			test:    "Multiple SNSTopic",
			dirName: "aws_sns_topic_multiple",
			mocks: func(client *repository2.MockSNSRepository, alerter *mocks.AlerterInterface) {
				client.On("ListAllTopics").Return([]*sns.Topic{
					{TopicArn: awssdk.String("arn:aws:sns:eu-west-3:526954929923:user-updates-topic")},
					{TopicArn: awssdk.String("arn:aws:sns:eu-west-3:526954929923:user-updates-topic2")},
					{TopicArn: awssdk.String("arn:aws:sns:eu-west-3:526954929923:user-updates-topic3")},
				}, nil)
			},
			err: nil,
		},
		{
			test:    "cannot list SNSTopic",
			dirName: "aws_sns_topic_empty",
			mocks: func(client *repository2.MockSNSRepository, alerter *mocks.AlerterInterface) {
				awsError := awserr.NewRequestFailure(awserr.New("AccessDeniedException", "", errors.New("")), 403, "")
				client.On("ListAllTopics").Return(nil, awsError)

				alerter.On("SendAlert", resourceaws.AwsSnsTopicResourceType, alerts.NewRemoteAccessDeniedAlert(common2.RemoteAWSTerraform, remoteerr.NewResourceListingErrorWithType(awsError, resourceaws.AwsSnsTopicResourceType, resourceaws.AwsSnsTopicResourceType), alerts.EnumerationPhase)).Return()
			},
			err: nil,
		},
	}

	schemaRepository := testresource.InitFakeSchemaRepository("aws", "3.19.0")
	resourceaws.InitResourcesMetadata(schemaRepository)
	factory := terraform3.NewTerraformResourceFactory(schemaRepository)
	deserializer := resource.NewDeserializer(factory)

	for _, c := range cases {
		t.Run(c.test, func(tt *testing.T) {
			shouldUpdate := c.dirName == *goldenfile.Update

			sess := session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))

			scanOptions := ScannerOptions{Deep: true}
			providerLibrary := terraform3.NewProviderLibrary()
			remoteLibrary := common2.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			fakeRepo := &repository2.MockSNSRepository{}
			c.mocks(fakeRepo, alerter)

			var repo repository2.SNSRepository = fakeRepo
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
				repo = repository2.NewSNSRepository(sess, cache.New(0))
			}

			remoteLibrary.AddEnumerator(aws2.NewSNSTopicEnumerator(repo, factory))
			remoteLibrary.AddDetailsFetcher(resourceaws.AwsSnsTopicResourceType, common2.NewGenericDetailsFetcher(resourceaws.AwsSnsTopicResourceType, provider, deserializer))

			testFilter := &enumeration.MockFilter{}
			testFilter.On("IsTypeIgnored", mock.Anything).Return(false)

			s := NewScanner(remoteLibrary, alerter, scanOptions, testFilter)
			got, err := s.Resources()
			assert.Equal(tt, c.err, err)
			if err != nil {
				return
			}
			test.TestAgainstGoldenFile(got, resourceaws.AwsSnsTopicResourceType, c.dirName, provider, deserializer, shouldUpdate, tt)
			alerter.AssertExpectations(tt)
			fakeRepo.AssertExpectations(tt)
		})
	}
}

func TestSNSTopicPolicyScan(t *testing.T) {

	cases := []struct {
		test    string
		dirName string
		mocks   func(*repository2.MockSNSRepository, *mocks.AlerterInterface)
		err     error
	}{
		{
			test:    "no SNS Topic policy",
			dirName: "aws_sns_topic_policy_empty",
			mocks: func(client *repository2.MockSNSRepository, alerter *mocks.AlerterInterface) {
				client.On("ListAllTopics").Return([]*sns.Topic{}, nil)
			},
			err: nil,
		},
		{
			test:    "Multiple SNSTopicPolicy",
			dirName: "aws_sns_topic_policy_multiple",
			mocks: func(client *repository2.MockSNSRepository, alerter *mocks.AlerterInterface) {
				client.On("ListAllTopics").Return([]*sns.Topic{
					{TopicArn: awssdk.String("arn:aws:sns:us-east-1:526954929923:my-topic-with-policy")},
					{TopicArn: awssdk.String("arn:aws:sns:us-east-1:526954929923:my-topic-with-policy2")},
				}, nil)
			},
			err: nil,
		},
		{
			test:    "cannot list SNSTopic",
			dirName: "aws_sns_topic_policy_topic_list",
			mocks: func(client *repository2.MockSNSRepository, alerter *mocks.AlerterInterface) {
				awsError := awserr.NewRequestFailure(awserr.New("AccessDeniedException", "", errors.New("")), 403, "")
				client.On("ListAllTopics").Return(nil, awsError)

				alerter.On("SendAlert", resourceaws.AwsSnsTopicPolicyResourceType, alerts.NewRemoteAccessDeniedAlert(common2.RemoteAWSTerraform, remoteerr.NewResourceListingErrorWithType(awsError, resourceaws.AwsSnsTopicPolicyResourceType, resourceaws.AwsSnsTopicResourceType), alerts.EnumerationPhase)).Return()
			},
			err: nil,
		},
	}

	schemaRepository := testresource.InitFakeSchemaRepository("aws", "3.19.0")
	resourceaws.InitResourcesMetadata(schemaRepository)
	factory := terraform3.NewTerraformResourceFactory(schemaRepository)
	deserializer := resource.NewDeserializer(factory)

	for _, c := range cases {
		t.Run(c.test, func(tt *testing.T) {
			shouldUpdate := c.dirName == *goldenfile.Update

			sess := session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))

			scanOptions := ScannerOptions{Deep: true}
			providerLibrary := terraform3.NewProviderLibrary()
			remoteLibrary := common2.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			fakeRepo := &repository2.MockSNSRepository{}
			c.mocks(fakeRepo, alerter)

			var repo repository2.SNSRepository = fakeRepo
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
				repo = repository2.NewSNSRepository(sess, cache.New(0))
			}

			remoteLibrary.AddEnumerator(aws2.NewSNSTopicPolicyEnumerator(repo, factory))
			remoteLibrary.AddDetailsFetcher(resourceaws.AwsSnsTopicPolicyResourceType, common2.NewGenericDetailsFetcher(resourceaws.AwsSnsTopicPolicyResourceType, provider, deserializer))

			testFilter := &enumeration.MockFilter{}
			testFilter.On("IsTypeIgnored", mock.Anything).Return(false)

			s := NewScanner(remoteLibrary, alerter, scanOptions, testFilter)
			got, err := s.Resources()
			assert.Equal(tt, c.err, err)
			if err != nil {
				return
			}
			test.TestAgainstGoldenFile(got, resourceaws.AwsSnsTopicPolicyResourceType, c.dirName, provider, deserializer, shouldUpdate, tt)
			alerter.AssertExpectations(tt)
			fakeRepo.AssertExpectations(tt)
		})
	}
}

func TestSNSTopicSubscriptionScan(t *testing.T) {

	cases := []struct {
		test    string
		dirName string
		mocks   func(*repository2.MockSNSRepository, *mocks.AlerterInterface)
		err     error
	}{
		{
			test:    "no SNS Topic Subscription",
			dirName: "aws_sns_topic_subscription_empty",
			mocks: func(client *repository2.MockSNSRepository, alerter *mocks.AlerterInterface) {
				client.On("ListAllSubscriptions").Return([]*sns.Subscription{}, nil)
			},
			err: nil,
		},
		{
			test:    "Multiple SNSTopic Subscription",
			dirName: "aws_sns_topic_subscription_multiple",
			mocks: func(client *repository2.MockSNSRepository, alerter *mocks.AlerterInterface) {
				client.On("ListAllSubscriptions").Return([]*sns.Subscription{
					{SubscriptionArn: awssdk.String("arn:aws:sns:us-east-1:526954929923:user-updates-topic2:c0f794c5-a009-4db4-9147-4c55959787fa")},
					{SubscriptionArn: awssdk.String("arn:aws:sns:us-east-1:526954929923:user-updates-topic:b6e66147-2b31-4486-8d4b-2a2272264c8e")},
				}, nil)
			},
			err: nil,
		},
		{
			test:    "Multiple SNSTopic Subscription with one pending and one incorrect",
			dirName: "aws_sns_topic_subscription_multiple",
			mocks: func(client *repository2.MockSNSRepository, alerter *mocks.AlerterInterface) {
				client.On("ListAllSubscriptions").Return([]*sns.Subscription{
					{SubscriptionArn: awssdk.String("PendingConfirmation"), Endpoint: awssdk.String("TEST")},
					{SubscriptionArn: awssdk.String("Incorrect"), Endpoint: awssdk.String("INCORRECT")},
					{SubscriptionArn: awssdk.String("arn:aws:sns:us-east-1:526954929923:user-updates-topic2:c0f794c5-a009-4db4-9147-4c55959787fa")},
					{SubscriptionArn: awssdk.String("arn:aws:sns:us-east-1:526954929923:user-updates-topic:b6e66147-2b31-4486-8d4b-2a2272264c8e")},
				}, nil)

				alerter.On("SendAlert", "aws_sns_topic_subscription.PendingConfirmation", aws2.NewWrongArnTopicAlert("PendingConfirmation", awssdk.String("TEST"))).Return()

				alerter.On("SendAlert", "aws_sns_topic_subscription.Incorrect", aws2.NewWrongArnTopicAlert("Incorrect", awssdk.String("INCORRECT"))).Return()
			},
			err: nil,
		},
		{
			test:    "cannot list SNSTopic subscription",
			dirName: "aws_sns_topic_subscription_list",
			mocks: func(client *repository2.MockSNSRepository, alerter *mocks.AlerterInterface) {
				awsError := awserr.NewRequestFailure(awserr.New("AccessDeniedException", "", errors.New("")), 403, "")
				client.On("ListAllSubscriptions").Return(nil, awsError)

				alerter.On("SendAlert", resourceaws.AwsSnsTopicSubscriptionResourceType, alerts.NewRemoteAccessDeniedAlert(common2.RemoteAWSTerraform, remoteerr.NewResourceListingErrorWithType(awsError, resourceaws.AwsSnsTopicSubscriptionResourceType, resourceaws.AwsSnsTopicSubscriptionResourceType), alerts.EnumerationPhase)).Return()
			},
			err: nil,
		},
	}

	schemaRepository := testresource.InitFakeSchemaRepository("aws", "3.19.0")
	resourceaws.InitResourcesMetadata(schemaRepository)
	factory := terraform3.NewTerraformResourceFactory(schemaRepository)
	deserializer := resource.NewDeserializer(factory)

	for _, c := range cases {
		t.Run(c.test, func(tt *testing.T) {
			shouldUpdate := c.dirName == *goldenfile.Update

			sess := session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))

			scanOptions := ScannerOptions{Deep: true}
			providerLibrary := terraform3.NewProviderLibrary()
			remoteLibrary := common2.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			fakeRepo := &repository2.MockSNSRepository{}
			c.mocks(fakeRepo, alerter)

			var repo repository2.SNSRepository = fakeRepo
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
				repo = repository2.NewSNSRepository(sess, cache.New(0))
			}

			remoteLibrary.AddEnumerator(aws2.NewSNSTopicSubscriptionEnumerator(repo, factory, alerter))
			remoteLibrary.AddDetailsFetcher(resourceaws.AwsSnsTopicSubscriptionResourceType, common2.NewGenericDetailsFetcher(resourceaws.AwsSnsTopicSubscriptionResourceType, provider, deserializer))

			testFilter := &enumeration.MockFilter{}
			testFilter.On("IsTypeIgnored", mock.Anything).Return(false)

			s := NewScanner(remoteLibrary, alerter, scanOptions, testFilter)
			got, err := s.Resources()
			assert.Equal(tt, err, c.err)
			if err != nil {
				return
			}
			test.TestAgainstGoldenFile(got, resourceaws.AwsSnsTopicSubscriptionResourceType, c.dirName, provider, deserializer, shouldUpdate, tt)
			alerter.AssertExpectations(tt)
			fakeRepo.AssertExpectations(tt)
		})
	}
}
