package repository

import (
	"context"
	"fmt"
	"github.com/snyk/driftctl/enumeration/remote/azurerm/common"
	"github.com/snyk/driftctl/enumeration/remote/cache"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresql"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PostgresqlRespository interface {
	ListAllServers() ([]*armpostgresql.Server, error)
	ListAllDatabasesByServer(server *armpostgresql.Server) ([]*armpostgresql.Database, error)
}

type postgresqlServersClientImpl struct {
	client *armpostgresql.ServersClient
}

type postgresqlServersClient interface {
	List(context.Context, *armpostgresql.ServersListOptions) (armpostgresql.ServersListResponse, error)
}

func (c postgresqlServersClientImpl) List(ctx context.Context, options *armpostgresql.ServersListOptions) (armpostgresql.ServersListResponse, error) {
	return c.client.List(ctx, options)
}

type postgresqlDatabaseClientImpl struct {
	client *armpostgresql.DatabasesClient
}

type postgresqlDatabaseClient interface {
	ListByServer(context.Context, string, string, *armpostgresql.DatabasesListByServerOptions) (armpostgresql.DatabasesListByServerResponse, error)
}

func (c postgresqlDatabaseClientImpl) ListByServer(ctx context.Context, resGroup string, serverName string, options *armpostgresql.DatabasesListByServerOptions) (armpostgresql.DatabasesListByServerResponse, error) {
	return c.client.ListByServer(ctx, resGroup, serverName, options)
}

type postgresqlRepository struct {
	serversClient  postgresqlServersClient
	databaseClient postgresqlDatabaseClient
	cache          cache.Cache
}

func NewPostgresqlRepository(cred azcore.TokenCredential, options *arm.ClientOptions, config common.AzureProviderConfig, cache cache.Cache) *postgresqlRepository {
	return &postgresqlRepository{
		postgresqlServersClientImpl{client: armpostgresql.NewServersClient(config.SubscriptionID, cred, options)},
		postgresqlDatabaseClientImpl{client: armpostgresql.NewDatabasesClient(config.SubscriptionID, cred, options)},
		cache,
	}
}

func (s *postgresqlRepository) ListAllServers() ([]*armpostgresql.Server, error) {
	cacheKey := "postgresqlListAllServers"

	defer s.cache.Unlock(cacheKey)
	if v := s.cache.GetAndLock(cacheKey); v != nil {
		return v.([]*armpostgresql.Server), nil
	}

	res, err := s.serversClient.List(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	s.cache.Put(cacheKey, res.Value)
	return res.Value, nil
}

func (s *postgresqlRepository) ListAllDatabasesByServer(server *armpostgresql.Server) ([]*armpostgresql.Database, error) {
	res, err := azure.ParseResourceID(*server.ID)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("postgresqlListAllDatabases_%s_%s", res.ResourceGroup, *server.Name)
	if v := s.cache.Get(cacheKey); v != nil {
		return v.([]*armpostgresql.Database), nil
	}

	result, err := s.databaseClient.ListByServer(context.Background(), res.ResourceGroup, *server.Name, nil)
	if err != nil {
		return nil, err
	}

	s.cache.Put(cacheKey, result.Value)
	return result.Value, nil
}
