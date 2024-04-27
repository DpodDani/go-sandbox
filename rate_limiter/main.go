package main

import (
	"log"
	"time"
)

func main() {
	log.Println("Hello world - from ratelimiter")

	tb := NewTokenBucket(
		4,                            // capacity
		time.Duration(time.Second*1), // refill rate
		2,                            // refill amount
	)

	for i := 1; i < 50; i++ {
		_, err := tb.SendRequest(i)
		if err != nil {
			log.Println("Requested failed")
		} else {
			log.Println("Request succeeded")
		}
		time.Sleep(time.Millisecond * 100)
	}
}
