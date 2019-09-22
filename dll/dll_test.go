package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItem(t *testing.T) {
	x := Item{123, nil, nil}
	y := Item{321, &x, nil}
	x.next = &y
	assert.Equal(t, y.prev, &x)
}

func TestNil(t *testing.T) {
	lst := List{}
	assert.Equal(t, lst.Len(), 0)
	t.Log(lst.First(), lst.Last())
	assert.Nil(t, lst.First())
	assert.Nil(t, lst.Last())
}

func TestPush(t *testing.T) {
	lst := List{}
	lst.PushFront(1)
	lst.PushBack(2)
	lst.PushFront(3)
	first := lst.First()
	middle := first.Next()
	last := lst.Last()
	assert.Equal(t, first.Value(), 3)
	assert.Equal(t, middle.Value(), 1)
	assert.Equal(t, last.Value(), 2)
	assert.Equal(t, middle.Next(), lst.Last())
	assert.Equal(t, middle.Prev(), lst.First())
}

func TestRemove(t *testing.T) {
	lst := List{}
	lst.PushBack(1)
	// remove the only item
	lst.Remove(lst.First())
	assert.Nil(t, lst.First())
	assert.Nil(t, lst.Last())
	lst.PushBack(2)
	pitem := lst.First()
	lst.PushFront(3)
	lst.PushBack(4)
	// remove from the middle
	lst.Remove(pitem)
	assert.Equal(t, (*lst.First()).Value(), 3)
	assert.Equal(t, (*lst.Last()).Value(), 4)
	assert.Equal(t, lst.Len(), 2)
}
