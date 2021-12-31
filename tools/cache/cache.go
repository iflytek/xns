package cache

import (
	"sync"
	"time"
)

type elem struct {
	deadline time.Time
	value    interface{}
}

type Cache struct {
	data map[interface{}]*elem
	mu   sync.Mutex
}

func NewCache()*Cache{
	return &Cache{
		data: map[interface{}]*elem{},
	}
}

func (c *Cache) Set(key, val interface{}, timeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.data[key]
	if ok {
		e.value = val
		e.deadline = time.Now().Add(timeout)
		return
	}
	c.data[key] = &elem{
		deadline: time.Now().Add(timeout),
		value:    val,
	}
}

func (c *Cache) Get(key interface{}) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.data[key]
	if ok {
		if time.Since(e.deadline) > 0 {
			delete(c.data, key)
			return nil, false
		}
		return e.value,true
	}
	return nil,false
}
