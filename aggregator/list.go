package aggregator

import (
	"context"

	"github.com/nejtr0n/simple-dimple/generator"
)

type List []chan<- generator.Message

func (l List) Send(ctx context.Context, message generator.Message) {
	for _, aggregatorCh := range l {
		select {
		case <-ctx.Done():
			return
		case aggregatorCh <- message:
		}
	}
}

func (l List) Close() {
	for _, aggregatorCh := range l {
		close(aggregatorCh)
	}
}
