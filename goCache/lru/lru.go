package lru

import "container/list"

//Cache Cache
type Cache struct {
	maxBytes  int64
	nbytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

//Value Value
type Value interface {
	Len() int
}

//New Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//Get look up for key
func (c *Cache) Get(key string) (Value, bool) {
	if el, ok := c.cache[key]; ok {
		c.ll.MoveToFront(el)
		kv := el.Value.(*entry)
		return kv.value, true
	}

	return nil, false
}

//RemoveOldest remove nodes
func (c *Cache) RemoveOldest() {
	el := c.ll.Back()

	if el != nil {
		c.ll.Remove(el)
		kv := el.Value.(*entry)

		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())

		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

//Add add node
func (c *Cache) Add(key string, value Value) {
	if el, ok := c.cache[key]; ok {
		c.ll.MoveToFront(el)

		kv := el.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		el := c.ll.PushFront(&entry{key, value})
		c.cache[key] = el
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.nbytes > c.maxBytes {
		c.RemoveOldest()
	}
}

//Len Len
func (c *Cache) Len() int {
	return c.ll.Len()
}
