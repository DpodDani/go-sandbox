package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenBucketRefill(t *testing.T) {
	tb := TokenBucket{
		Capacity:        4,
		RefillRate:      time.Duration(time.Millisecond * 50),
		RefillAmount:    2,
		availableTokens: 0,
	}

	require.Equal(t, uint16(0), tb.availableTokens)

	go tb.run()

	// premature check to ensure refill rate has begun
	time.Sleep(time.Millisecond * 10)
	assert.True(t, tb.availableTokens > uint16(1))
	assert.NotEqual(t, tb.Capacity, tb.availableTokens)

	// 2nd check to ensure full refill has completed
	time.Sleep(time.Millisecond * 100)
	require.Equal(t, tb.Capacity, tb.availableTokens)
}

func TestTokenBucketAvailableTokenDecrement(t *testing.T) {
	tb := TokenBucket{
		Capacity:        4,
		RefillRate:      time.Duration(time.Second * 1),
		RefillAmount:    2,
		availableTokens: 4,
	}

	go tb.run()

	tb.SendRequest(struct{}{})
	require.Equal(t, uint16(3), tb.availableTokens)

	time.Sleep(time.Second * 1)
	assert.Equal(t, uint16(3), tb.availableTokens, "Expected no change to available tokens.")

	tb.SendRequest(struct{}{})
	require.Equal(t, uint16(2), tb.availableTokens)

	time.Sleep(time.Millisecond * 1500) // give time for refill to happen
	assert.Equal(t, tb.Capacity, tb.availableTokens, "Expected available tokens to match capacity")
}
