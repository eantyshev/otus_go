package main

import "fmt"

type Item struct {
	value      interface{}
	prev, next *Item
}

func (item Item) Value() interface{} {
	return item.value
}

func (item Item) Prev() *Item {
	return item.prev
}

func (item Item) Next() *Item {
	return item.next
}

type List struct {
	head, tail *Item
	length     int
}

func (l List) Len() int {
	return l.length
}

func (l List) First() *Item {
	return l.head
}

func (l List) Last() *Item {
	return l.tail
}

func (l *List) PushFront(v interface{}) {
	item := Item{}
	item.value = v
	item.next = l.head
	if l.head != nil {
		l.head.prev = &item
	}
	l.head = &item
	if l.tail == nil {
		l.tail = &item
	}
	l.length++
}

func (l *List) PushBack(v interface{}) {
	item := Item{}
	item.value = v
	item.prev = (*l).tail
	if l.tail != nil {
		l.tail.next = &item
	}
	l.tail = &item
	if l.head == nil {
		l.head = &item
	}
	l.length++
}

func (l *List) Remove(pitem *Item) {
	if pitem.prev != nil {
		pitem.prev.next = pitem.next
	}
	if pitem.next != nil {
		pitem.next.prev = pitem.prev
	}
	if l.head == pitem {
		l.head = pitem.next
	}
	if l.tail == pitem {
		l.tail = pitem.prev
	}
	l.length--
}

func main() {
	lst := List{}
	lst.PushFront(1)
	lst.PushFront(2)
	lst.PushBack(3)
	fmt.Println(lst, lst.First(), lst.Last())
}
