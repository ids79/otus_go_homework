package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	_, isItem := c.Get(key)
	c.Lock()
	defer c.Unlock()
	if !isItem {
		if c.queue.Len() == c.capacity {
			cItem := c.queue.Back().Value.(cacheItem)
			delete(c.items, cItem.key)
			c.queue.Remove(c.queue.Back())
		}
		newElem := cacheItem{key, value}
		c.items[key] = c.queue.PushFront(newElem)
	} else {
		c.items[key].Value = cacheItem{key, value}
	}
	return isItem
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	lItem, isItem := c.items[key]
	if isItem {
		c.queue.MoveToFront(lItem)
		cItem := c.queue.Front().Value.(cacheItem)
		return cItem.value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
