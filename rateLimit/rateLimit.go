package rateLimit

import (
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
	"log"
	"time"
)

type RateLimit struct {
	client *rate.Limiter
	ctx context.Context
}

func Init(limitPerSec int) *RateLimit {
	return &RateLimit {
		client: rate.NewLimiter(rate.Every(time.Second), limitPerSec),
		ctx: context.Background(),
	}
}

// allow request
func (r RateLimit) Allow() bool {
	return r.client.Allow()
}

// pool request
func (r RateLimit) Stack() {
	if e := r.client.Wait(r.ctx); e != nil {
		log.Println(e)
	}
}