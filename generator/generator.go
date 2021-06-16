package generator

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

var (
	timeAfter  = time.After
	randomizer = randomInt
)

type Worker func(ctx context.Context, duration time.Duration) <-chan Message

func NewWorker(source string, base int, step int) Worker {
	return func(ctx context.Context, duration time.Duration) <-chan Message {
		out := make(chan Message)
		go func() {
			defer close(out)
			generationComplete := timeAfter(duration)
			for {
				select {
				case <-generationComplete:
					return
				case <-ctx.Done():
					return
				default:
					base = randomizer(base, step)
					out <- Message{
						Id:    source,
						Value: base,
					}
				}
			}
		}()
		return out
	}
}

type Generator chan (<-chan Message)

func (g Generator) Add(ctx context.Context, duration int, producer Worker) {
	g <- producer(ctx, time.Duration(duration)*time.Second)
}

func (g Generator) Start(ctx context.Context, sendTimeout int) <-chan Message {
	var wg sync.WaitGroup
	out := make(chan Message)
	output := func(c <-chan Message) {
		defer wg.Done()
		for m := range c {
			select {
			case <-ctx.Done():
				return
			case <-timeAfter(time.Duration(sendTimeout) * time.Second):
				return
			case out <- m:
			}
		}
	}
	wg.Add(len(g))
	for c := range g {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func randomInt(base, step int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(step) + base
}
