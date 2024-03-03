package repository_import

import (
	"context"
	"hezzl_t/app/internal/repository"
	clickhouse_repo "hezzl_t/app/internal/repository/clickhouse"
	nats_repo "hezzl_t/app/internal/repository/nats"
	"hezzl_t/app/internal/repository/postgres"
	redis_repo "hezzl_t/app/internal/repository/redis"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/go-redis/redis"
	"github.com/nats-io/nats.go"
)

type RepositoryImports struct {
	repository.Good
	repository.Cash
	repository.Broker
	repository.Log
}

func NewRepositoryImports(client *redis.Client, n *nats.Conn, conn driver.Conn, ctx context.Context) *RepositoryImports {
	return &RepositoryImports{
		postgres.NewGoodRepo(),
		redis_repo.NewCashRepo(client),
		nats_repo.NewNatsLogRepo(n),
		clickhouse_repo.NewLogRepo(conn, ctx),
	}
}
