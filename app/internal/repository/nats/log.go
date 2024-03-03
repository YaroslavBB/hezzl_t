package nats_repo

import (
	"encoding/json"
	"hezzl_t/app/internal/entity"
	"hezzl_t/app/internal/entity/global"
	"hezzl_t/app/internal/repository"

	"github.com/nats-io/nats.go"
)

type natsLogRepo struct {
	nats *nats.Conn
}

func NewNatsLogRepo(n *nats.Conn) repository.Broker {
	return &natsLogRepo{
		nats: n,
	}
}

func (r *natsLogRepo) SendLog(g entity.Good) (err error) {
	data, err := json.Marshal(g)
	if err != nil {
		return
	}

	err = r.nats.Publish(global.NatsSubj, data)
	return
}
