package collections

import "sync"

// Cache is a concurrency safe generic k/v data store of items.
type Cache[K comparable, V any] struct {
	items map[K]V
	mu    sync.RWMutex
}

// NewCache returns a newly initialized Cache object mapping keys of type comparable,
// to their respective values of type any.
func NewCache[K comparable, V any](m map[K]V) Cache[K, V] {
	return Cache[K, V]{
		items: m,
		mu:    sync.RWMutex{},
	}
}

// Get returns an item from the Cache if it is present along with true.
// If the specified item is not present, a nil value is sent along with false.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.items[key]
	if ok {
		return value, true
	}
	return value, false
}

// Set will update an item in the Cache if there is already an existing value.
// Otherwise, a new value will be inserted into the cache.
func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

// SetIfAbsent will add an item to the Cache if there is no existing value for the specified key.
func (c *Cache[K, V]) SetIfAbsent(key K, value V) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.items[key]; ok {
		return false
	}
	c.items[key] = value
	return true
}

// Remove will delete an item from the cache for the specified key.
func (c *Cache[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Keys returns a slice containing the keys for the Cache's data store.
func (c *Cache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]K, len(c.items))

	counter := 0
	for k := range c.items {
		keys[counter] = k
		counter++
	}
	return keys
}

// Has returns true if there is a value with the specified key in the Cache, otherwise false is returned.
func (c *Cache[K, V]) Has(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if _, ok := c.items[key]; ok {
		return true
	}
	return false
}

// Clear purges the current items in the Cache by reinitializing the map.
func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[K]V)
}

// Count returns the number of items currently in the Cache.
func (c *Cache[K, V]) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Items returns a slice containing a copy of the items currently in the Cache.
func (c *Cache[K, V]) Items() []V {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items := make([]V, len(c.items))

	counter := 0
	for _, val := range c.items {
		items[counter] = val
		counter++
	}
	return items
}

// Empty returns true if the Cache is empty or false otherwise.
func (c *Cache[K, V]) Empty() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items) == 0
}

// Pop will return an item from the Cache if it's present and also remove the item from the Cache.
// If an item is removed from the cache true is returned, otherwise false is returned.
func (c *Cache[K, V]) Pop(key K) (V, bool) {
	var (
		value V
		ok    bool
	)

	c.mu.Lock()
	defer c.mu.Unlock()

	if value, ok = c.items[key]; ok {
		delete(c.items, key)
		return value, true
	}
	return value, false
}
