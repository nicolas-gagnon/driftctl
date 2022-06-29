package remote

import (
	"testing"

	"github.com/snyk/driftctl/enumeration"
	"github.com/snyk/driftctl/enumeration/remote/alerts"
	"github.com/snyk/driftctl/enumeration/remote/cache"
	common2 "github.com/snyk/driftctl/enumeration/remote/common"
	remoteerr "github.com/snyk/driftctl/enumeration/remote/error"
	google2 "github.com/snyk/driftctl/enumeration/remote/google"
	"github.com/snyk/driftctl/enumeration/remote/google/repository"
	terraform3 "github.com/snyk/driftctl/enumeration/terraform"

	"github.com/snyk/driftctl/enumeration/resource"
	googleresource "github.com/snyk/driftctl/enumeration/resource/google"
	"github.com/snyk/driftctl/mocks"

	testgoogle "github.com/snyk/driftctl/test/google"
	testresource "github.com/snyk/driftctl/test/resource"
	terraform2 "github.com/snyk/driftctl/test/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestGoogleBigtableInstance(t *testing.T) {

	cases := []struct {
		test             string
		assertExpected   func(t *testing.T, got []*resource.Resource)
		response         []*assetpb.Asset
		responseErr      error
		setupAlerterMock func(alerter *mocks.AlerterInterface)
		wantErr          error
	}{
		{
			test:     "no instance",
			response: []*assetpb.Asset{},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 0)
			},
		},
		{
			test: "one instance returned",
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 1)
				assert.Equal(t, "projects/cloudskiff-dev-elie/instances/tf-instance", got[0].ResourceId())
				assert.Equal(t, "google_bigtable_instance", got[0].ResourceType())
			},
			response: []*assetpb.Asset{
				{
					AssetType: "bigtableadmin.googleapis.com/Instance",
					Name:      "//bigtable.googleapis.com/projects/cloudskiff-dev-elie/instances/tf-instance",
					Resource: &assetpb.Resource{
						Data: func() *structpb.Struct {
							v, err := structpb.NewStruct(map[string]interface{}{
								"name": "projects/cloudskiff-dev-elie/instances/tf-instance",
							})
							if err != nil {
								t.Fatal(err)
							}
							return v
						}(),
					},
				},
			},
		},
		{
			test: "one instance without resource data",
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 0)
			},
			response: []*assetpb.Asset{
				{
					AssetType: "bigtableadmin.googleapis.com/Instance",
					Name:      "//bigtable.googleapis.com/projects/cloudskiff-dev-elie/instances/tf-instance",
				},
				{
					AssetType: "bigtableadmin.googleapis.com/Instance",
					Name:      "//bigtable.googleapis.com/projects/cloudskiff-dev-elie/instances/tf-instance-2",
					Resource: &assetpb.Resource{
						Data: func() *structpb.Struct {
							v, err := structpb.NewStruct(map[string]interface{}{})
							if err != nil {
								t.Fatal(err)
							}
							return v
						}(),
					},
				},
			},
		},
		{
			test: "cannot list instances",
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 0)
			},
			responseErr: status.Error(codes.PermissionDenied, "The caller does not have permission"),
			setupAlerterMock: func(alerter *mocks.AlerterInterface) {
				alerter.On(
					"SendAlert",
					"google_bigtable_instance",
					alerts.NewRemoteAccessDeniedAlert(
						common2.RemoteGoogleTerraform,
						remoteerr.NewResourceListingError(
							status.Error(codes.PermissionDenied, "The caller does not have permission"),
							"google_bigtable_instance",
						),
						alerts.EnumerationPhase,
					),
				).Once()
			},
		},
	}

	providerVersion := "3.78.0"
	schemaRepository := testresource.InitFakeSchemaRepository("google", providerVersion)
	googleresource.InitResourcesMetadata(schemaRepository)
	factory := terraform3.NewTerraformResourceFactory(schemaRepository)

	for _, c := range cases {
		t.Run(c.test, func(tt *testing.T) {
			scanOptions := ScannerOptions{}
			providerLibrary := terraform3.NewProviderLibrary()
			remoteLibrary := common2.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			if c.setupAlerterMock != nil {
				c.setupAlerterMock(alerter)
			}

			assetClient, err := testgoogle.NewFakeAssertServerWithList(c.response, c.responseErr)
			if err != nil {
				tt.Fatal(err)
			}

			realProvider, err := terraform2.InitTestGoogleProvider(providerLibrary, providerVersion)
			if err != nil {
				tt.Fatal(err)
			}

			repo := repository.NewAssetRepository(assetClient, realProvider.GetConfig(), cache.New(0))

			remoteLibrary.AddEnumerator(google2.NewGoogleBigTableInstanceEnumerator(repo, factory))

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
			if c.assertExpected != nil {
				c.assertExpected(t, got)
			}
		})
	}
}

func TestGoogleBigtableTable(t *testing.T) {

	cases := []struct {
		test             string
		assertExpected   func(t *testing.T, got []*resource.Resource)
		response         []*assetpb.Asset
		responseErr      error
		setupAlerterMock func(alerter *mocks.AlerterInterface)
		wantErr          error
	}{
		{
			test:     "no table",
			response: []*assetpb.Asset{},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 0)
			},
		},
		{
			test: "one resource returned",
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 1)
				assert.Equal(t, "projects/cloudskiff-dev-elie/instances/tf-instance/tables/tf-table", got[0].ResourceId())
				assert.Equal(t, "google_bigtable_table", got[0].ResourceType())
			},
			response: []*assetpb.Asset{
				{
					AssetType: "bigtableadmin.googleapis.com/Table",
					Name:      "//bigtable.googleapis.com/projects/cloudskiff-dev-elie/instances/tf-instance/tables/tf-table",
					Resource: &assetpb.Resource{
						Data: func() *structpb.Struct {
							v, err := structpb.NewStruct(map[string]interface{}{
								"name": "projects/cloudskiff-dev-elie/instances/tf-instance/tables/tf-table",
							})
							if err != nil {
								t.Fatal(err)
							}
							return v
						}(),
					},
				},
			},
		},
		{
			test: "one resource without resource data",
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 0)
			},
			response: []*assetpb.Asset{
				{
					AssetType: "bigtableadmin.googleapis.com/Table",
					Name:      "//bigtable.googleapis.com/projects/cloudskiff-dev-elie/instances/tf-instance/tables/tf-table",
				},
			},
		},
		{
			test: "cannot list cloud functions",
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 0)
			},
			responseErr: status.Error(codes.PermissionDenied, "The caller does not have permission"),
			setupAlerterMock: func(alerter *mocks.AlerterInterface) {
				alerter.On(
					"SendAlert",
					"google_bigtable_table",
					alerts.NewRemoteAccessDeniedAlert(
						common2.RemoteGoogleTerraform,
						remoteerr.NewResourceListingError(
							status.Error(codes.PermissionDenied, "The caller does not have permission"),
							"google_bigtable_table",
						),
						alerts.EnumerationPhase,
					),
				).Once()
			},
		},
	}

	providerVersion := "3.78.0"
	schemaRepository := testresource.InitFakeSchemaRepository("google", providerVersion)
	googleresource.InitResourcesMetadata(schemaRepository)
	factory := terraform3.NewTerraformResourceFactory(schemaRepository)

	for _, c := range cases {
		t.Run(c.test, func(tt *testing.T) {
			scanOptions := ScannerOptions{}
			providerLibrary := terraform3.NewProviderLibrary()
			remoteLibrary := common2.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			if c.setupAlerterMock != nil {
				c.setupAlerterMock(alerter)
			}

			assetClient, err := testgoogle.NewFakeAssertServerWithList(c.response, c.responseErr)
			if err != nil {
				tt.Fatal(err)
			}

			realProvider, err := terraform2.InitTestGoogleProvider(providerLibrary, providerVersion)
			if err != nil {
				tt.Fatal(err)
			}

			repo := repository.NewAssetRepository(assetClient, realProvider.GetConfig(), cache.New(0))

			remoteLibrary.AddEnumerator(google2.NewGoogleBigtableTableEnumerator(repo, factory))

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
			if c.assertExpected != nil {
				c.assertExpected(t, got)
			}
		})
	}
}
