package redis_repo

import (
	"encoding/json"
	"fmt"
	"hezzl_t/app/internal/entity"
	"hezzl_t/app/internal/entity/global"
	"hezzl_t/app/internal/repository"
	"time"

	"github.com/go-redis/redis"
)

type cashRepo struct {
	client *redis.Client
}

func NewCashRepo(client *redis.Client) repository.Cash {
	return &cashRepo{
		client: client,
	}
}

func (r *cashRepo) CashGood(key string, value entity.Good) (err error) {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.client.Set(key, data, time.Minute).Err()
	return
}

func (r *cashRepo) FindGood(key string) (g entity.Good, err error) {
	data, err := r.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			err = global.ErrNoData
			return
		}
		return
	}

	err = json.Unmarshal([]byte(data), &g)
	return
}

func (r *cashRepo) DelGood(key string) (err error) {
	err = r.client.Del(key).Err()
	return
}

func (r *cashRepo) CashGoodInfo(key entity.GoodInfoCashKey, value entity.GoodInfo) (err error) {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	keyByte, err := json.Marshal(key)
	if err != nil {
		return err
	}

	keyStr := fmt.Sprintf("%s%s", global.GoodInfoCashKey, string(keyByte))

	err = r.client.Set(keyStr, data, time.Minute).Err()
	return
}

func (r *cashRepo) FindGoodInfo(key entity.GoodInfoCashKey) (data entity.GoodInfo, err error) {
	keyByte, err := json.Marshal(key)
	if err != nil {
		return
	}

	keyStr := fmt.Sprintf("%s%s", global.GoodInfoCashKey, string(keyByte))

	g, err := r.client.Get(keyStr).Result()
	if err != nil {
		if err == redis.Nil {
			err = global.ErrNoData
			return
		}
		return
	}

	err = json.Unmarshal([]byte(g), &data)
	return
}

func (r *cashRepo) DelGoodInfo() (err error) {
	var (
		keys   []string
		cursor uint64
		stop   = false
	)

	for !stop {
		keys, cursor, err = r.client.Scan(cursor, fmt.Sprintf("%s*", global.GoodInfoCashKey), 10).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			err = r.client.Del(key).Err()
			if err != nil {
				return
			}
		}

		if cursor == 0 {
			stop = true
		}
	}
	return
}
