package cache

import (
	"strconv"
	"sync"
)

//Token ...
type Token string

//Adapter ...
type Adapter struct {
	mu    *sync.Mutex
	cache Cache
}

//Set ...
func (cw *Adapter) Set(key, val, exptimeStr string) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	exptime, err := strconv.Atoi(exptimeStr)
	if err != nil {
		return ClientErrorReply
	}
	cw.cache.Set(key, val, exptime)
	return StoredReply
}

//Add ...
func (cw *Adapter) Add(key, val, exptimeStr string) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	exptime, err := strconv.Atoi(exptimeStr)
	if err != nil {
		return ClientErrorReply
	}
	if cw.cache.Exists(key) {
		return NotStoredReply
	}
	cw.cache.Set(key, val, exptime)
	return StoredReply
}

//Replace ...
func (cw *Adapter) Replace(key, val, exptimeStr string) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	exptime, err := strconv.Atoi(exptimeStr)
	if err != nil {
		return ClientErrorReply
	}
	if cw.cache.Exists(key) {
		cw.cache.Set(key, val, exptime)
		return StoredReply
	}
	return NotStoredReply
}

//Append ...
func (cw *Adapter) Append(key, val, exptimeStr string) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	return cw.appendPrependHelper(key, val, exptimeStr, true)
}

//Prepend ...
func (cw *Adapter) Prepend(key, val, exptimeStr string) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	return cw.appendPrependHelper(key, val, exptimeStr, false)
}

func (cw *Adapter) appendPrependHelper(key, val, exptimeStr string, isAppend bool) Reply {
	currVal, exists := cw.cache.Get(key)
	if exists == false {
		return NotStoredReply
	}
	exptime, err := strconv.Atoi(exptimeStr)
	if err != nil {
		return ClientErrorReply
	}
	if isAppend {
		cw.cache.Set(key, currVal+val, exptime)
	} else { //is prepend
		cw.cache.Set(key, val+currVal, exptime)
	}

	return StoredReply
}

//Increment ...
func (cw *Adapter) Increment(key, numStr string) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	return cw.incrDecrHelper(key, numStr, true)
}

//Decrement ...
func (cw *Adapter) Decrement(key, numStr string) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	return cw.incrDecrHelper(key, numStr, false)
}

//Increment ...
func (cw *Adapter) incrDecrHelper(key, val string, isAddition bool) Reply {
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
func (cw *Adapter) CompareAndSwap(key, val, exptimeStr string, casKey Token) Reply {
	// exptime, err := strconv.Atoi(exptimeStr)
	// if err != nil {
	// 	return ClientErrorReply
	// }
	return NotImplementedReply
}

//Get ...
func (cw *Adapter) Get(key string) (Reply, string) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	val, exists := cw.cache.Get(key)
	if exists == false {
		return NotFoundReply, ""
	}
	return ValueReply, val
}

//GetEntryPlusToken ...
func (cw *Adapter) GetEntryPlusToken(key string) (Reply, string, Token) {
	return NotImplementedReply, "", ""
}

//Delete ...
func (cw *Adapter) Delete(key string) Reply {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	if cw.cache.Exists(key) {
		cw.cache.Delete(key)
		return DeletedReply
	}
	return NotFoundReply
}

//Clear ...
func (cw *Adapter) Clear() Reply {
	return NotImplementedReply
}

//Stats ...
func (cw *Adapter) Stats() Reply {
	return NotImplementedReply
}
