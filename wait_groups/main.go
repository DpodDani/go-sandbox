package main

import (
	"context"
	"log"
	"sync"
	"time"
)

func main() {
	log.Println("Hello world")

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Infinitely put apples, oranges and peaches onto their own channels
	fruitSpawner := func(fruit string, stream chan string) {
		for {
			select {
			case <-ctx.Done():
				return
			case stream <- fruit:
			}
		}
	}

	appleStream := make(chan string)
	go fruitSpawner("apple", appleStream)

	orangeStream := make(chan string)
	go fruitSpawner("orange", orangeStream)

	peachStream := make(chan string)
	go fruitSpawner("peach", peachStream)

	// Create one function that spawns 3 more functions that consume from the apple channel
	wg.Add(1)
	go MultipleConsumers(ctx, appleStream, &wg)

	// Create two functions that consume from the oranges and peaches channels respectively
	wg.Add(2)
	go SingleConsumer(ctx, orangeStream, &wg)
	go SingleConsumer(ctx, peachStream, &wg)

	wg.Wait()
}

func MultipleConsumers(parentCtx context.Context, in <-chan string, parentWg *sync.WaitGroup) {
	var wg sync.WaitGroup

	consumerWork := func(ctx context.Context, stream <-chan string) {
		for {
			select {
			case fruit := <-stream:
				log.Println(fruit)
			case <-ctx.Done():
				wg.Done()
				return
			}
		}
	}

	ctx, cancel := context.WithTimeout(parentCtx, time.Second*1)
	defer cancel()

	wg.Add(1)
	go consumerWork(ctx, in)

	wg.Add(1)
	go consumerWork(ctx, in)

	wg.Add(1)
	go consumerWork(ctx, in)

	wg.Wait()
	parentWg.Done()
}

func SingleConsumer(ctx context.Context, in <-chan string, wg *sync.WaitGroup) {
	for {
		select {
		case fruit := <-in:
			log.Println(fruit)
		case <-ctx.Done():
			wg.Done()
			return
		}
	}
}
