package remote

import (
	"testing"

	"github.com/snyk/driftctl/enumeration"
	azurerm2 "github.com/snyk/driftctl/enumeration/remote/azurerm"
	repository2 "github.com/snyk/driftctl/enumeration/remote/azurerm/repository"
	"github.com/snyk/driftctl/enumeration/remote/cache"
	common2 "github.com/snyk/driftctl/enumeration/remote/common"
	remoteerr "github.com/snyk/driftctl/enumeration/remote/error"
	terraform3 "github.com/snyk/driftctl/enumeration/terraform"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/pkg/errors"
	"github.com/snyk/driftctl/enumeration/resource"
	resourceazure "github.com/snyk/driftctl/enumeration/resource/azurerm"
	"github.com/snyk/driftctl/mocks"

	"github.com/snyk/driftctl/test"
	"github.com/snyk/driftctl/test/goldenfile"
	testresource "github.com/snyk/driftctl/test/resource"
	terraform2 "github.com/snyk/driftctl/test/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAzurermCompute_Image(t *testing.T) {

	dummyError := errors.New("this is an error")

	tests := []struct {
		test           string
		mocks          func(*repository2.MockComputeRepository, *mocks.AlerterInterface)
		assertExpected func(t *testing.T, got []*resource.Resource)
		wantErr        error
	}{
		{
			test: "no images",
			mocks: func(repository *repository2.MockComputeRepository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllImages").Return([]*armcompute.Image{}, nil)
			},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 0)
			},
		},
		{
			test: "error listing images",
			mocks: func(repository *repository2.MockComputeRepository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllImages").Return(nil, dummyError)
			},
			wantErr: remoteerr.NewResourceListingError(dummyError, resourceazure.AzureImageResourceType),
		},
		{
			test: "multiple images including an invalid ID",
			mocks: func(repository *repository2.MockComputeRepository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllImages").Return([]*armcompute.Image{
					{
						Resource: armcompute.Resource{
							ID:   to.StringPtr("/subscriptions/4e411884-65b0-4911-bc80-52f9a21942a2/resourceGroups/testgroup/providers/Microsoft.Compute/images/image1"),
							Name: to.StringPtr("image1"),
						},
					},
					{
						Resource: armcompute.Resource{
							ID:   to.StringPtr("/subscriptions/4e411884-65b0-4911-bc80-52f9a21942a2/resourceGroups/testgroup/providers/Microsoft.Compute/images/image2"),
							Name: to.StringPtr("image2"),
						},
					},
					{
						Resource: armcompute.Resource{
							ID:   to.StringPtr("/invalid-id/image3"),
							Name: to.StringPtr("image3"),
						},
					},
				}, nil)
			},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 2)

				assert.Equal(t, got[0].ResourceId(), "/subscriptions/4e411884-65b0-4911-bc80-52f9a21942a2/resourceGroups/testgroup/providers/Microsoft.Compute/images/image1")
				assert.Equal(t, got[0].ResourceType(), resourceazure.AzureImageResourceType)

				assert.Equal(t, got[1].ResourceId(), "/subscriptions/4e411884-65b0-4911-bc80-52f9a21942a2/resourceGroups/testgroup/providers/Microsoft.Compute/images/image2")
				assert.Equal(t, got[1].ResourceType(), resourceazure.AzureImageResourceType)
			},
		},
	}

	providerVersion := "2.71.0"
	schemaRepository := testresource.InitFakeSchemaRepository("azurerm", providerVersion)
	resourceazure.InitResourcesMetadata(schemaRepository)
	factory := terraform3.NewTerraformResourceFactory(schemaRepository)

	for _, c := range tests {
		t.Run(c.test, func(tt *testing.T) {
			scanOptions := ScannerOptions{}
			remoteLibrary := common2.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			fakeRepo := &repository2.MockComputeRepository{}
			c.mocks(fakeRepo, alerter)

			remoteLibrary.AddEnumerator(azurerm2.NewAzurermImageEnumerator(fakeRepo, factory))

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

func TestAzurermCompute_SSHPublicKey(t *testing.T) {

	dummyError := errors.New("this is an error")

	tests := []struct {
		test    string
		dirName string
		mocks   func(*repository2.MockComputeRepository, *mocks.AlerterInterface)
		wantErr error
	}{
		{
			test:    "no public key",
			dirName: "azurerm_ssh_public_key_empty",
			mocks: func(repository *repository2.MockComputeRepository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllSSHPublicKeys").Return([]*armcompute.SSHPublicKeyResource{}, nil)
			},
		},
		{
			test:    "error listing public keys",
			dirName: "azurerm_ssh_public_key_empty",
			mocks: func(repository *repository2.MockComputeRepository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllSSHPublicKeys").Return(nil, dummyError)
			},
			wantErr: remoteerr.NewResourceListingError(dummyError, resourceazure.AzureSSHPublicKeyResourceType),
		},
		{
			test:    "multiple public keys",
			dirName: "azurerm_ssh_public_key_multiple",
			mocks: func(repository *repository2.MockComputeRepository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllSSHPublicKeys").Return([]*armcompute.SSHPublicKeyResource{
					{
						Resource: armcompute.Resource{
							ID:   to.StringPtr("/subscriptions/7bfb2c5c-7308-46ed-8ae4-fffa356eb406/resourceGroups/TESTRESGROUP/providers/Microsoft.Compute/sshPublicKeys/example-key"),
							Name: to.StringPtr("example-key"),
						},
					},
					{
						Resource: armcompute.Resource{
							ID:   to.StringPtr("/subscriptions/7bfb2c5c-7308-46ed-8ae4-fffa356eb406/resourceGroups/TESTRESGROUP/providers/Microsoft.Compute/sshPublicKeys/example-key2"),
							Name: to.StringPtr("example-key2"),
						},
					},
				}, nil)
			},
		},
	}

	providerVersion := "2.71.0"
	schemaRepository := testresource.InitFakeSchemaRepository("azurerm", providerVersion)
	resourceazure.InitResourcesMetadata(schemaRepository)
	factory := terraform3.NewTerraformResourceFactory(schemaRepository)
	deserializer := resource.NewDeserializer(factory)

	for _, c := range tests {
		t.Run(c.test, func(tt *testing.T) {
			shouldUpdate := c.dirName == *goldenfile.Update

			scanOptions := ScannerOptions{Deep: true}
			providerLibrary := terraform3.NewProviderLibrary()
			remoteLibrary := common2.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			fakeRepo := &repository2.MockComputeRepository{}
			c.mocks(fakeRepo, alerter)

			var repo repository2.ComputeRepository = fakeRepo
			providerVersion := "2.71.0"
			realProvider, err := terraform2.InitTestAzureProvider(providerLibrary, providerVersion)
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
				cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{})
				if err != nil {
					t.Fatal(err)
				}
				clientOptions := &arm.ClientOptions{}
				repo = repository2.NewComputeRepository(cred, clientOptions, realProvider.GetConfig(), cache.New(0))
			}

			remoteLibrary.AddEnumerator(azurerm2.NewAzurermSSHPublicKeyEnumerator(repo, factory))
			remoteLibrary.AddDetailsFetcher(resourceazure.AzureSSHPublicKeyResourceType, common2.NewGenericDetailsFetcher(resourceazure.AzureSSHPublicKeyResourceType, provider, deserializer))

			testFilter := &enumeration.MockFilter{}
			testFilter.On("IsTypeIgnored", mock.Anything).Return(false)

			s := NewScanner(remoteLibrary, alerter, scanOptions, testFilter)
			got, err := s.Resources()
			assert.Equal(tt, c.wantErr, err)

			if err != nil {
				return
			}
			test.TestAgainstGoldenFile(got, resourceazure.AzureSSHPublicKeyResourceType, c.dirName, provider, deserializer, shouldUpdate, tt)
			alerter.AssertExpectations(tt)
			fakeRepo.AssertExpectations(tt)
		})
	}
}
