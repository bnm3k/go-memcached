package cache

import (
	"strconv"
	"sync"
)

//Token ...
type Token string

//Wrapper ...
type Wrapper struct {
	mu    *sync.Mutex
	cache Cache
}

//Set ...
func (cw *Wrapper) Set(key, val string, exptime int) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	cw.cache.Set(key, val, exptime)
	return StoredReply
}

//Add ...
func (cw *Wrapper) Add(key, val string, exptime int) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	if cw.cache.Exists(key) {
		return NotStoredReply
	}
	cw.cache.Set(key, val, exptime)
	return StoredReply
}

//Replace ...
func (cw *Wrapper) Replace(key, val string, exptime int) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	if cw.cache.Exists(key) {
		cw.cache.Set(key, val, exptime)
		return StoredReply
	}
	return NotStoredReply
}

//Append ...
func (cw *Wrapper) Append(key, val string, exptime int) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	return cw.appendPrependHelper(key, val, exptime, true)
}

//Prepend ...
func (cw *Wrapper) Prepend(key, val string, exptime int) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	return cw.appendPrependHelper(key, val, exptime, false)
}

func (cw *Wrapper) appendPrependHelper(key, val string, exptime int, isAppend bool) Reply {
	currVal, exists := cw.cache.Get(key)
	if exists == false {
		return NotStoredReply
	}
	if isAppend {
		cw.cache.Set(key, currVal+val, exptime)
	} else { //is prepend
		cw.cache.Set(key, val+currVal, exptime)
	}

	return StoredReply
}

//Increment ...
func (cw *Wrapper) Increment(key, numStr string) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	return cw.incrDecrHelper(key, numStr, true)
}

//Decrement ...
func (cw *Wrapper) Decrement(key, numStr string) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	return cw.incrDecrHelper(key, numStr, false)
}

//Increment ...
func (cw *Wrapper) incrDecrHelper(key, val string, isAddition bool) Reply {
	currVal, exists := cw.cache.Get(key)
	if exists == false {
		return NotFoundReply
	}
	opNum, err := strconv.Atoi(val)
	if err != nil {
		return ClientErrorReply
	}

	valNum, err := strconv.Atoi(currVal)
	if err != nil {
		return ClientErrorReply
	}
	var result int
	if isAddition {
		result = valNum + opNum
	} else {
		result = valNum - opNum
	}
	cw.cache.Set(key, strconv.Itoa(result), 0)
	if err != nil {
		return NotStoredReply
	}
	return StoredReply
}

//CompareAndSwap ...
func (cw *Wrapper) CompareAndSwap(key, val string, exptime int, casKey Token) Reply {
	return NotImplementedReply
}

//Get ...
func (cw *Wrapper) Get(key string) (Reply, string) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	val, exists := cw.cache.Get(key)
	if exists == false {
		return NotFoundReply, ""
	}
	return ValueReply, val
}

//GetEntryPlusToken ...
func (cw *Wrapper) GetEntryPlusToken(key string) (Reply, string, Token) {
	return NotImplementedReply, "", ""
}

//Delete ...
func (cw *Wrapper) Delete(key string) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	if cw.cache.Exists(key) {
		cw.cache.Delete(key)
		return DeletedReply
	}
	return NotFoundReply
}

//Clear ...
func (cw *Wrapper) Clear() Reply {
	return NotImplementedReply
}

//Stats ...
func (cw *Wrapper) Stats() Reply {
	return NotImplementedReply
}
