package customCache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
)

type Redis struct {
	client *redis.Client
	context context.Context
}

type Option struct {
	Host string
	Port int
	Password string
	Db int
	TimeoutMs struct {
		Read int
		Write int
	}
}

func Init(option Option) (*Redis, error) {
	c := redis.NewClient(&redis.Options{
		Addr: option.Host + ":" + strconv.Itoa(option.Port),
		Password: option.Password,
		DB: option.Db,
		WriteTimeout: time.Millisecond * time.Duration(option.TimeoutMs.Write),
		ReadTimeout:  time.Millisecond * time.Duration(option.TimeoutMs.Read),
	})

	ctx := context.Background()

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

