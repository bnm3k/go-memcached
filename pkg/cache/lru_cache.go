package cache

import (
	"container/list"
	"math"
	"time"
)

// LruCache contains an LRU LruCache
type LruCache struct {
	kv      map[string]*list.Element
	lruList *list.List
	max     int // Max items present, zero for unlimited
}

// node maps a value to a key
type node struct {
	key    string
	value  string
	expire int64 // Unix time
}

// NewLRUCache returns an empty LRUCache
func NewLRUCache(max int) *LruCache {
	if max < 1 {
		max = math.MaxInt64
	}
	c := &LruCache{
		kv:      make(map[string]*list.Element),
		lruList: list.New(),
		max:     max,
	}
	return c
}

// Exists returns true if entry with given key exists, else false
func (c *LruCache) Exists(key string) bool {
	_, exists := c.kv[key]
	return exists
}

// Set entry from given key-value plus add expiry
func (c *LruCache) Set(key string, value string, exptime int) {
	current, exists := c.kv[key]
	var expire int64 = 0
	if exptime > 0 {
		expire = time.Now().Unix() + int64(exptime)
	}

	if exists == false {
		//add new entry
		c.kv[key] = c.lruList.PushFront(&node{
			key:    key,
			value:  value,
			expire: expire,
		})
		if c.lruList.Len() > c.max {
			lruKey := (c.lruList.Back().Value).(*node).key
			c.Delete(lruKey)
		}
	} else {
		//update current entry
		//only update expire val if exptime g.t. 0
		current.Value.(*node).value = value
		if exptime > 0 {
			current.Value.(*node).expire = expire
		}
		c.lruList.MoveToFront(current)
	}
}

// Delete entry with given key
func (c *LruCache) Delete(key string) {
	current, exists := c.kv[key]
	if exists == true {
		c.lruList.Remove(current)
		delete(c.kv, key)
	}
}

// Get a key
func (c *LruCache) Get(key string) (string, bool) {

	current, exists := c.kv[key]
	if exists {
		expire := int64(current.Value.(*node).expire)
		if expire == 0 || expire > time.Now().Unix() {
			c.lruList.MoveToFront(current)
			return current.Value.(*node).value, true
		}
		// remove expired entry instead of returning it
		c.Delete(key)
	}
	return "", false
}
