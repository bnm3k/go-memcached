package lfulrt

import (
	"container/list"
	"errors"
	"math"
	"time"
)

/*
A single bucket stores all the keys with an equal number of
frequency access. A dictionary [frequencySet] is used for fast
lookup of random keys. Such lookups occur when a key is being used
hence needs to be bumped up to the next bucket. An accompanying linked-list
is used for popping off the least-recently used key. When a key is added to the bucket,
its also added to the head of the linked list. Hence, at any given instance, the tail
of the list counts as the least-recently-used for that bucket. The dictionary's value
is a pointer to the accompanying linked-list node of the key so that during both
random removal and popping the lru-key, the operations are O(1).
Note that adding is also O(1) since it's adding to the head of the list and the dict
*/
type bucket struct {
	frequencySet map[string]*list.Element
	lruList      *list.List
}

func newBucket() *bucket {
	return &bucket{
		frequencySet: make(map[string]*list.Element),
		lruList:      list.New(),
	}
}

func (b *bucket) add(key string) {
	elem := b.lruList.PushFront(key)
	b.frequencySet[key] = elem
}

func (b *bucket) remove(key string) {
	elem, isPresent := b.frequencySet[key]
	if isPresent {
		b.lruList.Remove(elem)
		delete(b.frequencySet, key)
	}

}

func (b *bucket) isEmpty() bool {
	return len(b.frequencySet) == 0
}

func (b *bucket) popLRU() (string, bool) {
	if len(b.frequencySet) == 0 {
		return "", false
	}
	lastElem := b.lruList.Back()
	key := lastElem.Value.(string)
	b.lruList.Remove(lastElem)
	delete(b.frequencySet, key)
	return key, true
}

/*
LfuLrtCache encompasses both the key-value map and a lfuList that's
used to track the frequencies of use of the keys.
An insertion does not count as use, hence frequency of use at that
point is zero.
On the other hand, an update and get count as a use.
For each use, the key is bumped up to the next bucket. Read section
on bucket to see why this operation is O(1)
Note, bumping up entails deleting from the current bucket and adding
to the next bucket.
On max capacity, the least-frequently-used key plus its value are evicted
However, if multiple keys are in the same frequency bucket, then the
least-recently-used key is evicted. See section on bucket to see how
it keeps track of the lru.
Note however that the eviction operation could potentially be O(f)
where f is the max frequency. In future, in order to mitigate this,
an additional field will be added to the LfuLrtCache struct to keep track
of the index of the bucket where the lfu key should be
*/

type payload struct {
	frequency int
	expire    int64
	value     string
}

//LfuLrtCache ...
type LfuLrtCache struct {
	lfuList []*bucket
	kvStore map[string]payload
	max     int
}

//Constructor ...
func Constructor(max int) *LfuLrtCache {
	if max < 1 {
		max = math.MaxInt64
	}
	lfuList := make([]*bucket, 1)
	lfuList[0] = newBucket()
	return &LfuLrtCache{
		lfuList: lfuList,
		kvStore: make(map[string]payload),
		max:     max,
	}
}

// Exists returns true if entry with given key exists, else false
func (c *LfuLrtCache) Exists(key string) bool {
	_, isPresent := c.kvStore[key]
	return isPresent
}

// returns true and deletes entry if is expired, else false
func (c *LfuLrtCache) checkIfExpired(key string, entry payload) bool {
	if entry.expire != 0 && entry.expire <= time.Now().Unix() {
		c.lfuList[entry.frequency].remove(key)
		delete(c.kvStore, key)
		return true
	}
	return false
}

// Set ...
func (c *LfuLrtCache) Set(key, value string, exptime int) {
	//add new entry
	entry, isPresent := c.kvStore[key]
	var expire int64 = 0
	if exptime > 0 {
		expire = time.Now().Unix() + int64(exptime)
	}
	if isPresent { //is update
		if isExpired := c.checkIfExpired(key, entry); isExpired {
			return
		}
		c.updateFrequency(key, entry)
		// during update, only update expire val if exptime g.t. 0
		if exptime > 0 {
			entry.expire = expire
		}
	} else { //new entry
		c.evictExtra()
		entry.frequency = 0
		entry.expire = expire
		c.lfuList[0].add(key)
	}
	entry.value = value
	c.kvStore[key] = entry

}

func (c *LfuLrtCache) updateFrequency(key string, entry payload) {
	c.lfuList[entry.frequency].remove(key)
	entry.frequency++
	c.kvStore[key] = entry
	if entry.frequency == len(c.lfuList) {
		c.lfuList = append(c.lfuList, newBucket())
	}
	c.lfuList[entry.frequency].add(key)
}

//Get ...
func (c *LfuLrtCache) Get(key string) (string, bool) {
	entry, isPresent := c.kvStore[key]
	if isPresent == false {
		return "", false
	}
	if isExpired := c.checkIfExpired(key, entry); isExpired {
		return "", false
	}
	c.updateFrequency(key, entry)
	return entry.value, true
}

var errorEvicting = errors.New("Error on eviction, incorrect LfuLrtCache state")

func (c *LfuLrtCache) evictExtra() error {
	if len(c.kvStore) >= c.max {
		for _, bucket := range c.lfuList {
			keyToEvict, isNotEmpty := bucket.popLRU()
			if isNotEmpty {
				delete(c.kvStore, keyToEvict)
				return nil
			}
		}
		return errorEvicting
	}
	return nil
}

// Delete entry with given key
func (c *LfuLrtCache) Delete(key string) {
	entry, isPresent := c.kvStore[key]
	if isPresent {
		c.lfuList[entry.frequency].remove(key)
		delete(c.kvStore, key)
	}
}
