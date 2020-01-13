package lfu

import (
	"math"
	"time"
)

type payload struct {
	frequency int
	value     string
	expire    int64 // Unix time
}

var emptyVal struct{}

type set map[string]struct{}

//LfuCache ...
type LfuCache struct {
	lfuList  []set
	kvStore  map[string]payload
	capacity int
}

//Constructor ...
func Constructor(capacity int) LfuCache {
	if capacity < 1 {
		capacity = math.MaxInt64
	}
	lfuList := make([]set, 1)
	lfuList[0] = make(map[string]struct{})
	return LfuCache{
		lfuList:  lfuList,
		kvStore:  make(map[string]payload),
		capacity: capacity,
	}
}

// Exists returns true if entry with given key exists, else false
func (c *LfuCache) Exists(key string) bool {
	_, isPresent := c.kvStore[key]
	return isPresent
}

// Set entry from given key-value plus add expiry
func (c *LfuCache) Set(key, value string, exptime int) {
	//add new entry
	entry, isPresent := c.kvStore[key]
	var expire int64 = 0
	if exptime > 0 {
		expire = time.Now().Unix() + int64(exptime)
	}
	if isExpired := c.checkIfExpired(key, entry); isExpired {
		return
	}
	if isPresent { //is update
		c.updateFrequency(key, entry)
		// during update, only update expire val if exptime g.t. 0
		if exptime > 0 {
			entry.expire = expire
		}
	} else { //new entry
		c.evictExtra()
		entry.frequency = 0
		entry.expire = expire
		c.lfuList[0][key] = emptyVal
	}
	entry.value = value
	c.kvStore[key] = entry
}

func (c *LfuCache) evictExtra() {
	if len(c.kvStore) >= c.capacity {
		for _, bucket := range c.lfuList {
			if len(bucket) == 0 {
				continue
			}
			for keyToEvict := range bucket {
				delete(bucket, keyToEvict)
				delete(c.kvStore, keyToEvict)
				return
			}
		}
	}
}

// Get entry by given key
func (c *LfuCache) Get(key string) (string, bool) {
	entry, isPresent := c.kvStore[key]
	isExpired := c.checkIfExpired(key, entry)
	if isPresent == false || isExpired {
		return "", false
	}
	c.updateFrequency(key, entry)
	return entry.value, true
}

// returns true and deletes entry if is expired, else false
func (c *LfuCache) checkIfExpired(key string, entry payload) bool {
	if entry.expire != 0 && entry.expire <= time.Now().Unix() {
		bucket := c.lfuList[entry.frequency]
		delete(bucket, key)
		delete(c.kvStore, key)
		return true
	}
	return false
}

func (c *LfuCache) updateFrequency(key string, entry payload) {
	delete(c.lfuList[entry.frequency], key)
	entry.frequency++
	c.kvStore[key] = entry
	if entry.frequency == len(c.lfuList) {
		c.lfuList = append(c.lfuList, make(map[string]struct{}))
	}
	c.lfuList[entry.frequency][key] = emptyVal
}

// Delete entry with given key
func (c *LfuCache) Delete(key string) {
	entry, isPresent := c.kvStore[key]
	if isPresent == true {
		bucket := c.lfuList[entry.frequency]
		delete(bucket, key)
		delete(c.kvStore, key)
	}
}
