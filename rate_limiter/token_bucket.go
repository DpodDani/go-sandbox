package main

import (
	"context"
	"fmt"
	"time"
)

type RateLimiter interface {
	SendRequest(req any) (bool, error)
}

type TokenBucket struct {
	Ctx          context.Context
	Capacity     uint16
	RefillRate   time.Duration
	RefillAmount uint16

	availableTokens uint16
}

func NewTokenBucket(ctx context.Context, capacity uint16, refillRate time.Duration, refillAmount uint16) *TokenBucket {
	tb := &TokenBucket{
		Ctx:             ctx,
		Capacity:        capacity,
		RefillRate:      refillRate,
		RefillAmount:    refillAmount,
		availableTokens: 0,
	}
	go tb.run()

	// wait until token bucket is filled up
	for {
		if tb.availableTokens == tb.Capacity {
			break
		}
	}

	return tb
}

func (tb *TokenBucket) refillBucket() {
	if tb.availableTokens+tb.RefillAmount <= tb.Capacity {
		tb.availableTokens += tb.RefillAmount
	}
	time.Sleep(tb.RefillRate)
}

func (tb *TokenBucket) run() {
	for {
		select {
		case <-tb.Ctx.Done():
			return
		default:
			tb.refillBucket()
		}
	}
}

func (tb *TokenBucket) SendRequest(req any) (bool, error) {
	if tb.availableTokens == 0 {
		return false, fmt.Errorf("no available tokens. Rejecting request: %+v", req)
	}
	tb.availableTokens -= 1
	return true, nil
}
