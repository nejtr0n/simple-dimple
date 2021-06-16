package calculator

import (
	"github.com/nejtr0n/simple-dimple/storage"
)

func NewCalculator() *calculator {
	return &calculator{
		data: make(map[string][]int),
	}
}

type calculator struct {
	data map[string][]int
}

func (c *calculator) Push(source string, value int) {
	c.data[source] = append(c.data[source], value)
}

func (c calculator) Calculate() []storage.Summary {
	var summary []storage.Summary
	for source, items := range c.data {
		summary = append(summary, storage.Summary{
			Source: source,
			Value:  c.average(items),
		})
	}
	return summary
}

func (c calculator) average(items []int) int {
	if len(items) == 0 {
		return 0
	}
	var sum int
	for _, item := range items {
		sum += item
	}
	return sum / len(items)
}
