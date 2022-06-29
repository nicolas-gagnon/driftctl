package remote

import (
	"testing"

	"github.com/snyk/driftctl/enumeration"
	azurerm2 "github.com/snyk/driftctl/enumeration/remote/azurerm"
	repository2 "github.com/snyk/driftctl/enumeration/remote/azurerm/repository"
	"github.com/snyk/driftctl/enumeration/remote/common"
	remoteerr "github.com/snyk/driftctl/enumeration/remote/error"
	"github.com/snyk/driftctl/enumeration/terraform"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresql"
	"github.com/pkg/errors"
	"github.com/snyk/driftctl/enumeration/resource"
	resourceazure "github.com/snyk/driftctl/enumeration/resource/azurerm"
	"github.com/snyk/driftctl/mocks"

	testresource "github.com/snyk/driftctl/test/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAzurermPostgresqlServer(t *testing.T) {

	dummyError := errors.New("this is an error")

	tests := []struct {
		test           string
		mocks          func(*repository2.MockPostgresqlRespository, *mocks.AlerterInterface)
		assertExpected func(t *testing.T, got []*resource.Resource)
		wantErr        error
	}{
		{
			test: "no postgres server",
			mocks: func(repository *repository2.MockPostgresqlRespository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllServers").Return([]*armpostgresql.Server{}, nil)
			},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 0)
			},
		},
		{
			test: "error listing postgres servers",
			mocks: func(repository *repository2.MockPostgresqlRespository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllServers").Return(nil, dummyError)
			},
			wantErr: remoteerr.NewResourceListingError(dummyError, resourceazure.AzurePostgresqlServerResourceType),
		},
		{
			test: "multiple postgres servers",
			mocks: func(repository *repository2.MockPostgresqlRespository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllServers").Return([]*armpostgresql.Server{
					{
						TrackedResource: armpostgresql.TrackedResource{
							Resource: armpostgresql.Resource{
								ID:   to.StringPtr("server1"),
								Name: to.StringPtr("server1"),
							},
						},
					},
					{
						TrackedResource: armpostgresql.TrackedResource{
							Resource: armpostgresql.Resource{
								ID:   to.StringPtr("server2"),
								Name: to.StringPtr("server2"),
							},
						},
					},
				}, nil)
			},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 2)

				assert.Equal(t, got[0].ResourceId(), "server1")
				assert.Equal(t, got[0].ResourceType(), resourceazure.AzurePostgresqlServerResourceType)

				assert.Equal(t, got[1].ResourceId(), "server2")
				assert.Equal(t, got[1].ResourceType(), resourceazure.AzurePostgresqlServerResourceType)
			},
		},
	}

	providerVersion := "2.71.0"
	schemaRepository := testresource.InitFakeSchemaRepository("azurerm", providerVersion)
	resourceazure.InitResourcesMetadata(schemaRepository)
	factory := terraform.NewTerraformResourceFactory(schemaRepository)

	for _, c := range tests {
		t.Run(c.test, func(tt *testing.T) {

			scanOptions := ScannerOptions{}
			remoteLibrary := common.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			fakeRepo := &repository2.MockPostgresqlRespository{}
			c.mocks(fakeRepo, alerter)

			var repo repository2.PostgresqlRespository = fakeRepo

			remoteLibrary.AddEnumerator(azurerm2.NewAzurermPostgresqlServerEnumerator(repo, factory))

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

func TestAzurermPostgresqlDatabase(t *testing.T) {

	dummyError := errors.New("this is an error")

	tests := []struct {
		test           string
		mocks          func(*repository2.MockPostgresqlRespository, *mocks.AlerterInterface)
		assertExpected func(t *testing.T, got []*resource.Resource)
		wantErr        error
	}{
		{
			test: "no postgres database",
			mocks: func(repository *repository2.MockPostgresqlRespository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllServers").Return([]*armpostgresql.Server{}, nil)
			},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 0)
			},
		},
		{
			test: "error listing postgres servers",
			mocks: func(repository *repository2.MockPostgresqlRespository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllServers").Return(nil, dummyError)
			},
			wantErr: remoteerr.NewResourceListingErrorWithType(dummyError, resourceazure.AzurePostgresqlDatabaseResourceType, resourceazure.AzurePostgresqlServerResourceType),
		},
		{
			test: "error listing postgres databases",
			mocks: func(repository *repository2.MockPostgresqlRespository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllServers").Return([]*armpostgresql.Server{
					{
						TrackedResource: armpostgresql.TrackedResource{
							Resource: armpostgresql.Resource{
								ID:   to.StringPtr("/subscriptions/7bfb2c5c-7308-46ed-8ae4-fffa356eb406/resourceGroups/api-rg-pro/providers/Microsoft.DBforPostgreSQL/servers/postgresql-server-8791542"),
								Name: to.StringPtr("postgresql-server-8791542"),
							},
						},
					},
				}, nil).Once()

				repository.On("ListAllDatabasesByServer", mock.IsType(&armpostgresql.Server{})).Return(nil, dummyError).Once()
			},
			wantErr: remoteerr.NewResourceListingError(dummyError, resourceazure.AzurePostgresqlDatabaseResourceType),
		},
		{
			test: "multiple postgres databases",
			mocks: func(repository *repository2.MockPostgresqlRespository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllServers").Return([]*armpostgresql.Server{
					{
						TrackedResource: armpostgresql.TrackedResource{
							Resource: armpostgresql.Resource{
								ID:   to.StringPtr("/subscriptions/7bfb2c5c-7308-46ed-8ae4-fffa356eb406/resourceGroups/api-rg-pro/providers/Microsoft.DBforPostgreSQL/servers/postgresql-server-8791542"),
								Name: to.StringPtr("postgresql-server-8791542"),
							},
						},
					},
				}, nil).Once()

				repository.On("ListAllDatabasesByServer", mock.IsType(&armpostgresql.Server{})).Return([]*armpostgresql.Database{
					{
						ProxyResource: armpostgresql.ProxyResource{
							Resource: armpostgresql.Resource{
								ID:   to.StringPtr("db1"),
								Name: to.StringPtr("db1"),
							},
						},
					},
					{
						ProxyResource: armpostgresql.ProxyResource{
							Resource: armpostgresql.Resource{
								ID:   to.StringPtr("db2"),
								Name: to.StringPtr("db2"),
							},
						},
					},
				}, nil).Once()
			},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 2)

				assert.Equal(t, got[0].ResourceId(), "db1")
				assert.Equal(t, got[0].ResourceType(), resourceazure.AzurePostgresqlDatabaseResourceType)

				assert.Equal(t, got[1].ResourceId(), "db2")
				assert.Equal(t, got[1].ResourceType(), resourceazure.AzurePostgresqlDatabaseResourceType)
			},
		},
	}

	providerVersion := "2.71.0"
	schemaRepository := testresource.InitFakeSchemaRepository("azurerm", providerVersion)
	resourceazure.InitResourcesMetadata(schemaRepository)
	factory := terraform.NewTerraformResourceFactory(schemaRepository)

	for _, c := range tests {
		t.Run(c.test, func(tt *testing.T) {

			scanOptions := ScannerOptions{}
			remoteLibrary := common.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			fakeRepo := &repository2.MockPostgresqlRespository{}
			c.mocks(fakeRepo, alerter)

			var repo repository2.PostgresqlRespository = fakeRepo

			remoteLibrary.AddEnumerator(azurerm2.NewAzurermPostgresqlDatabaseEnumerator(repo, factory))

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
