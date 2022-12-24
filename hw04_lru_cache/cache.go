package hw04lrucache

import "sync"

type Key string

type Cache[T ItemType] interface {
	Set(key Key, value T) bool
	Get(key Key) (T, bool)
	Clear()
}

type lruCache[T ItemType] struct {
	sync.Mutex
	capacity int
	queue    List[T]
	items    map[Key]*ListItem[T]
}

func (c *lruCache[T]) Set(key Key, value T) bool {
	c.Lock()
	defer c.Unlock()
	_, isItem := c.get(key)
	if !isItem {
		if c.queue.Len() == c.capacity {
			delete(c.items, c.queue.Back().key)
			c.queue.Remove(c.queue.Back())
		}
		cItem := c.queue.PushFront(value)
		cItem.key = key
		c.items[key] = cItem
	} else {
		c.items[key].Value = value
	}
	return isItem
}

func (c *lruCache[T]) Get(key Key) (T, bool) {
	c.Lock()
	defer c.Unlock()
	return c.get(key)
}

func (c *lruCache[T]) get(key Key) (T, bool) {
	lItem, isItem := c.items[key]
	if isItem {
		c.queue.MoveToFront(lItem)
		return c.queue.Front().Value, true
	}
	var t T
	return t, false
}

func (c *lruCache[T]) Clear() {
	c.Lock()
	defer c.Unlock()
	c.items = make(map[Key]*ListItem[T], c.capacity)
	c.queue = NewList[T]()
}

func NewCache[T ItemType](capacity int) Cache[T] {
	return &lruCache[T]{
		capacity: capacity,
		queue:    NewList[T](),
		items:    make(map[Key]*ListItem[T], capacity),
	}
}
