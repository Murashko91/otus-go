package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    *sync.Mutex
}

type ListItemData struct {
	value interface{}
	key   Key
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	item, ok := c.items[key]
	itemData := ListItemData{value, key}
	if ok {
		item.Value = itemData
		c.queue.MoveToFront(item)
	} else {
		if len(c.items) == c.capacity {
			lastItem := c.queue.Back()
			lastItemData := lastItem.Value.(ListItemData)
			lastItemKey := lastItemData.key
			delete(c.items, lastItemKey)
			c.queue.Remove(lastItem)
		}

		c.items[key] = c.queue.PushFront(itemData)
	}
	c.mutex.Unlock()

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()

	item, ok := c.items[key]

	if ok {
		itemData := item.Value.(ListItemData)
		c.queue.MoveToFront(item)
		c.mutex.Unlock()
		return itemData.value, true
	}
	c.mutex.Unlock()

	return nil, false
}

func (c *lruCache) Clear() {
	c.mutex.Lock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.mutex.Unlock()
}

func NewCache(capacity int) Cache {
	var mutex sync.Mutex
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mutex:    &mutex,
	}
}
