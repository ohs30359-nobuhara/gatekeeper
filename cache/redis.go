package customCache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type Redis struct {
	client *redis.Client
	context context.Context
}

func Init(ctx context.Context) (*Redis, error) {
	c := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
		WriteTimeout: time.Millisecond * 1,
	})

	if e := c.Ping(ctx).Err(); e != nil {
		log.Println(e.Error())
		return nil, e
	}

	return &Redis{
		client: c,
		context: ctx,
	}, nil
}

func (r Redis) Get(k string) (string, bool) {
	v := r.client.Get(r.context, k).Val()

	if len(v) == 0 {
		return v, false
	}

	return v, true
}

func (r Redis) Set(k string, v string, ex time.Duration) {
	e := r.client.Set(r.context, k, v, ex).Err()

	if e != nil {
		log.Fatal(e)
	}
}

