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
	defer c.mutex.Unlock()

	item, ok := c.items[key]
	itemData := ListItemData{value, key}

	if ok {
		item.Value = itemData
		c.queue.MoveToFront(item)
	} else {
		if len(c.items) == c.capacity {
			lastItem := c.queue.Back()
			lastItemKey := lastItem.Value.(ListItemData).key
			delete(c.items, lastItemKey)
			c.queue.Remove(lastItem)
		}

		c.items[key] = c.queue.PushFront(itemData)
	}

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.items[key]

	if ok {
		itemData := item.Value.(ListItemData)
		c.queue.MoveToFront(item)
		return itemData.value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mutex:    &sync.Mutex{},
	}
}
