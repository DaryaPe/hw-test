package hw04lrucache

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
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l lruCache) Set(key Key, value interface{}) bool {
	if l.items[key] != nil {
		l.items[key].Value = item{value: value, key: key}
		l.queue.MoveToFront(l.items[key])
		return true
	}
	if l.capacity == l.queue.Len() {
		oldKey := l.queue.Back().Value.(item).key
		delete(l.items, oldKey)
		l.queue.Remove(l.queue.Back())
	}
	l.items[key] = l.queue.PushFront(item{value: value, key: key})
	return false
}

func (l lruCache) Get(key Key) (interface{}, bool) {
	if l.items[key] == nil {
		return nil, false
	}
	l.queue.MoveToFront(l.items[key])
	return l.items[key].Value.(item).value, true
}

func (l lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

type item struct {
	value interface{}
	key   Key
}
