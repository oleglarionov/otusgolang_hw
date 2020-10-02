package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	items    map[Key]*cacheItem
	queue    List
	mux      sync.Mutex
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mux.Lock()
	defer l.mux.Unlock()

	_, keyExists := l.items[key]
	if keyExists {
		l.queue.Remove(l.items[key].qi)
	} else if len(l.items) == l.capacity {
		backQi := l.queue.Back()
		delete(l.items, backQi.Value.(*cacheItem).key)
		l.queue.Remove(backQi)
	}

	i := &cacheItem{
		value: value,
		key:   key,
	}
	l.queue.PushFront(i)
	i.qi = l.queue.Front()
	l.items[key] = i

	return keyExists
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mux.Lock()
	defer l.mux.Unlock()

	i, keyExists := l.items[key]
	if !keyExists {
		return nil, false
	}
	l.queue.MoveToFront(i.qi)
	return i.value, true
}

func (l *lruCache) Clear() {
	l.mux.Lock()
	defer l.mux.Unlock()

	l.items = make(map[Key]*cacheItem, l.capacity)
	l.queue = NewList()
}

type cacheItem struct {
	value interface{}
	key   Key
	qi    *listItem
}

func NewCache(capacity int) Cache {
	items := make(map[Key]*cacheItem, capacity)
	queue := NewList()

	return &lruCache{
		capacity: capacity,
		items:    items,
		queue:    queue,
	}
}
