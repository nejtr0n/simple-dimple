package aggregator

import (
	"context"
	"sync"
	"time"

	"github.com/nejtr0n/simple-dimple/calculator"
	"github.com/nejtr0n/simple-dimple/generator"
	"github.com/nejtr0n/simple-dimple/storage"
)

var (
	timeAfter = time.After
)

func NewAggregator(ctx context.Context, duration int) Aggregator {
	return func(wg *sync.WaitGroup, result chan<- storage.Summary) chan<- generator.Message {
		aggregator := make(chan generator.Message)
		go func() {
			defer wg.Done()
			for {
				d, done := aggregateIteration(ctx, aggregator, time.Duration(duration)*time.Second)
				for _, summary := range d {
					result <- summary
				}
				if done {
					return
				}
			}
		}()
		return aggregator
	}
}

type Aggregator func(*sync.WaitGroup, chan<- storage.Summary) chan<- generator.Message

func aggregateIteration(ctx context.Context, pipe <-chan generator.Message, duration time.Duration) (summary []storage.Summary, done bool) {
	calc, timer := calculator.NewCalculator(), timeAfter(duration)
	for {
		select {
		case <-ctx.Done():
			return calc.Calculate(), true
		case <-timer:
			return calc.Calculate(), false
		case message, ok := <-pipe:
			if ok == false {
				return calc.Calculate(), true
			}
			calc.Push(message.Id, message.Value)
		}
	}
}
