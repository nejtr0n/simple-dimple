package queue

import (
	"context"
	"sync"

	"github.com/nejtr0n/simple-dimple/generator"
)

type Queue chan generator.Message

func (q Queue) Merge(ctx context.Context, chs ...<-chan generator.Message) {
	var wg sync.WaitGroup
	output := func(c <-chan generator.Message) {
		defer wg.Done()
		for n := range c {
			select {
			case <-ctx.Done():
				return
			case q <- n:
			}
		}
	}
	wg.Add(len(chs))
	for _, c := range chs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(q)
	}()
}
