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
	mu       sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	item, has := l.items[key]
	if has {
		l.items[key].Value = itemCache{value: value, key: key}
		l.queue.MoveToFront(l.items[key])
		return true
	}
	if l.capacity == l.queue.Len() {
		item = l.queue.Back()
		delete(l.items, item.Value.(itemCache).key)
		l.queue.Remove(item)
	}
	l.items[key] = l.queue.PushFront(itemCache{value: value, key: key})
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.queue.Len() == 0 {
		return nil, false
	}
	if item, has := l.items[key]; has {
		l.queue.MoveToFront(item)
		return item.Value.(itemCache).value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

type itemCache struct {
	value interface{}
	key   Key
}
