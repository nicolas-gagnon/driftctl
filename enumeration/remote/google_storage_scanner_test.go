package remote

import (
	"context"
	"testing"

	"github.com/snyk/driftctl/enumeration"
	"github.com/snyk/driftctl/enumeration/remote/alerts"
	"github.com/snyk/driftctl/enumeration/remote/cache"
	common2 "github.com/snyk/driftctl/enumeration/remote/common"
	remoteerr "github.com/snyk/driftctl/enumeration/remote/error"
	google2 "github.com/snyk/driftctl/enumeration/remote/google"
	repository2 "github.com/snyk/driftctl/enumeration/remote/google/repository"
	terraform3 "github.com/snyk/driftctl/enumeration/terraform"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"github.com/snyk/driftctl/enumeration/resource"
	googleresource "github.com/snyk/driftctl/enumeration/resource/google"
	"github.com/snyk/driftctl/mocks"

	"github.com/snyk/driftctl/test"
	"github.com/snyk/driftctl/test/goldenfile"
	testgoogle "github.com/snyk/driftctl/test/google"
	testresource "github.com/snyk/driftctl/test/resource"
	terraform2 "github.com/snyk/driftctl/test/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGoogleStorageBucket(t *testing.T) {

	cases := []struct {
		test             string
		dirName          string
		response         []*assetpb.ResourceSearchResult
		responseErr      error
		setupAlerterMock func(alerter *mocks.AlerterInterface)
		wantErr          error
	}{
		{
			test:     "no storage buckets",
			dirName:  "google_storage_bucket_empty",
			response: []*assetpb.ResourceSearchResult{},
			wantErr:  nil,
		},
		{
			test:    "multiples storage buckets",
			dirName: "google_storage_bucket",
			response: []*assetpb.ResourceSearchResult{
				{
					AssetType:   "storage.googleapis.com/Bucket",
					DisplayName: "driftctl-unittest-1",
				},
				{
					AssetType:   "storage.googleapis.com/Bucket",
					DisplayName: "driftctl-unittest-2",
				},
				{
					AssetType:   "storage.googleapis.com/Bucket",
					DisplayName: "driftctl-unittest-3",
				},
			},
			wantErr: nil,
		},
		{
			test:        "cannot list storage buckets",
			dirName:     "google_storage_bucket_empty",
			responseErr: status.Error(codes.PermissionDenied, "The caller does not have permission"),
			setupAlerterMock: func(alerter *mocks.AlerterInterface) {
				alerter.On(
					"SendAlert",
					"google_storage_bucket",
					alerts.NewRemoteAccessDeniedAlert(
						common2.RemoteGoogleTerraform,
						remoteerr.NewResourceListingError(
							status.Error(codes.PermissionDenied, "The caller does not have permission"),
							"google_storage_bucket",
						),
						alerts.EnumerationPhase,
					),
				).Once()
			},
			wantErr: nil,
		},
	}

	providerVersion := "3.78.0"
	resType := resource.ResourceType(googleresource.GoogleStorageBucketResourceType)
	schemaRepository := testresource.InitFakeSchemaRepository("google", providerVersion)
	googleresource.InitResourcesMetadata(schemaRepository)
	factory := terraform3.NewTerraformResourceFactory(schemaRepository)
	deserializer := resource.NewDeserializer(factory)

	for _, c := range cases {
		t.Run(c.test, func(tt *testing.T) {
			shouldUpdate := c.dirName == *goldenfile.Update

			scanOptions := ScannerOptions{Deep: true}
			providerLibrary := terraform3.NewProviderLibrary()
			remoteLibrary := common2.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			if c.setupAlerterMock != nil {
				c.setupAlerterMock(alerter)
			}

			var assetClient *asset.Client
			if !shouldUpdate {
				var err error
				assetClient, err = testgoogle.NewFakeAssetServer(c.response, c.responseErr)
				if err != nil {
					tt.Fatal(err)
				}
			}

			realProvider, err := terraform2.InitTestGoogleProvider(providerLibrary, providerVersion)
			if err != nil {
				tt.Fatal(err)
			}
			provider := terraform2.NewFakeTerraformProvider(realProvider)
			provider.WithResponse(c.dirName)

			// Replace mock by real resources if we are in update mode
			if shouldUpdate {
				ctx := context.Background()
				assetClient, err = asset.NewClient(ctx)
				if err != nil {
					tt.Fatal(err)
				}
				err = realProvider.Init()
				if err != nil {
					tt.Fatal(err)
				}
				provider.ShouldUpdate()
			}

			repo := repository2.NewAssetRepository(assetClient, realProvider.GetConfig(), cache.New(0))

			remoteLibrary.AddEnumerator(google2.NewGoogleStorageBucketEnumerator(repo, factory))
			remoteLibrary.AddDetailsFetcher(resType, common2.NewGenericDetailsFetcher(resType, provider, deserializer))

			testFilter := &enumeration.MockFilter{}
			testFilter.On("IsTypeIgnored", mock.Anything).Return(false)

			s := NewScanner(remoteLibrary, alerter, scanOptions, testFilter)
			got, err := s.Resources()
			assert.Equal(tt, err, c.wantErr)
			if err != nil {
				return
			}
			alerter.AssertExpectations(tt)
			testFilter.AssertExpectations(tt)
			test.TestAgainstGoldenFile(got, resType.String(), c.dirName, provider, deserializer, shouldUpdate, tt)
		})
	}
}

func TestGoogleStorageBucketIAMMember(t *testing.T) {

	cases := []struct {
		test                  string
		dirName               string
		assetRepositoryMock   func(assetRepository *repository2.MockAssetRepository)
		storageRepositoryMock func(storageRepository *repository2.MockStorageRepository)
		responseErr           error
		setupAlerterMock      func(alerter *mocks.AlerterInterface)
		wantErr               error
	}{
		{
			test:    "no storage buckets",
			dirName: "google_storage_bucket_member_empty",
			assetRepositoryMock: func(assetRepository *repository2.MockAssetRepository) {
				assetRepository.On("SearchAllBuckets").Return([]*assetpb.ResourceSearchResult{}, nil)
			},
			wantErr: nil,
		},
		{
			test:    "multiples storage buckets, no bindings",
			dirName: "google_storage_bucket_member_empty",
			assetRepositoryMock: func(assetRepository *repository2.MockAssetRepository) {
				assetRepository.On("SearchAllBuckets").Return([]*assetpb.ResourceSearchResult{
					{
						AssetType:   "storage.googleapis.com/Bucket",
						DisplayName: "dctlgstoragebucketiambinding-1",
					},
					{
						AssetType:   "storage.googleapis.com/Bucket",
						DisplayName: "dctlgstoragebucketiambinding-2",
					},
				}, nil)
			},
			storageRepositoryMock: func(storageRepository *repository2.MockStorageRepository) {
				storageRepository.On("ListAllBindings", "dctlgstoragebucketiambinding-1").Return(map[string][]string{}, nil)
				storageRepository.On("ListAllBindings", "dctlgstoragebucketiambinding-2").Return(map[string][]string{}, nil)
			},
			wantErr: nil,
		},
		{
			test:    "Cannot list bindings",
			dirName: "google_storage_bucket_member_listing_error",
			assetRepositoryMock: func(assetRepository *repository2.MockAssetRepository) {
				assetRepository.On("SearchAllBuckets").Return([]*assetpb.ResourceSearchResult{
					{
						AssetType:   "storage.googleapis.com/Bucket",
						DisplayName: "dctlgstoragebucketiambinding-1",
					},
				}, nil)
			},
			storageRepositoryMock: func(storageRepository *repository2.MockStorageRepository) {
				storageRepository.On("ListAllBindings", "dctlgstoragebucketiambinding-1").Return(
					map[string][]string{},
					errors.New("googleapi: Error 403: driftctl-acc-circle@driftctl-qa-1.iam.gserviceaccount.com does not have storage.buckets.getIamPolicy access to the Google Cloud Storage bucket., forbidden"))
			},
			setupAlerterMock: func(alerter *mocks.AlerterInterface) {
				alerter.On(
					"SendAlert",
					"google_storage_bucket_iam_member",
					alerts.NewRemoteAccessDeniedAlert(
						common2.RemoteGoogleTerraform,
						remoteerr.NewResourceListingError(
							errors.New("googleapi: Error 403: driftctl-acc-circle@driftctl-qa-1.iam.gserviceaccount.com does not have storage.buckets.getIamPolicy access to the Google Cloud Storage bucket., forbidden"),
							"google_storage_bucket_iam_member",
						),
						alerts.EnumerationPhase,
					),
				).Once()
			},
			wantErr: nil,
		},
		{
			test:    "multiples storage buckets, multiple bindings",
			dirName: "google_storage_bucket_member_listing_multiple",
			assetRepositoryMock: func(assetRepository *repository2.MockAssetRepository) {
				assetRepository.On("SearchAllBuckets").Return([]*assetpb.ResourceSearchResult{
					{
						AssetType:   "storage.googleapis.com/Bucket",
						DisplayName: "dctlgstoragebucketiambinding-1",
					},
					{
						AssetType:   "storage.googleapis.com/Bucket",
						DisplayName: "dctlgstoragebucketiambinding-2",
					},
				}, nil)
			},
			storageRepositoryMock: func(storageRepository *repository2.MockStorageRepository) {
				storageRepository.On("ListAllBindings", "dctlgstoragebucketiambinding-1").Return(map[string][]string{
					"roles/storage.admin":        {"user:elie.charra@cloudskiff.com"},
					"roles/storage.objectViewer": {"user:william.beuil@cloudskiff.com"},
				}, nil)

				storageRepository.On("ListAllBindings", "dctlgstoragebucketiambinding-2").Return(map[string][]string{
					"roles/storage.admin":        {"user:william.beuil@cloudskiff.com"},
					"roles/storage.objectViewer": {"user:elie.charra@cloudskiff.com"},
				}, nil)
			},
			wantErr: nil,
		},
	}

	providerVersion := "3.78.0"
	resType := resource.ResourceType(googleresource.GoogleStorageBucketIamMemberResourceType)
	schemaRepository := testresource.InitFakeSchemaRepository("google", providerVersion)
	googleresource.InitResourcesMetadata(schemaRepository)
	factory := terraform3.NewTerraformResourceFactory(schemaRepository)
	deserializer := resource.NewDeserializer(factory)

	for _, c := range cases {
		t.Run(c.test, func(tt *testing.T) {
			repositoryCache := cache.New(100)

			shouldUpdate := c.dirName == *goldenfile.Update

			scanOptions := ScannerOptions{Deep: true}
			providerLibrary := terraform3.NewProviderLibrary()
			remoteLibrary := common2.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			if c.setupAlerterMock != nil {
				c.setupAlerterMock(alerter)
			}

			storageRepo := &repository2.MockStorageRepository{}
			if c.storageRepositoryMock != nil {
				c.storageRepositoryMock(storageRepo)
			}
			var storageRepository repository2.StorageRepository = storageRepo
			if shouldUpdate {
				storageClient, err := storage.NewClient(context.Background())
				if err != nil {
					panic(err)
				}
				storageRepository = repository2.NewStorageRepository(storageClient, repositoryCache)
			}

			assetRepo := &repository2.MockAssetRepository{}
			if c.assetRepositoryMock != nil {
				c.assetRepositoryMock(assetRepo)
			}
			var assetRepository repository2.AssetRepository = assetRepo

			realProvider, err := terraform2.InitTestGoogleProvider(providerLibrary, providerVersion)
			if err != nil {
				tt.Fatal(err)
			}
			provider := terraform2.NewFakeTerraformProvider(realProvider)
			provider.WithResponse(c.dirName)

			remoteLibrary.AddEnumerator(google2.NewGoogleStorageBucketIamMemberEnumerator(assetRepository, storageRepository, factory))

			testFilter := &enumeration.MockFilter{}
			testFilter.On("IsTypeIgnored", mock.Anything).Return(false)

			s := NewScanner(remoteLibrary, alerter, scanOptions, testFilter)
			got, err := s.Resources()
			assert.Equal(tt, c.wantErr, err)
			if err != nil {
				return
			}
			alerter.AssertExpectations(tt)
			testFilter.AssertExpectations(tt)
			test.TestAgainstGoldenFile(got, resType.String(), c.dirName, provider, deserializer, shouldUpdate, tt)
		})
	}
}
