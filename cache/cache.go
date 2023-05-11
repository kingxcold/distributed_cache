package cache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

func New() *Cache {
	return &Cache{data: make(map[string][]byte)}
}

func (c *Cache) Set(key []byte, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[string(key)] = value

	go func() {
		<-time.After(ttl)
		delete(c.data, string(key))
	}()

	return nil
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.data[string(key)]
	if !ok {
		return []byte(""), fmt.Errorf("key (%s) not found\n", key)
	}
	return value, nil
}

func (c *Cache) Has(key []byte) bool {
	c.lock.RLock()
	defer c.lock.Unlock()
	_, ok := c.data[string(key)]
	if !ok {
		return false
	}
	return true
}

func (c *Cache) Delete(key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.data, string(key))
	return nil
}
