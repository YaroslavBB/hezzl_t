package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"hezzl_t/app/internal/entity"
	"hezzl_t/app/internal/entity/global"
	repository_import "hezzl_t/app/internal/repository/import"

	"github.com/nats-io/nats.go"
)

type LoggerUsecase struct {
	ri  *repository_import.RepositoryImports
	n   *nats.Conn
	ctx context.Context
}

func NewLoggerUsecase(ri *repository_import.RepositoryImports, n *nats.Conn, ctx context.Context) *LoggerUsecase {
	return &LoggerUsecase{
		ri:  ri,
		n:   n,
		ctx: ctx,
	}
}

func (u *LoggerUsecase) Log() error {
	var (
		buffer = make([]entity.Good, 0, 10)
		g      entity.Good
	)

	for {
		sub, err := u.n.SubscribeSync(global.NatsSubj)
		if err != nil {
			return err
		}

		defer sub.Unsubscribe()

		msg, err := sub.NextMsgWithContext(u.ctx)
		if err != nil {
			return err
		}

		err = json.Unmarshal(msg.Data, &g)
		if err != nil {
			return err
		}

		buffer = append(buffer, g)

		if len(buffer) > 10 {
			for _, item := range buffer {
				err = u.ri.Log.LogGood(item)
				if err != nil {
					fmt.Println(err)
					return err
				}
			}

			buffer = make([]entity.Good, 0, 10)
		}
	}
}
