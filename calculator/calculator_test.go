package calculator

import (
	"testing"

	"github.com/nejtr0n/simple-dimple/storage"
	"github.com/stretchr/testify/assert"
)

func TestCalculator(t *testing.T) {
	tests := []struct {
		name  string
		items []struct {
			source string
			value  int
		}
		expected []storage.Summary
	}{
		{
			name:     "zero average with empty data",
			items:    nil,
			expected: nil,
		},
		{
			name: "single value equals avg",
			items: []struct {
				source string
				value  int
			}{
				{
					source: "test",
					value:  1,
				},
			},
			expected: []storage.Summary{
				{
					Source: "test",
					Value:  1,
				},
			},
		},
		{
			name: "multiple value single source",
			items: []struct {
				source string
				value  int
			}{
				{
					source: "test",
					value:  2,
				},
				{
					source: "test",
					value:  4,
				},
			},
			expected: []storage.Summary{
				{
					Source: "test",
					Value:  3,
				},
			},
		},
		{
			name: "multiple value multiple source",
			items: []struct {
				source string
				value  int
			}{
				{
					source: "a",
					value:  2,
				},
				{
					source: "a",
					value:  4,
				},
				{
					source: "b",
					value:  5,
				},
				{
					source: "b",
					value:  6,
				},
			},
			expected: []storage.Summary{
				{
					Source: "a",
					Value:  3,
				},
				{
					Source: "b",
					Value:  5,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewCalculator()
			for _, item := range tt.items {
				calc.Push(item.source, item.value)
			}
			assert.ElementsMatch(t, tt.expected, calc.Calculate())
		})
	}

}
