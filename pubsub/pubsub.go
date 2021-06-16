package pubsub

import (
	"context"
	"sync"

	"github.com/nejtr0n/simple-dimple/aggregator"
	"github.com/nejtr0n/simple-dimple/generator"
	"github.com/nejtr0n/simple-dimple/storage"
)

func NewPubSub() *pubSub {
	return &pubSub{
		aggregators: make(map[string]aggregator.List),
		result:      make(chan storage.Summary),
	}
}

type pubSub struct {
	mu          sync.RWMutex
	wg          sync.WaitGroup
	aggregators map[string]aggregator.List
	result      chan storage.Summary
}

func (q *pubSub) Subscribe(sourceIds []string, handler aggregator.Aggregator) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.wg.Add(len(sourceIds))
	for _, sourceId := range sourceIds {
		q.aggregators[sourceId] = append(q.aggregators[sourceId], handler(&q.wg, q.result))
	}
}

func (q *pubSub) Process(ctx context.Context, queue <-chan generator.Message) <-chan storage.Summary {
	go func() {
		defer close(q.result)
		for message := range queue {
			list, ok := q.aggregators[message.Id]
			if ok {
				list.Send(ctx, message)
			}
		}
		q.close()
		q.wg.Wait()
	}()

	return q.result
}

func (q *pubSub) close() {
	for _, aggregatorChs := range q.aggregators {
		aggregatorChs.Close()
	}
}
