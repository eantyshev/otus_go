package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const N = 2

func TestTopN(t *testing.T) {
	cases := []struct {
		Text   string
		Result []WordStat
	}{
		{"qwe qwe\nother", []WordStat{{"qwe", 2}, {"other", 1}}},
		// more than N distinct words
		{"1 1 1 2 2 3", []WordStat{{"1", 3}, {"2", 2}}},
		// some unicode & empty separators
		{" \t йцук\nйцук\n  \n\n\n", []WordStat{{"йцук", 2}}},
	}
	for _, c := range cases {
		assert.Equal(t, c.Result, TopN(c.Text, 2))
	}
}
