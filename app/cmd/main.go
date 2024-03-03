package main

import (
	"context"
	"fmt"
	"hezzl_t/app/config"
	"hezzl_t/app/external"
	repository_import "hezzl_t/app/internal/repository/import"
	"hezzl_t/app/internal/usecase"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

func main() {
	fmt.Println("running")

	ctx := context.Background()

	config, err := config.NewConfig(os.Getenv("CONF_PATH"))
	if err != nil {
		log.Fatalln(err)
	}

	postgresDB, err := sqlx.Open("postgres", config.DbConn())
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDB.Close()

	clickhouseDb, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{config.ClockhouseConn()},
		Auth: clickhouse.Auth{
			Database: config.Clickhouse.Database,
			Username: config.Clickhouse.User,
			Password: config.Clickhouse.Password,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer clickhouseDb.Close()

	err = clickhouseDb.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(os.Getenv("MIGRATION_LOCK")); err != nil {
		if err := exec.Command("bash", os.Getenv("INIT_PG_DB")).Run(); err != nil {
			log.Fatal(err)
		}

		if err := InitClickhouseDb(ctx, clickhouseDb); err != nil {
			log.Fatal(err)
		}

		_, err = os.Create(os.Getenv("MIGRATION_LOCK"))
		if err != nil {
			log.Fatal(err)
		}
	}

	redis := redis.NewClient(&redis.Options{
		Addr: config.RedisConn(),
	})
	defer redis.Close()

	_, err = redis.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	natsConn, err := nats.Connect(config.NastConn())
	if err != nil {
		log.Fatal(err)
	}
	defer natsConn.Close()

	i := repository_import.NewRepositoryImports(redis, natsConn, clickhouseDb, ctx)
	goodUsecase := usecase.NewGoodUsecase(i)
	handler := external.NewGoodHandlers(*goodUsecase, postgresDB)
	logUsecase := usecase.NewLoggerUsecase(i, natsConn, ctx)

	go logUsecase.Log()

	mux := http.NewServeMux()
	mux.HandleFunc("/good/create", handler.CreateGood)
	mux.HandleFunc("/good/update", handler.UpdateGood)
	mux.HandleFunc("/good", handler.FindGood)
	mux.HandleFunc("/good/remove", handler.RemoveGood)
	mux.HandleFunc("/goods/list", handler.LoadGoodInfo)
	mux.HandleFunc("/good/reprioritiize", handler.UpdatePriority)

	http.ListenAndServe(config.ServerIP(), mux)
}

func InitClickhouseDb(ctx context.Context, clickhouseDb driver.Conn) (err error) {
	if err = clickhouseDb.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS logs(
		ID INT,
		ProjectID INT,
		Name String,
		Description String,
		Priority INT,
		Removed Boolean,
		EventTime DateTime
	) engine = MergeTree()
	PRIMARY KEY (ID); 
		`); err != nil {
		return
	}

	if err = clickhouseDb.Exec(ctx, `
	ALTER TABLE logs
    ADD INDEX id_index (ID) TYPE minmax;
		`); err != nil {
		return
	}

	if err = clickhouseDb.Exec(ctx, `
	ALTER TABLE logs
	ADD INDEX project_id_index (ProjectID) TYPE minmax;
		`); err != nil {
		return
	}

	if err = clickhouseDb.Exec(ctx, `
	ALTER TABLE logs
	ADD INDEX name_index (Name) TYPE minmax;
		`); err != nil {
		return
	}

	return
}
