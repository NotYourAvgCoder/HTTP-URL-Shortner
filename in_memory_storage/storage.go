package inmemorystorage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisDatabase struct {
	client   *redis.Client
	addr     string
	password string
	dbOption int
}

func InitializeRedisDB(addr string, password string, dbOptions int) *RedisDatabase {
	return &RedisDatabase{
		addr:     addr,
		password: password,
		dbOption: dbOptions,
	}
}

func (rd *RedisDatabase) Connect() error {
	log.Info().Msg("trying to connect to redis server")
	client := redis.NewClient(&redis.Options{
		Addr:     rd.addr,
		Password: rd.password, // no password set
		DB:       rd.dbOption, // use default DB
	})

	resp := client.ConfigSet(context.Background(), "save", "5 1")

	if resp.Err() != nil {
		return fmt.Errorf("error while setting save config : %v", resp.Err())
	}

	resp = client.Ping(context.Background())
	if resp.Err() != nil {
		return fmt.Errorf("error while pinging redis server : %v", resp.Err())
	}
	log.Info().Msg("successfully pinged redis server !!!")

	rd.client = client
	return nil
}

func (rd *RedisDatabase) Insert(key string, value string) error {
	status := rd.client.Set(context.Background(), key, value, 0)
	return status.Err()
}

func (rd *RedisDatabase) Get(key string) (string, error) {
	return rd.client.Get(context.Background(), key).Result()
}
