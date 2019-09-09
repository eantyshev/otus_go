package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnpack(t *testing.T) {
	assert.Equal(t, Unpack("a4bc2d5e"), "aaaabccddddde")
	assert.Equal(t, Unpack("abcd"), "abcd")
	assert.Equal(t, Unpack("45"), "")
	assert.Equal(t, Unpack("й5л10"), "йййййлллллллллл") //unicode symbols, several digit in a number
	assert.Equal(t, Unpack(`qwe\4\5`), "qwe45")
	assert.Equal(t, Unpack(`qwe\45`), "qwe44444")
	assert.Equal(t, Unpack(`qwe\\5`), `qwe\\\\\`)
	assert.Equal(t, Unpack(`\1qwe\\`), `1qwe\`)
	assert.Equal(t, Unpack(`\12qwe`), `11qwe`)
	assert.Equal(t, Unpack("ф0fa"), "fa")
}
